package Storage

import (
	"github.com/brokenbydefault/Nanollet/ProofWork"
)

type FrontierStore []byte
type PoWStore []byte

var Frontier []byte
var PrecomputedPoW []byte

func UpdateFrontier(hash []byte) {
	Frontier = hash
	if !ProofWork.IsValidProof(Frontier, PrecomputedPoW) {
		PrecomputedPoW = ProofWork.GenerateProof(Frontier)
	}
}
