package Block

import (
	"github.com/brokenbydefault/Nanollet/ProofWork"
)

func (s *SendBlock) CreateProof() []byte {
	if ProofWork.IsValidProof(s.Previous, s.Work) {
		return s.Work
	}
	s.Work = ProofWork.GenerateProof(s.Previous)
	return s.Work
}

func (s *ReceiveBlock) CreateProof() []byte {
	if ProofWork.IsValidProof(s.Previous, s.Work) {
		return s.Work
	}
	s.Work = ProofWork.GenerateProof(s.Previous)
	return s.Work
}

func (s *OpenBlock) CreateProof() []byte {
	if ProofWork.IsValidProof(s.Account, s.Work) {
		return s.Work
	}
	s.Work = ProofWork.GenerateProof(s.Account)
	return s.Work
}

func (s *ChangeBlock) CreateProof() []byte {
	if ProofWork.IsValidProof(s.Previous, s.Work) {
		return s.Work
	}
	s.Work = ProofWork.GenerateProof(s.Previous)
	return s.Work
}
