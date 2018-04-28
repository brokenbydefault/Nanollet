package RPCClient

import (
	"crypto/subtle"
	"errors"
	"github.com/brokenbydefault/Nanollet/RPC/rpctypes"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func Subscribe(c rpctypes.Connection, pk Wallet.PublicKey) (err error) {
	req := SubscribeRequest{
		PublicKey: pk,
		DefaultRequest: DefaultRequest{
			Action: "subscribe",
			App:    "nanosub",
		},
	}

	var resp Subscription
	err = c.SendRequestJSON(&req, &resp)

	if subtle.ConstantTimeCompare(resp.PublicKey, pk) == 0 {
		return errors.New("not subscribed")
	}

	return err
}

func Unsubscribe(c rpctypes.Connection) (err error) {
	req := SubscribeRequest{
		DefaultRequest: DefaultRequest{
			Action: "unsubscribe",
			App:    "nanosub",
		},
	}

	var resp Subscription
	err = c.SendRequestJSON(&req, &resp)

	if resp.PublicKey != nil {
		return errors.New("still subscribed")
	}

	return err
}
