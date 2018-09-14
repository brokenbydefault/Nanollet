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
	Encode(dst []byte) (n int, err error)
	Decode(rHeader *Header, src []byte) (err error)

	ModifyHeader(h *Header)
}

type PacketTCP interface {
	Encode(dst io.Writer) (err error)
	Decode(rHeader *Header, src io.Reader) (err error)

	ModifyHeader(h *Header)
}

func EncodePacketUDP(lHeader Header, packet PacketUDP) []byte {
	dst := make([]byte, PackageSize)
	packet.ModifyHeader(&lHeader)

	nH, err := lHeader.Encode(dst)
	if err != nil {
		return nil
	}

	nP, err := packet.Encode(dst[nH:])
	if err != nil {
		return nil
	}

	return dst[:nH+nP]
}

func EncodePacketTCP(lHeader Header, packet PacketTCP, writer io.Writer) (err error) {
	packet.ModifyHeader(&lHeader)

	_, err = lHeader.Write(writer)
	if err != nil {
		return err
	}

	err = packet.Encode(writer)
	if err != nil {
		return err
	}

	return nil
}

func DecodeResponsePacketTCP(messageType MessageType) PacketTCP {
	switch messageType {
	case FrontierReq:
		return &FrontierReqPackageResponse{}
	case BulkPull:
		return &BulkPullPackageResponse{}
	case BulkPullBlocks:
		panic("not supported")
	case BulkPullAccount:
		return &BulkPullAccountPackageResponse{}
	default:
		panic("not supported")
	}
}
