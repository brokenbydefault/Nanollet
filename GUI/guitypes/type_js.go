// +build js

package guitypes

import (
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"honnef.co/go/js/dom"
)

type Page interface {
	OnContinue(w dom.Document, action string)
	OnView(w dom.Document)
	Name() string
}

type Application interface {
	Pages() []Page
	Name() string
	Display() Front.HTMLPAGE
	HaveSidebar() bool
}

type Sector struct{}
type App struct{}
