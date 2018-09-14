package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"time"
	"github.com/brokenbydefault/Nanollet/Util"
	"errors"
	"golang.org/x/crypto/blake2b"
	"encoding/binary"
)

//@TODO Support votes by blocks

type ConfirmACKPackage struct {
	PublicKey Wallet.PublicKey
	Signature Wallet.Signature
	Sequence  [8]byte
	Hashes    []Block.BlockHash
}

const (
	ConfirmACKPackageSizeMin = 32 + 64 + 8 + 32
	ConfirmACKPackageSizeMax = MessageSize
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
)

var votePrefix = []byte("vote ")

func NewConfirmACKPackage(sk *Wallet.SecretKey, txs ...Block.Transaction) (packet *ConfirmACKPackage) {
	packet = &ConfirmACKPackage{
		PublicKey: sk.PublicKey(),
		Sequence:  newSequence(),
	}

	for _, tx := range txs {
		packet.Hashes = append(packet.Hashes, tx.Hash())
	}

	sig, err := sk.Sign(Util.CreateHash(32, votePrefix, concatHashes(packet.Hashes), packet.Sequence[:]))
	if err != nil {
		return
	}

	packet.Signature = sig

	return packet
}

func (p *ConfirmACKPackage) Encode(dst []byte) (n int, err error) {
	if p == nil {
		return 0, nil
	}

	if len(dst) < ConfirmACKPackageSizeMax {
		return 0, ErrDestinationLenghtNotEnough
	}

	if err != nil {
		return 0, err
	}

	n += copy(dst[n:], p.PublicKey[:])
	n += copy(dst[n:], p.Signature[:])
	n += copy(dst[n:], p.Sequence[:])
	n += copy(dst[n:], concatHashes(p.Hashes))

	return n, nil
}

func (p *ConfirmACKPackage) Decode(rHeader *Header, src []byte) (err error) {
	if p == nil {
		return
	}

	if rHeader == nil {
		return ErrInvalidHeaderParameters
	}

	if l := len(src); l < ConfirmACKPackageSizeMin || l > ConfirmACKPackageSizeMax {
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
			p.Hashes = append(p.Hashes, Block.NewBlockHash(src[i:]))
		}

		if ok := p.PublicKey.IsValidSignature(Util.CreateHash(32, votePrefix, src[bi:], p.Sequence[:]), &p.Signature); !ok {
			return ErrInvalidSignature
		}
	} else {
		tx, _, err := Block.NewTransaction(rHeader.ExtensionType.GetBlockType())
		if err != nil {
			return err
		}

		if err = tx.Decode(src[bi:]); err != nil {
			return err
		}

		p.Hashes = []Block.BlockHash{tx.Hash()}
		if ok := p.PublicKey.IsValidSignature(Util.CreateHash(32, p.Hashes[0][:], p.Sequence[:]), &p.Signature); !ok {
			return ErrInvalidSignature
		}
	}

	return nil
}

func (p *ConfirmACKPackage) ModifyHeader(h *Header) {
	h.MessageType = ConfirmACK
	h.ExtensionType.Add(ExtensionType(Block.NotABlock) << 8)
}

func concatHashes(hashes []Block.BlockHash) (b []byte) {
	for _, h := range hashes {
		b = append(b, h[:]...)
	}

	return b
}

func newSequence() (b [8]byte) {
	binary.PutVarint(b[:], time.Now().Unix())
	return b
}
