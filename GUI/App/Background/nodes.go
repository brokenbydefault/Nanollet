package Background

import (
	"time"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
)

func UpdateNodeCount(w *DOM.Window) {
	for range time.Tick(10 * time.Second) {
		DOM.UpdateNodesCount(w)
	}
}
