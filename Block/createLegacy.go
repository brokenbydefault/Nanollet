package Block

import (
	"crypto/subtle"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Config"
)

func CreateSignedSendBlock(sk Wallet.SecretKey, sending, balance *Numbers.RawAmount, previous []byte, destination Wallet.Address) (BlockTransaction, error) {
	blk, err := CreateSignedUniversalSendBlock(sk, Config.DefaultRepresentative, balance, sending, previous, destination)
	if err != nil {
		return nil, err
	}

	return blk.SwitchToUniversalBlock().SwitchTo(Send), nil
}

func CreateSignedOpenBlock(sk Wallet.SecretKey, source []byte) (BlockTransaction, error) {
	blk, err := CreateSignedUniversalOpenBlock(sk, Config.DefaultRepresentative, Numbers.NewRaw(), source)
	if err != nil {
		return nil, err
	}

	return blk.SwitchToUniversalBlock().SwitchTo(Open), nil
}

func CreateSignedReceiveBlock(sk Wallet.SecretKey, source, previous []byte) (BlockTransaction, error) {
	blk, err := CreateSignedUniversalReceiveBlock(sk, Config.DefaultRepresentative, Numbers.NewRaw(), Numbers.NewRaw(), previous, source)
	if err != nil {
		return nil, err
	}

	return blk.SwitchToUniversalBlock().SwitchTo(Receive), nil
}

func CreateSignedChangeBlock(sk Wallet.SecretKey, previous []byte, representative Wallet.Address) (BlockTransaction, error) {
	blk, err := CreateSignedUniversalChangeBlock(sk, representative, Numbers.NewRaw(), previous)
	if err != nil {
		return nil, err
	}

	return blk.SwitchToUniversalBlock().SwitchTo(Change), nil
}

func CreateSignedReceiveOrOpenBlock(sk Wallet.SecretKey, source, previous []byte) (blk BlockTransaction, err error) {
	pk, err := sk.PublicKey()
	if err != nil {
		return blk, err
	}

	if previous == nil || subtle.ConstantTimeCompare(previous, pk) == 1 {
		return CreateSignedOpenBlock(sk, source)
	}

	return CreateSignedReceiveBlock(sk, source, previous)
}
