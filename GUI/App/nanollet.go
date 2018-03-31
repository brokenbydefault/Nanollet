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
		&PageList{},
	}
}

type PageWallet guitypes.Sector

func (c *PageWallet) Name() string {
	return "wallet"
}

func (c *PageWallet) OnView(w *window.Window) {
	// no-op
}

func (c *PageWallet) OnContinue(w *window.Window) {
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

	raw, err := Numbers.NewHumanFromString(whole+"."+decimal, Numbers.MegaXRB).ConvertToRawAmount()
	if err != nil {
		return
	}

	if !Storage.Amount.Subtract(raw).IsValid() {
		DOM.UpdateNotification(w, "The given amount is higher than the maximum")
		return
	}

	blk, err := Block.CreateSignedSendBlock(&Storage.SK, raw, Storage.Amount, Storage.Frontier, &address)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem creating a block")
		return
	}

	err = Background.PublishBlockToQueue(blk, raw)
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

func (c *PageReceive) OnContinue(w *window.Window) {
	// no-op
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

		txdiv, _ := sciter.CreateElement("div", "")
		txdiv.SetAttr("class", tx.Type)
		txbox.Append(txdiv)

		tpdiv, _ := sciter.CreateElement("div", "")
		tpdiv.SetAttr("class", "type")
		txdiv.Append(tpdiv)

		amm, err := Numbers.NewHumanFromRaw(tx.Amount).ConvertToBase(Numbers.MegaXRB, 6)
		if err != nil {
			return
		}

		ammdiv, _ := sciter.CreateElement("div", amm)
		ammdiv.SetAttr("class", "amount")
		txdiv.Append(ammdiv)

		if i == 4 {
			break
		}
	}

}

func (c *PageList) OnContinue(w *window.Window) {
	//no-op
}
