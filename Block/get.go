package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/ProofWork"
)

func (d *DefaultBlock) SetWork(w ProofWork.Work) {
	d.PoW = w
}

func (d *DefaultBlock) SetSignature(s Wallet.Signature) {
	d.Signature = s
}

func (d *DefaultBlock) GetSubType() BlockType {
	if d.subType == Invalid {
		return d.mainType
	}
	return d.subType
}

func (s *SendBlock) GetType() BlockType {
	return Send
}

func (s *ReceiveBlock) GetType() BlockType {
	return Receive
}

func (s *OpenBlock) GetType() BlockType {
	return Open
}

func (s *ChangeBlock) GetType() BlockType {
	return Change
}

func (u *UniversalBlock) GetType() BlockType {
	return State
}

func (s *SendBlock) SetFrontier(h BlockHash) {
	s.Previous = h
}

func (s *ReceiveBlock) SetFrontier(h BlockHash) {
	s.Previous = h
}

func (s *OpenBlock) SetFrontier(h BlockHash) {
	s.Source = h
}

func (s *ChangeBlock) SetFrontier(h BlockHash) {
	s.Previous = h
}

func (u *UniversalBlock) SetFrontier(h BlockHash) {
	copy(u.Previous[:], h[:])
}

func (s *SendBlock) SetBalance(n *Numbers.RawAmount) {
	s.Balance = n
}

func (s *ReceiveBlock) SetBalance(n *Numbers.RawAmount) {
	// no-op
}

func (s *OpenBlock) SetBalance(n *Numbers.RawAmount) {
	// no-op
}

func (s *ChangeBlock) SetBalance(n *Numbers.RawAmount) {
	// no-op
}

func (u *UniversalBlock) SetBalance(n *Numbers.RawAmount) {
	u.Balance = n
}

func (s *SendBlock) GetTarget() (pk Wallet.PublicKey, hash BlockHash) {
	return s.Destination, hash
}

func (s *ReceiveBlock) GetTarget() (pk Wallet.PublicKey, hash BlockHash) {
	return pk, s.Source
}

func (s *OpenBlock) GetTarget() (pk Wallet.PublicKey, hash BlockHash) {
	return pk, s.Source
}

func (s *ChangeBlock) GetTarget() (pk Wallet.PublicKey, hash BlockHash) {
	return pk, hash
}

func (u *UniversalBlock) GetTarget() (destination Wallet.PublicKey, source BlockHash) {
	return Wallet.PublicKey(u.Link), u.Link
}
