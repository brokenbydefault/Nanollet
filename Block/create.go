package Block

import (
	"errors"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"crypto/subtle"
)

var (
	ErrEndBlock         = errors.New("not a block")
	ErrInvalidBlock     = errors.New("invalid block")
	ErrInvalidBlockType = errors.New("invalid type")
)

func NewTransaction(blockType BlockType) (blk Transaction, size int, err error) {
	switch blockType & 0x0F {
	case Send:
		blk = &SendBlock{}
		size = SendSize
	case Receive:
		blk = &ReceiveBlock{}
		size = ReceiveSize
	case Change:
		blk = &ChangeBlock{}
		size = ChangeSize
	case Open:
		blk = &OpenBlock{}
		size = OpenSize
	case State:
		blk = &UniversalBlock{}
		size = StateSize
	case Invalid:
		err = ErrInvalidBlock
	case NotABlock:
		err = ErrEndBlock
	default:
		err = ErrInvalidBlockType
	}

	return blk, size, err
}

func CreateUniversalSendBlock(sk *Wallet.SecretKey, representative Wallet.PublicKey, balance, sending *Numbers.RawAmount, previous BlockHash, destination Wallet.PublicKey) (Transaction, error) {
	blk := &UniversalBlock{
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        balance.Subtract(sending),
		Previous:       previous,
		Link:           BlockHash(destination),
	}

	if sk == nil {
		return blk, nil
	}

	return attachSignature(sk, blk)
}

func CreateUniversalOpenBlock(sk *Wallet.SecretKey, representative Wallet.PublicKey, receiving *Numbers.RawAmount, source BlockHash) (Transaction, error) {
	blk := &UniversalBlock{
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        receiving,
		//Previous:       NewBlockHash(nil),
		Link: source,
	}

	if sk == nil {
		return blk, nil
	}

	return attachSignature(sk, blk)
}

func CreateUniversalReceiveBlock(sk *Wallet.SecretKey, representative Wallet.PublicKey, balance, receiving *Numbers.RawAmount, previous BlockHash, source BlockHash) (Transaction, error) {
	blk := &UniversalBlock{
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        balance.Add(receiving),
		Previous:       previous,
		Link:           source,
	}

	if sk == nil {
		return blk, nil
	}

	return attachSignature(sk, blk)
}

func CreateUniversalChangeBlock(sk *Wallet.SecretKey, representative Wallet.PublicKey, balance *Numbers.RawAmount, previous BlockHash) (Transaction, error) {
	blk := &UniversalBlock{
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        balance,
		Previous:       previous,
		//Link:           NewBlockHash(nil),
	}

	if sk == nil {
		return blk, nil
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalReceiveOrOpenBlock(sk *Wallet.SecretKey, representative Wallet.PublicKey, balance, receiving *Numbers.RawAmount, previous BlockHash, source BlockHash) (Transaction, error) {
	pk := sk.PublicKey()

	if previous == NewBlockHash(nil) || subtle.ConstantTimeCompare(previous[:], pk[:]) == 1 {
		return CreateUniversalOpenBlock(sk, representative, receiving, source)
	}

	return CreateUniversalReceiveBlock(sk, representative, balance, receiving, previous, source)
}

func attachSignature(sk *Wallet.SecretKey, blk Transaction) (Transaction, error) {
	hash := blk.Hash()

	sig, err := sk.Sign(hash[:])
	if err != nil {
		return nil, err
	}

	blk.SetSignature(sig)

	return blk, err
}
