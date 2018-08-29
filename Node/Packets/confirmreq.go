package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
)

//@TODO Support votes by blocks

type ConfirmReqPackage struct {
	Transaction Block.Transaction
}

const (
	ConfirmReqPackageSizeMin = Block.ReceiveSize
	ConfirmReqPackageSizeMax = PackageSize
)

func NewConfirmReqPackage(tx Block.Transaction) (packet *ConfirmReqPackage) {
	return &ConfirmReqPackage{
		Transaction: tx,
	}
}

func (p *ConfirmReqPackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
	if p == nil {
		return
	}

	return p.Transaction.Encode()[1:]
}

func (p *ConfirmReqPackage) Decode(rHeader *Header, data []byte) (err error) {
	if p == nil {
		return
	}

	if rHeader == nil {
		return ErrInvalidHeaderParameters
	}

	p.Transaction, _, err = Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
	if err != nil {
		return err
	}

	if err := p.Transaction.Decode(data); err != nil {
		return err
	}

	return nil
}

func (p *ConfirmReqPackage) ModifyHeader(h *Header) {
	h.MessageType = ConfirmReq
	h.ExtensionType.Add(ExtensionType(uint8(p.Transaction.GetType())) << 8)
}
