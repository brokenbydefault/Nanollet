package DOM

import (
	"github.com/sciter-sdk/go-sciter"
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
