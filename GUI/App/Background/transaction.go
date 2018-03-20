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
	ammount   *Numbers.RawAmount
}

type queueup chan txs

var queue = make(queueup, 32)

func PublishBlocksToQueue(blk []Block.BlockTransaction, ammount ...*Numbers.RawAmount) error {
	var rt = make(WaitConfirmation)
	defer close(rt)

	if len(ammount) == 0 {
		ammount[0], _ = Numbers.NewRawFromString("0")
	}

	queue <- txs{
		blocks:    blk,
		returning: rt,
		ammount:   ammount[0],
	}

	return <-rt
}

func PublishBlockToQueue(blk Block.BlockTransaction, ammount ...*Numbers.RawAmount) error {
	return PublishBlocksToQueue([]Block.BlockTransaction{blk}, ammount...)
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
				balance = Storage.Amount.Subtract(tx.ammount)
			case "open":
				balance = Storage.Amount.Add(tx.ammount)
			case "receive":
				balance = Storage.Amount.Add(tx.ammount)
			default:
				balance = Storage.Amount
			}

			blk.SetFrontier(Storage.Frontier)
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

			Storage.SetAmount(balance)
			Storage.SetFrontier(hash)
			Storage.History.Add(blk.GetType(), tx.ammount, hash, "") // @TODO Set Address
		}

		tx.returning <- geterr
	}

}
