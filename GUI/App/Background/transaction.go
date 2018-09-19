package Background

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Node"
	"time"
	"errors"
)

var (
	ErrInsufficientVotes = errors.New("insufficient votes")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInvalidFrontier   = errors.New("invalid frontier")
)

type WaitConfirmation chan error

type txs struct {
	blocks    []Block.Transaction
	blockType Block.BlockType
	returning WaitConfirmation
	amount    *Numbers.RawAmount
}

var queue = make(chan txs, 64)

func init() {
	go listen()
}

func PublishBlocksToQueue(blk []Block.Transaction, blockType Block.BlockType, amounts ...*Numbers.RawAmount) error {
	var waitchan = make(WaitConfirmation)
	defer close(waitchan)

	// amount is optional, set default of 0 if missing
	var amount *Numbers.RawAmount
	if len(amounts) == 0 {
		amount = Numbers.NewMin()
	} else {
		amount = amounts[0]
	}

	queue <- txs{
		blocks:    blk,
		returning: waitchan,
		blockType: blockType,
		amount:    amount,
	}

	return <-waitchan
}

func PublishBlockToQueue(blk Block.Transaction, blockType Block.BlockType, amount ...*Numbers.RawAmount) error {
	return PublishBlocksToQueue([]Block.Transaction{blk}, blockType, amount...)
}

func listen() {

	for tx := range queue {
		var err error
		for _, blk := range tx.blocks {
			err = processBlock(blk, tx.blockType, tx.amount)
			if err != nil {
				break
			}
		}

		tx.returning <- err
	}

}

func processBlock(tx Block.Transaction, blockType Block.BlockType, amm *Numbers.RawAmount) error {
	var err error = nil
	var balance *Numbers.RawAmount

	switch blockType {
	case Block.Send:
		balance = Storage.AccountStorage.Balance.Subtract(amm)
	case Block.Open:
		balance = Storage.AccountStorage.Balance.Add(amm)
	case Block.Receive:
		balance = Storage.AccountStorage.Balance.Add(amm)
	default:
		balance = Storage.AccountStorage.Balance
	}

	if !balance.IsValid() {
		return ErrInvalidAmount
	}

	tx.SetFrontier(Storage.AccountStorage.Frontier)
	//@TODO Support pre-computed PoW, again.
	//blk.SetWork(Storage.RetrievePrecomputedPoW())

	tx.Work()
	tx.SetBalance(balance)

	hash := tx.Hash()
	sig, err := Storage.AccountStorage.SecretKey.Sign(hash[:])
	if err != nil {
		return err
	}

	tx.SetSignature(sig)

	if err = Node.PostBlock(Connection, tx); err != nil {
		return err
	}

	Storage.TransactionStorage.Add(tx)

	//@TODO improve if not reach the quorum
	if !waitVotesConfirmation(tx, 2*time.Minute) {
		Storage.TransactionStorage.Remove(tx)
		return ErrInsufficientVotes
	}

	Storage.AccountStorage.Balance = tx.GetBalance()
	Storage.AccountStorage.Frontier = tx.Hash()

	return nil
}

func waitVotesConfirmation(tx Block.Transaction, duration time.Duration) bool {
	start := time.Now()
	hash := tx.Hash()

	for range time.Tick(2 * time.Second) {
		if Storage.TransactionStorage.IsConfirmed(&hash, &Storage.Configuration.Account.Quorum) {
			return true
		}

		Node.RequestVotes(Connection, tx)

		if time.Since(start) > duration {
			break
		}
	}

	return false
}
