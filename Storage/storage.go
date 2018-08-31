package Storage

import (
	"github.com/shibukawa/configdir"
	"github.com/brokenbydefault/Nanollet/Config"
	"os"
)

var (
	Permanent *configdir.Config
)

func init() {
	Permanent = configdir.New("BrokenByDefault", Config.Configuration().DefaultFolder).QueryFolders(configdir.Global)[0]

	os.UserCacheDir()
}
