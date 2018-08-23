package Storage

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/ProofWork"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var Representative Wallet.Address

var Frontier Block.BlockHash
var LastBlock Block.UniversalBlock

func UpdateFrontier(blk Block.Transaction) {
//	LastBlock = *blk.SwitchToUniversalBlock()
	Frontier = LastBlock.SwitchTo(blk.GetType()).Hash()
}

var precomputedPoW = make(chan []byte, 1)
var lastPoW ProofWork.Work

func UpdatePoW() {
	hash := Frontier
	if Frontier == nil {
		hash = Block.BlockHash(PK)
	}

	if lastPoW.IsValid(Frontier) {
		precomputedPoW <- lastPoW
		return
	}

	precomputedPoW <- ProofWork.GenerateProof(hash)
}

func RetrievePrecomputedPoW() []byte {
	lastPoW = <-precomputedPoW
	return lastPoW
}
