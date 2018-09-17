// +build !js

package App

import (
	"strings"

	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/OpenCAP"
)

type NanolletApp guitypes.App

func (c *NanolletApp) Name() string {
	return "nanollet"
}

func (c *NanolletApp) HaveSidebar() bool {
	return true
}

func (c *NanolletApp) Display() Front.HTMLPAGE {
	return Front.HTMLNanollet
}

func (c *NanolletApp) Pages() []guitypes.Page {
	return []guitypes.Page{
		&PageWallet{},
		&PageReceive{},
		&PageRepresentative{},
		&PageList{},
	}
}

type PageWallet guitypes.Sector

func (c *PageWallet) Name() string {
	return "send"
}

func (c *PageWallet) OnView(w *window.Window) {
	// no-op
}

func (c *PageWallet) OnContinue(w *window.Window, _ string) {
	page := DOM.SetSector(c)

	var errs []error
	whole, err := page.GetStringValue(w, ".whole")
	errs = append(errs, err)

	decimal, err := page.GetStringValue(w, ".decimal")
	errs = append(errs, err)

	addrOrAlias, err := page.GetStringValue(w, ".address")
	errs = append(errs, err)

	if Util.CheckError(errs) != nil {
		return
	}

	if addrOrAlias == "" || (whole == "" && decimal == "") {
		// Empty values not return errors
		return
	}

	var dest Wallet.PublicKey
	switch {
	case Wallet.Address(addrOrAlias).IsValid():
		if dest, err = Wallet.Address(addrOrAlias).GetPublicKey(); err != nil {
			DOM.UpdateNotification(w, "The address is wrong")
			return
		}
	case OpenCAP.Address(addrOrAlias).IsValid():
		if dest, err = OpenCAP.Address(addrOrAlias).GetPublicKey(); err != nil {
			DOM.UpdateNotification(w, "The address was not found")
			return
		}
	default:
		DOM.UpdateNotification(w, "The address invalid or it's not supported")
		return
	}

	if !Util.StringIsNumeric(whole) || !Util.StringIsNumeric(decimal) {
		DOM.UpdateNotification(w, "The given amount is incorrect")
		return
	}

	amount, err := Numbers.NewHumanFromString(whole+"."+decimal, Numbers.MegaXRB).ConvertToRawAmount()
	if err != nil {
		return
	}

	if !Storage.AccountStorage.Balance.Subtract(amount).IsValid() {
		DOM.UpdateNotification(w, "The given amount is higher than the maximum")
		return
	}

	tx, err := Block.CreateUniversalSendBlock(&Storage.AccountStorage.SecretKey, Storage.AccountStorage.Representative, Storage.AccountStorage.Balance, amount, Storage.AccountStorage.Frontier, dest)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlockToQueue(tx, Block.Send, amount)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem sending a block")
		return
	}

	DOM.UpdateAmount(w)
	DOM.UpdateNotification(w, "Your payment was sent successfully.")
	page.ApplyForIt(w, ".whole, .decimal, .address", DOM.ClearValue)
}

type PageReceive guitypes.Sector

func (c *PageReceive) Name() string {
	return "receive"
}

func (c *PageReceive) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	textarea, _ := page.SelectFirstElement(w, ".address")
	textarea.SetValue(sciter.NewValue(string(Storage.AccountStorage.PublicKey.CreateAddress())))
	DOM.ReadOnlyElement(textarea)
}

func (c *PageReceive) OnContinue(w *window.Window, _ string) {
	// no-op
}

type PageRepresentative guitypes.Sector

func (c *PageRepresentative) Name() string {
	return "representative"
}

func (c *PageRepresentative) OnView(w *window.Window) {
	// no-op
}

func (c *PageRepresentative) OnContinue(w *window.Window, _ string) {
	page := DOM.SetSector(c)

	addr, err := page.GetStringValue(w, ".address")
	if err != nil {
		return
	}

	address := Wallet.Address(addr)
	if address == "" {
		return
	}

	representative, err := address.GetPublicKey()
	if err != nil || !address.IsValid() {
		DOM.UpdateNotification(w, "The given address is invalid")
		return
	}

	blk, err := Block.CreateUniversalChangeBlock(&Storage.AccountStorage.SecretKey, representative, Storage.AccountStorage.Balance, Storage.AccountStorage.Frontier)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlockToQueue(blk, Block.Change)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem sending a block")
		return
	}

	DOM.UpdateAmount(w)
	DOM.UpdateNotification(w, "Your representative was changed successfully.")
	page.ApplyForIt(w, ".address", DOM.ClearValue)
}

type PageList guitypes.Sector

func (c *PageList) Name() string {
	return "history"
}

func (c *PageList) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	balance, _ := Numbers.NewHumanFromRaw(Storage.AccountStorage.Balance).ConvertToBase(Numbers.MegaXRB, int(Numbers.MegaXRB))
	display, _ := page.SelectFirstElement(w, ".fullamount")
	display.SetValue(sciter.NewValue(balance))

	if Storage.TransactionStorage.Count() == 0 {
		return
	}

	txbox, _ := page.SelectFirstElement(w, ".txbox")
	DOM.ClearHTML(txbox)

	for i, tx := range Storage.TransactionStorage.GetByFrontier(Storage.AccountStorage.Frontier) {

		hashPrev := tx.GetPrevious()
		txPrev, _ := Storage.TransactionStorage.GetByHash(&hashPrev)

		txType := Block.GetSubType(tx, txPrev)
		txAmount := Block.GetAmount(tx, txPrev)

		humanAmount, err := Numbers.NewHumanFromRaw(txAmount).ConvertToBase(Numbers.MegaXRB, 6)
		if err != nil {
			return
		}

		txdiv := DOM.CreateElementAppendTo("div", "", "item", "", txbox)

		DOM.CreateElementAppendTo("div", strings.ToUpper(txType.String()), "type", "", txdiv)
		DOM.CreateElementAppendTo("div", humanAmount, "amount", "", txdiv)

		if i == 4 {
			break
		}
	}

}

func (c *PageList) OnContinue(w *window.Window, _ string) {
	//no-op
}
