package RPCClient

import (
	"github.com/brokenbydefault/Nanollet/RPC/rpctypes"
	"crypto/subtle"
	"errors"
	"github.com/brokenbydefault/MFA/mfatypes"
)

func SubscribeMFA(c rpctypes.Connection, pk [32]byte) (err error) {
	req := mfatypes.SubscribeRequest{
		PublicKey: pk[:],
		DefaultRequest: mfatypes.DefaultRequest{
			Action: "subscribe",
			App:    "nanomfa",
		},
	}

	var resp mfatypes.Subscription
	err = c.SendRequestJSON(&req, &resp)

	if subtle.ConstantTimeCompare(resp.PublicKey, pk[:]) == 0 {
		return errors.New("not subscribed")
	}

	return err
}

func UnsubscribeMFA(c rpctypes.Connection) (err error) {
	req := mfatypes.SubscribeRequest{
		DefaultRequest: mfatypes.DefaultRequest{
			Action: "unsubscribe",
			App:    "nanomfa",
		},
	}

	var resp mfatypes.Subscription
	err = c.SendRequestJSON(&req, &resp)

	if resp.PublicKey != nil {
		return errors.New("still subscribed")
	}

	return err
}

func SendToken(c rpctypes.Connection, destination []byte, env []byte) (err error) {
	req := mfatypes.EnvelopeRequest{
		Envelope:  env[:],
		PublicKey: destination[:],
		DefaultRequest: mfatypes.DefaultRequest{
			Action: "send",
			App:    "nanomfa",
		},
	}

	var resp mfatypes.Subscription
	err = c.SendRequestJSON(&req, &resp)
	if err != nil {
		return err
	}

	if resp.Error != "" {
		err = errors.New(resp.Error)
	}

	return err
}
