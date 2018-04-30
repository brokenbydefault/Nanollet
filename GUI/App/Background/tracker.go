package Background

import (
	"bytes"
	"encoding/json"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/sciter-sdk/go-sciter/window"
	"time"
)

var min, _ = Numbers.NewRawFromString("1000000000000000000000000") //@TODO Enable user to change the minimum amount

func StartAddress(w *window.Window) error {
	info, err := RPCClient.GetAccountInformation(Connectivity.Socket, Storage.PK.CreateAddress())

	if err != nil {
		if err.Error() == RPCClient.ErrNotOpenedAccount.Error() {
			info.Balance, _ = Numbers.NewRawFromString("0")
			info.Representative = Config.DefaultRepresentative
		} else {
			return err
		}
	}

	Storage.Representative = info.Representative
	Storage.Amount = info.Balance

	if info.Frontier != nil {
		frontierblk, err := RPCClient.GetBlockByHash(Connectivity.Socket, info.Frontier)
		if err != nil {
			return err
		}

		Storage.UpdateFrontier(frontierblk.SwitchToUniversalBlock())
	}

	go Storage.UpdatePoW()
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
	pends := make(chan []byte, 64)

	err := RPCClient.Subscribe(conn, Storage.PK)
	if err != nil {
		time.Sleep(2 * time.Second)
		realtimeupdate(w)
		return
	}

	go func(){
		if err := conn.ReceiveAllMessages(nil, pends); err != nil {
			time.Sleep(2 * time.Second)
			pendings(w)
			realtimeupdate(w)
			return
		}
	}()

	for p := range pends {
		pend := RPCClient.CallbackResponse{}

		err := json.Unmarshal(p, &pend)
		if err != nil {
			continue
		}

		// If the destination is not the currently public-key or already sent a received: skip
		if !bytes.Equal(pend.Destination, Storage.PK) || Storage.History.AlreadyReceived(pend.Hash) {
			continue
		}

		blk, err := Block.CreateSignedUniversalReceiveOrOpenBlock(Storage.SK, Storage.Representative, Storage.Amount, pend.Amount, Storage.Frontier, pend.Hash)
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
}

func pendings(w *window.Window) {
	pends, err := RPCClient.GetAccountPending(Connectivity.Socket, 1000, min, Storage.PK.CreateAddress())
	if err != nil {
		return
	}

	for _, pend := range pends {
		blk, err := Block.CreateSignedUniversalReceiveOrOpenBlock(Storage.SK, Storage.Representative, Storage.Amount, pend.Amount, Storage.Frontier, pend.Hash)
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
