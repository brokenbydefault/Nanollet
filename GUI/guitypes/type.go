package guitypes

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/GUI/Front"
)

type Page interface {
	OnContinue(w *window.Window)
	OnView(w *window.Window)
	Name() string
}

type Application interface {
	Pages() []Page
	Name() string
	Display() Front.HTMLPAGE
	HaveSidebar() bool
}

type Sector struct {}
type App struct {}