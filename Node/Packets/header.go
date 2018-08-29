package Packets

import (
	"errors"
	"github.com/brokenbydefault/Nanollet/Block"
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
}

func NewHeader() *Header {
	return &Header{
		MagicNumber:   82,
		NetworkType:   Beta, // Important
		VersionMax:    13,
		VersionUsing:  13,
		VersionMin:    13,
		MessageType:   Invalid,
		ExtensionType: 0,
	}
}

func (h *Header) Encode() (data []byte) {
	if h == nil {
		return
	}

	data = append(data, []byte{
		byte(h.MagicNumber),
		byte(h.NetworkType),
		byte(h.VersionMax),
		byte(h.VersionUsing),
		byte(h.VersionMin),
		byte(h.MessageType),
		byte(h.ExtensionType),
		byte(h.ExtensionType >> 8),
	}...)

	return data
}

func (h *Header) Decode(data []byte) (err error) {
	if len(data) < HeaderSize {
		return ErrInvalidHeaderSize
	}

	h.MagicNumber = data[0]
	h.NetworkType = NetworkType(data[1])
	h.VersionMax = data[2]
	h.VersionUsing = data[3]
	h.VersionMin = data[4]
	h.MessageType = MessageType(data[5])
	h.ExtensionType = ExtensionType(uint16(data[6]) | uint16(data[7])<<8)

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
