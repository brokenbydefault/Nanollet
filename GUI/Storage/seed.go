package Storage

import (
	"github.com/shibukawa/configdir"
	"errors"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Config"
)

var dirconfig *configdir.Config
var SEED Wallet.Seed

func Start() {
	if Config.IsDebugEnabled() {
		dirconfig = configdir.New("BrokenByDefault", "Nanollet-DEBUG").QueryFolders(configdir.Global)[0]
	}else {
		dirconfig = configdir.New("BrokenByDefault", "Nanollet").QueryFolders(configdir.Global)[0]
	}
}

func ExistSeed() bool {
	if dirconfig == nil {
		return false
	}
	return dirconfig.Exists("wallet.dat")
}

func SaveSeed(seed string) error {
	if dirconfig == nil {
		return errors.New("impossible to store the data")
	}
	return dirconfig.WriteFile("wallet.dat", []byte(seed))
}

func RetrieveSeed() ([]byte, error) {
	if dirconfig == nil {
		return nil, errors.New("impossible to retrieve the data")
	}
	return dirconfig.ReadFile("wallet.dat")
}