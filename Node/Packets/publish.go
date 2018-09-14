package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
)

type PushPackage struct {
	Transaction Block.Transaction
}

func NewPushPackage(transaction Block.Transaction) (packet *PushPackage) {
	return &PushPackage{
		transaction,
	}
}

func (p *PushPackage) Encode(dst []byte) (n int, err error)  {
	if p == nil {
		return
	}

	n += copy(dst, p.Transaction.Encode()[1:])

	return n, err
}

func (p *PushPackage) Decode(rHeader *Header, data []byte) (err error) {
	if p == nil {
		return nil
	}

	if rHeader == nil {
		return ErrInvalidHeaderParameters
	}

	p.Transaction, _, err = Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
	if err != nil {
		return err
	}

	return p.Transaction.Decode(data)
}

func (p *PushPackage) ModifyHeader(h *Header) {
	h.MessageType = Publish
	h.ExtensionType |= ExtensionType(uint8(p.Transaction.GetType())) << 8
}
