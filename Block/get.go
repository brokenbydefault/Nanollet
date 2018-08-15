package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func (d *DefaultBlock) SetWork(w []byte) {
	d.PoW = w
}

func (d *DefaultBlock) SetSignature(s []byte) {
	d.Signature = s
}

func (d *DefaultBlock) GetType() BlockType {
	return d.Type
}

func (d *DefaultBlock) GetSubType() BlockType {
	if d.SubType == "" {
		return d.Type
	}
	return d.SubType
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
	var hash [32]byte
	copy(hash[:], h)

	u.Previous = hash[:]
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

func (s *SendBlock) GetTarget() (Wallet.Address, BlockHash) {
	return s.Destination, nil
}

func (s *ReceiveBlock) GetTarget() (Wallet.Address, BlockHash) {
	return "", s.Source
}

func (s *OpenBlock) GetTarget() (Wallet.Address, BlockHash) {
	return "", s.Source
}

func (s *ChangeBlock) GetTarget() (Wallet.Address, BlockHash) {
	// no-op
	return "", nil
}

func (u *UniversalBlock) GetTarget() (destination Wallet.Address, source BlockHash) {
	return Wallet.PublicKey(u.Link).CreateAddress(), u.Link
}
