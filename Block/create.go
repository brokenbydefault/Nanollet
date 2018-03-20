package Block

import (
	"encoding/json"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"errors"
)

//@TODO
// Easier way to create blocks

func GetTypeFromJSON(jsn []byte) string {
	dblock := DefaultBlock{}

	err := json.Unmarshal(jsn, &dblock)
	if err != nil {
		return ""
	}

	return dblock.Type
}

//@TODO Create use JSONMarshal instead
func NewBlockFromJSON(jsn []byte) (blk UniversalBlock, err error) {
	blkSerialized := SerializedUniversalBlock{}

	err = json.Unmarshal(jsn, &blkSerialized)
	if err != nil {
		return
	}

	var errs []error
	blk.Type = blkSerialized.Type
	blk.Work, err = Util.UnsafeHexDecode(blkSerialized.Work)
	errs = append(errs, err)
	blk.Signature, err = Util.UnsafeHexDecode(blkSerialized.Signature)
	errs = append(errs, err)

	if blkSerialized.Account != "" {
		blk.Account, err = Wallet.Address(blkSerialized.Account).GetPublicKey()
		errs = append(errs, err)
	}

	if blkSerialized.Representative != "" {
		blk.Representative, err = Wallet.Address(blkSerialized.Representative).GetPublicKey()
		errs = append(errs, err)
	}

	blk.Amount, _ = Numbers.NewRawFromString("0")
	if blkSerialized.Amount != "" {
		blk.Amount, err = Numbers.NewRawFromHex(blkSerialized.Amount)
		errs = append(errs, err)
	}

	if blkSerialized.Balance != "" {
		blk.Balance, err = Numbers.NewRawFromHex(blkSerialized.Balance)
		errs = append(errs, err)
	}

	if blkSerialized.Previous != "" {
		blk.Previous, err = Util.UnsafeHexDecode(blkSerialized.Previous)
		errs = append(errs, err)
	}

	switch blkSerialized.Type {
	case "send":
		blk.Destination, err = Wallet.Address(blkSerialized.Destination).GetPublicKey()
		errs = append(errs, err)
	case "receive":
		fallthrough
	case "open":
		blk.Source = []byte(blkSerialized.Source)
	}

	err = Util.CheckError(errs)
	return
}

func CreateSignedSendBlock(sk *Wallet.SecretKey, sending, balance *Numbers.RawAmount, previous []byte, destination *Wallet.Address) (blk *SendBlock, err error) {
	var errs []error

	finalbalance := balance.Subtract(sending)
	if !finalbalance.IsValid() {
		err = errors.New("valid amount")
		return
	}

	blk = &SendBlock{}
	blk.Type = "send"
	blk.Balance = finalbalance
	blk.Previous = previous
	blk.Destination, err = destination.GetPublicKey()
	errs = append(errs, err)
	blk.Signature, err = sk.CreateSignature(blk.Hash())
	errs = append(errs, err)

	err = Util.CheckError(errs)
	return
}

func CreateSignedOpenBlock(sk *Wallet.SecretKey, source []byte) (blk *OpenBlock, err error) {
	var errs []error

	blk = &OpenBlock{}
	blk.Type = "open"
	blk.Source = source
	blk.Representative, err = sk.PublicKey()
	errs = append(errs, err)
	blk.Account, err = sk.PublicKey()
	errs = append(errs, err)
	blk.Signature, err = sk.CreateSignature(blk.Hash())
	errs = append(errs, err)

	err = Util.CheckError(errs)
	return
}

func CreateSignedReceiveBlock(sk *Wallet.SecretKey, source, previous []byte) (blk *ReceiveBlock, err error) {
	var errs []error

	blk = &ReceiveBlock{}
	blk.Type = "receive"
	blk.Source = source
	blk.Previous = previous
	blk.Signature, err = sk.CreateSignature(blk.Hash())
	errs = append(errs, err)

	err = Util.CheckError(errs)
	return
}

func CreateSignedReceiveOrOpenBlock(sk *Wallet.SecretKey, source, previous []byte) (blk BlockTransaction, err error) {
	if previous == nil {
		return CreateSignedOpenBlock(sk, source)
	}

	return CreateSignedReceiveBlock(sk, source, previous)
}

func CreateSignedChangeBlock(sk *Wallet.SecretKey, previous []byte, representative *Wallet.Address) (blk *ChangeBlock, err error) {
	var errs []error

	blk = &ChangeBlock{}
	blk.Type = "change"
	blk.Previous = previous
	blk.Representative, err = representative.GetPublicKey()
	errs = append(errs, err)
	blk.Signature, err = sk.CreateSignature(blk.Hash())
	errs = append(errs, err)

	err = Util.CheckError(errs)
	return
}
