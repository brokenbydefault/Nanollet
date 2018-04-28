// +build !js

package DOM

import (
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

type Page string

func SetPage(name string) Page {
	return Page(name)
}

func SetSector(pg guitypes.Page) Page {
	return Page(pg.Name())
}

func (p Page) SelectAllElement(w *window.Window, css string) ([]*sciter.Element, error) {
	el, err := w.GetRootElement()
	if err != nil {
		return nil, err
	}

	el, err = el.SelectFirst("[page=\"" + string(p) + "\"]")
	if err != nil {
		return nil, err
	}

	return el.Select(css)
}

func (p Page) SelectFirstElement(w *window.Window, css string) (*sciter.Element, error) {
	el, err := w.GetRootElement()
	if err != nil {
		return nil, err
	}

	el, err = el.SelectFirst("[page=\"" + string(p) + "\"]")
	if err != nil {
		return nil, err
	}

	return el.SelectFirst(css)
}

func SelectAllElement(w *window.Window, css string) ([]*sciter.Element, error) {
	el, err := w.GetRootElement()
	if err != nil {
		return nil, err
	}

	return el.Select(css)
}

func SelectFirstElement(w *window.Window, css string) (*sciter.Element, error) {
	el, err := w.GetRootElement()
	if err != nil {
		return nil, err
	}

	return el.SelectFirst(css)
}
