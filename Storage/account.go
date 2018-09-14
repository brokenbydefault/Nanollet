package Storage

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

var AccessStorage Access

type Access struct {
	Token    [32]byte
	Password []byte
	Seed     Wallet.Seed
}

var AccountStorage Account

type Account struct {
	SecretKey Wallet.SecretKey
	PublicKey Wallet.PublicKey

	Representative Wallet.PublicKey
	Frontier       Block.BlockHash
	Balance        *Numbers.RawAmount

	//PreComputedWork Block.Work
}
