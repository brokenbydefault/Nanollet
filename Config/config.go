package Config

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"net"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
)

type Config struct {
	DebugStatus bool

	DefaultRepresentative Wallet.Address
	DefaultAuthorities    []Wallet.PublicKey
	DefaultNodes          []string

	DefaultFolder string

	DefaultMinimumAmount *Numbers.RawAmount

	DefaultUDPNetwork *net.UDPAddr
	DefaultPeers []*Peer.Peer
}

var defaultConfig = Config{
	DebugStatus: false,

	DefaultRepresentative: Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95"),
	DefaultFolder:         "Nanollet",

	DefaultMinimumAmount: Numbers.NewRaw(),

	DefaultUDPNetwork: &net.UDPAddr{
		Port: 6000,
	},

	DefaultPeers: Peer.NewPeersFromString("rai-beta.raiblocks.net:54000"),
}

func Configuration() Config {
	return defaultConfig
}
