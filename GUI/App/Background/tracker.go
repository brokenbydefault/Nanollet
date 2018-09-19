package Background

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/App/DOM"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Node"
	"net"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"time"
)

var Connection Node.Node

func init() {
	Connection = &Node.Server{
		Peers:          &Storage.PeerStorage,
		Transactions:   &Storage.TransactionStorage,
		Header:         Storage.Configuration.Node.Header,
		PublishHandler: PublishHandler,
	}
	go Connection.Start()
}

func PublishHandler(srv *Node.Server, _ *net.UDPAddr, rHeader *Packets.Header, msg []byte) {
	packet := new(Packets.PushPackage)

	if err := packet.Decode(rHeader, msg); err != nil {
		return
	}

	dest, _ := packet.Transaction.GetTarget()
	acc := packet.Transaction.GetAccount()
	if acc != Storage.AccountStorage.PublicKey && dest != Storage.AccountStorage.PublicKey {
		return
	}

	if !packet.Transaction.IsValidPOW() {
		return
	}

	srv.Transactions.Add(packet.Transaction)
}

func StartAddress(w *DOM.Window) error {
	txs, err := Node.GetHistory(Connection, &Storage.AccountStorage.PublicKey, nil)
	if err != nil {
		return err
	}

	if len(txs) == 0 {
		Storage.AccountStorage.Frontier = Block.NewBlockHash(nil)
		Storage.AccountStorage.Representative = Storage.Configuration.Account.Representative
		Storage.AccountStorage.Balance = Numbers.NewMin()
	} else {

		if txs[0].GetType() == Block.State {
			Storage.AccountStorage.Frontier = txs[0].Hash()
			Storage.AccountStorage.Representative = txs[0].SwitchToUniversalBlock(nil, nil).Representative
			Storage.AccountStorage.Balance = txs[0].GetBalance()
		} else {
			balance, err := Node.GetBalance(Connection, &Storage.AccountStorage.PublicKey)
			if err == nil {
				return err
			}

			Storage.AccountStorage.Frontier = txs[0].Hash()
			for _, tx := range txs {
				if typ := tx.GetType(); typ == Block.Change || typ == Block.Open {
					Storage.AccountStorage.Representative = tx.SwitchToUniversalBlock(nil, nil).Representative
				}
			}
			Storage.AccountStorage.Balance = balance
		}

	}

	go realtimeUpdate(w)

	Storage.TransactionStorage.Add(txs...)
	DOM.UpdateAmount(w)

	go pending(w)

	return nil
}

func realtimeUpdate(w *DOM.Window) {
	for tx := range Storage.TransactionStorage.Listen() {
		tx := tx

		if dest, _ := tx.GetTarget(); dest != Storage.AccountStorage.PublicKey {
			continue
		}

		hash := tx.Hash()
		if tx, ok := Storage.TransactionStorage.GetByLinkHash(&hash); ok {
			hash, sig := tx.Hash(), tx.GetSignature()
			if Storage.AccountStorage.PublicKey.IsValidSignature(hash[:], &sig) {
				continue
			}
		}

		if !Storage.TransactionStorage.IsConfirmed(&hash, &Storage.Configuration.Account.Quorum) {
			DOM.UpdateNotification(w, "New payment identified, voting in progress.")
		}

		go acceptPending(w, tx)
	}
}

func acceptPending(w *DOM.Window, tx Block.Transaction) {
	hash := tx.Hash()

	if waitVotesConfirmation(tx, 6*time.Minute) {
		amount, err := Node.GetAmount(Connection, tx)
		if err != nil {
			return
		}

		blk, err := Block.CreateSignedUniversalReceiveOrOpenBlock(&Storage.AccountStorage.SecretKey, Storage.AccountStorage.Representative, Storage.AccountStorage.Balance, amount, Storage.AccountStorage.Frontier, hash)
		if err != nil {
			return
		}

		if err := PublishBlockToQueue(blk, Block.Receive, amount); err == nil {
			DOM.UpdateNotification(w, "You have received a new payment.")
			DOM.UpdateAmount(w)
		}
	}
}

func pending(w *DOM.Window) {
	txsPend, err := Node.GetPendings(Connection, &Storage.AccountStorage.PublicKey, Storage.Configuration.Account.MinimumAmount)
	if err != nil {
		return
	}

	Storage.TransactionStorage.Add(txsPend...)
}
