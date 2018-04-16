package DOM

import (
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"html"
)

func UpdateAmount(w *window.Window) error {
	hamm := Numbers.NewHumanFromRaw(Storage.Amount)

	for el, scale := range map[string]int{
		".ammount": 6,
	} {
		balance, err := hamm.ConvertToBase(Numbers.MegaXRB, scale)
		if err != nil {
			return err
		}

		display, err := SelectFirstElement(w, el)
		if err != nil {
			return err
		}

		err = display.SetValue(sciter.NewValue(balance))
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateNotification(w *window.Window, msg string) {
	box, _ := SelectFirstElement(w, "section.notification")

	nt := CreateElementAppendTo("button", html.EscapeString(msg), "notification", "", box)
	nt.OnClick(func() {
		nt.SetHtml(" ", sciter.SOH_REPLACE)
		nt.Clear()
	})
}
