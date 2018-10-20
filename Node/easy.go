package Node

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Node/Packets"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Block"
	"errors"
	"github.com/brokenbydefault/Nanollet/Util"
)

var (
	ErrNoOpenBlock        = errors.New("the open-block was not found")
	ErrInvalidSignature   = errors.New("signature doesn't match")
	ErrLegacyNotSupported = errors.New("legacy blocks not supported")
)

// GetInformation retrievers basic information (balance and frontier) of the account. It will request nodes of the
// network.
func GetInformation(c Node, pk *Wallet.PublicKey) (balance *Numbers.RawAmount, frontier Block.BlockHash, err error) {
	req := Packets.NewBulkPullAccountPackageRequest(*pk, Numbers.NewMax())

	packets, cancel := c.SendTCP(req, Packets.BulkPullAccount)
	for packet := range packets {
		p, ok := packet.(*Packets.BulkPullAccountPackageResponse)
		if !ok {
			continue
		}

		cancel()
		balance = p.Balance
		frontier = p.Frontier
	}

	return balance, frontier, nil
}

// GetBalance is a wrapper of GetInformation
func GetBalance(c Node, pk *Wallet.PublicKey) (balance *Numbers.RawAmount, err error) {
	balance, _, err = GetInformation(c, pk)
	return
}

func GetPendings(c Node, pk *Wallet.PublicKey, minAmount *Numbers.RawAmount) (txs []Block.Transaction, err error) {
	req := Packets.NewBulkPullAccountPackageRequest(*pk, minAmount)

	packets, _ := c.SendTCP(req, Packets.BulkPullAccount)

	pends := make(map[Block.BlockHash]Block.Transaction)
	for packet := range packets {
		p, ok := packet.(*Packets.BulkPullAccountPackageResponse)
		if !ok {
			continue
		}

		for _, hash := range p.Pending {
			if _, ok := pends[hash]; !ok {
				tx, err := GetBlock(c, &hash)
				if tx == nil || err != nil {
					break
				}

				if dest, _ := tx.GetTarget(); dest != *pk {
					break
				}

				if typ := tx.GetType(); typ != Block.State && typ != Block.Send {
					break
				}

				pends[hash] = tx
			}
		}
	}

	for _, tx := range pends {
		txs = append(txs, tx)
	}

	return txs, nil
}

func GetAmount(c Node, tx Block.Transaction) (*Numbers.RawAmount, error) {
	if typ := tx.GetType(); typ != Block.State && typ != Block.Send {
		return nil, ErrLegacyNotSupported
	}

	hashPrev := tx.GetPrevious()
	if Util.IsEmpty(hashPrev[:]) {
		return tx.GetBalance(), nil
	}

	txPrev, err := GetBlock(c, &hashPrev)
	if err != nil {
		return nil, err
	}

	if typ := txPrev.GetType(); typ != Block.State && typ != Block.Send {
		return nil, ErrLegacyNotSupported
	}

	return txPrev.GetBalance().Subtract(tx.GetBalance()), nil
}

func GetHistory(c Node, pk *Wallet.PublicKey, start *Block.BlockHash) (txs []Block.Transaction, err error) {
	if start == nil {
		h := Block.NewBlockHash(nil)
		start = &h
	}

	req := Packets.NewBulkPullPackageRequest(*pk, *start)

	packets, _ := c.SendTCP(req, Packets.BulkPull)
	for packet := range packets {
		p, ok := packet.(*Packets.BulkPullPackageResponse)
		if !ok {
			continue
		}

		l := len(p.Transactions)
		if l <= 0 {
			continue
		}

		// If the last block is not open
		if hashPrev := p.Transactions[l-1].GetPrevious(); !Util.IsEmpty(hashPrev[:]) {
			continue
		}

		pktx := p.Transactions[l-1].GetAccount()
		if pktx != *pk {
			continue
		}

		lastPreviousHash := p.Transactions[0].Hash()
		for _, tx := range p.Transactions {
			hash, sig := tx.Hash(), tx.GetSignature()

			if (!pk.IsValidSignature(hash[:], &sig) && !Block.IsEpoch(tx)) || hash != lastPreviousHash {
				continue
			}

			lastPreviousHash = tx.GetPrevious()
		}

		if len(p.Transactions) >= len(txs) {
			txs, err = p.Transactions, nil
		}
	}

	return txs, err
}

func GetBlock(c Node, hash *Block.BlockHash) (tx Block.Transaction, err error) {
	req := Packets.NewBulkPullPackageRequest(Wallet.PublicKey(*hash), *hash)

	packets, cancel := c.SendTCP(req, Packets.BulkPull)

	for packet := range packets {
		p, ok := packet.(*Packets.BulkPullPackageResponse)
		if !ok || len(p.Transactions) != 1 {
			continue
		}

		if p.Transactions[0].Hash() == *hash {
			cancel()
		}

		tx = p.Transactions[0]
		return
	}

	return tx, err
}

func PostBlock(c Node, tx Block.Transaction) (err error) {
	req := Packets.NewPushPackage(tx)

	if err := c.SendUDP(req); err != nil {
		return err
	}

	return nil
}

func RequestVotes(c Node, tx Block.Transaction) (err error) {
	req := Packets.NewConfirmReqPackage(tx)

	if err := c.SendUDP(req); err != nil {
		return err
	}

	return nil
}

