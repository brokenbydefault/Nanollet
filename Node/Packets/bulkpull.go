package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"bufio"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/Inkeliz/blakEd25519"
	"golang.org/x/crypto/blake2b"
)

type BulkPullPackageResponse struct {
	Transactions []Block.Transaction
}

type BulkPullPackageRequest struct {
	PublicKey Wallet.PublicKey
	End       Block.BlockHash
}

func NewBulkPullResponse(txs []Block.Transaction) (packet *BulkPullPackageResponse) {
	return &BulkPullPackageResponse{
		Transactions: txs,
	}
}

func NewBulkPullPackageRequest(pk Wallet.PublicKey, end Block.BlockHash) (packet *BulkPullPackageRequest) {
	return &BulkPullPackageRequest{
		PublicKey: pk,
		End:       end,
	}
}

func (p *BulkPullPackageRequest) Encode(_ *Header, dst io.Writer) (err error) {
	if p == nil {
		return
	}

	if _, err = dst.Write(p.PublicKey[:]); err != nil {
		return err
	}

	if _, err = dst.Write(p.End[:]); err != nil {
		return err
	}

	return nil
}

func (p *BulkPullPackageRequest) Decode(_ *Header, src io.Reader) (err error) {
	if p == nil {
		return
	}

	if n, err := src.Read(p.PublicKey[:]); n != blakEd25519.PublicKeySize || err != nil {
		return ErrInvalidMessageSize
	}

	if n, err := src.Read(p.End[:]); n != blake2b.Size256 || err != nil {
		return ErrInvalidMessageSize
	}

	return nil
}

func (p *BulkPullPackageResponse) Encode(_ *Header, dst io.Writer) (err error) {
	if p == nil {
		return
	}

	for _, tx := range p.Transactions {
		if _, err = dst.Write(tx.Encode()); err != nil {
			return err
		}
	}

	if _, err = dst.Write([]byte{byte(Block.NotABlock)}); err != nil {
		return err
	}

	return nil
}

func (p *BulkPullPackageResponse) Decode(_ *Header, src io.Reader) (err error) {
	reader := bufio.NewReader(src)

	for {
		blockType, err := reader.ReadByte()
		if err != nil {
			return err
		}

		tx, size, err := Block.NewTransaction(Block.BlockType(blockType))
		if err != nil {
			if err == Block.ErrEndBlock {
				return nil
			}
			return err
		}

		btx := make([]byte, size)

		if _, err = reader.Read(btx); err != nil {
			return err
		}

		if err = tx.Decode(btx); err != nil {
			return err
		}

		p.Transactions = append(p.Transactions, tx)
	}
}

func (p *BulkPullPackageResponse) ModifyHeader(h *Header) {
	h.SetRemoveHeader(true)
}

func (p *BulkPullPackageRequest) ModifyHeader(h *Header) {
	h.MessageType = BulkPull
}
