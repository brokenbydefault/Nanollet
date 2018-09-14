package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"encoding"
)

type Transaction interface {
	Encode() (data []byte)
	encoding.BinaryMarshaler
	Decode(data []byte) (err error)
	encoding.BinaryUnmarshaler
	SwitchToUniversalBlock(previousBlock *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock

	Work() Work
	Hash() BlockHash

	GetBalance() *Numbers.RawAmount
	GetType() BlockType
	GetTarget() (destination Wallet.PublicKey, source BlockHash)
	GetPrevious() BlockHash
	GetWork() Work
	GetSignature() Wallet.Signature

	SetWork(pow Work)
	SetSignature(sig Wallet.Signature)
	SetFrontier(hash BlockHash)
	SetBalance(balance *Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	PoW       Work             `json:"work"`
	Signature Wallet.Signature `json:"signature"`
	hash      BlockHash
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
