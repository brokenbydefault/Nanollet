package TwoFactor

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/TwoFactor/Ephemeral"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
)

func TestNewRequesterServer(t *testing.T) {

	var received bool

	computer := Ephemeral.NewEphemeral()
	_, smartphone, _ := Wallet.GenerateRandomKeyPair()

	seedfy, _ := NewSeedFY()
	token, _ := NewToken(seedfy.String(), []byte("123456"))

	request, response := NewRequesterServer(&computer, []Wallet.PublicKey{smartphone.PublicKey()})
	go func() {
		envelope := <-response
		received = true

		if envelope.Capsule.Token != token {
			t.Error("token not equal")
		}

	}()

	if err := ReplyRequest(&smartphone, token, request); err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)

	if !received {
		t.Error("package not received/sent")
	}
}
