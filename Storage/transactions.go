package Storage

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"math/rand"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"sync"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
)

var TransactionStorage = TransactionBox{
	list:      new(sync.Map),
	listeners: make([]chan *Transaction, 0, 1),
}

type Transaction struct {
	Block.Transaction
	Date  time.Time
	votes *sync.Map // map[Wallet.PublicKey]uint64 // map[PublicKey Of Voter]Sequence
}

type TransactionBox struct {
	list      *sync.Map // map[Block.BlockHash]Block.Transaction
	listeners []chan *Transaction
}

func (h *TransactionBox) Listen() <-chan *Transaction {
	c := make(chan *Transaction)

	if h == nil {
		return c
	}

	h.listeners = append(h.listeners, c)

	return c
}

func (h *TransactionBox) Count() (len int) {
	if h == nil || h.list == nil {
		return 0
	}

	h.list.Range(func(_, _ interface{}) bool {
		len++
		return true
	})

	return len
}

func (h *TransactionBox) IsConfirmed(hash *Block.BlockHash, quorum *Peer.Quorum) (valid bool) {
	if h == nil || h.list == nil || quorum == nil {
		return false
	}

	t, ok := h.GetByHash(hash)
	if !ok || t.votes == nil {
		return false
	}

	if txs, ok := h.GetByPreviousHash(hash); ok {
		for _, tx := range txs {
			hash := tx.Hash()

			if h.IsConfirmed(&hash, quorum) {
				return true
			}
		}
	}

	previous := t.GetPrevious()
	txs, _ := h.GetByPreviousHash(&previous)

	var possibleWinner, winner = make(map[Block.BlockHash]int), Block.BlockHash{}
	for _, pk := range quorum.PublicKeys {

		// For each authorities we need to get him vote, if any. Since one representative can change the vote,
		// we need to compare the sequence of all transactions voted by him.
		for _, tx := range txs {
			hash := tx.Hash()

			higherSeq := uint64(0)
			if seq, ok := t.votes.Load(pk); ok && seq.(uint64) > higherSeq {
				possibleWinner[hash]++

				if possibleWinner[hash] > possibleWinner[winner] {
					winner = hash
				}
			}
		}
	}

	votesNeeded := quorum.Calc(len(possibleWinner))
	if winner == *hash && possibleWinner[*hash] >= votesNeeded {
		return true
	}

	return false
}

func (h *TransactionBox) WaitConfirmation(quorum *Peer.Quorum, timeout time.Duration, hashes ...*Block.BlockHash) (winner *Block.BlockHash, ok bool) {
	if len(hashes) == 0 || quorum == nil {
		return nil, false
	}

	c := make(chan *Block.BlockHash, len(hashes))
	for _, hash := range hashes {
		if hash == nil {
			continue
		}

		go func(t *TransactionBox, h *Block.BlockHash, ch chan *Block.BlockHash) {
			if t.IsConfirmed(h, quorum) {
				c <- hash
			}
		}(h, hash, c)
	}

	select {
	case hash := <-c:
		return hash, true
	case <-time.After(timeout):
		return nil, false
	}
}

func (h *TransactionBox) AddVotes(hash *Block.BlockHash, pk *Wallet.PublicKey, seq uint64) {
	if h == nil {
		return
	}

	t, ok := h.GetByHash(hash)
	if !ok {
		return
	}

	if t.votes == nil {
		t.votes = new(sync.Map)
	}

	t.votes.LoadOrStore(*pk, seq)
}

func (h *TransactionBox) GetByHash(hash *Block.BlockHash) (item *Transaction, ok bool) {
	if h == nil || Util.IsEmpty(hash[:]) {
		return
	}

	value, ok := h.list.Load(*hash)
	if !ok {
		return nil, ok
	}

	item, ok = value.(*Transaction)
	if !ok {
		return nil, ok
	}

	return item, ok
}

func (h *TransactionBox) GetByLinkHash(hash *Block.BlockHash) (item *Transaction, ok bool) {
	if h == nil {
		return
	}

	h.list.Range(func(key, value interface{}) bool {
		item, ok = value.(*Transaction)
		if ok {
			if _, src := item.GetTarget(); src == *hash {
				return false
			}
		}

		item, ok = nil, false
		return true
	})

	return item, ok
}

func (h *TransactionBox) GetByPreviousHash(hash *Block.BlockHash) (items []*Transaction, ok bool) {
	if h == nil {
		return
	}

	h.list.Range(func(key, value interface{}) bool {
		if item, ok := value.(*Transaction); ok && item.GetPrevious() == *hash {
			items = append(items, item)
		}

		return true
	})

	if len(items) == 0 {
		return nil, false
	}

	return items, true
}

func (h *TransactionBox) GetByFrontier(hash Block.BlockHash) (items []*Transaction) {
	if h == nil {
		return
	}

	var (
		tx *Transaction
		ok bool
	)

	for {
		tx, ok = h.GetByHash(&hash)
		if !ok || Util.IsEmpty(hash[:]) {
			break
		}

		items = append(items, tx)
		hash = tx.GetPrevious()
	}

	return items
}

func (h *TransactionBox) GetAll() (items []Block.Transaction) {
	if h == nil {
		return
	}

	h.list.Range(func(key, value interface{}) bool {
		item, ok := value.(*Transaction)
		if ok {
			items = append(items, item)
		}

		return true
	})

	return items
}

func (h *TransactionBox) GetRandom(n int) (items []Block.Transaction) {
	if h == nil {
		return
	}

	list := h.GetAll()

	l := len(list)
	if l == 0 {
		return
	}

	if n <= 0 || n > l {
		n = l
	}

	var random = map[int]int{}
	for i := 0; i < n; i++ {
		g := rand.Intn(l)
		random[g] = g
	}

	for i := range random {
		if val, ok := list[i].(Block.Transaction); ok {
			items = append(items, val)
		}
	}

	return
}

func (h *TransactionBox) Add(txs ...Block.Transaction) {
	if h == nil {
		return
	}

	for _, tx := range txs {
		hash := tx.Hash()
		t := &Transaction{
			Transaction: tx,
			Date:        time.Now(),
			votes:       new(sync.Map),
		}

		if _, old := h.list.LoadOrStore(hash, t); !old {
			h.notifyListeners(t)
		}

	}
}

func (h *TransactionBox) Remove(txs ...Block.Transaction) {
	if h == nil {
		return
	}

	for _, tx := range txs {
		h.list.Delete(tx.Hash())
	}
}

func (h *TransactionBox) notifyListeners(tx *Transaction) {
	if h == nil {
		return
	}

	for _, l := range h.listeners {
		l <- tx
	}
}
