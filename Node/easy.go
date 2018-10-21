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
//
// This function can be maliciously affected, the node can reply with any arbitrary information.
//
// It's highly recommended to use the GetHistory instead GetInformation, which also provides
// the balance (if state/send) and the frontier.
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

// GetBalance is a wrapper of GetInformation, which retrieves a balance of a specific pk.
//
// It's highly recommended to use the GetHistory instead GetInformation, which also provides
// the balance (if state/send) and the frontier.
func GetBalance(c Node, pk *Wallet.PublicKey) (balance *Numbers.RawAmount, err error) {
	balance, _, err = GetInformation(c, pk)
	return
}

// GetPendings get all blocks which the given pk don't publish the respective "receive". This function can lie,
// since malicious nodes can send any block.
//
// It's highly recommended to request vote for each pending block before sent a "receive" for that block.
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

// GetAmount get the amount sent or received in the tx. It will subtract the amount of
// tx.Balance - txPrevious.Balance, it only works if the previous blocks are also a State and Send blocks.
//
// This function don't request votes, however if the given tx is already valid the previous tx will be also valid,
// since the hash is considered unique.
func GetAmount(c Node, tx Block.Transaction) (*Numbers.RawAmount, error) {
	if typ := tx.GetType(); typ != Block.State && typ != Block.Send {
		return nil, ErrLegacyNotSupported
	}

	hashPrev := tx.GetPrevious()
	if Util.IsEmpty(hashPrev[:]) {
		return tx.GetBalance(), nil
	}

	txPrev, err := GetBlock(c, &hashPrev)
	if txPrev ==nil || err != nil {
		return nil, err
	}

	if typ := txPrev.GetType(); typ != Block.State && typ != Block.Send {
		return nil, ErrLegacyNotSupported
	}

	return txPrev.GetBalance().Subtract(tx.GetBalance()), nil
}

// GetHistory retrieves all blocks for given pk. It will return the largest chain received from
// network. Since all blocks should be signed by the pk there's no way to forge the block without
// the knowing of the secret-key.
//
// This function don't request for vote, you need to request votes to verify if the block is valid on the network.
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

			if  hash != lastPreviousHash || (!pk.IsValidSignature(hash[:], &sig) && !Block.IsEpoch(tx)) {
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

// GetMultiplesHistory retrieves all blocks for given pk. It will all chains received from
// network, instead of the largest (GetHistory).
//
// This function don't request for vote, you need to request votes to verify if the block is valid on the network.
// For usual cases the `GetHistory` should be used.
func GetMultiplesHistory(c Node, pk *Wallet.PublicKey, start *Block.BlockHash) (txs map[Block.BlockHash][]Block.Transaction, err error) {
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

		// If already have this chain
		if _, ok := txs[p.Transactions[0].Hash()]; ok  {
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

			if hash != lastPreviousHash || (!pk.IsValidSignature(hash[:], &sig) && !Block.IsEpoch(tx)) {
				continue
			}

			lastPreviousHash = tx.GetPrevious()
		}

		txs[p.Transactions[0].Hash()] = p.Transactions
	}

	return txs, err
}

// GetBlock retrieves a single block from the network, using the hash. Since the hash is unique,
// there's no room for a malicious node forge the block, since the block most have the same hash.
//
// This function don't request for vote, you need to request votes to verify if the block is valid on the network.
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

		return p.Transactions[0], nil
	}

	return tx, ErrTCPNotAvailable
}

// PostBlock sends a `publish` for the given tx, publishing the block to
// network.
func PostBlock(c Node, tx Block.Transaction) (err error) {
	req := Packets.NewPushPackage(tx)

	if err := c.SendUDP(req); err != nil {
		return err
	}

	return nil
}

// RequestVotes sends a `confirm_req` for the given tx. It returns error if
// impossible to send the package. The votes are received by the ConfirmReqHandler.
func RequestVotes(c Node, tx Block.Transaction) (err error) {
	req := Packets.NewConfirmReqPackage(tx)

	if err := c.SendUDP(req); err != nil {
		return err
	}

	return nil
}

