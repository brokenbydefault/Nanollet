// +build js

package DOM

import (
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
	"honnef.co/go/js/dom"
	"errors"
)

var (
	ErrInvalidElement = errors.New("invalid element")
)

type Page string

func SetPage(name string) Page {
	return Page(name)
}

func SetSector(pg guitypes.Page) Page {
	return Page(pg.Name())
}

func (p Page) SelectAllElement(w dom.Document, css string) ([]dom.Element, error) {
	el := w.QuerySelector("[page=\"" + string(p) + "\"]")
	if el == nil {
		return nil, ErrInvalidElement
	}

	els := el.QuerySelectorAll(css)
	if len(els) == 0 {
		return nil, ErrInvalidElement
	}

	return els, nil
}

func (p Page) SelectFirstElement(w dom.Document, css string) (dom.Element, error) {
	els, err := p.SelectAllElement(w, css)
	if err != nil {
		return nil, err
	}

	return els[0], nil
}

func SelectAllElement(w dom.Document, css string) ([]dom.Element, error) {
	els := w.QuerySelectorAll(css)
	if len(els) == 0 {
		return nil, ErrInvalidElement
	}

	return els, nil
}

func SelectFirstElement(w dom.Document, css string) (dom.Element, error) {
	els, err := SelectAllElement(w, css)
	if err != nil {
		return nil, ErrInvalidElement
	}

	return els[0], nil
}
