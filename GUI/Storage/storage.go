package Storage

import (
	"github.com/shibukawa/configdir"
	"github.com/brokenbydefault/Nanollet/Config"
)

var Permanent *configdir.Config

func Start() {
	name := "Nanollet"

	if Config.IsDebugEnabled() {
		name = "Nanollet-DEBUG"
	}

	Permanent = configdir.New("BrokenByDefault", name).QueryFolders(configdir.Global)[0]
}
