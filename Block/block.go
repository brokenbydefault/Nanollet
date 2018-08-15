package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/ProofWork"
)

type BlockTransaction interface {
	Serialize() ([]byte, error)
	SwitchToUniversalBlock() *UniversalBlock

	Work() ProofWork.Work
	Hash() BlockHash

	GetType() BlockType
	GetSubType() BlockType
	GetTarget() (destination Wallet.Address, source BlockHash)

	SetWork([]byte)
	SetSignature([]byte)
	SetFrontier(BlockHash)
	SetBalance(*Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	Type      BlockType        `json:"type"`
	SubType   BlockType        `json:"subtype,omitempty"`
	PoW       ProofWork.Work   `json:"work"`
	Signature Wallet.Signature `json:"signature"`
}

//--------------

type SendBlock struct {
	DefaultBlock
	Previous    BlockHash          `json:"previous"`
	Destination Wallet.Address     `json:"destination"`
	Balance     *Numbers.RawAmount `json:"balance"`
}

//--------------

type ReceiveBlock struct {
	DefaultBlock
	Previous BlockHash `json:"previous"`
	Source   BlockHash `json:"source"`
}

//--------------

type OpenBlock struct {
	DefaultBlock
	Account        Wallet.Address `json:"account"`
	Representative Wallet.Address `json:"representative"`
	Source         BlockHash      `json:"source"`
}

//--------------

type ChangeBlock struct {
	DefaultBlock
	Previous       BlockHash      `json:"previous"`
	Representative Wallet.Address `json:"representative"`
}

//--------------

type UniversalBlock struct {
	DefaultBlock
	Account        Wallet.Address     `json:"account"`
	Previous       BlockHash          `json:"previous"`
	Representative Wallet.Address     `json:"representative"`
	Balance        *Numbers.RawAmount `json:"balance"`
	Link           BlockHash          `json:"link"`

	Amount      *Numbers.RawAmount `json:"amount,omitempty"`
}
