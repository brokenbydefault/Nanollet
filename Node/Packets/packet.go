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
	ErrUnsupportedMessage         = errors.New("unsupported message type")
	ErrInvalidMessageSize         = errors.New("invalid message")
	ErrDestinationLenghtNotEnough = errors.New("dst are smaller than message")
)

type PacketUDP interface {
	Encode(rHeader *Header, dst []byte) (n int, err error)
	Decode(rHeader *Header, src []byte) (err error)

	ModifyHeader(h *Header)
}

type PacketTCP interface {
	Encode(rHeader *Header, dst io.Writer) (err error)
	Decode(rHeader *Header, src io.Reader) (err error)

	ModifyHeader(h *Header)
}

func EncodePacketUDP(lHeader, rHeader *Header, packet PacketUDP) []byte {
	if lHeader == nil {
		lHeader = NewHeader()
	}

	dst := make([]byte, PackageSize)
	packet.ModifyHeader(lHeader)

	nH, err := lHeader.Encode(dst)
	if err != nil {
		return nil
	}

	nP, err := packet.Encode(rHeader, dst[nH:])
	if err != nil {
		return nil
	}

	return dst[:nH+nP]
}
