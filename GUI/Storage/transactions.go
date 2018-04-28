package Storage

import (
	"bytes"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC"
)

var Amount *Numbers.RawAmount

type HistoryStore []RPCClient.SingleHistory

var History HistoryStore

func (h *HistoryStore) Set(hist []RPCClient.SingleHistory) {
	*h = hist
}

func (h *HistoryStore) ExistHash(hash Block.BlockHash) bool {
	for _, blk := range *h {
		if bytes.Equal(blk.Hash, hash) {
			return true
		}
	}

	return false
}

func (h *HistoryStore) AlreadyReceived(hash Block.BlockHash) bool {
	for _, blk := range *h {
		if bytes.Equal(blk.Source, hash) {
			return true
		}
	}

	return false
}

func (h *HistoryStore) Add(blk Block.BlockTransaction, amount *Numbers.RawAmount) {
	hist := RPCClient.SingleHistory{}
	hist.Type = blk.GetSubType()
	hist.Amount = amount
	hist.Destination, hist.Source = blk.GetTarget()
	hist.Hash = blk.Hash()

	*h = append([]RPCClient.SingleHistory{hist}, *h...)
}

func (h *HistoryStore) Next(page uint32) {
	//@TODO pagination
}
