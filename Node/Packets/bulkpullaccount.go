package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"bufio"
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/binary"
)

type BulkPullAccountPackageRequest struct {
	PublicKey            Wallet.PublicKey
	MinimumPendingAmount *Numbers.RawAmount
	Flags                uint8
}

type BulkPullAccountPackageResponse struct {
	Frontier Block.BlockHash
	Balance  *Numbers.RawAmount
	Pending  []Block.BlockHash
}

func NewBulkPullAccountPackageRequest(pk Wallet.PublicKey, minAmount *Numbers.RawAmount) (packet *BulkPullAccountPackageRequest) {
	return &BulkPullAccountPackageRequest{
		PublicKey:            pk,
		MinimumPendingAmount: minAmount,
		Flags:                0,
	}
}

func NewBulkPullAccountPackageResponse(frontier Block.BlockHash, balance *Numbers.RawAmount, pending []Block.BlockHash) (packet *BulkPullAccountPackageResponse) {
	return &BulkPullAccountPackageResponse{
		Frontier: frontier,
		Balance:  balance,
		Pending:  pending,
	}
}

func (p *BulkPullAccountPackageRequest) Encode(dst io.Writer) (err error) {
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

	p.MinimumPendingAmount = new(Numbers.RawAmount)
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

//@TODO (inkeliz) implement Encode
func (p *BulkPullAccountPackageResponse) Encode(dst io.Writer) (err error) {
	if p == nil {
		return
	}

	/**
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
	**/

	return nil
}

func (p *BulkPullAccountPackageResponse) Decode(_ *Header, src io.Reader) (err error) {
	if p == nil {
		return
	}

	buf := bufio.NewReaderSize(src, 48)

	if _, err := buf.Read(p.Frontier[:]); err != nil {
		return err
	}

	p.Balance = new(Numbers.RawAmount)
	if err := p.Balance.Read(buf); err != nil {
		return err
	}

	for {
		hash := Block.BlockHash{}
		if err := binary.Read(buf, binary.BigEndian, &hash); err != nil {
			return nil
		}

		amount := make([]byte, 16)
		if err := binary.Read(buf, binary.BigEndian, &amount); err != nil {
			return nil
		}

		if Util.IsEmpty(hash[:]) {
			return nil
		}

		p.Pending = append(p.Pending, hash)
	}

	return nil
}

func (p *BulkPullAccountPackageResponse) ModifyHeader(h *Header) {
	h.SetRemoveHeader(true)
}

func (p *BulkPullAccountPackageRequest) ModifyHeader(h *Header) {
	h.MessageType = BulkPullAccount
}
