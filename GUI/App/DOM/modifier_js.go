// +build js

package DOM

import (
	"honnef.co/go/js/dom"
)

func ApplyForAll(w dom.Document, css string, mod Modifier) error {
	els, err := SelectAllElement(w, css)
	if err != nil {
		return err
	}

	return applyForAll(els, mod)
}

func (p *Page) ApplyForAll(w dom.Document, css string, mod Modifier) error {
	els, err := p.SelectAllElement(w, css)
	if err != nil {
		return err
	}

	return applyForAll(els, mod)
}

func ApplyForIt(w dom.Document, css string, mod Modifier) error {
	it, err := SelectFirstElement(w, css)
	if err != nil {
		return err
	}

	return mod(it)
}

func (p *Page) ApplyForIt(w dom.Document, css string, mod Modifier) error {
	it, err := p.SelectFirstElement(w, css)
	if err != nil {
		return err
	}

	return mod(it)
}

func applyForAll(els []dom.Element, mod Modifier) error {
	for _, e := range els {
		if err := mod(e); err != nil {
			return err
		}
	}

	return nil
}

type Modifier func(el dom.Element) error

func DisableElement(el dom.Element) error {
	el.SetAttribute("disabled", "")
	return nil
}

func EnableElement(el dom.Element) error {
	el.RemoveAttribute("disabled")
	return nil
}

func HideElement(el dom.Element) error {
	el.SetAttribute("style", "display: none;")
	return nil
}

func ShowElement(el dom.Element) error {
	el.SetAttribute("style", "display: block;")
	return nil
}

func VisibleElement(el dom.Element) error {
	el.SetAttribute("style", "visibility: visible;")
	return nil
}

func InvisibleElement(el dom.Element) error {
	el.SetAttribute("style", "visibility: hidden;")
	return nil
}

func VisitedElement(el dom.Element) error {
	return nil
}

func UnvisitedElement(el dom.Element) error {
	return nil
}

func ReadOnlyElement(el dom.Element) error {
	el.SetAttribute("readonly", "")
	return nil
}

func WriteOnlyElement(el dom.Element) error {
	el.RemoveAttribute("readonly")
	return nil
}

func Checked(el dom.Element) error {
	el.SetAttribute("checked", "")
	return nil
}

func Unchecked(el dom.Element) error {
	el.RemoveAttribute("checked")
	return nil
}

func ClearValue(el dom.Element) error {
	el.SetNodeValue("")
	return nil
}

func ClearHTML(el dom.Element) error {
	el.SetTextContent("")
	return nil
}

func DestroyHTML(el dom.Element) error {
	return ClearHTML(el)
}
