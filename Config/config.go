package Config

import "github.com/brokenbydefault/Nanollet/Wallet"

var Debug = false
var DefaultRepresentative = Wallet.Address("xrb_1ywcdyz7djjdaqbextj4wh1db3wykze5ueh9wnmbgrcykg3t5k1se7zyjf95")

func IsDebugEnabled() bool {
	return Debug
}
