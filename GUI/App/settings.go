package App

import (
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Storage"
)

type SettingsApp struct{}

func (c *SettingsApp) Name() string {
	return "settings"
}

func (c *SettingsApp) HaveSidebar() bool {
	return true
}

func (c *SettingsApp) Pages() []DOM.Page {
	return []DOM.Page{
		&PageSeed{},
		//&PageAuthorities{},
	}
}

type PageSeed struct{}

func (c *PageSeed) Name() string {
	return "seed"
}

func (c *PageSeed) OnView(w *DOM.Window, dom *DOM.DOM) {
	seedbox, err := dom.SelectFirstElement(".seed")
	if err != nil {
		return
	}

	seedbox.SetValue(Storage.PermanentStorage.SeedFY.String())
}

func (c *PageSeed) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	//no-op
}
