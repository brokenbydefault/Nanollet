// +build !js

package App

import (
	"github.com/brokenbydefault/Nanofy"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"os"
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

func (c *PageSign) OnContinue(w *window.Window, _ string) {
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

	blks, err := Nanofy.NewNanofierVersion1().CreateBlock(file, Storage.SK, Storage.Representative, Storage.Amount, &Storage.LastBlock)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlocksToQueue(blks, )
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

func (c *PageVerify) OnContinue(w *window.Window, _ string) {
	page := DOM.SetSector(c)

	addr, _ := page.GetStringValue(w, ".address")
	filepath, err := page.GetStringValue(w, ".filepath")
	if addr == "" || filepath == "" || err != nil {
		return
	}

	pk, err := Wallet.Address(addr).GetPublicKey()
	if !Wallet.Address(addr).IsValid() || err != nil {
		DOM.UpdateNotification(w, "The given address is invalid")
		return
	}

	file, err := os.Open(filepath[7:])
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	statfile, err := file.Stat()
	if statfile.IsDir() || err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	filehash, err := Nanofy.CreateHash(file)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem hashing the file")
		return
	}

	retriveblk, err := RPCClient.GetBlockByFile(Connectivity.Socket, filehash, pk)
	if err != nil || retriveblk.Error != "" {
		DOM.UpdateNotification(w, "There was a problem retrieving the information")
		return
	}

	if !retriveblk.Exist {
		DOM.UpdateNotification(w, "Wrong! This address never had signed this file")
		return
	}

	var blks = make([]Block.UniversalBlock, 3)
	var previous = retriveblk.FlagHash
	for i := range blks {
		blks[i], err = RPCClient.GetBlockByHash(Connectivity.Socket, previous)
		if err != nil {
			DOM.UpdateNotification(w, "There was a problem retrieving the information"+err.Error())
			return
		}

		previous = blks[i].Previous
	}

	nanofier, err := Nanofy.NewNanofierFromFlagBlock(&blks[0])
	if err != nil {
		DOM.UpdateNotification(w, "Nanofy version not supported")
		return
	}

	if nanofier.VerifySignature(&pk, &blks[0], &blks[1], &blks[2], filehash) {
		DOM.UpdateNotification(w, "Correct! This address signs this given file")
	} else {
		DOM.UpdateNotification(w, "Wrong! This address never had signed this file")
	}

	namebox, _ := page.SelectFirstElement(w, ".name")
	namebox.SetHtml("Drop the file here", sciter.SIH_REPLACE_CONTENT)

	page.ApplyForIt(w, ".filepath", DOM.ClearValue)
}
