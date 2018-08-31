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

	SetWork(pow ProofWork.Work)
	SetSignature(sig Wallet.Signature)
	SetFrontier(hash BlockHash)
	SetBalance(balance *Numbers.RawAmount)
}

//--------------

type DefaultBlock struct {
	mainType BlockType
	subType   BlockType
	PoW       ProofWork.Work   `json:"work"`
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