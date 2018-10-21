// +build !js

package DOM

import "errors"

var (
	ErrInvalidElement = errors.New("invalid element")
)

func (el *Element) SelectAllElement(css string) ([]*Element, error) {
	e, err := el.el.Select(css)
	if e == nil {
		err = ErrInvalidElement
	}
	return NewElements(e), err
}

func (el *Element) SelectFirstElement(css string) (*Element, error) {
	e, err := el.el.SelectFirst(css)
	if e == nil {
		err = ErrInvalidElement
	}
	return NewElement(e), err
}

func (dom *DOM) SelectAllElement(css string) ([]*Element, error) {
	return dom.el.SelectAllElement(css)
}

func (dom *DOM) SelectFirstElement(css string) (*Element, error) {
	return dom.el.SelectFirstElement(css)
}