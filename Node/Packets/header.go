package Packets

import (
	"errors"
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
)

// NetworkType is a one-byte which defines the network which is connected to, such as Live or Test.
type NetworkType byte

const (
	Test NetworkType = iota + 65
	Beta
	Live
)

// MessageType is one-byte which says what type of message was received or sent, such as Publish or ConfirmReq.
type MessageType byte

const (
	Invalid         MessageType = iota
	NotType
	KeepAlive
	Publish
	ConfirmReq
	ConfirmACK
	BulkPull
	BulkPush
	FrontierReq
	BulkPullBlocks
	NodeHandshake
	BulkPullAccount
)

// Extensions is originally one uint16, in little-endian, which tells some specific behaviors and limitations of the
// node and block usage.
type ExtensionType uint16

const (
	ExtendedNode = ExtensionType(0x8000)
)

func (t ExtensionType) Is(expected ExtensionType) (ok bool) {
	return t&expected == expected
}

func (t *ExtensionType) Add(extension ExtensionType) {
	*t |= extension
}

func (t ExtensionType) GetBlockType() Block.BlockType {
	return Block.BlockType((t >> 8) & 0x0F)
}

var (
	ErrInvalidHeaderParameters = errors.New("invalid header parameters")
	ErrInvalidHeaderSize       = errors.New("invalid header size")
)

// HeaderSize represent the amount bytes used in the header
const HeaderSize = 8

type Header struct {
	MagicNumber   byte
	NetworkType   NetworkType
	VersionMax    byte
	VersionUsing  byte
	VersionMin    byte
	MessageType   MessageType
	ExtensionType ExtensionType

	removeHeader bool
}

func NewHeader() *Header {
	return &Header{
		MagicNumber:   82,
		NetworkType:   Live,
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Invalid,
		ExtensionType: 0,
		removeHeader:  false,
	}
}

func (h *Header) SetRemoveHeader(opt bool) {
	h.removeHeader = opt
}

func (h *Header) Encode(dst []byte) (n int, err error) {
	if h == nil {
		return
	}

	if len(dst) < HeaderSize {
		return 0, ErrDestinationLenghtNotEnough
	}

	if h.removeHeader {
		return 0, nil
	}

	n = copy(dst, []byte{
		byte(h.MagicNumber),
		byte(h.NetworkType),
		byte(h.VersionMax),
		byte(h.VersionUsing),
		byte(h.VersionMin),
		byte(h.MessageType),
		byte(h.ExtensionType),
		byte(h.ExtensionType >> 8),
	})

	return n, err
}

func (h *Header) Decode(src []byte) (err error) {
	if len(src) < HeaderSize {
		return ErrInvalidHeaderSize
	}

	h.MagicNumber = src[0]
	h.NetworkType = NetworkType(src[1])
	h.VersionMax = src[2]
	h.VersionUsing = src[3]
	h.VersionMin = src[4]
	h.MessageType = MessageType(src[5])
	h.ExtensionType = ExtensionType(uint16(src[6]) | uint16(src[7])<<8)

	if h.MagicNumber != []byte("R")[0] {
		return ErrInvalidHeaderParameters
	}

	//@TODO Verify network
	/**
	if h.NetworkType {
	}
	**/

	//@TODO Verify version
	/**
	if h.VersionUsing {
	}
	**/

	if h.MessageType >= BulkPullAccount {
		return ErrInvalidHeaderParameters
	}

	return nil
}

func (h *Header) Read(src io.Reader) (n int, err error) {
	b := make([]byte, HeaderSize)

	n, err = src.Read(b)
	if err != nil {
		return n, err
	}

	err = h.Decode(b)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (h *Header) Write(dst io.Writer) (n int, err error) {
	if h.removeHeader {
		return 0, nil
	}

	b := make([]byte, HeaderSize)
	if n, err = h.Encode(b); err != nil {
		return n, err
	}

	if n, err = dst.Write(b); err != nil {
		return n, err
	}

	return HeaderSize, nil
}
