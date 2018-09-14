package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"golang.org/x/crypto/blake2b"
)

var addressBase = Wallet.PublicKey{0x51, 0x14, 0xab, 0x7c, 0x6a, 0xd0, 0xd6, 0xc3, 0x14, 0xc5, 0xc2, 0x8e, 0x36, 0xb0, 0x8a, 0x65, 0x0a, 0xd4, 0x2b, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type Nanofier interface {
	CreateBlocks() (txs []Block.Transaction, err error)
	CreateSignatureBlock() (tx Block.Transaction, err error)
	CreateFlagBlock() (tx Block.Transaction, err error)

	IsValid() (ok bool)
	IsCorrectlyFormatted() (ok bool)

	Flag() Wallet.PublicKey
	Amount() *Numbers.RawAmount
}

func CreatePublicKey(version uint64) Wallet.PublicKey {
	pk := Wallet.NewPublicKey(addressBase[:])
	pk[24] = byte(version)
	pk[25] = byte(version >> 8)
	pk[26] = byte(version >> 16)
	pk[27] = byte(version >> 24)
	pk[28] = byte(version >> 32)
	pk[29] = byte(version >> 40)
	pk[30] = byte(version >> 48)
	pk[31] = byte(version >> 56)

	return pk
}

func CreateAddress(version uint64) Wallet.Address {
	return CreatePublicKey(version).CreateAddress()
}

func CreateHash(file io.Reader) (hash Wallet.PublicKey, err error) {
	blake, _ := blake2b.New(32, nil)

	_, err = io.Copy(blake, file)
	if err != nil {
		return hash, err
	}

	copy(hash[:], blake.Sum(nil))
	return hash, nil
}

func isSignatureValid(pk *Wallet.PublicKey, tx Block.Transaction) bool {
	if hash, sig := tx.Hash(), tx.GetSignature(); pk.IsValidSignature(hash[:], &sig) {
		return true
	}

	return false
}
