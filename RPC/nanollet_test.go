package RPCClient

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"testing"
)

func TestGetAccountBalance(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z")

	accountBalance, err := GetAccountBalance(Connectivity.Socket, addr)
	if err != nil {
		t.Error(err)
		return
	}

	n := Numbers.NewHumanFromRaw(accountBalance.Balance)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = n.ConvertToBase(Numbers.MegaXRB, 6)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestGetAccountBalanceInvalid(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1nanofy1the1next1tra111111111is1an1pkregister1u11111ashk5pbd")

	_, err := GetAccountBalance(Connectivity.Socket, addr)
	if err == nil {
		t.Error(err)
	}

	if err != ErrBadAddress {
		t.Error(err)
		return
	}

}

func TestGetMultipleAccountsBalance(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z")
	addr2 := Wallet.Address("xrb_1nanofy1the1next1transaction1is1an1pkregister1u11111ashk5pbd")

	accountBalance, err := GetMultiAccountsBalance(Connectivity.Socket, addr, addr2)
	if err != nil {
		t.Error(err)
		return
	}

	n := Numbers.NewHumanFromRaw(accountBalance[addr].Balance)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = n.ConvertToBase(Numbers.MegaXRB, 6)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAccountInformation(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z")
	_, err := GetAccountInformation(Connectivity.Socket, addr)
	if err != nil {
		t.Error(err)
		return
	}

	addr = Wallet.Address("xrb_1nanofy1the1next1transaction1is1an1pkregister1u11111ashk5pbd")
	_, err = GetAccountInformation(Connectivity.Socket, addr)
	if err.Error() != ErrNotOpenedAccount.Error() {
		t.Error(err)
		return
	}

}

func TestGetAccountHistory(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z")
	_, err := GetAccountHistory(Connectivity.Socket, 100, addr)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAccountHistoryNotExist(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1nanofy1the1next1transaction1is1an1pkregister1u11111ashk5pbd")
	_, err := GetAccountHistory(Connectivity.Socket, 100, addr)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestGetMultiAccountsPending(t *testing.T) {

	StartWebsocket()

	pk, _, _ := Wallet.GenerateRandomKeyPair()

	addr := pk.CreateAddress()
	addr2 := Wallet.Address("xrb_1nanofy8on8preceding8transaction11111111111111111111chcdnjcj")
	amm, err := Numbers.NewRawFromString("0")

	r, err := GetMultiAccountsPending(Connectivity.Socket, 25, amm, addr, addr2)
	if err != nil {
		t.Error(err)
		return
	}

	if len(r[addr2]) <= 0 {
		t.Error("error")
		return
	}
}

func TestGetMultiAccountsPendingOverLimit(t *testing.T) {

	StartWebsocket()

	pk, _, _ := Wallet.GenerateRandomKeyPair()

	addr := pk.CreateAddress()
	addr2 := Wallet.Address("xrb_1nanofy8on8preceding8transaction11111111111111111111chcdnjcj")
	amm, err := Numbers.NewRawFromString("0")

	_, err = GetMultiAccountsPending(Connectivity.Socket, 100, amm, addr, addr2)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAccountPending(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1nanofy8on8preceding8transaction11111111111111111111chcdnjcj")
	amm, err := Numbers.NewRawFromString("0")

	r, err := GetAccountPending(Connectivity.Socket, 25, amm, addr)
	if err != nil {
		t.Error(err)
		return
	}

	if len(r) <= 0 {
		t.Error("error")
		return
	}

	if r[0].Source == "" {
		t.Error("error")
		return
	}

}

func TestGetAccountPendingOverLimit(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_1nanofy8on8preceding8transaction11111111111111111111chcdnjcj")
	amm, err := Numbers.NewRawFromString("0")

	_, err = GetAccountPending(Connectivity.Socket, 100, amm, addr)
	if err != err {
		t.Error(err)
		return
	}

}

func TestGetAccountPendingInvalid(t *testing.T) {

	StartWebsocket()

	addr := Wallet.Address("xrb_3yxxyyrdeapnxe1dyxxxxxxxsyhxou4risfkzibe8bfdjj3663d7ppy1enb")
	amm, err := Numbers.NewRawFromString("0")

	_, err = GetAccountPending(Connectivity.Socket, 25, amm, addr)
	if err.Error() != ErrBadAddress.Error() {
		t.Error(err)
		return
	}

}

func TestSendBlock(t *testing.T) {

	StartWebsocket()

	pk, sk, _ := Wallet.RecoverKeyPairFromClassicalSeed("165F4C1726F8245BE38410D524A77D8B1CB2687CD5F74377475A44506DD90C74", 0)
	addr := pk.CreateAddress()

	info, err := GetAccountInformation(Connectivity.Socket, addr)
	if err != nil {
		t.Error(err)
		return
	}

	block := &Block.ChangeBlock{}
	block.Representative = addr
	block.Type = "change"
	block.Previous = info.Frontier
	block.Signature, _ = sk.CreateSignature(block.Hash())

	_, err = BroadcastBlock(Connectivity.Socket, block)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestSendBlock2(t *testing.T) {

	StartWebsocket()

	pk, sk, _ := Wallet.RecoverKeyPairFromClassicalSeed("165F4C1726F8245BE38410D524A77D8B1CB2687CD5F74377475A44506DD90C74", 0)
	addr := pk.CreateAddress()

	info, err := GetAccountInformation(Connectivity.Socket, addr)
	if err != nil {
		t.Error(err)
		return
	}

	min, _ := Numbers.NewRawFromString("0")
	pend, err := GetAccountPending(Connectivity.Socket, 10, min, addr)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range pend {
		block := &Block.ReceiveBlock{}
		block.Type = "receive"
		block.Previous = info.Frontier
		block.Source = v.Hash
		block.Signature, _ = sk.CreateSignature(block.Hash())
		block.PoW = block.Work()

		_, err := BroadcastBlock(Connectivity.Socket, block)
		if err != nil {
			t.Error(err)
			return
		}

	}

}

func TestGetBlockByStringHash(t *testing.T) {

	StartWebsocket()

	_, err := GetBlockByStringHash(Connectivity.Socket, "505726E8B8FABFB714823E38A7BCAF90D9DBA700B36E71DFE7836850B1070EB3")
	if err != nil {
		t.Error(err)
		return
	}

}
