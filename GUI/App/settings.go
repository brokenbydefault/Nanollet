package App

import (
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/sciter-sdk/go-sciter"
	"github.com/brokenbydefault/Nanollet/Storage"
)

type SettingsApp guitypes.App

func (c *SettingsApp) Name() string {
	return "settings"
}

func (c *SettingsApp) HaveSidebar() bool {
	return true
}

func (c *SettingsApp) Display() Front.HTMLPAGE {
	return Front.HTMLSettings
}

func (c *SettingsApp) Pages() []guitypes.Page {
	return []guitypes.Page{
		&PageSeed{},
		//&PageAuthorities{},
	}
}

type PageSeed guitypes.Sector

func (c *PageSeed) Name() string {
	return "seed"
}

func (c *PageSeed) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	seedbox, err := page.SelectFirstElement(w, ".seed")
	if err != nil {
		return
	}

	seedbox.SetValue(sciter.NewValue(Storage.PermanentStorage.SeedFY.String()))
}

func (c *PageSeed) OnContinue(w *window.Window, action string) {
	//no-op
}