package RPCClient

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
)

type DefaultRequest struct {
	Action string `json:"action"`
	App    string `json:"app,omitempty"`
}

type DefaultResponse struct {
	Error string `json:"error,omitempty"`
}

//--------------

type AccountsBalancesRequest struct {
	Accounts []Wallet.Address `json:"accounts"`
	DefaultRequest
}

type AccountBalance struct {
	Balance *Numbers.RawAmount
	Pending *Numbers.RawAmount
}

type AccountBalances struct {
	Balances map[Wallet.Address]AccountBalance
	DefaultResponse
}

type MultiplesAccountsBalance map[Wallet.Address]AccountBalance

//--------------

type AccountInformationRequest struct {
	Account Wallet.Address `json:"account"`
	Weight  bool           `json:"weight,omitempty"`
	Pending bool           `json:"pending,omitempty"`
	DefaultRequest
}

type AccountInformation struct {
	Frontier            Block.BlockHash    `json:"frontier"`
	OpenBlock           Block.BlockHash    `json:"open_block"`
	RepresentativeBlock Block.BlockHash    `json:"representative_block"`
	Balance             *Numbers.RawAmount `json:"balance"`
	BlockCount          uint64             `json:"block_count,string"`
	Pending             *Numbers.RawAmount `json:"pending,omitempty"`
	DefaultResponse
}

//--------------

type AccountHistoryRequest struct {
	Account Wallet.Address `json:"account"`
	Count   int            `json:"count"`
	Raw     bool           `json:"raw"`
	DefaultRequest
}

type SingleHistory struct {
	Hash        Block.BlockHash    `json:"hash"`
	Type        string             `json:"type"`
	SubType     string             `json:"subtype,omitempty"`
	Source      Block.BlockHash    `json:"source,omitempty"`
	Destination Wallet.Address     `json:"destination,omitempty"`
	Account     Wallet.Address     `json:"account"`
	Amount      *Numbers.RawAmount `json:"amount"`
}

type AccountHistory []SingleHistory

//--------------

type AccountsPendingRequest struct {
	Accounts  []Wallet.Address   `json:"accounts"`
	Count     int                `json:"count"`
	Threshold *Numbers.RawAmount `json:"threshold"`
	Source    bool               `json:"source"`
	DefaultRequest
}

type SinglePending struct {
	Hash   Block.BlockHash
	Amount *Numbers.RawAmount
	Source Wallet.Address
	DefaultResponse
}

type AccountPending []SinglePending
type AccountsPendingOriginal map[string]SinglePending
type MultiplesAccountsPending map[Wallet.Address]AccountPending

//--------------

type ProcessBlock struct {
	Hash Block.BlockHash
}

//--------------

type RetrieveBlock struct {
	Block Block.BlockTransaction
}
