package GUI

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/GUI/App"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/Config"
	"path/filepath"
)

func Start() {

	w, err := window.New(sciter.SW_MAIN|sciter.SW_RESIZEABLE|sciter.SW_OWNS_VM, sciter.NewRect(200, 200, 900, 600))
	if err != nil {
		panic(err)
	}

	if Config.IsDebugEnabled() {
		w.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)
	}

	w.LoadHtml(string(Front.HTMLBase), "/")
	w.SetCSS(Front.CSSStyle, "style.css", "text/css")

	App.InitApplication(w, &App.NanofyApp{})
	App.InitApplication(w, &App.NanolletApp{})
	App.InitApplication(w, &App.AccountApp{})

	App.ViewApplication(w, &App.AccountApp{})

	if Storage.Permanent.Exists("wallet.dat") {
		App.ViewPage(w, &App.PagePassword{})
	}

	w.Show()
	w.Run()

}

func Unpack() {
	if err := Storage.Permanent.WriteFile("sciter.link", Front.Sciter); err != nil {
		panic(err)
	}
	sciter.SetDLL(filepath.Join(Storage.Permanent.Path, "sciter.link"))
}
