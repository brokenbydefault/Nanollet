package Packets

import (
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"bufio"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/Inkeliz/blakEd25519"
	"encoding/binary"
)

type FrontierReqPackageRequest struct {
	PublicKey Wallet.PublicKey
	Age       uint32
	Count     uint32
}

func NewFrontierReqPackageRequest(pk Wallet.PublicKey, age, end uint32) (packet *FrontierReqPackageRequest) {
	return &FrontierReqPackageRequest{
		PublicKey: pk,
		Age:       age,
		Count:     end,
	}
}

func (p *FrontierReqPackageRequest) Encode(_ *Header, dst io.Writer) (err error) {
	if p == nil {
		return
	}

	if _, err = dst.Write(p.PublicKey[:]); err != nil {
		return err
	}

	if err = binary.Write(dst, binary.BigEndian, p.Age); err != nil {
		return err
	}

	if err = binary.Write(dst, binary.BigEndian, p.Count); err != nil {
		return err
	}

	return nil
}

func (p *FrontierReqPackageRequest) Decode(_ *Header, src io.Reader) (err error) {
	if p == nil {
		return
	}

	if n, err := src.Read(p.PublicKey[:]); n != blakEd25519.PublicKeySize || err != nil {
		return ErrInvalidMessageSize
	}

	b := make([]byte, 4)
	if n, err := src.Read(b); n != 8 || err != nil {
		return ErrInvalidMessageSize
	}

	p.Age = binary.BigEndian.Uint32(b)

	if n, err := src.Read(b); n != 8 || err != nil {
		return ErrInvalidMessageSize
	}

	p.Count = binary.BigEndian.Uint32(b)

	return nil
}

type Frontier struct {
	Account Wallet.PublicKey
	Hash    Block.BlockHash
}

type FrontierReqPackageResponse struct {
	Frontiers    []Frontier
	transactions []Block.Transaction
}

func NewFrontierReqPackageResponse(txs []Block.Transaction) (packet *FrontierReqPackageResponse) {
	return &FrontierReqPackageResponse{
		transactions: txs,
	}
}

func (p *FrontierReqPackageResponse) Encode(_ *Header, dst io.Writer) (err error) {
	if p == nil {
		return
	}

	for _, tx := range p.transactions {
		hash := tx.Hash()
		if _, err = dst.Write(hash[:]); err != nil {
			return err
		}
	}

	if _, err = dst.Write([]byte{byte(Block.NotABlock)}); err != nil {
		return err
	}

	return nil
}

func (p *FrontierReqPackageResponse) Decode(_ *Header, src io.Reader) (err error) {
	reader := bufio.NewReader(src)

	for {

		frontier := Frontier{}

		blockType, err := reader.ReadByte()
		if err != nil {
			return err
		}

		if blockType == byte(Block.Invalid) {
			return nil
		}

		if _, err = reader.Read(frontier.Account[:]); err != nil {
			return err
		}

		if _, err = reader.Read(frontier.Hash[:]); err != nil {
			return err
		}

		p.Frontiers = append(p.Frontiers, )
	}
}

func (p *FrontierReqPackageResponse) ModifyHeader(h *Header) {
	h.SetRemoveHeader(true)
}

func (p *FrontierReqPackageRequest) ModifyHeader(h *Header) {
	h.MessageType = FrontierReq
}
