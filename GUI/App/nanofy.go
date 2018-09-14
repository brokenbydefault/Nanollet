// +build !js

package App

import (
	"github.com/brokenbydefault/Nanollet/Nanofy"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"os"
	"github.com/brokenbydefault/Nanollet/Node"
	"github.com/brokenbydefault/Nanollet/Block"
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

	filePath, err := page.GetStringValue(w, ".filepath")
	if filePath == "" || err != nil {
		return
	}

	file, err := os.Open(filePath[7:])
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	stats, err := file.Stat()
	if err != nil || stats.IsDir() {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	previous, ok := Storage.TransactionStorage.GetByHash(&Storage.AccountStorage.Frontier)
	if !ok {
		DOM.UpdateNotification(w, "Previous block not found")
		return
	}

	nanofier, err := Nanofy.NewStateSigner(file, &Storage.AccountStorage.SecretKey, previous)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	blks, err := nanofier.CreateBlocks()
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlocksToQueue(blks, Block.Send, nanofier.Amount())
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem sending a block")
		return
	}

	DOM.UpdateAmount(w)
	DOM.UpdateNotification(w, "Your signature was sent successfully.")

	nameBox, _ := page.SelectFirstElement(w, ".name")
	nameBox.SetHtml("Drop the file here", sciter.SIH_REPLACE_CONTENT)

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
	filePath, err := page.GetStringValue(w, ".filepath")
	if addr == "" || filePath == "" || err != nil {
		return
	}

	pk, err := Wallet.Address(addr).GetPublicKey()
	if !Wallet.Address(addr).IsValid() || err != nil {
		DOM.UpdateNotification(w, "The given address is invalid")
		return
	}

	file, err := os.Open(filePath[7:])
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
		return
	}

	stats, err := file.Stat()
	if err != nil || stats.IsDir() {
		DOM.UpdateNotification(w, "Only files can be signed")
		return
	}

	txs, err := Node.GetHistory(Background.Connection, &pk, nil)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem retrieving the information")
		return
	}

	if Nanofy.VerifyFromHistory(file, pk, txs) {
		DOM.UpdateNotification(w, "Correct! This address signs this given file")
	} else {
		DOM.UpdateNotification(w, "Wrong! This address never had signed this file")
	}

	nameBox, _ := page.SelectFirstElement(w, ".name")
	nameBox.SetHtml("Drop the file here", sciter.SIH_REPLACE_CONTENT)

	page.ApplyForIt(w, ".filepath", DOM.ClearValue)
}
