package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/binary"
	"errors"
)

//@TODO Support votes by blocks

type ConfirmACKPackage struct {
	PublicKey    Wallet.PublicKey
	secretKey    Wallet.SecretKey
	Signature    Wallet.Signature
	Sequence     int64
	Transaction  Block.Transaction
	transactions []Block.Transaction
	Hashes       []Block.BlockHash
}

const (
	ConfirmACKPackageSizeMin = 32 + 64 + 8 + 32
	ConfirmACKPackageSizeMax = PackageSize
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
)

var votePrefix = []byte("vote ")

func NewConfirmACKPackage(sk Wallet.SecretKey, txs ...Block.Transaction) (packet *ConfirmACKPackage) {
	pk, _ := sk.PublicKey()

	return &ConfirmACKPackage{
		PublicKey:    pk,
		secretKey:    sk,
		Sequence:     time.Now().Unix(),
		transactions: txs,
	}
}

func (p *ConfirmACKPackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
	if p == nil || p.transactions == nil {
		return nil
	}

	data = make([]byte, ConfirmACKPackageSizeMax)

	copy(data[0:32], p.PublicKey)
	binary.LittleEndian.PutUint64(data[96:104], uint64(p.Sequence))

	var end = 104
	if len(p.transactions) > 1 {
		for _, tx := range p.transactions {
			end += copy(data[end:], tx.Hash())
		}

		sig, err := p.secretKey.CreateSignature(Util.CreateHash(32, votePrefix, data[104:end], data[96:104]))
		if err != nil {
			return nil
		}

		copy(data[32:96], sig)
	} else {
		end += copy(data[104:], p.transactions[0].Encode()[1:])

		sig, err := p.secretKey.CreateSignature(Util.CreateHash(32, p.transactions[0].Hash(), data[96:104]))
		if err != nil {
			return nil
		}

		copy(data[32:96], sig)
	}

	return data[:end]
}

func (p *ConfirmACKPackage) Decode(rHeader *Header, data []byte) (err error) {
	if p == nil {
		return
	}

	if rHeader == nil {
		return ErrInvalidHeaderParameters
	}

	if l := len(data); l <= ConfirmACKPackageSizeMin || l > ConfirmACKPackageSizeMax {
		return ErrInvalidMessageSize
	}

	p.PublicKey, p.Signature = make([]byte, 32), make([]byte, 64)

	bi := 0
	bi += copy(p.PublicKey, data[bi:bi+32])
	bi += copy(p.Signature, data[bi:bi+64])
	p.Sequence = int64(binary.LittleEndian.Uint64(data[bi:bi+8]))
	bi += 8

	if blktype := rHeader.ExtensionType.GetBlockType(); blktype == Block.NotABlock {
		if l := len(data) - bi; l <= 0 || l%32 != 0 {
			return ErrInvalidMessageSize
		}

		for bi = 104; bi < len(data); bi += 32 {
			p.Hashes = append(p.Hashes, data[bi:])
		}

		if ok := p.PublicKey.IsValidSignature(Util.CreateHash(32, votePrefix, data[104:], data[96:104]), p.Signature); !ok {
			return ErrInvalidSignature
		}
	} else {
		tx, _, err := Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
		if err != nil {
			return err
		}

		p.Transaction = tx

		if err = p.Transaction.Decode(data[bi:]); err != nil {
			return err
		}

		if ok := p.PublicKey.IsValidSignature(Util.CreateHash(32, p.Transaction.Hash(), data[96:104]), p.Signature); !ok {
			return ErrInvalidSignature
		}
	}

	return nil
}

func (p *ConfirmACKPackage) ModifyHeader(h *Header) {
	h.MessageType = ConfirmACK

	if p.Transaction != nil {
		h.ExtensionType.Add(ExtensionType(uint8(p.Transaction.GetType())) << 8)
	}

	if len(p.transactions) == 1 {
		h.ExtensionType.Add(ExtensionType(uint8(p.transactions[0].GetType())) << 8)
	}

	if p.Hashes != nil || len(p.transactions) > 1 {
		h.ExtensionType.Add(ExtensionType(Block.NotABlock) << 8)
	}
}
