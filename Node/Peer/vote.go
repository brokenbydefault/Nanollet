package Peer

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
)

type Quorum struct {
	PublicKeys []Wallet.PublicKey
	Common     int
	Fork       int
}

func (q *Quorum) CalcCommon() int {
	if q == nil {
		return (len(q.PublicKeys) * 50) / 100
	}

	return (len(q.PublicKeys) * q.Common) / 100
}

func (q *Quorum) CalcFork() int {
	if q == nil {
		return (len(q.PublicKeys) * 50) / 100
	}

	return (len(q.PublicKeys) * q.Fork) / 100
}

func (q *Quorum) Calc(possibleWinners int) int {
	if possibleWinners > 1 {
		return q.CalcFork()
	}

	return q.CalcCommon()
}