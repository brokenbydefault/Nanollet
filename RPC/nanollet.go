package RPCClient

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC/internal"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Util"
	"errors"
	"github.com/brokenbydefault/Nanollet/RPC/rpctypes"
	"encoding/json"
)

var ErrNotOpenedAccount = errors.New("Account not found")
var ErrBadAddress = errors.New("Bad account number")
var ErrEmpty = errors.New("empty")

//@TODO Use custom UnmarshalJSON/MarshalJSON and rewrite the entire code.
func NewAccountBalance(addr Wallet.Address) (AccountsBalancesRequest) {
	return NewMultiAccountsBalance(addr)
}

// GetAccountBalance will retrieve the amount and pending from given account address. It will
// return an non-nil error if can't reach the server or the address is invalid.
func GetAccountBalance(c rpctypes.Connection, addr Wallet.Address) (AccountBalance, error) {
	reqresp, err := GetMultiAccountsBalance(c, addr)
	if err != nil {
		return AccountBalance{}, err
	}

	return reqresp[addr], nil
}

func NewMultiAccountsBalance(addr ...Wallet.Address) (AccountsBalancesRequest) {
	return AccountsBalancesRequest{
		Accounts:       addr,
		DefaultRequest: defaultRequest("accounts_balances"),
	}
}

// GetMultiAccountsBalance will do the same from GetAccountBalance but will return an map, the index
// of the map will the address. Each value from the map is the respective balance/amount.
func GetMultiAccountsBalance(c rpctypes.Connection, addr ...Wallet.Address) (MultiplesAccountsBalance, error) {
	req := NewMultiAccountsBalance(addr...)

	reqresp := AccountBalances{}
	err := c.SendRequestJSON(&req, &reqresp)
	if err == nil && len(reqresp.Balances) == 0 {
		return nil, ErrBadAddress
	}

	return reqresp.Balances, err
}

func NewAccountInformation(addr Wallet.Address) (AccountInformationRequest) {
	return AccountInformationRequest{
		Account:        addr,
		Pending:        true,
		DefaultRequest: defaultRequest("account_info"),
	}
}

func GetAccountInformation(c rpctypes.Connection, addr Wallet.Address) (AccountInformation, error) {
	resp := AccountInformation{}
	req := NewAccountInformation(addr)

	err := c.SendRequestJSON(&req, &resp)
	if resp.Error != "" {
		return resp, errors.New(resp.Error)
	}

	return resp, err
}

func NewAccountHistory(limit int, addr Wallet.Address) (AccountHistoryRequest) {
	return AccountHistoryRequest{
		Account:        addr,
		Count:          limit,
		Raw:            true,
		DefaultRequest: defaultRequest("account_history"),
	}
}

func (d *AccountHistory) UnmarshalJSON(data []byte) (err error) {
	var def []SingleHistory

	var teststring string
	err = json.Unmarshal(data, &teststring)

	if err == nil {
		*d = def
		return nil
	}

	var testerror DefaultResponse
	json.Unmarshal(data, &testerror)
	if testerror.Error != "" {
		return errors.New(testerror.Error)
	}

	err = json.Unmarshal(data, &def)
	for i, hist := range def {
		if hist.Amount == nil {
			def[i].Amount, _ = Numbers.NewRawFromString("0")
		}
	}

	*d = def
	return err
}

// GetAccountHistory will retrieve all the history (Send/Receive) from given account address. It will
// return an non-nil error if can't reach the server or the address is not opened, not
// having an "OpenBlock".
func GetAccountHistory(c rpctypes.Connection, limit int, addr Wallet.Address) (resp AccountHistory, err error) {
	var reqresp struct {
		History AccountHistory
	}

	req := NewAccountHistory(limit, addr)
	err = c.SendRequestJSON(req, &reqresp)

	return reqresp.History, err
}

func NewAccountPending(limit int, minimum *Numbers.RawAmount, addr Wallet.Address) (AccountsPendingRequest) {
	return NewMultiAccountsPending(limit, minimum, addr)
}

func GetAccountPending(c rpctypes.Connection, limit int, minimum *Numbers.RawAmount, addr Wallet.Address) (AccountPending, error) {
	reqresp, err := GetMultiAccountsPending(c, limit, minimum, addr)
	if err != nil {
		return AccountPending{}, err
	}

	return reqresp[addr], nil
}

func NewMultiAccountsPending(limit int, minimum *Numbers.RawAmount, addr ...Wallet.Address) (AccountsPendingRequest) {
	return AccountsPendingRequest{
		Accounts:       addr,
		Threshold:      minimum,
		Count:          limit,
		Source:         true,
		DefaultRequest: defaultRequest("accounts_pending"),
	}
}

func (d *AccountsPendingOriginal) UnmarshalJSON(data []byte) (err error) {
	var def map[string]SinglePending

	var teststring string
	err = json.Unmarshal(data, &teststring)

	if err == nil {
		*d = def
		return nil
	}

	err = json.Unmarshal(data, &def)
	*d = def
	return err
}

func GetMultiAccountsPending(c rpctypes.Connection, limit int, minimum *Numbers.RawAmount, addr ...Wallet.Address) (MultiplesAccountsPending, error) {
	var reqresp struct {
		Blocks map[Wallet.Address]AccountsPendingOriginal
		DefaultResponse
	}

	resp := MultiplesAccountsPending{}
	req := NewMultiAccountsPending(limit, minimum, addr...)

	err := c.SendRequestJSON(&req, &reqresp)
	if reqresp.Error != "" {
		return resp, errors.New(reqresp.Error)
	}

	for addr, blks := range reqresp.Blocks {
		resp[Wallet.Address(addr)] = AccountPending{}

		for hash, blk := range blks {
			blk.Hash, _ = Util.UnsafeHexDecode(hash)
			resp[Wallet.Address(addr)] = append(resp[Wallet.Address(addr)], blk)
		}
	}

	return resp, err
}

// BroadcastBlock will perform the PoW and broadcast the block to the network
func BroadcastBlock(c rpctypes.Connection, block Block.BlockTransaction) (resp ProcessBlock, err error) {
	if block == nil {
		return resp, errors.New("invalid block")
	}

	block.CreateProof()

	blk, err := block.Serialize()
	if err != nil {
		return
	}

	req := internal.ProcessBlockRequest{}
	req.App = "nanollet"
	req.Action = "process"
	req.Block = string(blk)

	reqresp := internal.ProcessBlockResponse{}
	err = c.SendRequestJSON(&req, &reqresp)
	if err != nil {
		return
	}

	if reqresp.Hash == "" {
		return resp, errors.New("error to process the block")
	}

	h, err := Util.UnsafeHexDecode(reqresp.Hash)
	if err != nil {
		return
	}

	resp.Hash = h
	return
}

func NewBlockByStringHash(hash string) (req internal.RetrieveBlockRequest) {
	req.App = "nanollet"
	req.Action = "block"
	req.Hash = hash

	return
}

func GetBlockByStringHash(c rpctypes.Connection, hash string) (Block.UniversalBlock, error) {
	req := NewBlockByStringHash(hash)

	reqresp := internal.RetrieveBlockResponse{}
	err := c.SendRequestJSON(&req, &reqresp)
	if err != nil {
		return Block.UniversalBlock{}, err
	}

	return Block.NewBlockFromJSON([]byte(reqresp.Contents))
}

func GetBlockByHash(c rpctypes.Connection, hash []byte) (Block.UniversalBlock, error) {
	return GetBlockByStringHash(c, Util.UnsafeHexEncode(hash))
}

func defaultRequest(action string) DefaultRequest {
	return DefaultRequest{
		Action: action,
		App:    "nanollet",
	}
}
