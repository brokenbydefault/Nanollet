// +build js

package DOM

import (
	"honnef.co/go/js/dom"
)

func CreateElement(tag, value, class string, id string) dom.Element {
	root := dom.GetWindow().Document().(dom.HTMLDocument)

	el := root.CreateElement(tag)
	el.SetAttribute("class", class)
	el.SetAttribute("id", id)

	return el
}

func CreateElementAppendTo(tag, value, class string, id string, target dom.Element) dom.Element {
	el := CreateElement(tag, value, class, id)
	target.AppendChild(el)
	return el
}
