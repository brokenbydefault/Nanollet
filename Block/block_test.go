package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"testing"
)

func TestSendBlock_Hash(t *testing.T) {
	sb := SendBlock{}
	sb.Previous, _ = Util.UnsafeHexDecode("16E1238F60735B4ECF2CEFABED385B815D82FB9FB5FA6D1A7785CA8DB14B386C")
	sb.Balance, _ = Numbers.NewRawFromString("3000000000000000000000000")
	sb.Destination, _ = Wallet.Address("xrb_3yxxyyrdeapnxe1dyx1fha759syhxou4risfkzibe8bfdjj3663d7ppy1enb").GetPublicKey()
	sb.Work, _ = Util.UnsafeHexDecode("d000000000070080")

}
