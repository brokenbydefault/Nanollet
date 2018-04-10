package Background

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
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

	if len(amount) == 0 {
		amount[0], _ = Numbers.NewRawFromString("0")
	}

	queue <- txs{
		blocks:    blk,
		returning: rt,
		amount:    amount[0],
	}

	return <-rt
}

func PublishBlockToQueue(blk Block.BlockTransaction, amount ...*Numbers.RawAmount) error {
	return PublishBlocksToQueue([]Block.BlockTransaction{blk}, amount...)
}

func StartTransaction() {
	go execute()
}

func execute() {

	for tx := range queue {
		var geterr error = nil

		for _, blk := range tx.blocks {
			var err error = nil
			var balance *Numbers.RawAmount

			switch blk.GetType() {
			case "send":
				balance = Storage.Amount.Subtract(tx.amount)
			case "open":
				balance = Storage.Amount.Add(tx.amount)
			case "receive":
				balance = Storage.Amount.Add(tx.amount)
			default:
				balance = Storage.Amount
			}

			blk.SetFrontier(Storage.Frontier)
			blk.SetWork(Storage.PrecomputedPoW)
			blk.SetBalance(balance)

			hash := blk.Hash()
			sig, err := Storage.SK.CreateSignature(hash)
			if err != nil {
				geterr = err
				break
			}
			blk.SetSignature(sig)

			_, err = RPCClient.BroadcastBlock(Connectivity.Socket, blk)
			if err != nil {
				geterr = err
				break
			}

			Storage.Amount = balance
			Storage.History.Add(blk, tx.amount)
			Storage.Frontier = hash

			if len(queue) == 0 {
				go Storage.UpdateFrontier(hash)
			}
		}

		tx.returning <- geterr
	}

}
