package App

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"strconv"
)

type AccountApp guitypes.App

func (c *AccountApp) Name() string {
	return "account"
}

func (c *AccountApp) HaveSidebar() bool {
	return false
}

func (c *AccountApp) Display() Front.HTMLPAGE {
	return Front.HTMLAccount
}

func (c *AccountApp) Pages() []guitypes.Page {
	return []guitypes.Page{
		&PageIndex{},
		&PageGenerate{},
		&PageImport{},
		&PagePassword{},
		&PageAddress{},
	}
}

type PageIndex guitypes.Sector

func (c *PageIndex) Name() string {
	return "index"
}

func (c *PageIndex) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	gen, _ := page.SelectFirstElement(w, ".genSeed")
	gen.OnClick(func() {
		ViewPage(w, &PageGenerate{})
	})

	imp, _ := page.SelectFirstElement(w, ".importSeed")
	imp.OnClick(func() {
		ViewPage(w, &PageImport{})
	})
}

func (c *PageIndex) OnContinue(w *window.Window) {
	// no-op
}

type PageGenerate guitypes.Sector

func (c *PageGenerate) Name() string {
	return "generate"
}

func (c *PageGenerate) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	seed, err := Wallet.NewSeedFY(Wallet.V1, Wallet.Nanollet)
	if err != nil {
		return
	}

	textarea, _ := page.SelectFirstElement(w, ".seed")
	textarea.SetValue(sciter.NewValue(seed.Encode()))
	DOM.ReadOnlyElement(textarea)
}

func (c *PageGenerate) OnContinue(w *window.Window) {
	page := DOM.SetSector(c)
	seed, err := page.GetStringValue(w, ".seed")
	if seed == "" || err != nil {
		panic(err)
	}

	err = Storage.Permanent.WriteFile("wallet.dat", []byte(seed))
	if err != nil {
		panic(err)
	}

	ViewPage(w, &PagePassword{})
}

type PageImport guitypes.Sector

func (c *PageImport) Name() string {
	return "import"
}

func (c *PageImport) OnView(w *window.Window) {
	// no-op
}

func (c *PageImport) OnContinue(w *window.Window) {
	page := DOM.SetSector(c)

	seed, err := page.GetStringValue(w, ".seed")
	if seed == "" || err != nil {
		return
	}

	_, err = Wallet.ReadSeedFY(seed)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem interpreting your SEEDFY, it's wrong or isn't supported anymore")
		return
	}

	err = Storage.Permanent.WriteFile("wallet.dat", []byte(seed))
	if err != nil {
		panic(err)
	}

	ViewPage(w, &PagePassword{})
	page.ApplyForIt(w, ".seed", DOM.ClearValue)
}

type PagePassword guitypes.Sector

func (c *PagePassword) Name() string {
	return "password"
}

func (c *PagePassword) OnView(w *window.Window) {
	// no-op
}

func (c *PagePassword) OnContinue(w *window.Window) {
	page := DOM.SetSector(c)

	password, err := page.GetStringValue(w, ".password")
	if err != nil || len(password) < 8 {
		DOM.UpdateNotification(w, "There was a problem with your password, this is too short")
		return
	}

	seed, err := Storage.Permanent.ReadFile("wallet.dat")
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem reading your seed or it doesn't exist anymore")
		return
	}

	seedfy, err := Wallet.ReadSeedFY(string(seed))
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem interpreting your SEEDFY, this is incorrect or isn't supported")
		return
	}

	Storage.SEED = seedfy.RecoverSeed(password, nil)
	ViewPage(w, &PageAddress{})
	DOM.ApplyForIt(w, ".password", DOM.ClearValue)
}

type PageAddress guitypes.Sector

func (c *PageAddress) Name() string {
	return "address"
}

func (c *PageAddress) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	selectbox, err := page.SelectFirstElement(w, ".address")
	if err != nil {
		panic(err)
	}

	DOM.ClearHTML(selectbox)
	for i := uint32(0); i < 16; i++ {
		pk, _, err := Storage.SEED.CreateKeyPair(Wallet.Nano, i)
		if err != nil {
			panic(err)
		}

		opt, _ := sciter.CreateElement("option", string(pk.CreateAddress()))
		opt.SetAttr("value", strconv.FormatUint(uint64(i), 10))

		selectbox.Append(opt)
	}
}

func (c *PageAddress) OnContinue(w *window.Window) {
	page := DOM.SetSector(c)

	index, err := page.GetStringValue(w, ".address")
	if index == "" || err != nil {
		return
	}

	i, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		return
	}

	pk, sk, err := Storage.SEED.CreateKeyPair(Wallet.Nano, uint32(i))
	if err != nil {
		return
	}

	Storage.SK = sk
	Storage.PK = pk

	err = Background.StartAddress(w)
	if err != nil {
		DOM.UpdateNotification(w, "There was a critical problem connecting to our servers, please try again")
		return
	}

	Storage.SEED = nil
	page.ApplyForIt(w, ".address", DOM.ClearHTML)

	Background.StartTransaction()
	ViewApplication(w, &NanolletApp{})
}
