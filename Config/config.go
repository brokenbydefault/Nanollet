package Config

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

type Config struct {
	DebugStatus bool

	DefaultRepresentative Wallet.Address
	DefaultAuthorities    []Wallet.PublicKey
	DefaultNodes          []string

	DefaultFolder string

	DefaultMinimumAmount *Numbers.RawAmount
}

var Default = Config{
	DebugStatus: false,

	DefaultRepresentative: Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95"),
	DefaultFolder:         "Nanollet",

	DefaultMinimumAmount:  Numbers.NewRaw(),
}

func Configuration() Config {
	return Default
}
