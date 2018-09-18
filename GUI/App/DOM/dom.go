// +build !js

package DOM

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"
)

type Window struct {
	win  *window.Window
	root *DOM
}

type DOM struct {
	el *Element
}

type Element struct {
	el *sciter.Element
}

func NewWindow(w *window.Window) *Window {
	el, err := w.GetRootElement()
	if err != nil {
		panic(err)
	}

	return &Window{
		win:  w,
		root: &DOM{el: &Element{el: el}},
	}
}

func NewDOMApplication(app Application, w *Window) *DOM {
	el, err := w.win.GetRootElement()
	if err != nil {
		panic(err)
	}

	el, err = el.SelectFirst("[application=\"" + app.Name() + "\"]")
	if err != nil {
		panic(err)
	}

	return &DOM{el: &Element{el}}
}

func NewDOMPage(page Page, dom *DOM) *DOM {
	el, err := dom.el.el.SelectFirst("[page=\"" + page.Name() + "\"]")
	if err != nil {
		panic(err)
	}

	el.UID()

	return &DOM{el: &Element{el}}
}

func NewElement(sciterEl *sciter.Element) *Element {
	return &Element{
		sciterEl,
	}
}

func NewElements(sciterEls []*sciter.Element) []*Element {
	els := make([]*Element, len(sciterEls))

	for i, el := range sciterEls {
		els[i] = &Element{el,}
	}

	return els
}
