package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"bufio"
)

type BulkPullPackage struct {
	Transactions []Block.Transaction
}

func NewBulkPull(transactions []Block.Transaction) (packet *BulkPullPackage) {
	return &BulkPullPackage{
		transactions,
	}
}

func (p *BulkPullPackage) Encode(lHeader *Header, rHeader *Header, dst io.Writer) {
	for _, tx := range p.Transactions {
		dst.Write(tx.Encode())
	}
}

func (p *BulkPullPackage) Decode(rHeader *Header, src io.Reader) (err error) {
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

func (p *BulkPullPackage) ModifyHeader(h *Header) {
	h.MessageType = BulkPull
}
