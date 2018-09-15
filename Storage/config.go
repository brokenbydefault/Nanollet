package Storage

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Node/Peer"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
)

var Configuration Config = Live

type Config struct {
	DebugStatus bool

	Account struct {
		Representative Wallet.PublicKey
		Quorum         Peer.Quorum
		//DefaultPrefixes       []string
		MinimumAmount *Numbers.RawAmount
	}

	Node struct {
		Peers  []*Peer.Peer
		Header Packets.Header
		//DefaultMinimumWork int64
		//DefaultGenesis
	}

	Storage struct {
		Folder string
	}
}

var Live = Config{
	DebugStatus: false,

	Account: struct {
		Representative Wallet.PublicKey
		Quorum         Peer.Quorum
		MinimumAmount  *Numbers.RawAmount
	}{
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
				Wallet.Address("xrb_3pczxuorp48td8645bs3m6c3xotxd3idskrenmi65rbrga5zmkemzhwkaznh").MustGetPublicKey(),
				Wallet.Address("xrb_3rw4un6ys57hrb39sy1qx8qy5wukst1iiponztrz9qiz6qqa55kxzx4491or").MustGetPublicKey(),
				Wallet.Address("xrb_3jwrszth46rk1mu7rmb4rhm54us8yg1gw3ipodftqtikf5yqdyr7471nsg1k").MustGetPublicKey(),
				Wallet.Address("xrb_1niabkx3gbxit5j5yyqcpas71dkffggbr6zpd3heui8rpoocm5xqbdwq44oh").MustGetPublicKey(),
				Wallet.Address("xrb_1brainb3zz81wmhxndsbrjb94hx3fhr1fyydmg6iresyk76f3k7y7jiazoji").MustGetPublicKey(),
				Wallet.Address("xrb_1nanode8ngaakzbck8smq6ru9bethqwyehomf79sae1k7xd47dkidjqzffeg").MustGetPublicKey(),
				Wallet.Address("xrb_1tig1rio7iskejqgy6ap75rima35f9mexjazdqqquthmyu48118jiewny7zo").MustGetPublicKey(),
				Wallet.Address("xrb_16k5pimotz9zehjk795wa4qcx54mtusk8hc5mdsjgy57gnhbj3hj6zaib4ic").MustGetPublicKey(),
				Wallet.Address("xrb_3rropjiqfxpmrrkooej4qtmm1pueu36f9ghinpho4esfdor8785a455d16nf").MustGetPublicKey(),
				Wallet.Address("xrb_1i9ugg14c5sph67z4st9xk8xatz59xntofqpbagaihctg6ngog1f45mwoa54").MustGetPublicKey(),
				Wallet.Address("xrb_1x7biz69cem95oo7gxkrw6kzhfywq4x5dupw4z1bdzkb74dk9kpxwzjbdhhs").MustGetPublicKey(),
				Wallet.Address("xrb_1ninja7rh37ehfp9utkor5ixmxyg8kme8fnzc4zty145ibch8kf5jwpnzr3r").MustGetPublicKey(),
			},
			Common: 30,
			Fork:   50,
		},
		MinimumAmount: Numbers.NewRawFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD3, 0xC2, 0x1B, 0xCE, 0xCC, 0xED, 0xA1, 0x00, 0x00, 0x00}),
	},

	Node: struct {
		Peers  []*Peer.Peer
		Header Packets.Header
	}{
		Peers: Peer.NewPeersFromString(
			"rai.raiblocks.net:7075",
			"185.243.9.164:7075",
			"206.189.190.7:7075",
			"198.245.55.107:7075",
		),
		Header: Packets.Header{
			MagicNumber:  82,
			NetworkType:  Packets.Live,
			VersionMax:   14,
			VersionUsing: 14,
			VersionMin:   7,
		},
	},

	Storage: struct {
		Folder string
	}{
		Folder: "Nanollet",
	},
}

var Beta = Config{
	DebugStatus: false,

	Account: struct {
		Representative Wallet.PublicKey
		Quorum         Peer.Quorum
		MinimumAmount  *Numbers.RawAmount
	}{
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
	},

	Node: struct {
		Peers  []*Peer.Peer
		Header Packets.Header
	}{
		Peers: Peer.NewPeersFromString(
			"127.0.0.1:54000",
			"rai-beta.raiblocks.net:54000",
		),
		Header: Packets.Header{
			MagicNumber:  82,
			NetworkType:  Packets.Beta,
			VersionMax:   14,
			VersionUsing: 14,
			VersionMin:   7,
		},
	},

	Storage: struct {
		Folder string
	}{
		Folder: "Nanollet-DEBUG",
	},
}
