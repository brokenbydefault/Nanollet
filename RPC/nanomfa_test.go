package RPCClient

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/MFA"
)

func TestSendToken(t *testing.T) {
	Connectivity.Socket.StartWebsocket()

	receiver, _ := Util.SecureHexDecode("79C4051F9E7C8C04B421736FD401E41186A9578C73A600B8FD54A096A89F7E6C")

	seedfy, _ := MFA.NewSeedFY()
	seed, _ := MFA.ReadSeedFY(seedfy.Encode(), "")
	token, _ := MFA.RecoverToken(seed, 0)

	sender, _ := MFA.NewSender(MFA.GenerateDevice())
	env, _ := sender.CreateEnvelope(receiver[:], token)

	if err := SendToken(Connectivity.Socket, receiver[:], env); err != nil {
		t.Error(err)
	}
}
