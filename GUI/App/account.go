// +build !js

package App

import (
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"strconv"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/MFA"
	"image/color"
	"bytes"
	"encoding/json"
	"github.com/brokenbydefault/MFA/mfatypes"
	"github.com/brokenbydefault/Nanollet/Util"
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
		&PageMFA{},
	}
}

type PageIndex guitypes.Sector

func (c *PageIndex) Name() string {
	return "index"
}

func (c *PageIndex) OnView(w *window.Window) {
	// no-op
}

func (c *PageIndex) OnContinue(w *window.Window, action string) {
	switch action {
	case "genSeed":
		ViewPage(w, &PageGenerate{})
	case "importSeed":
		ViewPage(w, &PageImport{})
	}
}

type PageGenerate guitypes.Sector

func (c *PageGenerate) Name() string {
	return "generate"
}

func (c *PageGenerate) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	seed, err := Wallet.NewSeedFY(Wallet.V0, Wallet.Nanollet)
	if err != nil {
		return
	}

	textarea, _ := page.SelectFirstElement(w, ".seed")
	textarea.SetValue(sciter.NewValue(seed.Encode()))
	DOM.ReadOnlyElement(textarea)
}

func (c *PageGenerate) OnContinue(w *window.Window, _ string) {
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

func (c *PageImport) OnContinue(w *window.Window, _ string) {
	page := DOM.SetSector(c)

	seed, err := page.GetStringValue(w, ".seed")
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

func (c *PagePassword) OnContinue(w *window.Window, _ string) {
	page := DOM.SetSector(c)

	password, err := page.GetStringValue(w, ".password")
	if err != nil || len(password) < 8 {
		DOM.UpdateNotification(w, "There was a problem with your password, this is too short")
		return
	}
	DOM.ApplyForIt(w, ".password", DOM.ClearValue)

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

	need2FA, err := page.GetStringValue(w, ".ask2fa")
	if err == nil && need2FA != "" {
		Storage.PASSWORD = []byte(password)
		ViewPage(w, &PageMFA{})
		return
	}

	Storage.SEED = seedfy.RecoverSeed([]byte(password), nil)
	ViewPage(w, &PageAddress{})
}

type PageMFA guitypes.Sector

func (c *PageMFA) Name() string {
	return "mfa"
}

func (c *PageMFA) OnView(w *window.Window) {
	page := DOM.SetSector(c)

	keypinning, _ := Storage.Permanent.ReadFile("mfa.dat")

	MFAConn := Connectivity.NewSocket()
	envs := make(chan []byte, 16)

	receiver := MFA.NewReceiver()
	receiverPK := receiver.Ephemeral.PublicKey()

	qrcodespace, _ := page.SelectFirstElement(w, ".qrcode")
	DOM.ClearHTML(qrcodespace)
	DOM.CreateQRCodeAppendTo(Util.SecureHexEncode(receiverPK[:]), color.RGBA{220, 220, 223, 1}, 300, qrcodespace)

	// @TODO Improve code
	err := RPCClient.SubscribeMFA(MFAConn, receiverPK)
	if err != nil {
		DOM.UpdateNotification(w, "There was a critical problem connecting to our servers, please try again")
		return
	}

	// @TODO Remove goroutines here
	go func() {
		if err := MFAConn.ReceiveAllMessages(nil, envs); err != nil {
			return
		}
	}()

	go func() {
		defer MFAConn.CloseWebsocket()

		for e := range envs {
			jsn := mfatypes.CallbackResponse{}

			if err := json.Unmarshal(e, &jsn); err != nil {
				continue
			}

			token, sender, err := receiver.OpenEnvelope(jsn.Envelope)
			if err != nil || keypinning != nil && !bytes.Equal(sender, keypinning) {
				continue
			}

			Storage.TOKENMFA = token
			Storage.Permanent.WriteFile("mfa.dat", sender)
			receiver.Ephemeral = nil
			c.OnContinue(w, "")

			return
		}
	}()

	return
}

func (c *PageMFA) OnContinue(w *window.Window, _ string) {
	if Storage.TOKENMFA == nil {
		return
	}

	seed, err := Storage.Permanent.ReadFile("wallet.dat")
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem reading your seed or it doesn't exist anymore")
		return
	}

	seedfy, err := Wallet.ReadSeedFY(string(seed))
	if err != nil {
		DOM.UpdateNotification(w, "There was a problem reading your SEEDFY, this is incorrect or isn't supported")
		return
	}

	Storage.SEED = seedfy.RecoverSeed(Storage.PASSWORD, Storage.TOKENMFA)
	Storage.TOKENMFA = nil
	Storage.PASSWORD = ""

	ViewPage(w, &PageAddress{})
}

type PageAddress guitypes.Sector

const ADDRESS_PER_PAGE uint32 = 5

func (c *PageAddress) Name() string {
	return "address"
}

func (c *PageAddress) Position(w *window.Window) uint32 {
	page := DOM.SetSector(c)

	index, err := page.GetStringValue(w, ".address option")
	if index == "" || err != nil {
		return 0
	}

	i, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		return 0
	}

	return uint32(i)
}

func (c *PageAddress) UpdateList(w *window.Window, min, max uint32) {
	page := DOM.SetSector(c)

	selectbox, err := page.SelectFirstElement(w, ".address")
	if err != nil {
		panic(err)
	}

	value, _ := selectbox.GetValue()

	DOM.ClearHTML(selectbox)
	for i := min; i < max; i++ {
		pk, _, err := Storage.SEED.CreateKeyPair(Wallet.Nano, i)
		if err != nil {
			panic(err)
		}

		addr := string(pk.CreateAddress())

		opt := DOM.CreateElementAppendTo("option", addr[0:16]+" ... "+addr[48:64], "item", "", selectbox)
		opt.SetAttr("value", strconv.FormatUint(uint64(i), 10))

		if value.String() != "" && uint32(value.Int64()) == i {
			DOM.Checked(opt)
		}
	}
}

func (c *PageAddress) Next(w *window.Window) {
	pos := c.Position(w)
	if pos == 1<<32-1 {
		return
	}

	c.UpdateList(w, pos+ADDRESS_PER_PAGE, pos+(ADDRESS_PER_PAGE*2))
}

func (c *PageAddress) Previous(w *window.Window) {
	pos := c.Position(w)
	if pos == 0 {
		return
	}

	c.UpdateList(w, pos-ADDRESS_PER_PAGE, pos)
}

func (c *PageAddress) OnView(w *window.Window) {
	c.UpdateList(w, 0, 5)
}

func (c *PageAddress) OnContinue(w *window.Window, action string) {

	switch action {
	case "next":
		c.Next(w)
	case "previous":
		c.Previous(w)
	case "continue":
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

}
