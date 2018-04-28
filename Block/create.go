package Block

import (
	"encoding/json"
	"errors"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"crypto/subtle"
	"bytes"
)

func NewBlockFromJSON(jsn []byte) (UniversalBlock, error) {
	ub := UniversalBlock{
		Balance: Numbers.NewRaw(),
		Amount:  Numbers.NewRaw(),
	}

	// RPC INCONSISTENCY: All calls uses `balance` as numeric string, but the
	// `Send` (the legacy-one) uses hex.
	fixme := struct {
		Type BlockType
		Balance string
	}{}
	json.Unmarshal(jsn, &fixme)

	if fixme.Type == Send {
		amm, err := Numbers.NewRawFromHex(fixme.Balance)
		if err != nil {
			return ub, err
		}
		jsn = bytes.Replace(jsn, []byte(`"balance": "`+fixme.Balance+`"`), []byte(`"balance": "`+amm.ToString()+`"`), 1)
	}
	//////////////////////////////////////////////////////////////////////

	return ub, json.Unmarshal(jsn, &ub)
}

func parseAccount(errs []error, f func() (Wallet.PublicKey, error)) Wallet.Address {
	pk, err := f()
	errs = append(errs, err)

	return pk.CreateAddress()
}

func parseRepresentative(errs []error, addr Wallet.Address) Wallet.Address {
	if !addr.IsValid() {
		errs = append(errs, errors.New("invalid representative"))
	}

	return addr
}

func parseDestination(errs []error, f func() (Wallet.PublicKey, error)) Wallet.PublicKey {
	pk, err := f()
	errs = append(errs, err)

	return pk
}

func parseBalance(errs []error, n *Numbers.RawAmount) *Numbers.RawAmount {
	if ok := n.IsValid(); !ok {
		errs = append(errs, errors.New("valid amount"))
	}

	return n
}

func attachSignature(sk Wallet.SecretKey, blk BlockTransaction) (BlockTransaction, error) {
	sig, err := sk.CreateSignature(blk.Hash())
	if err != nil {
		return nil, err
	}

	blk.SetSignature(sig)
	return blk, err
}

func CreateSignedUniversalSendBlock(sk Wallet.SecretKey, representative Wallet.Address, balance, sending *Numbers.RawAmount, previous BlockHash, destination Wallet.Address) (BlockTransaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    State,
			SubType: Send,
		},
		Account:        parseAccount(errs, sk.PublicKey),
		Representative: parseRepresentative(errs, representative),
		Balance:        parseBalance(errs, balance.Subtract(sending)),
		Previous:       previous,
		Link:           BlockHash(parseDestination(errs, destination.GetPublicKey)),
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalOpenBlock(sk Wallet.SecretKey, representative Wallet.Address, receiving *Numbers.RawAmount, source BlockHash) (BlockTransaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    State,
			SubType: Open,
		},
		Account:        parseAccount(errs, sk.PublicKey),
		Representative: parseRepresentative(errs, representative),
		Balance:        parseBalance(errs, receiving),
		Previous:       make([]byte, 32),
		Link:           source,
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalReceiveBlock(sk Wallet.SecretKey, representative Wallet.Address, balance, receiving *Numbers.RawAmount, previous BlockHash, source BlockHash) (BlockTransaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    State,
			SubType: Receive,
		},
		Account:        parseAccount(errs, sk.PublicKey),
		Representative: parseRepresentative(errs, representative),
		Balance:        parseBalance(errs, balance.Add(receiving)),
		Previous:       previous,
		Link:           source,
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalChangeBlock(sk Wallet.SecretKey, representative Wallet.Address, balance *Numbers.RawAmount, previous BlockHash) (BlockTransaction, error) {
	var errs []error

	blk := &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    State,
			SubType: Change,
		},
		Account:        parseAccount(errs, sk.PublicKey),
		Representative: parseRepresentative(errs, representative),
		Balance:        parseBalance(errs, balance),
		Previous:       previous,
		Link:           make([]byte, 32),
	}

	if err := Util.CheckError(errs); err != nil {
		return nil, err
	}

	return attachSignature(sk, blk)
}

func CreateSignedUniversalReceiveOrOpenBlock(sk Wallet.SecretKey, representative Wallet.Address, balance, receiving *Numbers.RawAmount, previous BlockHash, source BlockHash) (BlockTransaction, error) {
	pk, err := sk.PublicKey()
	if err != nil {
		return nil, err
	}

	if previous == nil || subtle.ConstantTimeCompare(previous, pk) == 1 {
		return CreateSignedUniversalOpenBlock(sk, representative, receiving, source)
	}

	return CreateSignedUniversalReceiveBlock(sk, representative, balance, receiving, previous, source)
}
