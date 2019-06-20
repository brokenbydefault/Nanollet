package App

import (
	"github.com/Inkeliz/blakEd25519"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"strings"
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
		&PageAuthorities{},
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

	skbox, err := dom.SelectFirstElement(".sk")
	if err != nil {
		return
	}

	seedbox.SetValue(Storage.PersistentStorage.SeedFY.String())
	skbox.SetValue(Util.SecureHexEncode(Storage.AccountStorage.SecretKey[:blakEd25519.PublicKeySize]))
}

func (c *PageSeed) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	//no-op
}

type PageAuthorities struct{}

func (c *PageAuthorities) Name() string {
	return "authorities"
}

func (c *PageAuthorities) OnView(w *DOM.Window, dom *DOM.DOM) {
	autbox, err := dom.SelectFirstElement(".aut")
	if err != nil {
		return
	}


	var auts string
	for _, pk := range Storage.Configuration.Account.Quorum.PublicKeys {
		auts += string(pk.CreateAddress())
		auts += "\n"
	}

	autbox.SetValue(auts)
}

func (c *PageAuthorities) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	autbox, err := dom.SelectFirstElement(".aut")
	if err != nil {
		return
	}

	auts, err := autbox.GetStringValue()
	if err != nil {
		return
	}

	var pks []Wallet.PublicKey
	for _, addr := range strings.Split(strings.Replace(auts, "\r\n", "\n", -1), "\n") {
		address := Wallet.Address(strings.TrimSpace(addr))
		if address.IsValid() {
			pks = append(pks, address.MustGetPublicKey())
		}
	}

	Storage.Configuration.Account.Quorum.PublicKeys = pks
	Storage.PersistentStorage.Quorum = Storage.Configuration.Account.Quorum
	Storage.Engine.Save(&Storage.PersistentStorage)
}
