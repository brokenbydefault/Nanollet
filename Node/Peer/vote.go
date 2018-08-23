package Peer

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

type Quorum struct {
	Common int
	Fork   int
}

func (q *Quorum) CalcCommon(pks []Wallet.PublicKey) int {
	if q == nil {
		return  (len(pks) * 50) / 100
	}

	return (len(pks) * q.Common) / 100
}

func (q *Quorum) CalcFork(pks []Wallet.PublicKey) int {
	if q == nil {
		return (len(pks) * 50) / 100
	}

	return (len(pks) * q.Fork) / 100
}

type Votes struct {
	list map[string]map[string]vote
}

type vote struct {
	Sequence int64
	Hash     Block.BlockHash
}

func (v *Votes) Add(pk Wallet.PublicKey, seq int64, tx Block.Transaction) {
	previous := string(tx.SwitchToUniversalBlock(nil, nil).Previous)

	if vote, ok := v.list[previous][string(pk)]; !ok || vote.Sequence <= seq {
		v.list[previous][string(pk)] = vote{
			Sequence: seq,
			Hash:     vote.Hash,
		}
	}

}

func (v *Votes) Remove(previous Block.BlockHash) {
	if v == nil {
		return
	}

	delete(v.list, string(previous))
}

// @TODO Speedup and simply
func (v *Votes) Confirmed(pks []Wallet.PublicKey, quorum *Quorum) (txs []Block.BlockHash) {
	if v == nil {
		return
	}

	for _, vote := range v.list {
		result := map[string]int{}

		for _, pk := range pks {
			vote := vote[string(pk)]

			if _, ok := result[string(vote.Hash)]; ok {
				result[string(vote.Hash)]++
			} else {
				result[string(vote.Hash)] = 1
			}
		}

		qWinner, hWinner, vWinner := quorum.CalcCommon(pks), Block.BlockHash{}, 0
		if len(result) != 1 {
			qWinner = quorum.CalcFork(pks)
		}

		for hash, votes := range result {
			if votes > qWinner && votes > vWinner {
				hWinner, vWinner = Block.BlockHash(hash), votes
			}
		}

		if hWinner != nil {
			txs = append(txs, hWinner)
		}
	}

	return txs
}
