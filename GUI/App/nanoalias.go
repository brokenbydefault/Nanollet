package App

import (
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/NanoAlias"
	"github.com/brokenbydefault/Nanollet/Storage"
	"strings"
)

type NanoAliasApp struct{}

func (c *NanoAliasApp) Name() string {
	return "nanoalias"
}

func (c *NanoAliasApp) HaveSidebar() bool {
	return true
}

func (c *NanoAliasApp) Pages() []DOM.Page {
	return []DOM.Page{
		&PageRegister{},
	}
}

type PageRegister struct{}

func (c *PageRegister) Name() string {
	return "register"
}

func (c *PageRegister) OnView(w *DOM.Window, dom *DOM.DOM) {
	return
}

func (c *PageRegister) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	alias, err := dom.GetStringValueOf(".alias")
	if err != nil || alias == "" {
		return
	}

	alias = strings.TrimSpace(alias)
	alias = strings.ToLower(alias)

	if rune(alias[0]) != '@' {
		alias = "@" + alias
	}

	if ok := NanoAlias.Address(alias).IsValid(); !ok {
		DOM.UpdateNotification(w, "You can't use this alias. Invalid characters on the alias")
		return
	}

	previous, ok := Storage.TransactionStorage.GetByHash(&Storage.AccountStorage.Frontier)
	if !ok {
		DOM.UpdateNotification(w, "Can't find your last block, you need a opened account")
		return
	}

	if ok := NanoAlias.IsAvailable(alias); !ok {
		DOM.UpdateNotification(w, "You can't use this alias. It's already registered")
		return
	}

	defer DOM.UpdateAmount(w)
	if err := NanoAlias.Register(&Storage.AccountStorage.SecretKey, previous, alias); err != nil {
		DOM.UpdateNotification(w, "Impossible to register due to some error") //@TODO improve error
		return
	}

	DOM.UpdateNotification(w, "Complete! Try now, send using "+alias)
	return
}

type PageLookup struct{}

func (c *PageLookup) Name() string {
	return "lookup"
}

func (c *PageLookup) OnView(w *DOM.Window, dom *DOM.DOM) {
	return
}

func (c *PageLookup) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	panic("implement me")
}
