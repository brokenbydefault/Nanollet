package Block

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

type BlockTransaction interface {
	Serialize() ([]byte, error)
	CreateProof() []byte
	Hash() BlockHash

	GetType() string

	SetWork([]byte)
	SetSignature([]byte)
	SetFrontier(BlockHash)
	SetBalance(*Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	Type      string `json:"type" serialize:"string"`
	Work      []byte `json:"work" serialize:"hex"`
	Signature []byte `json:"signature" serialize:"hex"`
}

type SerializedDefaultBlock struct {
	Type      string `json:"type"`
	Work      string `json:"work" serialize:"hex"`
	Signature string `json:"signature" serialize:"hex"`
}

//--------------

type SendBlock struct {
	Previous    []byte
	Destination Wallet.PublicKey
	Balance     *Numbers.RawAmount
	DefaultBlock
}

type SerializedSendBlock struct {
	Previous    string `json:"previous" serialize:"hex"`
	Destination string `json:"destination"`
	Balance     string `json:"balance" serialize:"hex"`
	SerializedDefaultBlock
}

//--------------

type ReceiveBlock struct {
	Previous []byte
	Source   []byte
	DefaultBlock
}

type SerializedReceiveBlock struct {
	Previous string `json:"previous" serialize:"hex"`
	Source   string `json:"source" serialize:"hex"`
	SerializedDefaultBlock
}

//--------------

type OpenBlock struct {
	Source         []byte
	Representative Wallet.PublicKey
	Account        Wallet.PublicKey
	DefaultBlock
}

type SerializedOpenBlock struct {
	Source         string `json:"source" serialize:"hex"`
	Representative string `json:"representative" serialize:"hex"`
	Account        string `json:"account" serialize:"hex"`
	SerializedDefaultBlock
}

//--------------

type ChangeBlock struct {
	Previous       []byte
	Representative Wallet.PublicKey
	DefaultBlock
}

type SerializedChangeBlock struct {
	Previous       string `json:"previous" serialize:"hex"`
	Representative string `json:"representative" serialize:"hex"`
	SerializedDefaultBlock
}

//--------------

type SerializedUniversalBlock struct {
	Account        string `json:"account" serialize:"hex"`
	Previous       string `json:"previous" serialize:"hex"`
	Representative string `json:"representative" serialize:"hex"`
	Balance        string `json:"balance" serialize:"hex"`
	Amount         string `json:"amount" serialize:"hex"`

	Link string `json:"link" serialize:"hex"` // Link or Target?!

	Destination string `json:"destination"`
	Source      string `json:"source" serialize:"hex"`

	SerializedDefaultBlock
}

type UniversalBlock struct {
	Account        Wallet.PublicKey
	Previous       []byte
	Representative Wallet.PublicKey
	Balance        *Numbers.RawAmount
	Amount         *Numbers.RawAmount

	Destination Wallet.PublicKey `json:"destination"`
	Source      []byte           `json:"source" serialize:"hex"`

	DefaultBlock
}
