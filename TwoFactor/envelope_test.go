// +build !js

package TwoFactor

import (
	"testing"
	"bytes"
	"github.com/brokenbydefault/Nanollet/TwoFactor/Ephemeral"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func TestEncapsulateMFA(t *testing.T) {
	devicepk, devicesk, _ := Wallet.GenerateRandomKeyPair()

	computer := Ephemeral.NewEphemeral()
	smartphone := Ephemeral.NewEphemeral()

	seedfy, err := NewSeedFY()
	if err != nil {
		t.Error(err)
	}

	token, err := NewToken(seedfy.String(), []byte("123456789"))
	if err != nil {
		t.Error(err)
	}

	envelope := NewEnvelope(devicepk, smartphone.PublicKey(), computer.PublicKey(), token)
	envelope.Sign(&devicesk)
	envelope.Encrypt(&smartphone)
	smartphoneEnvelope, _ := envelope.MarshalBinary()

	rEnvelope := new(Envelope)
	rEnvelope.UnmarshalBinary(smartphoneEnvelope)
	rEnvelope.Decrypt(&computer)

	if !bytes.Equal(rEnvelope.Capsule.Token[:], token[:]) || bytes.Equal(envelope.Capsule.Token[:], token[:]) {
		t.Error("encryption wrong")
	}

	if !rEnvelope.IsValidSignature(nil) {
		t.Error("invalid signature")
	}
}
