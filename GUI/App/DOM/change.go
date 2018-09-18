// +build !js

package DOM

import (
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"fmt"
)

func UpdateAmount(w *Window) error {
	humanAmm := Numbers.NewHumanFromRaw(Storage.AccountStorage.Balance)

	for el, scale := range map[string]int{
		".ammount": 6,
	} {
		balance, err := humanAmm.ConvertToBase(Numbers.MegaXRB, scale)
		if err != nil {
			return err
		}

		display, err := w.root.SelectFirstElement(el)
		if err != nil {
			return err
		}

		err = display.SetValue(balance)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateNodesCount(w *Window) error {
	for _, el := range []string{".nodes"} {
		display, err := w.root.SelectFirstElement(el)
		if err != nil {
			return err
		}

		active, all := Storage.PeerStorage.CountActive()

		if err = display.SetValue(fmt.Sprintf("%d (%d)", active, all)); err != nil {
			return err
		}
	}

	return nil
}

func UpdateNotification(w *Window, msg string) {
	box, _ := w.root.SelectFirstElement("section.notification")

	nt := box.CreateElementWithAttr("button", msg, Attrs{"class": "notification"})
	nt.On(Click, func(_ string) {
		//	nt.SetHtml(" ", sciter.SOH_REPLACE)
		//	nt.Clear()
	})
}
