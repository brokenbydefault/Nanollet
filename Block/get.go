package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func (s *SendBlock) SetFrontier(h BlockHash) {
	s.Previous = h
}

func (s *ReceiveBlock) SetFrontier(h BlockHash) {
	s.Previous = h
}

func (s *OpenBlock) SetFrontier(h BlockHash) {
}

func (s *ChangeBlock) SetFrontier(h BlockHash) {
	s.Previous = h
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

func (d *DefaultBlock) SetWork(w []byte) {
	d.Work = w
}

func (d *DefaultBlock) SetSignature(s []byte) {
	d.Signature = s
}

func (d *DefaultBlock) GetType() string {
	return d.Type
}

func (s *SendBlock) GetTarget() (Wallet.Address, BlockHash) {
	return s.Destination.CreateAddress(), nil
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