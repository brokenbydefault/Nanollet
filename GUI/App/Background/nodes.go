package Background

import (
	"time"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/sciter-sdk/go-sciter/window"
)

func UpdateNodeCount(w *window.Window) {
	for range time.Tick(10 * time.Second) {
		DOM.UpdateNodesCount(w)
	}
}
