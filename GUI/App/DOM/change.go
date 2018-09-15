// +build !js

package DOM

import (
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"fmt"
)

func UpdateAmount(w *window.Window) error {
	humanAmm := Numbers.NewHumanFromRaw(Storage.AccountStorage.Balance)

	for el, scale := range map[string]int{
		".ammount": 6,
	} {
		balance, err := humanAmm.ConvertToBase(Numbers.MegaXRB, scale)
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

func UpdateNodesCount(w *window.Window) error {
	for _, el := range []string{".nodes"} {
		display, err := SelectFirstElement(w, el)
		if err != nil {
			return err
		}

		active, all := Storage.PeerStorage.CountActive()

		if err = display.SetValue(sciter.NewValue(fmt.Sprintf("%d (%d)", active, all))); err != nil {
			return err
		}
	}

	return nil
}


func UpdateNotification(w *window.Window, msg string) {
	box, _ := SelectFirstElement(w, "section.notification")

	nt := CreateElementAppendTo("button", msg, "notification", "", box)
	nt.OnClick(func() {
		nt.SetHtml(" ", sciter.SOH_REPLACE)
		nt.Clear()
	})
}
