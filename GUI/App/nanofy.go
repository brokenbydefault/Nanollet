package App

import (
	"github.com/brokenbydefault/Nanollet/Nanofy"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Node"
	"github.com/brokenbydefault/Nanollet/Block"
)

type NanofyApp struct{}

func (c *NanofyApp) Name() string {
	return "nanofy"
}

func (c *NanofyApp) HaveSidebar() bool {
	return true
}

func (c *NanofyApp) Pages() []DOM.Page {
	return []DOM.Page{
		&PageSign{},
		&PageVerify{},
	}
}

type PageSign struct{}

func (c *PageSign) Name() string {
	return "sign"
}

func (c *PageSign) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PageSign) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	file, err := dom.GetFileOf(".filepath")
	if err != nil {
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

	nameBox, _ := dom.SelectFirstElement(".name")
	nameBox.SetHTML("Drop the file here", DOM.InnerReplaceContent)

	dom.ApplyFor(".filepath", DOM.ClearValue)
}

type PageVerify struct{}

func (c *PageVerify) Name() string {
	return "verify"
}

func (c *PageVerify) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PageVerify) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	addr, _ := dom.GetStringValueOf(".address")
	if addr == "" {
		return
	}

	pk, err := Wallet.Address(addr).GetPublicKey()
	if !Wallet.Address(addr).IsValid() || err != nil {
		DOM.UpdateNotification(w, "The given address is invalid")
		return
	}

	file, err := dom.GetFileOf(".filepath")
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem opening the file")
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

	nameBox, _ := dom.SelectFirstElement(".name")
	nameBox.SetHTML("Drop the file here", DOM.InnerReplaceContent)

	dom.ApplyFor(".filepath", DOM.ClearValue)
}
