// +build js

package DOM

import "honnef.co/go/js/dom"

func (el *Element) Apply(mod Modifier) error {
	return mod(el.el)
}

func (el *Element) SetAttr(name, value string) error {
	el.el.SetAttribute(name, value)
	return nil
}

func (el *Element) SetText(text string) error {
	el.el.SetTextContent(text)
	return nil
}

func (el *Element) SetValue(value string) error {
	el.el.SetNodeValue(value)
	return nil
}

func (el *Element) SetHTML(html string, method ReplaceMethod) (err error) {
	switch method {
	case InnerReplaceContent:
		el.el.SetInnerHTML(html)
	case InnerPrepend:
		el.el.SetInnerHTML(el.el.InnerHTML() + html)
	case InnerAppend:
		el.el.SetInnerHTML(html + el.el.InnerHTML())
	case OuterReplace:
		el.el.SetOuterHTML(html)
	case OuterPrepend:
		el.el.SetOuterHTML(html + el.el.OuterHTML())
	case OuterAppend:
		el.el.SetOuterHTML(el.el.OuterHTML() + html)
	default:
		err = ErrInvalidReplaceMethod
	}

	return err
}

func (el *Element) On(method ActionMethod, f func(class string)) (err error) {
	class, _ := el.GetAttr("class")

	switch method {
	case Click:
		el.el.AddEventListener("click", false, func(_ dom.Event) {
			go f(class)
		})
	default:
		err = ErrInvalidActionMethod
	}

	return err
}

func (dom *DOM) ApplyForAll(css string, mod Modifier) error {
	els, err := dom.SelectAllElement(css)
	if err != nil {
		return err
	}

	for _, el := range els {
		if err := el.Apply(mod); err != nil {
			return err
		}
	}

	return nil
}

func (dom *DOM) ApplyFor(css string, mod Modifier) error {
	el, err := dom.SelectFirstElement(css)
	if err != nil {
		return err
	}

	return el.Apply(mod)
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
