package DOM

import (
	"strings"
)

func (w *Window) InitApplication(app Application) {
	w.StartApplication(app)

	if !app.HaveSidebar() {
		return
	}

	controlBar, _ := w.root.SelectFirstElement(".control")
	moduleButton := controlBar.CreateElement("button", "")

	moduleButton.CreateElementWithAttr("span", strings.Title(app.Name()), Attrs{"class": "title"})
	moduleButton.CreateElementWithAttr("span", "", Attrs{"class": "pointer"})

	moduleButton.On(Click, func(class string) {
		w.root.ApplyForAll(".control button", DisableElement)
		defer w.root.ApplyForAll(".control button", EnableElement)
		w.ViewApplication(app)
	})

	aside := controlBar.CreateElementWithAttr("aside", "", Attrs{"class": "application", "id": app.Name()})

	for _, p := range app.Pages() {
		page := p

		controlButton := aside.CreateElementWithAttr("button", "", Attrs{"class": strings.Title(page.Name())})
		block := controlButton.CreateElementWithAttr("span", "", Attrs{"class": "block"})

		block.CreateElementWithAttr("icon", "", Attrs{"class": "icon-" + page.Name()})
		block.CreateElementWithAttr("span", strings.Title(page.Name()), Attrs{"class": "title"})
		block.CreateElementWithAttr("span", "", Attrs{"class": "pointer"})

		controlButton.On(Click, func(class string) {
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

	w.root.ApplyForAll(".application#"+app.Name()+" button", ShowElement)

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
