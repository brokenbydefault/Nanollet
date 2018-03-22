package Background

import (
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Block"
	"time"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
)

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

	Storage.Frontier = info.Frontier
	Storage.Amount = info.Balance
	DOM.UpdateAmount(w)

	hist, err := RPCClient.GetAccountHistory(Connectivity.Socket, 25, Storage.PK.CreateAddress())
	if err != nil {
		return err
	}

	Storage.History.Set(hist)

	go startTracking(w)
	return nil
}

// @TODO Use the callback from the RPC instead of polling
// For now we need to polling the server each 10 seconds, far from optimal.
func startTracking(w *window.Window) {

	for range time.Tick(10 * time.Second) {

		min, _ := Numbers.NewRawFromString("1000000000000000000000000") //@TODO Enable user to change the minimum amount
		pends, err := RPCClient.GetAccountPending(Connectivity.Socket, 10, min, Storage.PK.CreateAddress())
		if err != nil {
			continue
		}

		if len(pends) > 0 {
			DOM.UpdateNotification(w, "New payments were detected in the network")
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

}
