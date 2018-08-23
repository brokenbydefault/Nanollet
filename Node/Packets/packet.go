package Packets

import (
	"errors"
	"io"
)

const (
	MessageSize = 508
	PackageSize = HeaderSize + MessageSize
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message type")
)

type PacketUDP interface {
	Encode(lHeader *Header, rHeader *Header) (data []byte)
	Decode(rHeader *Header, data []byte) (err error)

	ModifyHeader(h *Header)
}

type PacketTCP interface {
	Encode(lHeader *Header, rHeader *Header, dst io.Writer)
	Decode(rHeader *Header, src io.Reader) (err error)

	ModifyHeader(h *Header)
}

func DecodePacketUDP(data []byte) (header Header, packet PacketUDP, err error) {
	header.Decode(data)

	switch header.MessageType {
	case KeepAlive:
		packet = &KeepAlivePackage{}
	case NodeHandshake:
		packet = &HandshakePackage{}
	default:
		return header, packet, ErrUnsupportedMessage
	}

	return header, packet, nil
}
