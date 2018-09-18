// +build !js

package DOM

import (
	"github.com/sciter-sdk/go-sciter"
	"encoding/base64"
)

func (el *Element) CreateElement(tag, text string) *Element {
	e, _ := sciter.CreateElement(tag, text)

	el.el.Append(e)
	return NewElement(e)
}

type Attrs map[string]string

func (el *Element) CreateElementWithAttr(tag, text string, attrs Attrs) *Element {
	e := el.CreateElement(tag, text)
	for name, value := range attrs {
		e.SetAttr(name, value)
	}

	return e
}

func (el *Element) CreateQRCode(png []byte) *Element {
	return el.CreateElementWithAttr("img", "", Attrs{"src": "data:image/png;base64, " + base64.StdEncoding.EncodeToString(png)})
}