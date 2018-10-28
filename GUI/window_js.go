// +build js

package GUI

import (
	"github.com/brokenbydefault/Nanollet/GUI/App"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"honnef.co/go/js/dom"
)

func Start() {
	win := DOM.NewWindow(dom.GetWindow())
	win.InitApplication(new(App.NanolletApp))
	win.InitApplication(new(App.NanofyApp))
	win.InitApplication(new(App.AccountApp))
	win.InitApplication(new(App.NanoAliasApp))
	win.InitApplication(new(App.SettingsApp))
	win.ViewApplication(new(App.AccountApp))

	if Storage.PersistentStorage.SeedFY != *new(Wallet.SeedFY) {
		win.ViewPage(new(App.PagePassword))
	}

	go Background.UpdateNodeCount(win)
}
