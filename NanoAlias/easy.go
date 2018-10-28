package NanoAlias

import (
	"fmt"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/GUI/App/Background"
	"github.com/brokenbydefault/Nanollet/Node"
	"github.com/brokenbydefault/Nanollet/Storage"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
)

func IsAvailable(alias string) bool {
	if _, err := Address(alias).GetPublicKey(); err == ErrNotFound {
		return true
	}

	return false
}

func Register(sk *Wallet.SecretKey, previous Block.Transaction, alias string) error {
	aliasPK, aliasSK, err := Address(alias).GetAliasKey()
	if err != nil {
		return err
	}

	send, err := CreateSendBlock(aliasPK, sk, previous)
	if err != nil {
		return err
	}

	if err := Background.PublishBlockToQueue(send, Block.Send, Amount); err != nil {
		return err
	}

	txs, err := CreateAliasBlocks(&aliasSK, send)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		tx.Work()

		Storage.TransactionStorage.Add(tx)
		if err := Node.PostBlock(Background.Connection, tx); err != nil {
			return err
		}

		hash := tx.Hash()
		Node.RequestVotes(Background.Connection, tx)

		winner, ok := Storage.TransactionStorage.WaitConfirmation(&Storage.Configuration.Account.Quorum, 30 * time.Second, &hash)
		if !ok || *winner != hash {
			fmt.Print(ok, winner)
			return ErrUnconfirmedRegister
		}

	}

	return nil
}

func CreateSendBlock(aliasPK Wallet.PublicKey, sk *Wallet.SecretKey, previous Block.Transaction) (tx Block.Transaction, err error) {
	if previous.GetType() != Block.State {
		return nil, ErrUnsupportedBlock
	}

	return Block.CreateUniversalSendBlock(sk, previous.SwitchToUniversalBlock(nil, nil).Representative, Amount, previous.GetBalance(), previous.Hash(), aliasPK)
}

func CreateAliasBlocks(aliasSK *Wallet.SecretKey, sendTx Block.Transaction) (txs []Block.Transaction, err error) {
	if sendTx.GetType() != Block.State {
		return nil, ErrInvalidAliasBlock
	}

	txs = make([]Block.Transaction, 2)
	txs[0], err = Block.CreateUniversalOpenBlock(aliasSK, Representative, Amount, sendTx.Hash())
	if err != nil {
		return nil, err
	}

	txs[1], err = Block.CreateUniversalSendBlock(aliasSK, Representative, txs[0].GetBalance(), txs[0].GetBalance(), txs[0].Hash(), sendTx.GetAccount())
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func IsValidOpenBlock(txOpen Block.Transaction) bool {
	if txOpen.GetType() != Block.State {
		return false
	}

	if txOpen.GetBalance().Compare(Amount) != 0 {
		return false
	}

	if txOpen.SwitchToUniversalBlock(nil, nil).Representative != Representative {
		return false
	}

	return true
}
