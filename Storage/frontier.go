package Storage

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/ProofWork"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
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
	if Util.IsEmpty(Frontier[:]) {
		hash = Block.BlockHash(PK)
	}

	if lastPoW.IsValid(Frontier[:]) {
		precomputedPoW <- lastPoW[:]
		return
	}

	pow := ProofWork.GenerateProof(hash[:])
	precomputedPoW <- pow[:]
}

func RetrievePrecomputedPoW() []byte {
	lastPoW = ProofWork.NewWork(<-precomputedPoW)
	return lastPoW[:]
}
