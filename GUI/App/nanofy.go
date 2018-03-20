package App

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"os"
	"github.com/brokenbydefault/Nanofy"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/sciter-sdk/go-sciter"
)

type NanofyApp guitypes.App

func (c *NanofyApp) Name() string {
	return "nanofy"
}

func (c *NanofyApp) HaveSidebar() bool {
	return true
}

func (c *NanofyApp) Display() Front.HTMLPAGE {
	return Front.HTMLNanofy
}

func (c *NanofyApp) Pages() []guitypes.Page {
	return []guitypes.Page{
		&PageSign{},
		&PageVerify{},
	}
}

type PageSign guitypes.Sector

func (c *PageSign) Name() string {
	return "sign"
}

func (c *PageSign) OnView(w *window.Window) {
	// no-op
}

func (c *PageSign) OnContinue(w *window.Window) {
	page := DOM.SetSector(c)

	filepath, err := page.GetStringValue(w, ".filepath")
	if filepath == "" || err != nil {
		return
	}

	file, err := os.Open(filepath[7:])
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	blks, err := Nanofy.NV0.CreateFileBlocks(file, &Storage.SK, Storage.Amount, Storage.Frontier)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	tx := make([]Block.BlockTransaction, 2)
	for i, blk := range blks {
		tx[i] = blk
	}

	err = Background.PublishBlocksToQueue(tx, Nanofy.NV0.Value())
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem sending a block")
		return
	}

	DOM.UpdateAmount(w)
	DOM.UpdateNotification(w, "Your signature was sent successfully.")

	namebox, _ := page.SelectFirstElement(w, ".name")
	namebox.SetHtml("Drop the file here", sciter.SIH_REPLACE_CONTENT)

	page.ApplyForIt(w, ".filepath", DOM.ClearValue)
}

type PageVerify guitypes.Sector

func (c *PageVerify) Name() string {
	return "verify"
}

func (c *PageVerify) OnView(w *window.Window) {
	// no-op
}

func (c *PageVerify) OnContinue(w *window.Window) {
	page := DOM.SetSector(c)

	filepath, err := page.GetStringValue(w, ".filepath")
	if filepath == "" || err != nil {
		return
	}

	addr, _ := page.GetStringValue(w, ".address")
	if !Wallet.Address(addr).IsValid() {
		return
	}

	file, err := os.Open(filepath[7:])
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	filehash, err := Nanofy.CreateFileHash(file)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem hashing the file")
		return
	}

	pk, _ := Wallet.Address(addr).GetPublicKey()
	blks, err := RPCClient.GetBlockByFile(Connectivity.Socket, filehash, pk)
	if err != nil || blks.Error != "" {
		DOM.UpdateNotification(w, "There was a problem retrieving the information")
		return
	}

	if !blks.Exist {
		DOM.UpdateNotification(w, "Wrong! This address never had signed this file")
		return
	}

	flagblk, errb1 := RPCClient.GetBlockByHash(Connectivity.Socket, blks.FlagHash)
	sigblk, errb2 := RPCClient.GetBlockByHash(Connectivity.Socket, blks.SigHash)

	if errb1 != nil || errb2 != nil {
		DOM.UpdateNotification(w, "There was a problem retrieving the information")
		return
	}


	if Nanofy.VerifySignature(&pk, filehash, flagblk, sigblk) {
		DOM.UpdateNotification(w, "Correct! This address signs this given file")
	}else{
		DOM.UpdateNotification(w, "Wrong! This address never had signed this file")
	}

	namebox, _ := page.SelectFirstElement(w, ".name")
	namebox.SetHtml("Drop the file here", sciter.SIH_REPLACE_CONTENT)

	page.ApplyForIt(w, ".filepath", DOM.ClearValue)
}
