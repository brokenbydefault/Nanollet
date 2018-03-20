package Block

import "github.com/brokenbydefault/Nanollet/ProofWork"

//@TODO Prevent do PoW when one valid PoW is already set

func (s *SendBlock) CreateProof() []byte {
	s.Work = ProofWork.GenerateProof(s.Previous)
	return s.Work
}

func (s *ReceiveBlock) CreateProof() []byte {
	s.Work = ProofWork.GenerateProof(s.Previous)
	return s.Work
}

func (s *OpenBlock) CreateProof() []byte {
	s.Work = ProofWork.GenerateProof(s.Account)
	return s.Work
}

func (s *ChangeBlock) CreateProof() []byte {
	s.Work = ProofWork.GenerateProof(s.Previous)
	return s.Work
}