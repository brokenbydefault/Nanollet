package Storage

import (
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

type HistoryStore []RPCClient.SingleHistory

var History HistoryStore

func (h *HistoryStore) Set(hist []RPCClient.SingleHistory) {
	History = hist
}

func (h *HistoryStore) Add(types string, amount *Numbers.RawAmount, hash []byte, account Wallet.Address) {
	if types == "change" {
		return
	}

	hist := RPCClient.SingleHistory{}
	hist.Type = types
	hist.Amount = amount
	hist.Hash = hash
	hist.Account = account
	History = append([]RPCClient.SingleHistory{hist}, History...)
}

func (h *HistoryStore) Next(page uint32) () {
	//@TODO pagination
}
