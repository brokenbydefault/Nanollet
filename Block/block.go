package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/ProofWork"
)

type Transaction interface {
	Encode() (data []byte)
	Decode(data []byte) (err error)
	SwitchToUniversalBlock(previousBlock *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock

	Work() ProofWork.Work
	Hash() BlockHash

	GetType() BlockType
	GetSubType() BlockType
	GetTarget() (destination Wallet.PublicKey, source BlockHash)
	GetPrevious() BlockHash

	SetWork(pow []byte)
	SetSignature(sig []byte)
	SetFrontier(hash BlockHash)
	SetBalance(balance *Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	Type      BlockType        `json:"type"`
	SubType   BlockType        `json:"subtype,omitempty"`
	PoW       ProofWork.Work   `json:"work"`
	Signature Wallet.Signature `json:"signature"`
	hash      BlockHash
}

//--------------

type SendBlock struct {
	DefaultBlock
	Previous    BlockHash          `json:"previous"`
	Destination Wallet.PublicKey   `json:"destination"`
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
	Account        Wallet.PublicKey `json:"account"`
	Representative Wallet.PublicKey `json:"representative"`
	Source         BlockHash        `json:"source"`
}

//--------------

type ChangeBlock struct {
	DefaultBlock
	Previous       BlockHash        `json:"previous"`
	Representative Wallet.PublicKey `json:"representative"`
}

//--------------

type UniversalBlock struct {
	DefaultBlock
	Account        Wallet.PublicKey   `json:"account"`
	Previous       BlockHash          `json:"previous"`
	Representative Wallet.PublicKey   `json:"representative"`
	Balance        *Numbers.RawAmount `json:"balance"`
	Link           BlockHash          `json:"link"`
}
