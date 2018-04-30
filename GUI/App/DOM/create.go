// +build !js

package DOM

import (
	"github.com/sciter-sdk/go-sciter"
	"github.com/skip2/go-qrcode"
	"image/color"
	"encoding/base64"
)

func CreateElement(tag, value, class string, id string) *sciter.Element {
	el, _ := sciter.CreateElement(tag, value)
	el.SetAttr("class", class)
	el.SetAttr("id", id)
	return el
}

func CreateElementAppendTo(tag, value, class string, id string, target *sciter.Element) *sciter.Element {
	el := CreateElement(tag, value, class, id)
	target.Append(el)
	return el
}

func CreateQRCodeAppendTo(text string, color color.RGBA, size int, target *sciter.Element) *sciter.Element {
	qr, _ := qrcode.New(text, qrcode.Highest)
	qr.BackgroundColor = color

	png, _ := qr.PNG(size)

	el := CreateElementAppendTo("img", "", "", "", target)
	el.SetAttr("src", "data:image/png;base64, "+base64.StdEncoding.EncodeToString(png))

	return el
}
