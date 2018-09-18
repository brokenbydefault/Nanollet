// +build !js

package DOM

func (el *Element) SelectAllElement(css string) ([]*Element, error) {
	e, err := el.el.Select(css)
	return NewElements(e), err
}

func (el *Element) SelectFirstElement(css string) (*Element, error) {
	e, err := el.el.SelectFirst(css)
	return NewElement(e), err
}

func (dom *DOM) SelectAllElement(css string) ([]*Element, error) {
	return dom.el.SelectAllElement(css)
}

func (dom *DOM) SelectFirstElement(css string) (*Element, error) {
	return dom.el.SelectFirstElement(css)
}