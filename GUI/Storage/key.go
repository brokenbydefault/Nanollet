package Storage

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
)

var SK Wallet.SecretKey
var PK Wallet.PublicKey

func SetPrivateKey(sk Wallet.SecretKey, pk Wallet.PublicKey) {
	SK = sk
	PK = pk
}
