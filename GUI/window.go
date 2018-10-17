package GUI

import (
	"github.com/brokenbydefault/Nanollet/GUI/App"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"path/filepath"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
)

func init() {
	if err := Storage.ArbitraryStorage.WriteFile("sciter.link", Front.Sciter); err != nil {
		panic(err)
	}

	sciter.SetDLL(filepath.Join(Storage.ArbitraryStorage.Path, "sciter.link"))
}

func Start() {

	w, err := window.New(sciter.SW_MAIN|sciter.SW_RESIZEABLE|sciter.SW_TITLEBAR|sciter.SW_CONTROLS|sciter.SW_GLASSY|sciter.SW_OWNS_VM, sciter.NewRect(200, 200, 900, 600))
	if err != nil {
		panic(err)
	}
	w.SetTitle("Nanollet")

	if Storage.Configuration.DebugStatus {
		w.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)
	}

	w.LoadHtml(string(Front.HTML), "/")
	w.SetCSS(Front.CSSStyle, "style.css", "text/css")

	win := DOM.NewWindow(w)
	win.InitApplication(new(App.NanolletApp))
	win.InitApplication(new(App.NanofyApp))
	win.InitApplication(new(App.AccountApp))
	win.InitApplication(new(App.SettingsApp))
	win.ViewApplication(new(App.AccountApp))

	if Storage.PermanentStorage.SeedFY != *new(Wallet.SeedFY) {
		win.ViewPage(new(App.PagePassword))
	}

	go Background.UpdateNodeCount(win)

	w.Show()
	w.Run()

}
