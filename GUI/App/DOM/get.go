package DOM

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"
	"strings"
)



func (p *Page) GetStringValue(w *window.Window, css string) (result string, err error) {
	input, err := p.SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getValue(input)
}

func GetStringValue(w *window.Window, css string) (result string, err error) {
	input, err := SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getValue(input)
}

func getValue(el *sciter.Element) (result string, err error) {
	value, err := el.GetValue()
	if err != nil {
		return
	}

	return strings.TrimSpace(value.String()), nil
}