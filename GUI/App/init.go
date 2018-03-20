package App

import (
	"github.com/sciter-sdk/go-sciter/window"
	"strings"
	"github.com/sciter-sdk/go-sciter"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/GUI/guitypes"
)

func InitApplication(w *window.Window, app guitypes.Application) {
	elem, err := DOM.SelectFirstElement(w, ".dynamic")
	elem.SetHtml(string(app.Display()), sciter.SIH_APPEND_AFTER_LAST)

	StartApplication(w, app)

	menubutton, err := DOM.SelectFirstElement(w, ".apps #"+strings.ToLower(app.Name()))
	if err != nil {
		return
	}

	menubutton.OnClick(func() {
		ViewApplication(w, app)
	})


	controlbar, _ := DOM.SelectFirstElement(w, ".control")
	control, _ := sciter.CreateElement("aside", "")
	control.SetAttr("application", app.Name())
	controlbar.Append(control)

	for _, p := range app.Pages() {
		page := p

		controlbutton, _ := sciter.CreateElement("button", "")
		controlbutton.SetAttr("class", p.Name())
		control.Append(controlbutton)

		icon, _ := sciter.CreateElement("icon", "")
		icon.SetAttr("class", "icon-"+p.Name())
		controlbutton.Append(icon)
		controlbutton.OnClick(func() {
			ViewPage(w, page)
		})
	}
}

func ViewApplication(w *window.Window, app guitypes.Application) error {
	DOM.ApplyForAll(w, "[application], [page]", DOM.HideElement)

	if app.HaveSidebar() {
		el, _ := DOM.SelectFirstElement(w, "body")
		el.SetAttr("class", "")
	}

	DOM.ApplyForAll(w, "[application=\""+app.Name()+"\"]", DOM.ShowElement)

	return ViewPage(w, app.Pages()[0])
}

func StartApplication(w *window.Window, app guitypes.Application) error {

	for _, p := range app.Pages() {
		sector := DOM.SetSector(p)
		page := p

		continuebtn, err := sector.SelectFirstElement(w, ".continue")
		if err == nil {
			continuebtn.OnClick(func() {
				go func() {
					sector.ApplyForIt(w, ".continue", DOM.DisableElement)
					page.OnContinue(w)
					sector.ApplyForIt(w, ".continue", DOM.EnableElement)
				}()
			})
		}

	}

	return nil
}

func ViewPage(w *window.Window, page guitypes.Page) error {

	DOM.ApplyForAll(w, ".control button", DOM.UnvisitedElement)
	DOM.ApplyForIt(w, ".control button."+page.Name(), DOM.VisitedElement)

	DOM.ApplyForAll(w, "[page]", DOM.HideElement)
	page.OnView(w)

	DOM.ApplyForIt(w, "[page=\""+page.Name()+"\"]", DOM.ShowElement)

	return nil
}
