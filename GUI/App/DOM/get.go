// +build !js

package DOM

import (
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

func (p *Page) GetStringValue(w *window.Window, css string) (result string, err error) {
	input, err := p.SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getString(input)
}

func GetStringValue(w *window.Window, css string) (result string, err error) {
	input, err := SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getString(input)
}

func (p *Page) GetBytesValue(w *window.Window, css string) (result []byte, err error) {
	input, err := p.SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getBytes(input)
}

func GetBytesValue(w *window.Window, css string) (result []byte, err error) {
	input, err := SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getBytes(input)
}

func getBytes(el *sciter.Element) (result []byte, err error) {
	value, err := el.GetValue()
	if err != nil {
		return nil, err
	}

	//return value.Bytes(), nil
	return []byte(value.String()), nil
}

func getString(el *sciter.Element) (result string, err error) {
	value, err := el.GetValue()
	if err != nil {
		return "", err
	}

	return value.String(), nil
}
