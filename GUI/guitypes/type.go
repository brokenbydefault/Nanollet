package guitypes

import (
	"github.com/brokenbydefault/Nanollet/GUI/Front"
	"github.com/sciter-sdk/go-sciter/window"
)

type Page interface {
	OnContinue(w *window.Window, action string)
	OnView(w *window.Window)
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
