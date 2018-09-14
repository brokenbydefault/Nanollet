package TwoFactor

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"errors"
)

var (
	NotDesignedUse = errors.New("seedfy not intended to be used in MFA")
)

type Token [32]byte

func NewSeedFY() (Wallet.SeedFY, error) {
	return Wallet.NewCustomFY(Wallet.V0, Wallet.MFA, 2, 10)
}

func NewToken(seedfy string, pass []byte) (token Token, err error) {
	sf, err := Wallet.ReadSeedFY(seedfy)
	if err != nil {
		return token, err
	}

	if !sf.IsValid(Wallet.V0, Wallet.MFA) {
		return token, NotDesignedUse
	}

	seed := sf.RecoverSeed(pass, nil)

	_, sk, err := seed.CreateKeyPair(Wallet.Base, 0)
	if err != nil {
		return token, err
	}

	copy(token[:], sk[:32])

	return token, nil
}