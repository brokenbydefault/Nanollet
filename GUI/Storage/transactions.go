package Storage

import (
	"bytes"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC"
	"math/rand"
	"time"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var Amount *Numbers.RawAmount

var TransactionStorage TransactionsShelter

type TransactionsShelter struct {
	Unconfirmed TransactionBox
	Confirmed   TransactionBox
	Pending     TransactionBox
	Votes       Peer.Votes
}

type Transaction struct {
	Block.Transaction
	date time.Time
}

type TransactionBox struct {
	list      map[string]*Transaction
	listeners []chan *Transaction
}

func NewTransactionsShelter(pk Wallet.PublicKey) (ts *TransactionsShelter) {
	ts = new(TransactionsShelter)
	ts.KeepUpdated(pk)

	return ts
}

func (ts *TransactionsShelter) KeepUpdated(pk Wallet.PublicKey) {
	for range time.Tick(10 * time.Second) {

		for _, hash := range ts.Votes.Confirmed(Config.Configuration().DefaultAuthorities, nil) {

			tx, ok := ts.Unconfirmed.GetByHash(hash)
			if !ok || time.Since(tx.date) < time.Second*10 {
				continue
			}

			for _, tx := range ts.Unconfirmed.GetByFrontier(hash) {
				if dest, _ := tx.GetTarget(); bytes.Equal(dest, pk) {
					if _, ok := ts.Confirmed.GetByLinkHash(tx.Hash()); !ok {
						ts.Pending.Add(tx)
						continue
					}
				}

				ts.Confirmed.Add(tx)
				ts.Unconfirmed.Remove(tx)
			}
		}
	}
}


func (ts *TransactionsShelter) GetByHash(b Block.BlockHash) (tx Block.Transaction, t int, ok bool){
	if ts == nil {
		return
	}
	for i, txs := range []TransactionBox{ts.Unconfirmed, ts.Pending, ts.Confirmed} {
		if tx, ok = txs.GetByHash(b); ok {
			return tx, i, ok
		}
	}

	return nil, 0, false
}

func (ts *TransactionsShelter) GetByLinkHash(b Block.BlockHash) (tx Block.Transaction, t int, ok bool) {
	if ts == nil {
		return
	}

	for i, txs := range []TransactionBox{ts.Unconfirmed, h.Pending, st.Confirmed} {
		if tx, ok = txs.GetByLinkHash(b); ok {
			return tx, i, ok
		}
	}

	return nil, 0,false
}

func (ts *TransactionsShelter) GetByPreviousHash(b Block.BlockHash) (tx Block.Transaction, t int, ok bool) {
	if ts == nil {
		return
	}

	for i, txs := range []TransactionBox{ts.Unconfirmed, ts.Pending, ts.Confirmed} {
		if tx, ok = txs.GetByPreviousHash(b); ok {
			return tx, i, ok
		}
	}

	return nil, 0,false
}

func (h *TransactionBox) Listen() <-chan *Transaction {
	c := make(chan *Transaction)

	if h == nil {
		return c
	}

	h.listeners = append(h.listeners, c)

	return c
}

func (h *TransactionBox) notifyListeners(new *Transaction) {
	if h == nil {
		return
	}

	for _, l := range h.listeners {
		l <- new
	}
}

func (h *TransactionBox) GetByHash(b Block.BlockHash) (item *Transaction, ok bool) {
	if h == nil || b == nil  {
		return
	}

	item, ok = h.list[string(b)]
	return
}

func (h *TransactionBox) GetByLinkHash(b Block.BlockHash) (item Block.Transaction, ok bool) {
	if h == nil {
		return
	}

	for _, tx := range h.list {
		if tx.Transaction.GetType() == Block.Receive {
			if _, src := tx.GetTarget(); bytes.Equal(src, b) {
				return tx, true
			}
		}
	}

	return nil, false
}

func (h *TransactionBox) GetByPreviousHash(b Block.BlockHash) (item Block.Transaction, ok bool) {
	if h == nil {
		return
	}

	for _, tx := range h.list {
		if bytes.Equal(tx.GetPrevious(), b) {
			return tx, true
		}
	}

	return nil, false
}

func (h *TransactionBox) GetByFrontier(b Block.BlockHash) (items []Block.Transaction) {
	for {

		tx, ok := h.GetByHash(b)
		if  !ok {
			break
		}

		b = tx.GetPrevious()
		items = append(items, tx)
	}

	return items
}

func (h *TransactionBox) GetAll() (items []Block.Transaction) {
	for _, tx := range h.list {
		items = append(items, tx)
	}

	return items
}

func (h *TransactionBox) GetRange(min, max uint32) (items []Block.Transaction, ok bool) {
	if min <= 0 || max >= uint32(len(h.list)) {
		return nil, false
	}

	var i uint32 = 0
	for _, tx := range h.list {
		if i >= min && i < max {
			items = append(items, tx)
		}
		i++
	}

	return items, true
}

func (h *TransactionBox) Next(last, perPage uint32) (items []Block.Transaction, ok bool) {
	return h.GetRange(last, last+perPage)
}

func (h *TransactionBox) Previous(last, perPage uint32) (items []Block.Transaction, ok bool) {
	return h.GetRange(last-(perPage*2), last-perPage)
}

func (h *TransactionBox) GetRandom(n int) (items []*Transaction) {
	if h == nil {
		return
	}

	l := len(h.list)
	if l == 0 {
		return
	}

	var random = map[int]int{}
	for i := 0; i < n; i++ {
		n := rand.Intn(l - 1)
		random[n] = n
	}

	var i = 0
	for _, item := range h.list {
		if _, ok := random[i]; ok {
			items = append(items, item)
		}
		i++
	}

	return
}

func (h *TransactionBox) Add(transactions ...Block.Transaction) {
	if h == nil {
		return
	}

	for _, tx := range transactions {
		hash := tx.Hash()

		if _, ok := h.GetByHash(hash); !ok {
			h.list[string(hash)] = &Transaction{
				Transaction: tx,
				date:        time.Now(),
			}

			h.notifyListeners(h.list[string(hash)])

		}
	}
}

func (h *TransactionBox) Remove(txs ...Block.Transaction) {
	if h == nil {
		return
	}

	for _, tx := range txs {
		hash := tx.Hash()

		if _, ok := h.GetByHash(hash); !ok {
			delete(h.list, string(hash))
		}
	}
}

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

func (h *HistoryStore) Add(blk Block.Transaction, amount *Numbers.RawAmount) {
	hist := RPCClient.SingleHistory{}
	hist.Type = blk.GetSubType()
	hist.Amount = amount
	//	hist.Destination, hist.Source = blk.GetTarget()
	hist.Hash = blk.Hash()

	*h = append([]RPCClient.SingleHistory{hist}, *h...)
}

func (h *HistoryStore) Next(page uint32) {
	//@TODO pagination
}
