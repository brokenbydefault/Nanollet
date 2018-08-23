package Packets

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/Inkeliz/blakEd25519"
	"bytes"
)

type HandshakePackage struct {
	Challenge [32]byte
	PublicKey Wallet.PublicKey
	Signature Wallet.Signature
}

const (
	nodeHandshakeChallengeSize = 32
	nodeHandshakeResponseSize  = blakEd25519.PublicKeySize + blakEd25519.SignatureSize
	NodeHandshakePackageSize   = nodeHandshakeChallengeSize + nodeHandshakeResponseSize
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
		packet.Signature, _ = sk.CreateSignature(rChallenge)
	}

	return packet
}

func (p *HandshakePackage) Encode(lHeader *Header, rHeader *Header) (data []byte) {
	if p == nil {
		*p = HandshakePackage{}
	}

	p.ModifyHeader(lHeader)
	bi, data := lHeader.Encode(NodeHandshakePackageSize)

	if bytes.Equal(p.Challenge[:], make([]byte, 32)) == false {
		bi += copy(data[bi:], p.Challenge[0:])
	}

	if p.Signature != nil && p.PublicKey != nil {
		bi += copy(data[bi:], p.PublicKey)
		bi += copy(data[bi:], p.Signature)
	}

	return data[:bi]
}

func (p *HandshakePackage) Decode(rHeader *Header, data []byte) (err error) {
	if p == nil {
		*p = HandshakePackage{}
	}

	bi := 0

	if rHeader.ExtensionType.Is(Challenge) {
		bi += copy(p.Challenge[:], data[bi:32])
	}

	if rHeader.ExtensionType.Is(Response) {
		p.PublicKey, p.Signature = make([]byte, 32), make([]byte, 64)

		bi += copy(p.PublicKey, data[bi:bi+32])
		bi += copy(p.Signature, data[bi:bi+64])
	}

	return nil
}

func (p *HandshakePackage) ModifyHeader(h *Header) {
	h.MessageType = NodeHandshake

	if bytes.Equal(p.Challenge[:], make([]byte, 32)) == false {
		h.ExtensionType |= Challenge
	}

	if p.Signature != nil && p.PublicKey != nil {
		h.ExtensionType |= Response
	}

}
