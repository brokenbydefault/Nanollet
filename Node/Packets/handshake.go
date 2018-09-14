package Packets

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/Inkeliz/blakEd25519"
	"github.com/brokenbydefault/Nanollet/Util"
)

type HandshakePackage struct {
	Challenge [32]byte
	PublicKey Wallet.PublicKey
	Signature Wallet.Signature
}

const (
	NodeHandshakePackageSizeMin = 32
	NodeHandshakePackageSizeMax = NodeHandshakePackageSizeMin + blakEd25519.PublicKeySize + blakEd25519.SignatureSize
)

const (
	_         ExtensionType = iota
	Challenge
	Response
)

func NewHandshakePackage(lChallenge []byte, rChallenge []byte) (packet *HandshakePackage) {
	packet = new(HandshakePackage)

	if lChallenge != nil {
		copy(packet.Challenge[:], lChallenge)
	}

	if rChallenge != nil {
		pk, sk, _ := Wallet.GenerateRandomKeyPair()

		packet.PublicKey = pk
		sig, err := sk.Sign(rChallenge)
		if err != nil {
			return packet
		}
		packet.Signature = sig
	}

	return packet
}

func (p *HandshakePackage) Encode(dst []byte) (n int, err error) {
	if p == nil {
		return
	}

	if len(dst) < NodeHandshakePackageSizeMax {
		return 0, ErrDestinationLenghtNotEnough
	}

	if !Util.IsEmpty(p.Challenge[:]) {
		n += copy(dst[n:], p.Challenge[:])
	}

	if !Util.IsEmpty(p.Signature[:]) && !Util.IsEmpty(p.PublicKey[:]) {
		n += copy(dst[n:], p.PublicKey[:])
		n += copy(dst[n:], p.Signature[:])
	}

	return n, err
}

func (p *HandshakePackage) Decode(rHeader *Header, src []byte) (err error) {
	if p == nil {
		return
	}

	if rHeader == nil {
		return ErrInvalidHeaderParameters
	}

	if l := len(src); l > NodeHandshakePackageSizeMax || l < NodeHandshakePackageSizeMin {
		return ErrInvalidMessageSize
	}

	bi := 0
	if rHeader.ExtensionType.Is(Challenge) {
		bi += copy(p.Challenge[:], src[bi:bi+32])
	}

	if rHeader.ExtensionType.Is(Response) {
		bi += copy(p.PublicKey[:], src[bi:bi+32])
		bi += copy(p.Signature[:], src[bi:bi+64])
	}

	return nil
}

func (p *HandshakePackage) ModifyHeader(h *Header) {
	h.MessageType = NodeHandshake

	if !Util.IsEmpty(p.Challenge[:]) {
		h.ExtensionType.Add(Challenge)
	}

	if !Util.IsEmpty(p.Signature[:]) && !Util.IsEmpty(p.PublicKey[:]) {
		h.ExtensionType.Add(Response)
	}

}
