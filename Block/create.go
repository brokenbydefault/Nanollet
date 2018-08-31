package Block

import (
	"errors"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
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
		blk = &UniversalBlock{DefaultBlock: DefaultBlock{subType: blockType & 0xF0}}
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

func parseAccount(errs []error, f func() (Wallet.PublicKey, error)) Wallet.PublicKey {
	pk, err := f()
	errs = append(errs, err)

	return pk
}

func parseBalance(errs []error, n *Numbers.RawAmount) *Numbers.RawAmount {
	if ok := n.IsValid(); !ok {
		errs = append(errs, errors.New("valid amount"))
	}

	return n
}

func attachSignature(sk Wallet.SecretKey, blk Transaction) (Transaction, error) {
	hashb := blk.Hash()

	sig, err := sk.CreateSignature(hashb[:])
	if err != nil {
		return nil, err
	}

	blk.SetSignature(sig)
	return blk, err
}

func CreateSignedUniversalSendBlock(sk Wallet.SecretKey, representative Wallet.PublicKey, balance, sending *Numbers.RawAmount, previous BlockHash, destination Wallet.PublicKey) (Transaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			mainType: State,
			subType:  Send,
		},
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        parseBalance(errs, balance.Subtract(sending)),
		Previous:       previous,
		Link:           BlockHash(destination),
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalOpenBlock(sk Wallet.SecretKey, representative Wallet.PublicKey, receiving *Numbers.RawAmount, source BlockHash) (Transaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			mainType: State,
			subType:  Open,
		},
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        parseBalance(errs, receiving),
		Previous:       NewBlockHash(nil),
		Link:           source,
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalReceiveBlock(sk Wallet.SecretKey, representative Wallet.PublicKey, balance, receiving *Numbers.RawAmount, previous BlockHash, source BlockHash) (Transaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			mainType: State,
			subType:  Receive,
		},
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        parseBalance(errs, balance.Add(receiving)),
		Previous:       previous,
		Link:           source,
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalChangeBlock(sk Wallet.SecretKey, representative Wallet.PublicKey, balance *Numbers.RawAmount, previous BlockHash) (Transaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			mainType:    State,
			subType: Change,
		},
		Account:        sk.PublicKey(),
		Representative: representative,
		Balance:        parseBalance(errs, balance),
		Previous:       previous,
		Link:           NewBlockHash(nil),
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalReceiveOrOpenBlock(sk Wallet.SecretKey, representative Wallet.PublicKey, balance, receiving *Numbers.RawAmount, previous BlockHash, source BlockHash) (Transaction, error) {
	pk := sk.PublicKey()

	if previous == NewBlockHash(nil) || subtle.ConstantTimeCompare(previous[:], pk[:]) == 1 {
		return CreateSignedUniversalOpenBlock(sk, representative, receiving, source)
	}

	return CreateSignedUniversalReceiveBlock(sk, representative, balance, receiving, previous, source)
}
