package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/binary"
)

//@TODO Support votes by blocks

type ConfirmACKPackage struct {
	PublicKey   Wallet.PublicKey
	secretKey   Wallet.SecretKey
	Signature   Wallet.Signature
	Sequence    int64
	Transaction Block.Transaction
	Hashes      []Block.BlockHash
}

const (
	ConfirmACKPackageSizeMin = 32 + 64 + 8 + 32
	ConfirmACKPackageSizeMax = PackageSize
)

var votePrefix = []byte("vote ")

func NewConfirmACKPackage(sk Wallet.SecretKey, hashes []Block.BlockHash) (packet *ConfirmACKPackage) {
	pk, _ := sk.PublicKey()

	return &ConfirmACKPackage{
		PublicKey: pk,
		secretKey: sk,
		Sequence:  time.Now().Unix(),
		Hashes:    hashes,
	}
}

func (p *ConfirmACKPackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
	if p == nil || p.Hashes == nil || rHeader.VersionUsing < 13 {
		return nil
	}

	if p.Hashes != nil {

		l := len(p.Hashes)
		data = make([]byte, ConfirmACKPackageSizeMin+(32*l))

		copy(data[0:32], p.PublicKey)
		binary.LittleEndian.PutUint64(data[96:104], uint64(p.Sequence))
		for i, h := range p.Hashes {
			copy(data[104+(i*32):], h)
		}

		sig, err := p.secretKey.CreateSignature(Util.CreateHash(32, votePrefix, data[104+(l*32):], data[96:104]))
		if err != nil {
			return nil
		}

		copy(data[32:96], sig)
	}

	if p.Transaction != nil {
		//@TODO Support votes by block
	}

	/**
   if !rHeader.ExtensionType.Is(ExtendedNode) {
	   if tx.GetType() == Block.State && tx.GetSubType() != Block.Invalid {
		   tx = tx.SwitchToUniversalBlock(nil, nil).SwitchTo(tx.GetSubType())
	   }

	   return tx.Encode()[1:]
   }
   **/

	return data
}

func (p *ConfirmACKPackage) Decode(rHeader *Header, data []byte) (err error) {

	if bt := rHeader.ExtensionType.GetBlockType(); bt == Block.NotABlock {

	} else {

		tx, size, err := Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
		if err != nil {
			return err
		}

		p.Transaction = tx

		if err = p.Transaction.Decode(data[104:size]); err != nil {
			return err
		}
	}

	return nil
}

func (p *ConfirmACKPackage) ModifyHeader(h *Header) {
	h.MessageType = ConfirmACK

	if p.Transaction != nil {
		h.ExtensionType.Add(ExtensionType(uint8(p.Transaction.GetType())&0x0F) << 8)
	}

	if p.Hashes != nil {
		h.ExtensionType.Add(ExtensionType(Block.NotABlock))
	}
}
