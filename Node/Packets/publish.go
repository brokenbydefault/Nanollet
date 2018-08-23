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

func (p *PushPackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
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

func (p *PushPackage) Decode(rHeader *Header, data []byte) (err error) {
	p.Transaction, _, err = Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
	if err != nil {
		return err
	}

	return p.Transaction.Decode(data)
}

func (p *PushPackage) ModifyHeader(h *Header) {
	h.MessageType = Publish
	h.ExtensionType |= ExtensionType(uint8(p.Transaction.GetType())&0x0F) << 8
}
