package Node

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
	"time"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

func TestGetBalance(t *testing.T) {
	c := NewServer(Packets.Header{
		MagicNumber:   82,
		NetworkType:   Packets.Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Packets.Invalid,
		ExtensionType: 0,
	})
	c.Start()
	time.Sleep(2 * time.Second)

	pkx := Wallet.Address("xrb_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est").MustGetPublicKey()

	amount, err := GetBalance(c, &pkx)
	if err != nil {
		t.Error(err)
		return
	}

	if amount.Compare(Numbers.NewMin()) == 0 || !amount.IsValid() {
		t.Error("invalid amount")
	}

}

func TestGetBlock(t *testing.T) {
	c := NewServer(Packets.Header{
		MagicNumber:   82,
		NetworkType:   Packets.Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Packets.Invalid,
		ExtensionType: 0,
	})
	c.Start()
	time.Sleep(1 * time.Second)

	hash := Block.NewBlockHash(Util.SecureHexMustDecode("8FF14E8A184F43B63B048C5D20862B7A9D91DA1275BC2F77E8149633B157BCB8"))
	tx, err := GetBlock(c, &hash)
	if err != nil {
		t.Error(err)
		return
	}

	if tx == nil {
		t.Error("no block found")
		return
	}
}

func TestGetHistory(t *testing.T) {
	c := NewServer(Packets.Header{
		MagicNumber:   82,
		NetworkType:   Packets.Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Packets.Invalid,
		ExtensionType: 0,
	})
	c.Start()
	time.Sleep(1 * time.Second)

	pk := Wallet.Address("xrb_3tz9pdfskx934ce36cf6h17uspp4hzsamr5hk7u1wd6em1gfsnb618hfsafc").MustGetPublicKey()

	txs, err := GetHistory(c, &pk, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if len(txs) <= 0 {
		t.Error("blocks not found")
		return
	}

	open := txs[len(txs)-1].Hash()
	if open != Block.NewBlockHash(Util.SecureHexMustDecode("C4977221708A90790665432F51CDE3B1E248F876448B35E7EBEF9285036D90C0")) {
		t.Error("open not found")
		return
	}

}

func TestGetPendings(t *testing.T) {
	c := NewServer(Packets.Header{
		MagicNumber:   82,
		NetworkType:   Packets.Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Packets.Invalid,
		ExtensionType: 0,
	})
	c.Start()
	time.Sleep(1 * time.Second)

	pk := Wallet.Address("xrb_1nanofy8on8preceding8transaction11111111411111111111pqdyc8af").MustGetPublicKey()

	txs, err := GetPendings(c, &pk, Numbers.NewMin())
	if err != nil {
		t.Error(err)
		return
	}

	if len(txs) < 10 {
		t.Error("not enough information")
		return
	}

}
