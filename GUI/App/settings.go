package App

import (
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/Inkeliz/go-sciter/window"
)

type SettingsApp guitypes.App

func (c *SettingsApp) Name() string {
	return "nanofy"
}

func (c *SettingsApp) HaveSidebar() bool {
	return true
}

func (c *SettingsApp) Display() Front.HTMLPAGE {
	return Front.HTMLNanofy
}

func (c *SettingsApp) Pages() []guitypes.Page {
	return []guitypes.Page{
		&PageSign{},
		&PageVerify{},
	}
}

type PageSeed guitypes.Sector

func (c *PageSeed) Name() string {
	return "sign"
}

func (c *PageSeed) OnView(w *window.Window) {
	// no-op
}

func (c *PageSeed) OnContinue(w *window.Window, _ string) {
	// no-op
}
