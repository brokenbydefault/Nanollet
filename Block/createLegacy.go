package Block

import (
	"crypto/subtle"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var (
	DefaultRepresentative, _ = Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95").GetPublicKey()
)

func CreateSendBlock(sk *Wallet.SecretKey, sending, balance *Numbers.RawAmount, previous BlockHash, destination Wallet.PublicKey) (Transaction, error) {
	blk, err := CreateUniversalSendBlock(sk, Wallet.NewPublicKey(nil), balance, sending, previous, destination)
	if err != nil {
		return nil, err
	}

	legacy := blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Send)

	if sk == nil {
		return legacy, nil
	}

	return attachSignature(sk, legacy)
}

func CreateOpenBlock(sk *Wallet.SecretKey, source BlockHash) (Transaction, error) {
	blk, err := CreateUniversalOpenBlock(sk, Wallet.NewPublicKey(nil), Numbers.NewRaw(), source)
	if err != nil {
		return nil, err
	}

	legacy := blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Open)

	if sk == nil {
		return legacy, nil
	}

	return attachSignature(sk, legacy)
}

func CreateReceiveBlock(sk *Wallet.SecretKey, source, previous BlockHash) (Transaction, error) {
	blk, err := CreateUniversalReceiveBlock(sk, Wallet.NewPublicKey(nil), Numbers.NewRaw(), Numbers.NewRaw(), previous, source)
	if err != nil {
		return nil, err
	}

	legacy := blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Receive)

	if sk == nil {
		return legacy, nil
	}

	return attachSignature(sk, legacy)
}

func CreateChangeBlock(sk *Wallet.SecretKey, previous BlockHash, representative Wallet.PublicKey) (Transaction, error) {
	blk, err := CreateUniversalChangeBlock(sk, representative, Numbers.NewRaw(), previous)
	if err != nil {
		return nil, err
	}

	legacy := blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Change)

	if sk == nil {
		return legacy, nil
	}

	return attachSignature(sk, legacy)
}

func CreateReceiveOrOpenBlock(sk *Wallet.SecretKey, source, previous BlockHash) (blk Transaction, err error) {
	pk := sk.PublicKey()

	if previous == NewBlockHash(nil) || subtle.ConstantTimeCompare(previous[:], pk[:]) == 1 {
		return CreateOpenBlock(sk, source)
	}

	return CreateReceiveBlock(sk, source, previous)
}
