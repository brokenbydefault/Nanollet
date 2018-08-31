package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"bufio"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Util"
)

type BulkPullAccountPackageRequest struct {
	PublicKey            Wallet.PublicKey
	MinimumPendingAmount *Numbers.RawAmount
	Flags                uint8
}

type BulkPullAccountPackageResponse struct {
	Frontier Block.BlockHash
	Balance  *Numbers.RawAmount
	Pending  []PendingInformation
}

type PendingInformation struct {
	Hash   Block.BlockHash
	Amount *Numbers.RawAmount
}

func NewBulkPullAccountPackageRequest(pk Wallet.PublicKey, minAmount *Numbers.RawAmount) (packet *BulkPullAccountPackageRequest) {
	return &BulkPullAccountPackageRequest{
		PublicKey:            pk,
		MinimumPendingAmount: minAmount,
	}
}

func NewBulkPullAccountPackageResponse(frontier Block.BlockHash, balance *Numbers.RawAmount, pending []PendingInformation) (packet *BulkPullAccountPackageResponse) {
	return &BulkPullAccountPackageResponse{
		Frontier: frontier,
		Balance:  balance,
		Pending:  pending,
	}
}

func (p *BulkPullAccountPackageRequest) Encode(_ *Header, dst io.Writer) (err error) {
	if p == nil {
		return
	}

	if _, err = dst.Write(p.PublicKey[:]); err != nil {
		return err
	}

	if _, err = dst.Write(p.MinimumPendingAmount.ToBytes()[:]); err != nil {
		return err
	}

	if _, err = dst.Write([]byte{byte(p.Flags)}); err != nil {
		return err
	}

	return nil
}

func (p *BulkPullAccountPackageRequest) Decode(_ *Header, src io.Reader) (err error) {
	if p == nil {
		return
	}

	if _, err := src.Read(p.PublicKey[:]); err != nil {
		return ErrInvalidMessageSize
	}

	if err := p.MinimumPendingAmount.Read(src); err != nil {
		return ErrInvalidMessageSize
	}

	i := [1]byte{}
	if n, err := src.Read(i[:]); n != 1 || err != nil {
		return ErrInvalidMessageSize
	}

	p.Flags = i[0]

	return nil
}

func (p *BulkPullAccountPackageResponse) Encode(_ *Header, dst io.Writer) (err error) {
	if p == nil {
		return
	}

	for _, tx := range p.Pending {
		if _, err = dst.Write(tx.Hash[:]); err != nil {
			return err
		}
		if err = tx.Amount.Write(dst); err != nil {
			return err
		}
	}

	if _, err = dst.Write([]byte{byte(Block.NotABlock)}); err != nil {
		return err
	}

	return nil
}

func (p *BulkPullAccountPackageResponse) Decode(_ *Header, src io.Reader) (err error) {
	reader := bufio.NewReader(src)

	if _, err := reader.Read(p.Frontier[:]); err != nil {
		return err
	}

	if err := p.Balance.Read(reader); err != nil {
		return err
	}

	for {
		pending := PendingInformation{}

		if _, err := reader.Read(pending.Hash[:]); err != nil {
			return err
		}

		if err := pending.Amount.Read(reader); err != nil {
			return err
		}

		if Util.IsEmpty(pending.Hash[:]) {
			return nil
		}

		p.Pending = append(p.Pending, pending)
	}

	return nil
}

func (p *BulkPullAccountPackageResponse) ModifyHeader(h *Header) {
	h.SetRemoveHeader(true)
}

func (p *BulkPullAccountPackageRequest) ModifyHeader(h *Header) {
	h.MessageType = BulkPullAccount
}
