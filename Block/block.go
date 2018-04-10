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
	GetTarget() (destination Wallet.Address, source BlockHash)

	SetWork([]byte)
	SetSignature([]byte)
	SetFrontier(BlockHash)
	SetBalance(*Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	Type      string `json:"type"`
	Work      []byte `json:"work"`
	Signature []byte `json:"signature"`
}

type SerializedDefaultBlock struct {
	Type      string `json:"type"`
	Work      string `json:"work"`
	Signature string `json:"signature"`
}

//--------------

type SendBlock struct {
	Previous    []byte
	Destination Wallet.PublicKey
	Balance     *Numbers.RawAmount
	DefaultBlock
}

type SerializedSendBlock struct {
	Previous    string `json:"previous"`
	Destination string `json:"destination"`
	Balance     string `json:"balance"`
	SerializedDefaultBlock
}

//--------------

type ReceiveBlock struct {
	Previous []byte
	Source   []byte
	DefaultBlock
}

type SerializedReceiveBlock struct {
	Previous string `json:"previous"`
	Source   string `json:"source"`
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
	Source         string `json:"source"`
	Representative string `json:"representative"`
	Account        string `json:"account"`
	SerializedDefaultBlock
}

//--------------

type ChangeBlock struct {
	Previous       []byte
	Representative Wallet.PublicKey
	DefaultBlock
}

type SerializedChangeBlock struct {
	Previous       string `json:"previous"`
	Representative string `json:"representative"`
	SerializedDefaultBlock
}

//--------------

type SerializedUniversalBlock struct {
	Account        string `json:"account"`
	Previous       string `json:"previous"`
	Representative string `json:"representative"`
	Balance        string `json:"balance"`
	Amount         string `json:"amount"`

	Link           string `json:"target"`
	LinkAccount	   string `json:"link_as_account"`

	Destination string `json:"destination"`
	Source      string `json:"source"`

	SerializedDefaultBlock
}

type UniversalBlock struct {
	Account        Wallet.PublicKey
	Previous       []byte
	Representative Wallet.PublicKey
	Balance        *Numbers.RawAmount
	Amount         *Numbers.RawAmount

	Link           []byte
	LinkAccount   Wallet.Address `json:"link_as_account"`

	Destination Wallet.PublicKey
	Source      []byte

	DefaultBlock
}
