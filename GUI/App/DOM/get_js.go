// +build js

package DOM

import (
	"honnef.co/go/js/dom"
	"strings"
)

func (p *Page) GetStringValue(w dom.Document, css string) (result string, err error) {
	input, err := p.SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getValue(input)
}

func GetStringValue(w dom.Document, css string) (result string, err error) {
	input, err := SelectFirstElement(w, css)
	if err != nil {
		return
	}

	return getValue(input)
}

func getValue(el dom.Element) (result string, err error) {
	if strings.ToUpper(el.TagName()) == "TEXTAREA" {
		return el.(*dom.HTMLTextAreaElement).Value, nil
	}

	return el.(*dom.HTMLInputElement).Value, nil
}
