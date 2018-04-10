package RPCClient

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

type SubscribeRequest struct {
	PublicKey Wallet.PublicKey `json:"pk"`
	DefaultRequest
}

type UnsubscribeRequest struct {
	DefaultRequest
}

type Subscription struct {
	PublicKey Wallet.PublicKey `json:"pk"`
}

type CallbackResponse struct {
	Hash Block.BlockHash
	Origin Wallet.PublicKey
	Destination Wallet.PublicKey
	Amount *Numbers.RawAmount
	Block []byte

	DefaultResponse
}
