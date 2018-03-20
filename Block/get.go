package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
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
}

func (s *OpenBlock) SetBalance(n *Numbers.RawAmount) {
}

func (s *ChangeBlock) SetBalance(n *Numbers.RawAmount) {
}


func (d *DefaultBlock) GetType() string {
	return d.Type
}

func (d *DefaultBlock) SetWork(w []byte) {
	d.Work = w
}

func (d *DefaultBlock) SetSignature(s []byte) {
	d.Signature = s
}
