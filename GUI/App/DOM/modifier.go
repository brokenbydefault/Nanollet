// +build !js

package DOM

import (
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

func ApplyForAll(w *window.Window, css string, mod Modifier) error {
	els, err := SelectAllElement(w, css)
	if err != nil {
		return err
	}

	return applyForAll(els, mod)
}

func (p *Page) ApplyForAll(w *window.Window, css string, mod Modifier) error {
	els, err := p.SelectAllElement(w, css)
	if err != nil {
		return err
	}

	return applyForAll(els, mod)
}

func ApplyForIt(w *window.Window, css string, mod Modifier) error {
	it, err := SelectFirstElement(w, css)
	if err != nil {
		return err
	}

	return mod(it)
}

func (p *Page) ApplyForIt(w *window.Window, css string, mod Modifier) error {
	it, err := p.SelectFirstElement(w, css)
	if err != nil {
		return err
	}

	return mod(it)
}

func applyForAll(els []*sciter.Element, mod Modifier) error {
	for _, e := range els {
		if err := mod(e); err != nil {
			return err
		}
	}

	return nil
}

type Modifier func(el *sciter.Element) error

func DisableElement(el *sciter.Element) error {
	return el.SetState(sciter.STATE_DISABLED, 0, true)
}

func EnableElement(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_DISABLED, true)
}

func HideElement(el *sciter.Element) error {
	return el.SetStyle("display", "none")
}

func ShowElement(el *sciter.Element) error {
	return el.SetStyle("display", "block")
}

func VisibleElement(el *sciter.Element) error {
	return el.SetStyle("visibility", "visible")
}

func InvisibleElement(el *sciter.Element) error {
	return el.SetStyle("visibility", "hidden")
}

func VisitedElement(el *sciter.Element) error {
	return el.SetState(sciter.STATE_VISITED, 0, true)
}

func UnvisitedElement(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_VISITED, true)
}

func ReadOnlyElement(el *sciter.Element) error {
	return el.SetState(sciter.STATE_READONLY, 0, true)
}

func WriteOnlyElement(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_READONLY, true)
}

func Checked(el *sciter.Element) error {
	return el.SetState(sciter.STATE_CHECKED, 0, true)
}

func Unchecked(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_CHECKED, true)
}

func ClearValue(el *sciter.Element) error {
	return el.SetValue(sciter.NewValue())
}

func ClearHTML(el *sciter.Element) error {
	return el.SetHtml(" ", sciter.SIH_REPLACE_CONTENT)
}

func DestroyHTML(el *sciter.Element) error {
	el.SetHtml(" ", sciter.SOH_REPLACE)
	return el.Clear()
}
