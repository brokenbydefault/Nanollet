package DOM

import (
	"strings"
	"fmt"
)

func (w *Window) InitApplication(app Application) {
	w.StartApplication(app)

	if !app.HaveSidebar() {
		button, err := w.root.SelectFirstElement(`.control button[id="`+ strings.Title(app.Name()) +`"]`)
		if err == nil {
			DestroyHTML(button.el)
		}

		return
	}

	button, err := w.root.SelectFirstElement(`.control button[id="`+ strings.Title(app.Name()) +`"]`)
	if err != nil {
		panic(fmt.Sprintf("element %s was not found", `.control button[id="`+ strings.Title(app.Name()) +`"]`))
	}

	button.On(Click, func(class string) {
		w.root.ApplyForAll(".control button", DisableElement)
		defer w.root.ApplyForAll(".control button", EnableElement)
		w.ViewApplication(app)
	})

	for _, p := range app.Pages() {
		page := p

		pagebutton, err := w.root.SelectFirstElement(".control aside button."+strings.Title(page.Name()))
		if err != nil {
			panic(fmt.Sprintf("element %s was not found", `".control aside button.`+strings.Title(page.Name())))
		}

		pagebutton.On(Click, func(class string) {
			w.root.ApplyForAll(".control button", DisableElement)
			defer w.root.ApplyForAll(".control button", EnableElement)
			w.ViewPage(page)
		})
	}
}

func (w *Window) ViewApplication(app Application) error {
	w.root.ApplyForAll(".application button, [page]", HideElement)

	if app.HaveSidebar() {
		el, _ := w.root.SelectFirstElement("body")
		el.SetAttr("class", "")
	}

	w.root.ApplyForAll(".application#"+strings.Title(app.Name())+" button", ShowElement)

	return w.ViewPage(app.Pages()[0])
}

func (w *Window) StartApplication(app Application) {
	domAPP := NewDOMApplication(app, w)

	for _, p := range app.Pages() {
		page := p
		dom := NewDOMPage(page, domAPP)

		btns, err := dom.SelectAllElement(`button, input[type="submit"]`)
		if err != nil {
			return
		}

		for _, btn := range btns {
			btn.On(Click, func(class string) {
				defer dom.ApplyForAll(`button, input[type="submit"]`, EnableElement)
				dom.ApplyForAll(`button, input[type="submit"]`, DisableElement)
				page.OnContinue(w, dom, class)
			})
		}
	}

	return
}

func (w *Window) ViewPage(page Page) error {
	w.root.ApplyForAll(".control button", UnvisitedElement)
	w.root.ApplyFor(".control button."+strings.Title(page.Name()), VisitedElement)

	w.root.ApplyForAll("[page]", HideElement)

	dom := NewDOMPage(page, w.root)
	page.OnView(w, dom)
	dom.el.Apply(ShowElement)

	return nil
}
