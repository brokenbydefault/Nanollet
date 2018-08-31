package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"errors"
	"golang.org/x/crypto/blake2b"
)

//@TODO Support votes by blocks

type ConfirmACKPackage struct {
	PublicKey    Wallet.PublicKey
	secretKey    Wallet.SecretKey
	Signature    Wallet.Signature
	Sequence     [8]byte
	Transaction  Block.Transaction
	transactions []Block.Transaction
	Hashes       []Block.BlockHash
}

const (
	ConfirmACKPackageSizeMin = 32 + 64 + 8 + 32
	ConfirmACKPackageSizeMax = MessageSize
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
)

var votePrefix = []byte("vote ")

func NewConfirmACKPackage(sk Wallet.SecretKey, txs ...Block.Transaction) (packet *ConfirmACKPackage) {
	packet = &ConfirmACKPackage{
		PublicKey:    sk.PublicKey(),
		secretKey:    sk,
		transactions: txs,
	}

	copy(packet.Sequence[:], Util.UintToBytes(uint64(time.Now().Unix()), Util.BigEndian))
	return packet
}

func (p *ConfirmACKPackage) Encode(rHeader *Header, dst []byte) (n int, err error) {
	if p == nil || p.transactions == nil {
		return 0, nil
	}

	if len(dst) < ConfirmACKPackageSizeMax {
		return 0, ErrDestinationLenghtNotEnough
	}

	var blk []byte
	if len(p.transactions) > 1 {
		for _, tx := range p.transactions {
			hash := tx.Hash()
			blk = append(blk, hash[:]...)
		}

		p.Signature, err = p.secretKey.CreateSignature(Util.CreateHash(32, votePrefix, blk, p.Sequence[:]))
	} else {
		blk = p.transactions[0].Encode()[1:]

		hash := p.transactions[0].Hash()
		p.Signature, err = p.secretKey.CreateSignature(Util.CreateHash(32, hash[:], p.Sequence[:]))
	}

	if err != nil {
		return 0, err
	}

	n += copy(dst[n:], p.PublicKey[:])
	n += copy(dst[n:], p.Signature[:])
	n += copy(dst[n:], p.Sequence[:])
	n += copy(dst[n:], blk)

	return n, nil
}

func (p *ConfirmACKPackage) Decode(rHeader *Header, src []byte) (err error) {
	if p == nil {
		return
	}

	if rHeader == nil {
		return ErrInvalidHeaderParameters
	}

	if l := len(src); l <= ConfirmACKPackageSizeMin || l > ConfirmACKPackageSizeMax {
		return ErrInvalidMessageSize
	}

	bi := 0
	bi += copy(p.PublicKey[:], src[bi:32])
	bi += copy(p.Signature[:], src[bi:bi+64])
	bi += copy(p.Sequence[:], src[bi:bi+8])

	if blktype := rHeader.ExtensionType.GetBlockType(); blktype == Block.NotABlock {
		l := len(src)
		if l-bi <= 0 || (l-bi)%blake2b.Size256 != 0 {
			return ErrInvalidMessageSize
		}

		for i := bi; i < l; i += blake2b.Size256 {
			p.Hashes = append(p.Hashes, Block.NewBlockHash(src[bi:]))
		}

		if ok := p.PublicKey.IsValidSignature(Util.CreateHash(32, votePrefix, src[bi:], p.Sequence[:]), p.Signature); !ok {
			return ErrInvalidSignature
		}
	} else {
		p.Transaction, _, err = Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
		if err != nil {
			return err
		}

		if err = p.Transaction.Decode(src[bi:]); err != nil {
			return err
		}

		hash := p.Transaction.Hash()
		if ok := p.PublicKey.IsValidSignature(Util.CreateHash(32, hash[:], p.Sequence[:]), p.Signature); !ok {
			return ErrInvalidSignature
		}
	}

	return nil
}

func (p *ConfirmACKPackage) ModifyHeader(h *Header) {
	h.MessageType = ConfirmACK

	if p.Transaction != nil || len(p.transactions) == 1 {
		if p.Transaction != nil {
			h.ExtensionType.Add(ExtensionType(uint8(p.Transaction.GetType())) << 8)
		} else {
			h.ExtensionType.Add(ExtensionType(uint8(p.transactions[0].GetType())) << 8)
		}
	}

	if p.Hashes != nil || len(p.transactions) > 1 {
		h.ExtensionType.Add(ExtensionType(Block.NotABlock) << 8)
	}
}
