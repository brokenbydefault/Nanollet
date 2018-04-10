package Background

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"encoding/json"
	"bytes"
)

var min, _ = Numbers.NewRawFromString("1000000000000000000000000") //@TODO Enable user to change the minimum amount
var pends = make(chan []byte, 10)

func StartAddress(w *window.Window) error {
	info, err := RPCClient.GetAccountInformation(Connectivity.Socket, Storage.PK.CreateAddress())

	if err != nil {
		if err.Error() == RPCClient.ErrNotOpenedAccount.Error() {
			info.Balance, _ = Numbers.NewRawFromString("0")
			info.Frontier = nil
		} else {
			return err
		}
	}

	Storage.UpdateFrontier(info.Frontier)
	Storage.Amount = info.Balance
	DOM.UpdateAmount(w)

	hist, err := RPCClient.GetAccountHistory(Connectivity.Socket, 1000, Storage.PK.CreateAddress())
	if err != nil {
		return err
	}

	Storage.History.Set(hist)

	go func() {
		pendings(w)
		realtimeupdate(w)
	}()
	return nil
}

func realtimeupdate(w *window.Window) {
	conn := Connectivity.NewSocket()

	err := RPCClient.Subscribe(conn, Storage.PK)
	if err != nil {
		panic(err)
	}

	go conn.ReceiveAllMessages(nil, pends)

	pend := RPCClient.CallbackResponse{}

	for p := range pends {
		err := json.Unmarshal(p, &pend)
		if err != nil {
			continue
		}

		// If the destination is not the currently public-key or already sent a received: skip
		if !bytes.Equal(pend.Destination, Storage.PK) || Storage.History.AlreadyReceived(pend.Hash) {
			continue
		}

		blk, err := Block.CreateSignedReceiveOrOpenBlock(&Storage.SK, pend.Hash, Storage.Frontier)
		if err != nil || Storage.History.ExistHash(blk.Hash()) {
			continue
		}

		err = PublishBlockToQueue(blk, pend.Amount)
		if err != nil {
			continue
		}

		DOM.UpdateAmount(w)
		DOM.UpdateNotification(w, "You had received a new payment")
	}

	pendings(w)
	realtimeupdate(w)
}

func pendings(w *window.Window) {
	pends, err := RPCClient.GetAccountPending(Connectivity.Socket, 1000, min, Storage.PK.CreateAddress())
	if err != nil {
		return
	}

	for _, pend := range pends {
		blk, err := Block.CreateSignedReceiveOrOpenBlock(&Storage.SK, pend.Hash, Storage.Frontier)
		if err != nil {
			continue
		}

		err = PublishBlockToQueue(blk, pend.Amount)
		if err != nil {
			continue
		}

		DOM.UpdateAmount(w)
		DOM.UpdateNotification(w, "You had received a new payment")
	}
}
