package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"errors"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

type BlockType byte

const (
	Invalid   BlockType = iota
	NotABlock
	Send
	Receive
	Open
	Change
	State
)

const (
	SendSize            = 152
	SendExtendedSize    = SendSize + 1
	SendHashableSize    = SendExtendedSize - 64 - 8
	ReceiveSize         = 136
	ReceiveExtendedSize = ReceiveSize + 1
	ReceiveHashableSize = ReceiveExtendedSize - 64 - 8
	OpenSize            = 168
	OpenExtendedSize    = OpenSize + 1
	OpenHashableSize    = OpenExtendedSize - 64 - 8
	ChangeSize          = 136
	ChangeExtendedSize  = ChangeSize + 1
	ChangeHashableSize  = ChangeExtendedSize - 64 - 8
	StateSize           = 216
	StateExtendedSize   = StateSize + 1
	StateHashableSize   = StateExtendedSize - 64 - 8
)

var (
	ErrInvalidSize = errors.New("invalid size")
)

func (u *UniversalBlock) Encode() (data []byte) {
	data = make([]byte, StateExtendedSize)

	// Enhanced node send the original type of the block in the first 4 bits, if it's not a original Universal Block.
	// Keep in mind that it's not compatible with the reference node (by Raiblocks Team). You must convert the type of
	// the block to the original one in that case.
	// If the first 4 bytes are zeroes, the SubType meaning "Invalid", which means that the type of the block
	// is originally state-block.
	data[0] = uint8(u.DefaultBlock.SubType)<<4 | uint8(State)

	copy(data[1:33], u.Account)
	copy(data[33:65], u.Previous)
	copy(data[65:97], u.Representative)
	copy(data[97:113], u.Balance.ToBytes())
	copy(data[113:145], u.Link)
	copy(data[145:209], u.DefaultBlock.Signature)
	copy(data[209:217], u.DefaultBlock.PoW)

	return data
}

func (u *UniversalBlock) Decode(data []byte) (err error) {
	if u == nil {
		*u = UniversalBlock{}
	}

	i := 0

	switch len(data) {
	case StateSize:
		u.DefaultBlock.Type = State
	case StateExtendedSize:
		u.DefaultBlock.Type = State
		u.DefaultBlock.SubType = BlockType(data[0] & 0xF0)
		i += 1
	default:
		return ErrInvalidSize
	}

	i = copy(u.Account, data[i:i+32])
	i = copy(u.Previous, data[i:i+32])
	i = copy(u.Representative, data[i:i+32])
	i, err = u.Balance.Copy(data[i:i+16])
	i = copy(u.Link, data[i:i+32])
	i = copy(u.DefaultBlock.Signature, data[i:i+64])
	i = copy(u.DefaultBlock.PoW, data[i:i+8])

	return err
}

func (u *UniversalBlock) SwitchTo(t BlockType) Transaction {
	switch t {
	case Open:
		return &OpenBlock{
			DefaultBlock: DefaultBlock{
				Type:      Open,
				SubType:   Open,
				Signature: u.Signature,
				PoW:       u.PoW,
			},
			Representative: u.Representative,
			Account:        u.Account,
			Source:         u.Link,
		}
	case Receive:
		return &ReceiveBlock{
			DefaultBlock: DefaultBlock{
				Type:      Receive,
				SubType:   Receive,
				Signature: u.Signature,
				PoW:       u.PoW,
			},
			Previous: u.Previous,
			Source:   u.Link,
		}
	case Send:
		return &SendBlock{
			DefaultBlock: DefaultBlock{
				Type:      Send,
				SubType:   Send,
				Signature: u.Signature,
				PoW:       u.PoW,
			},
			Previous:    u.Previous,
			Destination: Wallet.PublicKey(u.Link),
			Balance:     u.Balance,
		}
	case Change:
		return &ChangeBlock{
			DefaultBlock: DefaultBlock{
				Type:      Change,
				SubType:   Change,
				Signature: u.Signature,
				PoW:       u.PoW,
			},
			Previous:       u.Previous,
			Representative: u.Representative,
		}
	}

	return u
}

func (u *UniversalBlock) SwitchToUniversalBlock(_ *UniversalBlock, _ *Numbers.RawAmount) *UniversalBlock {
	return u
}

func (s *SendBlock) Encode() (data []byte) {
	data = make([]byte, 0, SendExtendedSize)

	data[0] = uint8(s.DefaultBlock.Type)
	copy(data[1:33], s.Previous)
	copy(data[33:65], s.Destination)
	copy(data[65:81], s.Balance.ToBytes())
	copy(data[81:145], s.DefaultBlock.Signature)
	copy(data[145:153], s.DefaultBlock.PoW)

	return data
}

func (s *SendBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = SendBlock{}
	}

	i := len(data) - SendSize
	if i > 1 {
		return ErrInvalidSize
	}

	s.DefaultBlock.Type = Send
	s.DefaultBlock.SubType = Send

	i = copy(s.Previous, data[i:i+32])
	i = copy(s.Destination, data[i:i+32])
	i, err = s.Balance.Copy(data[i:i+16])
	i = copy(s.DefaultBlock.Signature, data[i:i+64])
	i = copy(s.DefaultBlock.PoW, data[i:i+8])

	return err
}

func (s *SendBlock) SwitchToUniversalBlock(prevBlock *UniversalBlock, _ *Numbers.RawAmount) *UniversalBlock {
	if prevBlock == nil {
		prevBlock = &UniversalBlock{
			Account:        make([]byte, 32),
			Representative: make([]byte, 32),
		}
	}

	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:      State,
			SubType:   Send,
			Signature: s.DefaultBlock.Signature,
			PoW:       s.DefaultBlock.PoW,
		},
		Account:        prevBlock.Account,
		Previous:       s.Previous,
		Representative: prevBlock.Representative,
		Balance:        s.Balance,
		Link:           BlockHash(s.Destination),
	}
}

func (s *ReceiveBlock) Encode() (data []byte) {
	data = make([]byte, 0, ReceiveExtendedSize)

	data[0] = uint8(s.DefaultBlock.Type)
	copy(data[1:33], s.Previous)
	copy(data[33:65], s.Source)
	copy(data[65:129], s.DefaultBlock.Signature)
	copy(data[129:137], s.DefaultBlock.PoW)

	return data
}

func (s *ReceiveBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = ReceiveBlock{}
	}

	i := len(data) - ReceiveSize
	if i > 1 {
		return ErrInvalidSize
	}

	s.DefaultBlock.Type = Receive
	s.DefaultBlock.SubType = Receive

	i = copy(s.Previous, data[i:i+32])
	i = copy(s.Source, data[i:i+32])
	i = copy(s.DefaultBlock.Signature, data[i:i+64])
	i = copy(s.DefaultBlock.PoW, data[i:i+8])

	return err
}

func (s *ReceiveBlock) SwitchToUniversalBlock(prevBlock *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock {
	if prevBlock == nil {
		prevBlock = &UniversalBlock{
			Account:        make([]byte, 32),
			Representative: make([]byte, 32),
			Balance:        Numbers.NewRaw(),
		}
	}

	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:      State,
			SubType:   Receive,
			PoW:       s.DefaultBlock.PoW,
			Signature: s.DefaultBlock.Signature,
		},
		Account:        prevBlock.Account,
		Previous:       s.Previous,
		Representative: prevBlock.Representative,
		Balance:        prevBlock.Balance.Add(amm),
		Link:           s.Source,
	}
}

func (s *OpenBlock) Encode() (data []byte) {
	data = make([]byte, 0, OpenExtendedSize)

	data[0] = uint8(s.DefaultBlock.Type)
	copy(data[1:33], s.Source)
	copy(data[33:65], s.Representative)
	copy(data[65:97], s.Account)
	copy(data[97:161], s.DefaultBlock.Signature)
	copy(data[161:169], s.DefaultBlock.PoW)

	return data
}

func (s *OpenBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = OpenBlock{}
	}

	i := len(data) - OpenSize
	if i > 1 {
		return ErrInvalidSize
	}

	s.DefaultBlock.Type = Open
	s.DefaultBlock.SubType = Open

	i = copy(s.Source, data[i:i+32])
	i = copy(s.Representative, data[i:i+32])
	i = copy(s.Account, data[i:i+32])
	i = copy(s.DefaultBlock.Signature, data[i:i+64])
	i = copy(s.DefaultBlock.PoW, data[i:i+8])

	return err
}

func (s *OpenBlock) SwitchToUniversalBlock(_ *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock {
	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:      State,
			SubType:   Open,
			PoW:       s.DefaultBlock.PoW,
			Signature: s.DefaultBlock.Signature,
		},
		Account:        s.Account,
		Previous:       make([]byte, 32),
		Representative: s.Representative,
		Balance:        amm,
		Link:           s.Source,
	}
}

func (s *ChangeBlock) Encode() (data []byte) {
	data = make([]byte, 0, ChangeExtendedSize)

	data[0] = uint8(s.DefaultBlock.Type)
	copy(data[1:33], s.Previous)
	copy(data[33:65], s.Representative)
	copy(data[65:129], s.DefaultBlock.Signature)
	copy(data[129:137], s.DefaultBlock.PoW)

	return data
}

func (s *ChangeBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = ChangeBlock{}
	}

	i := len(data) - ChangeSize
	if i > 1 {
		return ErrInvalidSize
	}

	s.DefaultBlock.Type = Receive
	s.DefaultBlock.SubType = Receive

	i = copy(s.Previous, data[i:i+32])
	i = copy(s.Representative, data[i:i+32])
	i = copy(s.DefaultBlock.Signature, data[i:i+64])
	i = copy(s.DefaultBlock.PoW, data[i:i+8])

	return err
}

func (s *ChangeBlock) SwitchToUniversalBlock(prevBlock *UniversalBlock, _ *Numbers.RawAmount) *UniversalBlock {
	if prevBlock == nil {
		prevBlock = &UniversalBlock{
			Account: make([]byte, 32),
			Balance: Numbers.NewRaw(),
		}
	}

	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:      State,
			SubType:   Change,
			PoW:       s.DefaultBlock.PoW,
			Signature: s.DefaultBlock.Signature,
		},
		Account:        prevBlock.Account,
		Previous:       s.Previous,
		Representative: s.Representative,
		Balance:        prevBlock.Balance,
		Link:           make([]byte, 32),
	}
}
