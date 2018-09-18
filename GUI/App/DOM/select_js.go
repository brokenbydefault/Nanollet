// +build js

package DOM

import (
	"errors"
)

var (
	ErrInvalidElement = errors.New("invalid element")
)

func (el *Element) SelectAllElement(css string) ([]*Element, error) {
	e := el.el.QuerySelectorAll(css)
	if len(e) == 0 {
		return nil, ErrInvalidElement
	}

	return NewElements(e), nil
}

func (el *Element) SelectFirstElement(css string) (*Element, error) {
	e, err := el.SelectAllElement(css)
	if err != nil {
		return nil, err
	}

	return e[0], nil
}

func (dom *DOM) SelectAllElement(css string) ([]*Element, error) {
	return dom.el.SelectAllElement(css)
}

func (dom *DOM) SelectFirstElement(css string) (*Element, error) {
	return dom.el.SelectFirstElement(css)
}