package App

import (
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"strconv"
	"image/color"
	"github.com/brokenbydefault/Nanollet/TwoFactor"
	"github.com/brokenbydefault/Nanollet/TwoFactor/Ephemeral"
	"github.com/brokenbydefault/Nanollet/Util"
)

type AccountApp struct{}

func (c *AccountApp) Name() string {
	return "account"
}

func (c *AccountApp) HaveSidebar() bool {
	return false
}

func (c *AccountApp) Pages() []DOM.Page {
	return []DOM.Page{
		&PageToS{},
		&PageIndex{},
		&PageGenerate{},
		&PageImport{},
		&PagePassword{},
		&PageAddress{},
		&PageMFA{},
	}
}

type PageToS struct{}

func (c *PageToS) Name() string {
	return "tos"
}

func (c *PageToS) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PageToS) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	if action == "accept" {
		w.ViewPage(new(PageIndex))
	}
}

type PageIndex struct{}

func (c *PageIndex) Name() string {
	return "index"
}

func (c *PageIndex) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PageIndex) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {
	switch action {
	case "genSeed":
		w.ViewPage(new(PageGenerate))
	case "importSeed":
		w.ViewPage(new(PageImport))
	}
}

type PageGenerate struct{}

func (c *PageGenerate) Name() string {
	return "generate"
}

func (c *PageGenerate) OnView(w *DOM.Window, dom *DOM.DOM) {
	seedfy, err := Wallet.NewSeedFY(Wallet.V0, Wallet.Nanollet)
	if err != nil {
		return
	}

	textarea, _ := dom.SelectFirstElement(".seed")
	textarea.SetValue(seedfy.String())
	textarea.Apply(DOM.ReadOnlyElement)
}

func (c *PageGenerate) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	seedHex, err := dom.GetStringValueOf(".seed")
	if seedHex == "" || err != nil {
		panic(err)
	}

	sf, err := Wallet.ReadSeedFY(seedHex)
	if err != nil {
		panic(err)
	}

	Storage.PersistentStorage.AddSeedFY(sf)

	w.ViewPage(new(PagePassword))
}

type PageImport struct{}

func (c *PageImport) Name() string {
	return "import"
}

func (c *PageImport) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PageImport) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	seed, err := dom.GetStringValueOf(".seed")
	if seed == "" || err != nil {
		return
	}

	sf, err := Wallet.ReadSeedFY(seed)
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem interpreting your SEEDFY, it's wrong or isn't supported anymore")
		return
	}

	if ok := sf.IsValid(Wallet.Version(sf.Version), Wallet.Nanollet); !ok {
		DOM.UpdateNotification(w, "There was a problem interpreting your SEEDFY, it's wrong or isn't supported anymore")
		return
	}

	Storage.PersistentStorage.AddSeedFY(sf)

	w.ViewPage(new(PagePassword))
	dom.ApplyFor(".seed", DOM.ClearValue)
}

type PagePassword struct{}

func (c *PagePassword) Name() string {
	return "password"
}

func (c *PagePassword) OnView(w *DOM.Window, dom *DOM.DOM) {
	// no-op
}

func (c *PagePassword) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	password, err := dom.GetBytesValueOf(".password")
	if err != nil || len(password) < 8 {
		DOM.UpdateNotification(w, "There was a problem with your password, this is too short")
		return
	}

	seedfy := Storage.PersistentStorage.SeedFY

	need2FA, err := dom.GetStringValueOf(".ask2fa")
	if err == nil && need2FA != "" {
		Storage.AccessStorage.Password = password
		w.ViewPage(new(PageMFA))
		return
	}

	Storage.AccessStorage.Seed = seedfy.RecoverSeed(password, nil)
	w.ViewPage(new(PageAddress))

	dom.ApplyFor(".password", DOM.ClearValue)
}

type PageMFA struct{}

func (c *PageMFA) Name() string {
	return "mfa"
}

func (c *PageMFA) OnView(w *DOM.Window, dom *DOM.DOM) {
	sk := Ephemeral.NewEphemeral()
	requester, response, err := TwoFactor.NewRequesterServer(&sk, Storage.PersistentStorage.AllowedKeys)
	if err != nil {
		panic(err)
	}

	qr, err := requester.QRCode(300, color.RGBA{220, 220, 223, 1})
	if err != nil {
		panic(err)
	}

	qrSpace, _ := dom.SelectFirstElement(".qrcode")
	qrSpace.Apply(DOM.ClearHTML)
	qrSpace.CreateQRCode(qr)

	go func() {
		for resp := range response {
			//@TODO Notify the user to allow or not the key
			Storage.PersistentStorage.AddAllowedKey(resp.Capsule.Device)
			Storage.AccessStorage.Token = resp.Capsule.Token

			c.OnContinue(w, dom, "")
			break
		}
	}()

	return
}

func (c *PageMFA) OnContinue(w *DOM.Window, dom *DOM.DOM, _ string) {
	if Util.IsEmpty(Storage.AccessStorage.Token[:]) {
		return
	}

	seedfy := Storage.PersistentStorage.SeedFY

	Storage.AccessStorage.Seed = seedfy.RecoverSeed(Storage.AccessStorage.Password, Storage.AccessStorage.Token[:])
	copy(Storage.AccessStorage.Token[:], make([]byte, len(Storage.AccessStorage.Token)))
	copy(Storage.AccessStorage.Password[:], make([]byte, len(Storage.AccessStorage.Password)))

	w.ViewPage(new(PageAddress))
}

type PageAddress struct{}

const AddressesPerPage uint32 = 5

func (c *PageAddress) Name() string {
	return "address"
}

func (c *PageAddress) Position(dom *DOM.DOM) uint32 {
	index, err := dom.GetStringValueOf(".address option")
	if index == "" || err != nil {
		return 0
	}

	i, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		return 0
	}

	return uint32(i)
}

func (c *PageAddress) UpdateList(dom *DOM.DOM, min, max uint32) {
	selectbox, err := dom.SelectFirstElement(".address")
	if err != nil {
		panic(err)
	}

	value, _ := selectbox.GetStringValue()
	selectbox.Apply(DOM.ClearHTML)

	for i := min; i < max; i++ {
		pk, _, err := Storage.AccessStorage.Seed.CreateKeyPair(Wallet.Nano, i)
		if err != nil {
			panic(err)
		}

		addr := string(pk.CreateAddress())

		opt := selectbox.CreateElementWithAttr("option", addr[0:16]+" ... "+addr[48:64], DOM.Attrs{
			"class": "item",
			"value": strconv.FormatUint(uint64(i), 10),
		})

		if intVal, err := strconv.ParseUint(value, 10, 32); err != nil && value != "" && uint32(intVal) == i {
			opt.Apply(DOM.Checked)
		}
	}
}

func (c *PageAddress) Next(dom *DOM.DOM) {
	pos := c.Position(dom)
	if pos == 1<<32-1 {
		return
	}

	c.UpdateList(dom, pos+AddressesPerPage, pos+(AddressesPerPage*2))
}

func (c *PageAddress) Previous(dom *DOM.DOM) {
	pos := c.Position(dom)
	if pos == 0 {
		return
	}

	c.UpdateList(dom, pos-AddressesPerPage, pos)
}

func (c *PageAddress) OnView(w *DOM.Window, dom *DOM.DOM) {
	c.UpdateList(dom, 0, 5)
}

func (c *PageAddress) OnContinue(w *DOM.Window, dom *DOM.DOM, action string) {

	switch action {
	case "next":
		c.Next(dom)
	case "previous":
		c.Previous(dom)
	case "continue":
		index, err := dom.GetStringValueOf(".address")
		if index == "" || err != nil {
			return
		}

		i, err := strconv.ParseUint(index, 10, 32)
		if err != nil {
			return
		}

		pk, sk, err := Storage.AccessStorage.Seed.CreateKeyPair(Wallet.Nano, uint32(i))
		if err != nil {
			return
		}

		Storage.AccountStorage.SecretKey = sk
		Storage.AccountStorage.PublicKey = pk

		err = Background.StartAddress(w)
		if err != nil {
			DOM.UpdateNotification(w, "There was a critical problem connecting to our servers, please try again")
			return
		}

		Storage.AccessStorage.Seed = nil
		dom.ApplyFor(".address", DOM.ClearHTML)

		w.ViewApplication(new(NanolletApp))
	}

}
