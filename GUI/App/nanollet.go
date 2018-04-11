package App

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/sciter-sdk/go-sciter"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"strings"
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

	var errors []error
	whole, err := page.GetStringValue(w, ".whole")
	errors = append(errors, err)

	decimal, err := page.GetStringValue(w, ".decimal")
	errors = append(errors, err)

	addr, err := page.GetStringValue(w, ".address")
	errors = append(errors, err)

	if Util.CheckError(errors) != nil {
		return
	}

	address := Wallet.Address(addr)
	if address == "" || (whole == "" && decimal == "") {
		return
	}

	if !address.IsValid() {
		DOM.UpdateNotification(w, "The given address is invalid")
		return
	}

	if !Util.StringIsNumeric(whole) || !Util.StringIsNumeric(decimal) {
		DOM.UpdateNotification(w, "The given amount is incorrect")
		return
	}

	amm, err := Numbers.NewHumanFromString(whole+"."+decimal, Numbers.MegaXRB).ConvertToRawAmount()
	if err != nil {
		return
	}

	if !Storage.Amount.Subtract(amm).IsValid() {
		DOM.UpdateNotification(w, "The given amount is higher than the maximum")
		return
	}

	blk, err := Block.CreateSignedSendBlock(&Storage.SK, amm, Storage.Amount, Storage.Frontier, &address)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlockToQueue(blk, amm)
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
	textarea.SetValue(sciter.NewValue(string(Storage.PK.CreateAddress())))
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

	if !address.IsValid() {
		DOM.UpdateNotification(w, "The given address is invalid")
		return
	}

	blk, err := Block.CreateSignedChangeBlock(&Storage.SK, Storage.Frontier, &address)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlockToQueue(blk)
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

	balance, _ := Numbers.NewHumanFromRaw(Storage.Amount).ConvertToBase(Numbers.MegaXRB, int(Numbers.MegaXRB))
	display, _:= page.SelectFirstElement(w, ".fullamount")
	display.SetValue(sciter.NewValue(balance))

	if len(Storage.History) == 0 {
		return
	}

	txbox, _ := page.SelectFirstElement(w, ".txbox")
	DOM.ClearHTML(txbox)

	for i, hist := range Storage.History {
		tx := hist

		amm, err := Numbers.NewHumanFromRaw(tx.Amount).ConvertToBase(Numbers.MegaXRB, 6)
		if err != nil {
			return
		}

		blktype := tx.Type
		if tx.SubType != "" {
			blktype = tx.SubType
		}

		txdiv := DOM.CreateElementAppendTo("div", "", "item", "", txbox)

		DOM.CreateElementAppendTo("div", strings.ToUpper(blktype), "type", "", txdiv)
		DOM.CreateElementAppendTo("div", amm, "amount", "", txdiv)

		if i == 4 {
			break
		}
	}

}

func (c *PageList) OnContinue(w *window.Window, _ string) {
	//no-op
}
