// +build !js

package DOM

import (
	"github.com/sciter-sdk/go-sciter"
)

func (el *Element) Apply(mod Modifier) error {
	return mod(el.el)
}

func (el *Element) SetAttr(name, value string) error {
	return el.el.SetAttr(name, value)
}

func (el *Element) SetText(text string) error {
	return el.el.SetText(text)
}

func (el *Element) SetValue(value string) error {
	return el.el.SetValue(sciter.NewValue(value))
}

func (el *Element) SetHTML(html string, method ReplaceMethod) error {
	return el.el.SetHtml(html, sciter.SET_ELEMENT_HTML(method))
}

func (el *Element) On(method ActionMethod, f func(class string)) (err error) {
	class, _ := el.GetAttr("class")

	switch method {
	case Click:
		el.el.OnClick(func() {
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

type Modifier func(el *sciter.Element) error

func DisableElement(el *sciter.Element) error {
	return el.SetState(sciter.STATE_DISABLED, 0, true)
}

func EnableElement(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_DISABLED, true)
}

func HideElement(el *sciter.Element) error {
	return el.SetStyle("display", "none")
}

func ShowElement(el *sciter.Element) error {
	return el.SetStyle("display", "block")
}

func VisibleElement(el *sciter.Element) error {
	return el.SetStyle("visibility", "visible")
}

func InvisibleElement(el *sciter.Element) error {
	return el.SetStyle("visibility", "hidden")
}

func VisitedElement(el *sciter.Element) error {
	return el.SetState(sciter.STATE_VISITED, 0, true)
}

func UnvisitedElement(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_VISITED, true)
}

func ReadOnlyElement(el *sciter.Element) error {
	return el.SetState(sciter.STATE_READONLY, 0, true)
}

func WriteOnlyElement(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_READONLY, true)
}

func Checked(el *sciter.Element) error {
	return el.SetState(sciter.STATE_CHECKED, 0, true)
}

func Unchecked(el *sciter.Element) error {
	return el.SetState(0, sciter.STATE_CHECKED, true)
}

func ClearValue(el *sciter.Element) error {
	return el.SetValue(sciter.NewValue())
}

func ClearHTML(el *sciter.Element) error {
	return el.SetHtml(" ", sciter.SIH_REPLACE_CONTENT)
}

func DestroyHTML(el *sciter.Element) error {
	el.SetHtml(" ", sciter.SOH_REPLACE)
	return el.Clear()
}
