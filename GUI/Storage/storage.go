package Storage

import (
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/shibukawa/configdir"
)

var Permanent *configdir.Config

func Start() {
	name := "Nanollet"

	if Config.IsDebugEnabled() {
		name = "Nanollet-DEBUG"
	}

	Permanent = configdir.New("BrokenByDefault", name).QueryFolders(configdir.Global)[0]
}
