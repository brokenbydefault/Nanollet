// +build !js

package App

import (
	"strings"

	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/OpenCAP"
	"image/color"
)

type NanolletApp struct{}

func (c *NanolletApp) Name() string {
	return "nanollet"
}

func (c *NanolletApp) HaveSidebar() bool {
	return true
}

func (c *NanolletApp) Pages() []DOM.Page {
	return []DOM.Page{
		&PageWallet{},
		&PageReceive{},
		&PageList{},
		&PageRepresentative{},
	}
}

type PageWallet struct{}

func (c *PageWallet) Name() string {
	return "send"
}

func (c *PageWallet) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PageWallet) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	var errs []error
	whole, err := dom.GetStringValueOf(".whole")
	errs = append(errs, err)

	decimal, err := dom.GetStringValueOf(".decimal")
	errs = append(errs, err)

	addrOrAlias, err := dom.GetStringValueOf(".address")
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
	dom.ApplyFor(".whole, .decimal, .address", DOM.ClearValue)
}

type PageReceive struct{}

func (c *PageReceive) Name() string {
	return "receive"
}

func (c *PageReceive) OnView(w *DOM.Window, dom *DOM.DOM) {
	addr := Storage.AccountStorage.PublicKey.CreateAddress()

	textarea, _ := dom.SelectFirstElement(".address")
	textarea.SetValue(string(addr))

	textarea.Apply(DOM.ReadOnlyElement)

	qr, err := addr.QRCode(175, color.RGBA{220, 220, 223, 1})
	if err != nil {
		return
	}

	qrSpace, err := dom.SelectFirstElement(".qrcode")
	if err != nil {
		return
	}
	qrSpace.Apply(DOM.ClearHTML)
	qrSpace.CreateQRCode(qr)
}

func (c *PageReceive) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	// no-op
}

type PageRepresentative struct{}

func (c *PageRepresentative) Name() string {
	return "representative"
}

func (c *PageRepresentative) OnView(w *DOM.Window, dom *DOM.DOM) {
	current, err := dom.SelectFirstElement(".currentRepresentative")
	if err != nil {
		return
	}

	current.SetValue(string(Storage.AccountStorage.Representative.CreateAddress()))
}

func (c *PageRepresentative) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	addr, err := dom.GetStringValueOf(".address")
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

	current, err := dom.SelectFirstElement("currentRepresentative")
	if err != nil {
		return
	}

	current.SetText(string(Storage.AccountStorage.Representative.CreateAddress()))

	dom.ApplyFor(".address", DOM.ClearValue)
}

type PageList struct{}

func (c *PageList) Name() string {
	return "history"
}

func (c *PageList) OnView(w *DOM.Window, dom *DOM.DOM) {
	balance, _ := Numbers.NewHumanFromRaw(Storage.AccountStorage.Balance).ConvertToBase(Numbers.MegaXRB, int(Numbers.MegaXRB))
	display, _ := dom.SelectFirstElement(".fullamount")
	display.SetValue(balance)

	if Storage.TransactionStorage.Count() == 0 {
		return
	}

	txbox, _ := dom.SelectFirstElement(".txbox")
	txbox.Apply(DOM.ClearHTML)

	for i, tx := range Storage.TransactionStorage.GetByFrontier(Storage.AccountStorage.Frontier) {

		hashPrev := tx.GetPrevious()
		txPrev, _ := Storage.TransactionStorage.GetByHash(&hashPrev)

		txType := Block.GetSubType(tx, txPrev)
		txAmount := Block.GetAmount(tx, txPrev)

		humanAmount, err := Numbers.NewHumanFromRaw(txAmount).ConvertToBase(Numbers.MegaXRB, 6)
		if err != nil {
			return
		}

		txdiv := txbox.CreateElementWithAttr("div", "", DOM.Attrs{"class": "item"})

		txdiv.CreateElementWithAttr("div", strings.ToUpper(txType.String()), DOM.Attrs{"class": "type"})
		txdiv.CreateElementWithAttr("div", humanAmount, DOM.Attrs{"class": "amount"})

		if i == 4 {
			break
		}
	}

}

func (c *PageList) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	//no-op
}
