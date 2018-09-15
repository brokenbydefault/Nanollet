package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
)

func (d *DefaultBlock) SetWork(w Work) {
	d.PoW = w
}

func (d *DefaultBlock) GetWork() Work {
	return d.PoW
}

func (d *DefaultBlock) SetSignature(s Wallet.Signature) {
	d.Signature = s
}

func (d *DefaultBlock) GetSignature() Wallet.Signature {
	return d.Signature
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

func (s *SendBlock) GetAccount() (pk Wallet.PublicKey) {
	return pk
}

func (s *ReceiveBlock) GetAccount() (pk Wallet.PublicKey) {
	return pk
}

func (s *OpenBlock) GetAccount() (pk Wallet.PublicKey) {
	return s.Account
}

func (s *ChangeBlock) GetAccount() (pk Wallet.PublicKey) {
	return pk
}

func (u *UniversalBlock) GetAccount() (pk Wallet.PublicKey) {
	return u.Account
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

func (s *SendBlock) GetBalance() *Numbers.RawAmount {
	return Numbers.NewRawFromBytes(s.Balance.ToBytes())
}

func (s *ReceiveBlock) GetBalance() *Numbers.RawAmount {
	// no-op
	return nil
}

func (s *OpenBlock) GetBalance() *Numbers.RawAmount {
	// no-op
	return nil
}

func (s *ChangeBlock) GetBalance() *Numbers.RawAmount {
	// no-op
	return nil
}

func (u *UniversalBlock) GetBalance() *Numbers.RawAmount {
	return Numbers.NewRawFromBytes(u.Balance.ToBytes())
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

func GetSubType(tx, txPrevious Transaction) BlockType {
	if tx.GetType() != State {
		return tx.GetType()
	}

	if hashPrev := tx.GetPrevious(); Util.IsEmpty(hashPrev[:]) {
		return Open
	}

	if dest, source := tx.GetTarget(); Util.IsEmpty(dest[:]) && Util.IsEmpty(source[:]) {
		return Change
	}

	if tx.GetBalance().Compare(txPrevious.GetBalance()) == 1 {
		return Receive
	} else {
		return Send
	}

	return Invalid
}

func GetAmount(tx, txPrevious Transaction) *Numbers.RawAmount {
	switch GetSubType(tx, txPrevious) {
	case Open:
		return tx.GetBalance()
	case Change:
		return Numbers.NewMin()
	case Send:
		return tx.GetBalance().Subtract(txPrevious.GetBalance()).Abs()
	case Receive:
		return tx.GetBalance().Subtract(txPrevious.GetBalance()).Abs()
	}

	return Numbers.NewMin()
}
