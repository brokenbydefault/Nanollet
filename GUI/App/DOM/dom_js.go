// +build js

package DOM

import (
	"honnef.co/go/js/dom"
)

type Window struct {
	win  *dom.Window
	root *DOM
}

type DOM struct {
	el *Element
}

type Element struct {
	el dom.Element
}

func NewWindow(w dom.Window) *Window {
	return &Window{
		win:  &w,
		root: &DOM{el: &Element{el: w.Document().QuerySelector("body")}},
	}
}

func NewDOMApplication(app Application, w *Window) *DOM {
	el, err := w.root.SelectFirstElement("[application=\"" + app.Name() + "\"]")
	if err != nil {
		panic(err)
	}

	return &DOM{el: el}
}

func NewDOMPage(page Page, dom *DOM) *DOM {
	el, err := dom.SelectFirstElement("[page=\"" + page.Name() + "\"]")
	if err != nil {
		panic(err)
	}

	return &DOM{el: el}
}

func NewElement(jsEl dom.Element) *Element {
	return &Element{jsEl}
}

func NewElements(jsEls []dom.Element) []*Element {
	els := make([]*Element, len(jsEls))

	for i, el := range jsEls {
		els[i] = &Element{el}
	}

	return els
}
