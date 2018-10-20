package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"encoding"
)

var (
	GenesisAccount = Wallet.Address("xrb_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3").MustGetPublicKey()
	BurnAccount    = Wallet.Address("xrb_1111111111111111111111111111111111111111111111111111hifc8npp").MustGetPublicKey()
	Epoch          = [32]byte{0x65, 0x70, 0x6f, 0x63, 0x68, 0x20, 0x76, 0x31, 0x20, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
)

type Transaction interface {
	Encode() (data []byte)
	encoding.BinaryMarshaler
	Decode(data []byte) (err error)
	encoding.BinaryUnmarshaler
	SwitchToUniversalBlock(previousBlock *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock

	Work() Work
	Hash() BlockHash

	GetAccount() Wallet.PublicKey
	GetBalance() *Numbers.RawAmount
	GetType() BlockType
	GetTarget() (destination Wallet.PublicKey, source BlockHash)
	GetPrevious() BlockHash
	GetWork() Work
	GetSignature() Wallet.Signature

	IsValidPOW() bool
	//IsValidSignature() bool

	SetWork(pow Work)
	SetSignature(sig Wallet.Signature)
	SetFrontier(hash BlockHash)
	SetBalance(balance *Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	PoW       Work             `json:"work"`
	Signature Wallet.Signature `json:"signature"`
}

//--------------

type SendBlock struct {
	Previous    BlockHash          `json:"previous"`
	Destination Wallet.PublicKey   `json:"destination"`
	Balance     *Numbers.RawAmount `json:"balance"`
	DefaultBlock
}

//--------------

type ReceiveBlock struct {
	Previous BlockHash `json:"previous"`
	Source   BlockHash `json:"source"`
	DefaultBlock
}

//--------------

type OpenBlock struct {
	Account        Wallet.PublicKey `json:"account"`
	Representative Wallet.PublicKey `json:"representative"`
	Source         BlockHash        `json:"source"`
	DefaultBlock
}

//--------------

type ChangeBlock struct {
	Previous       BlockHash        `json:"previous"`
	Representative Wallet.PublicKey `json:"representative"`
	DefaultBlock
}

//--------------

type UniversalBlock struct {
	Account        Wallet.PublicKey   `json:"account"`
	Previous       BlockHash          `json:"previous"`
	Representative Wallet.PublicKey   `json:"representative"`
	Balance        *Numbers.RawAmount `json:"balance"`
	Link           BlockHash          `json:"link"`
	DefaultBlock
}
