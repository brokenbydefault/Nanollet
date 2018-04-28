package Background

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
)

type WaitConfirmation chan error

type txs struct {
	blocks    []Block.BlockTransaction
	returning WaitConfirmation
	amount    *Numbers.RawAmount
}

type queueup chan txs

var queue = make(queueup, 32)

func PublishBlocksToQueue(blk []Block.BlockTransaction, amount ...*Numbers.RawAmount) error {
	var rt = make(WaitConfirmation)
	defer close(rt)

	// amount is optional, set default of 0 if missing
	var amm *Numbers.RawAmount
	if len(amount) == 0 {
		amm, _ = Numbers.NewRawFromString("0")
	} else {
		amm = amount[0]
	}

	queue <- txs{
		blocks:    blk,
		returning: rt,
		amount:    amm,
	}

	return <-rt
}

func PublishBlockToQueue(blk Block.BlockTransaction, amount ...*Numbers.RawAmount) error {
	return PublishBlocksToQueue([]Block.BlockTransaction{blk}, amount...)
}

func StartTransaction() {
	go listen()
}

func listen() {

	for tx := range queue {
		var err error
		for _, blk := range tx.blocks {
			err = processBlock(blk, tx.amount)
			if err != nil {
				break
			}
		}

		tx.returning <- err
	}

}

func processBlock(blk Block.BlockTransaction, amm *Numbers.RawAmount) error {
	defer func() {
		go Storage.UpdatePoW()
	}()

	var err error = nil
	var balance *Numbers.RawAmount

	switch blk.GetSubType() {
	case Block.Send:
		balance = Storage.Amount.Subtract(amm)
	case Block.Open:
		balance = Storage.Amount.Add(amm)
	case Block.Receive:
		balance = Storage.Amount.Add(amm)
	default:
		balance = Storage.Amount
	}

	blk.SetFrontier(Storage.Frontier)
	blk.SetWork(Storage.RetrievePrecomputedPoW())
	blk.SetBalance(balance)

	hash := blk.Hash()
	sig, err := Storage.SK.CreateSignature(hash)
	if err != nil {
		return err
	}
	blk.SetSignature(sig)

	_, err = RPCClient.BroadcastBlock(Connectivity.Socket, blk)
	if err != nil {
		return err
	}

	Storage.Amount = balance
	Storage.History.Add(blk, amm)
	Storage.UpdateFrontier(blk)

	return nil
}
