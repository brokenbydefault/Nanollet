package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/binary"
)

//@TODO Support votes by blocks

type ConfirmReqPackage struct {
	Transaction Block.Transaction
}

const (
	ConfirmACKPackageSizeMin = 32 + 64 + 8 + 32
	ConfirmACKPackageSizeMax = PackageSize
)

var votePrefix = []byte("vote ")

func NewConfirmReqPackage(tx Block.Transaction) (packet *ConfirmReqPackage) {
	return &ConfirmReqPackage{
		Transaction: tx,
	}
}

func (p *ConfirmReqPackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
	tx := p.Transaction

	/**
   if !rHeader.ExtensionType.Is(ExtendedNode) {
	   if tx.GetType() == Block.State && tx.GetSubType() != Block.Invalid {
		   tx = tx.SwitchToUniversalBlock(nil, nil).SwitchTo(tx.GetSubType())
	   }

	   return tx.Encode()[1:]
   }
   **/

	return tx.Encode()[1:]
}

func (p *ConfirmReqPackage) Decode(rHeader *Header, data []byte) (err error) {
	p.Transaction, _, err = Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
	if err != nil {
		return err
	}

	return p.Transaction.Decode(data)
}

func (p *PushPackage) ModifyHeader(h *Header) {
	h.MessageType = ConfirmReq
	h.ExtensionType.Add(ExtensionType(uint8(p.Transaction.GetType())&0x0F) << 8)
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
