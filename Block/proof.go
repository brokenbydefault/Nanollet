package Block

import (
	"bytes"
	"github.com/brokenbydefault/Nanollet/ProofWork"
)

func (s *SendBlock) Work() ProofWork.Work {
	if !s.PoW.IsValid(s.Previous) {
		s.PoW = ProofWork.GenerateProof(s.Previous)
	}

	return s.PoW
}

func (s *ReceiveBlock) Work() ProofWork.Work {
	if !s.PoW.IsValid(s.Previous) {
		s.PoW = ProofWork.GenerateProof(s.Previous)
	}

	return s.PoW
}

func (s *OpenBlock) Work() ProofWork.Work {
	var previous, _ = s.Account.GetPublicKey()

	if !s.PoW.IsValid(previous) {
		s.PoW = ProofWork.GenerateProof(previous)
	}

	return s.PoW
}

func (s *ChangeBlock) Work() ProofWork.Work {
	if !s.PoW.IsValid(s.Previous) {
		s.PoW = ProofWork.GenerateProof(s.Previous)
	}

	return s.PoW
}

func (u *UniversalBlock) Work() ProofWork.Work {
	var previous []byte

	if u.Previous == nil || bytes.Equal(u.Previous, make([]byte, 32)) || u.SubType == Open {
		previous, _ = u.Account.GetPublicKey()
	}else{
		previous = u.Previous
	}

	if !u.PoW.IsValid(previous) {
		u.PoW = ProofWork.GenerateProof(previous)
	}

	return u.PoW
}
