package Block

import (
	"crypto/subtle"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var (
	DefaultRepresentative, _ = Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95").GetPublicKey()
)

func CreateSignedSendBlock(sk Wallet.SecretKey, sending, balance *Numbers.RawAmount, previous BlockHash, destination Wallet.PublicKey) (Transaction, error) {
	blk, err := CreateSignedUniversalSendBlock(sk, DefaultRepresentative, balance, sending, previous, destination)
	if err != nil {
		return nil, err
	}

	return attachSignature(sk, blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Send))
}

func CreateSignedOpenBlock(sk Wallet.SecretKey, source BlockHash) (Transaction, error) {
	blk, err := CreateSignedUniversalOpenBlock(sk, DefaultRepresentative, Numbers.NewRaw(), source)
	if err != nil {
		return nil, err
	}

	return attachSignature(sk, blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Open))
}

func CreateSignedReceiveBlock(sk Wallet.SecretKey, source, previous BlockHash) (Transaction, error) {
	blk, err := CreateSignedUniversalReceiveBlock(sk, DefaultRepresentative, Numbers.NewRaw(), Numbers.NewRaw(), previous, source)
	if err != nil {
		return nil, err
	}

	return attachSignature(sk, blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Receive))
}

func CreateSignedChangeBlock(sk Wallet.SecretKey, previous BlockHash, representative Wallet.PublicKey) (Transaction, error) {
	blk, err := CreateSignedUniversalChangeBlock(sk, representative, Numbers.NewRaw(), previous)
	if err != nil {
		return nil, err
	}

	return attachSignature(sk, blk.SwitchToUniversalBlock(nil, nil).SwitchTo(Change))
}

func CreateSignedReceiveOrOpenBlock(sk Wallet.SecretKey, source, previous BlockHash) (blk Transaction, err error) {
	pk := sk.PublicKey()

	if previous == NewBlockHash(nil) || subtle.ConstantTimeCompare(previous[:], pk[:]) == 1 {
		return CreateSignedOpenBlock(sk, source)
	}

	return CreateSignedReceiveBlock(sk, source, previous)
}
