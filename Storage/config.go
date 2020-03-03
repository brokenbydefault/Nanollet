package Storage

import (
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var Configuration = Live

type Config struct {
	DebugStatus bool

	Account ConfigAccount
	Node    ConfigNode
	Storage ConfigStorage
}

type ConfigAccount struct {
	Representative Wallet.PublicKey
	Quorum         Peer.Quorum
	//DefaultPrefixes       []string
	MinimumAmount *Numbers.RawAmount
	UnitAmount    Numbers.UnitBase
}

type ConfigNode struct {
	Peers  []*Peer.Peer
	Header Packets.Header
	//DefaultMinimumWork int64
	//DefaultGenesis
}

type ConfigStorage struct {
	Folder string
}

var Live = Config{
	DebugStatus: false,

	Account: ConfigAccount{
		Representative: Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95").MustGetPublicKey(),
		Quorum: Peer.Quorum{
			PublicKeys: []Wallet.PublicKey{
				Wallet.Address("xrb_3rw4un6ys57hrb39sy1qx8qy5wukst1iiponztrz9qiz6qqa55kxzx4491or").MustGetPublicKey(),
				Wallet.Address("xrb_3pczxuorp48td8645bs3m6c3xotxd3idskrenmi65rbrga5zmkemzhwkaznh").MustGetPublicKey(),
				Wallet.Address("xrb_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou").MustGetPublicKey(),
				Wallet.Address("xrb_1rs5rtyeo1owjt6cz9ypdkqyydq656kai8t35haiioapts39x96br5u4mbdw").MustGetPublicKey(),
				Wallet.Address("xrb_3arg3asgtigae3xckabaaewkx3bzsh7nwz7jkmjos79ihyaxwphhm6qgjps4").MustGetPublicKey(),
				Wallet.Address("xrb_1awsn43we17c1oshdru4azeqjz9wii41dy8npubm4rg11so7dx3jtqgoeahy").MustGetPublicKey(),
				Wallet.Address("xrb_3hd4ezdgsp15iemx7h81in7xz5tpxi43b6b41zn3qmwiuypankocw3awes5k").MustGetPublicKey(),
				Wallet.Address("xrb_1q3hqecaw15cjt7thbtxu3pbzr1eihtzzpzxguoc37bj1wc5ffoh7w74gi6p").MustGetPublicKey(),
				Wallet.Address("xrb_3rropjiqfxpmrrkooej4qtmm1pueu36f9ghinpho4esfdor8785a455d16nf").MustGetPublicKey(),
				Wallet.Address("xrb_1nanode8ngaakzbck8smq6ru9bethqwyehomf79sae1k7xd47dkidjqzffeg").MustGetPublicKey(),
				Wallet.Address("xrb_1hza3f7wiiqa7ig3jczyxj5yo86yegcmqk3criaz838j91sxcckpfhbhhra1").MustGetPublicKey(),
				Wallet.Address("xrb_1anrzcuwe64rwxzcco8dkhpyxpi8kd7zsjc1oeimpc3ppca4mrjtwnqposrs").MustGetPublicKey(),
				Wallet.Address("xrb_3dmtrrws3pocycmbqwawk6xs7446qxa36fcncush4s1pejk16ksbmakis78m").MustGetPublicKey(),
				Wallet.Address("xrb_1brainb3zz81wmhxndsbrjb94hx3fhr1fyydmg6iresyk76f3k7y7jiazoji").MustGetPublicKey(),
				Wallet.Address("xrb_1bj5cf9hkgkcspmn15day8cyn3hyaciufbba4rqmbnkmbdpjdmo9pwyatjoi").MustGetPublicKey(),
				Wallet.Address("xrb_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd").MustGetPublicKey(),
				Wallet.Address("xrb_1bananobjcrqugm87e8p3kxkhy7d1bzkty53n889iyunm83cp14rb9fin78p").MustGetPublicKey(),
				Wallet.Address("xrb_1n1hukyqred6yuch1xgtmdofe1bnc68eza733qmb6r19xo9us7qipbjujad1").MustGetPublicKey(),
				Wallet.Address("xrb_1gaysex8yymd5ef88hjqxt8xbjt63qz43cujrrzy4df9xb6zhf315csi35ww").MustGetPublicKey(),
			},
			Common: 50,
			Fork:   60,
		},
		MinimumAmount: Numbers.NewRawFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD3, 0xC2, 0x1B, 0xCE, 0xCC, 0xED, 0xA1, 0x00, 0x00, 0x00}),
		UnitAmount:    Numbers.MegaXRB,
	},

	Node: ConfigNode{
		Peers: Peer.NewPeersFromString(
			"rai.raiblocks.net:7075",
			"185.243.9.164:7075",
			"206.189.190.7:7075",
			"198.245.55.107:7075",
		),
		Header: Packets.Header{
			MagicNumber:  82,
			NetworkType:  Packets.Live,
			VersionMax:   0,
			VersionUsing: 255,
			VersionMin:   255,
		},
	},

	Storage: ConfigStorage{
		Folder: "Nanollet",
	},
}

var Beta = Config{
	DebugStatus: false,

	Account: ConfigAccount{
		Representative: Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95").MustGetPublicKey(),
		Quorum: Peer.Quorum{
			PublicKeys: []Wallet.PublicKey{
				Wallet.Address("xrb_3arg3asgtigae3xckabaaewkx3bzsh7nwz7jkmjos79ihyaxwphhm6qgjps4").MustGetPublicKey(),
				Wallet.Address("xrb_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou").MustGetPublicKey(),
				Wallet.Address("xrb_1q3hqecaw15cjt7thbtxu3pbzr1eihtzzpzxguoc37bj1wc5ffoh7w74gi6p").MustGetPublicKey(),
				Wallet.Address("xrb_3dmtrrws3pocycmbqwawk6xs7446qxa36fcncush4s1pejk16ksbmakis78m").MustGetPublicKey(),
				Wallet.Address("xrb_3hd4ezdgsp15iemx7h81in7xz5tpxi43b6b41zn3qmwiuypankocw3awes5k").MustGetPublicKey(),
				Wallet.Address("xrb_1awsn43we17c1oshdru4azeqjz9wii41dy8npubm4rg11so7dx3jtqgoeahy").MustGetPublicKey(),
				Wallet.Address("xrb_1anrzcuwe64rwxzcco8dkhpyxpi8kd7zsjc1oeimpc3ppca4mrjtwnqposrs").MustGetPublicKey(),
				Wallet.Address("xrb_1hza3f7wiiqa7ig3jczyxj5yo86yegcmqk3criaz838j91sxcckpfhbhhra1").MustGetPublicKey(),
			},
			Common: 30,
			Fork:   60,
		},
		MinimumAmount: Numbers.NewRawFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD3, 0xC2, 0x1B, 0xCE, 0xCC, 0xED, 0xA1, 0x00, 0x00, 0x00}),
		UnitAmount:    Numbers.MegaXRB,
	},

	Node: ConfigNode{
		Peers: Peer.NewPeersFromString(
			"127.0.0.1:54000",
			"rai-beta.raiblocks.net:54000",
		),
		Header: Packets.Header{
			MagicNumber:  82,
			NetworkType:  Packets.Beta,
			VersionMax:   0,
			VersionUsing: 255,
			VersionMin:   255,
		},
	},

	Storage: ConfigStorage{
		Folder: "Nanollet-DEBUG",
	},
}
