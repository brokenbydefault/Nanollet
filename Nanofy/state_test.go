package Nanofy

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Node"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Storage"
)

func TestVersion1_VerifyBlock(t *testing.T) {

	con := Node.NewServer(Storage.Configuration.Node.Header)
	con.Start()
	time.Sleep(1 * time.Second)

	flaghash := Block.NewBlockHash(Util.SecureHexMustDecode("A7CEB0E504E31D74CD87DB0990B9C44D8211987C464BE758C6007C6DA2D188FB"))
	flagblock, err := Node.GetBlock(con, &flaghash)
	if err != nil {
		panic(err)
	}

	sighash := Block.NewBlockHash(Util.SecureHexMustDecode("02C03EEFBFAC781971125D04EB0F28D510D7B303E86D1CAB22E0C86271862E9B"))
	sigblock, err := Node.GetBlock(con, &sighash)
	if err != nil {
		panic(err)
	}

	prevhash := Block.NewBlockHash(Util.SecureHexMustDecode("EF7A37D750DBB692F957C0EA2965B1F48EB2BC2504AE7AE07957904D83D5B267"))
	prevblock, err := Node.GetBlock(con, &prevhash)
	if err != nil {
		panic(err)
	}



	pk := Wallet.Address("xrb_31hbrc4zary87ardrg74pd6xy157z71t4b9edmrqbj5tqgxnriba5e7cf3o6").MustGetPublicKey()
	nanofier, err := NewStateVerifier(nil, &pk, prevblock, sigblock, flagblock)
	if err != nil  {
		t.Error(err)
		return 
	}

	if !nanofier.IsCorrectlyFormatted() {
		t.Error("one valid block was report as invalid")
	}

}
