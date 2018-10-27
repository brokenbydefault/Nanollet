package Nanofy

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Node"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Block"
	"time"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
)

func TestVersion0_VerifyBlock(t *testing.T) {

	con := Node.NewServer(Packets.Header{
		MagicNumber:   82,
		NetworkType:   Packets.Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Packets.Invalid,
		ExtensionType: 0,
	}, &Storage.PeerStorage, &Storage.TransactionStorage)

	Node.NewHandler(con).Start()

	time.Sleep(2 * time.Second)

	flaghash := Block.NewBlockHash(Util.SecureHexMustDecode("AE95D7936A23D0671DD4E7E0736612F5304A18AD80B2827B2D69A3482A38F1EA"))
	flagblock, err := Node.GetBlock(con, &flaghash)
	if err != nil {
		t.Error(err)
		return
	}

	sighash := Block.NewBlockHash(Util.SecureHexMustDecode("7F0AEA1B6F2E9FD60B120DF01B0FBF4CC1B4B539A1D6B69DC50EAE81FE1A72E7"))
	sigblock, err := Node.GetBlock(con, &sighash)
	if err != nil {
		t.Error(err)
		return
	}
	
	pk := Wallet.Address("xrb_3w73pgb33ht1ws7hwaek5ywyjdteoj4qmcrzayiogpbabbo3i49dkerosn1z").MustGetPublicKey()
	nanofier, err := NewLegacyVerifier(nil, &pk, sigblock, flagblock)
	if err != nil {
		t.Error(err)
		return
	}

	if !nanofier.IsCorrectlyFormatted() {
		t.Error("one valid block was report as invalid")
	}

}

func TestVersion0_VerifyBlock_Invalid(t *testing.T) {
	con := Node.NewServer(Packets.Header{
		MagicNumber:   82,
		NetworkType:   Packets.Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Packets.Invalid,
		ExtensionType: 0,
	}, &Storage.PeerStorage, &Storage.TransactionStorage)

	Node.NewHandler(con).Start()
	time.Sleep(1 * time.Second)

	flaghash := Block.NewBlockHash(Util.SecureHexMustDecode("AE95D7936A23D0671DD4E7E0736612F5304A18AD80B2827B2D69A3482A38F1EA"))
	flagblock, err := Node.GetBlock(con, &flaghash)
	if err != nil {
		t.Error(err)
		return
	}

	sighash := Block.NewBlockHash(Util.SecureHexMustDecode("72DC2B79C307600FE8187521DB2C0AAA2929D6E10C1E3E3B058ACB6B617EB019"))
	sigblock, err := Node.GetBlock(con, &sighash)
	if err != nil {
		t.Error(err)
		return
	}

	pk := Wallet.Address("xrb_3w73pgb33ht1ws7hwaek5ywyjdteoj4qmcrzayiogpbabbo3i49dkerosn1z").MustGetPublicKey()
	nanofier, err := NewLegacyVerifier(nil, &pk, sigblock, flagblock)
	if err != nil {
		t.Error(err)
		return
	}

	if nanofier.IsCorrectlyFormatted() {
		t.Error("one valid block was report as invalid")
	}

}
