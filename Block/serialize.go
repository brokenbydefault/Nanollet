package Block

import (
	"github.com/brokenbydefault/Nanollet/Numbers"
	"errors"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
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

func (t BlockType) String() string {
	switch t {
	case Invalid:
		return "Invalid"
	case NotABlock:
		return "Not a block"
	case Send:
		return "Send"
	case Receive:
		return "Receive"
	case Open:
		return "Open"
	case Change:
		return "Change"
	case State:
		return "State"
	}

	return "Invalid"
}

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
	ErrInvalidSize      = errors.New("invalid size")
	ErrInvalidPoW       = errors.New("invalid PoW")
	ErrInvalidSignature = errors.New("invalid block signature")
)

func (u *UniversalBlock) MarshalBinary() ([]byte, error) {
	return u.Encode(), nil
}

func (u *UniversalBlock) UnmarshalBinary(b []byte) (error) {
	return u.Decode(b)
}

func (u *UniversalBlock) Encode() (data []byte) {
	data = make([]byte, StateExtendedSize)

	data[0] = uint8(u.GetType())
	copy(data[1:33], u.Account[:])
	copy(data[33:65], u.Previous[:])
	copy(data[65:97], u.Representative[:])
	copy(data[97:113], u.Balance.ToBytes())
	copy(data[113:145], u.Link[:])
	copy(data[145:209], u.DefaultBlock.Signature[:])
	copy(data[209:217], u.DefaultBlock.PoW[:])

	return data
}

func (u *UniversalBlock) Decode(data []byte) (err error) {
	if u == nil {
		return ErrEndBlock
	}

	blk, err := checkAndCopy(StateSize, data)
	if err != nil {
		return err
	}

	*u = UniversalBlock{
		Account:        Wallet.NewPublicKey(blk[0:32]),
		Previous:       NewBlockHash(blk[32:64]),
		Representative: Wallet.NewPublicKey(blk[64:96]),
		Balance:        Numbers.NewRawFromBytes(blk[96:112]),
		Link:           NewBlockHash(blk[112:144]),
		DefaultBlock: DefaultBlock{
			Signature: Wallet.NewSignature(blk[144:208]),
			PoW:       NewWork(blk[208:216]),
		},
	}

	hash := u.Hash()
	if !u.Account.IsValidSignature(hash[:], &u.Signature) && !IsEpoch(u) {
		return ErrInvalidSignature
	}

	return nil
}

func (u *UniversalBlock) SwitchTo(t BlockType) Transaction {
	switch t {
	case Open:
		return &OpenBlock{
			DefaultBlock: DefaultBlock{
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
				Signature: u.Signature,
				PoW:       u.PoW,
			},
			Previous: u.Previous,
			Source:   u.Link,
		}
	case Send:
		return &SendBlock{
			DefaultBlock: DefaultBlock{
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

func (s *SendBlock) MarshalBinary() ([]byte, error) {
	return s.Encode(), nil
}

func (s *SendBlock) UnmarshalBinary(b []byte) (error) {
	return s.Decode(b)
}

func (s *SendBlock) Encode() (data []byte) {
	data = make([]byte, SendExtendedSize)

	data[0] = uint8(s.GetType())
	copy(data[1:33], s.Previous[:])
	copy(data[33:65], s.Destination[:])
	copy(data[65:81], s.Balance.ToBytes())
	copy(data[81:145], s.DefaultBlock.Signature[:])
	copy(data[145:153], s.DefaultBlock.PoW[:])

	return data
}

func (s *SendBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = SendBlock{}
	}

	blk, err := checkAndCopy(SendSize, data)
	if err != nil {
		return err
	}

	*s = SendBlock{
		Previous:    NewBlockHash(blk[0:32]),
		Destination: Wallet.NewPublicKey(blk[32:64]),
		Balance:     Numbers.NewRawFromBytes(blk[64:80]),
		DefaultBlock: DefaultBlock{
			Signature: Wallet.NewSignature(blk[80:144]),
			PoW:       NewWork(Util.ReverseBytes(blk[144:152])),
		},
	}

	return err
}

func (s *SendBlock) SwitchToUniversalBlock(prevBlock *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock {
	if prevBlock == nil {
		prevBlock = &UniversalBlock{}
	}

	if amm == nil {
		amm = Numbers.NewRaw()
	}

	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Signature: s.DefaultBlock.Signature,
			PoW:       s.DefaultBlock.PoW,
		},
		Account:        prevBlock.Account,
		Previous:       s.Previous,
		Representative: prevBlock.Representative,
		Balance:        s.Balance.Subtract(amm),
		Link:           BlockHash(s.Destination),
	}
}

func (s *ReceiveBlock) MarshalBinary() ([]byte, error) {
	return s.Encode(), nil
}

func (s *ReceiveBlock) UnmarshalBinary(b []byte) (error) {
	return s.Decode(b)
}

func (s *ReceiveBlock) Encode() (data []byte) {
	data = make([]byte, ReceiveExtendedSize)

	data[0] = uint8(s.GetType())
	copy(data[1:33], s.Previous[:])
	copy(data[33:65], s.Source[:])
	copy(data[65:129], s.DefaultBlock.Signature[:])
	copy(data[129:137], s.DefaultBlock.PoW[:])

	return data
}

func (s *ReceiveBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = ReceiveBlock{}
	}

	blk, err := checkAndCopy(ReceiveSize, data)
	if err != nil {
		return err
	}

	*s = ReceiveBlock{
		Previous: NewBlockHash(blk[0:32]),
		Source:   NewBlockHash(blk[32:64]),
		DefaultBlock: DefaultBlock{
			Signature: Wallet.NewSignature(blk[64:128]),
			PoW:       NewWork(Util.ReverseBytes(blk[128:136])),
		},
	}

	return nil
}

func (s *ReceiveBlock) SwitchToUniversalBlock(prevBlock *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock {
	if prevBlock == nil {
		prevBlock = &UniversalBlock{}
	}

	if amm == nil {
		amm = Numbers.NewRaw()
	}

	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
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

func (s *OpenBlock) MarshalBinary() ([]byte, error) {
	return s.Encode(), nil
}

func (s *OpenBlock) UnmarshalBinary(b []byte) (error) {
	return s.Decode(b)
}

func (s *OpenBlock) Encode() (data []byte) {
	data = make([]byte, OpenExtendedSize)

	data[0] = uint8(s.GetType())
	copy(data[1:33], s.Source[:])
	copy(data[33:65], s.Representative[:])
	copy(data[65:97], s.Account[:])
	copy(data[97:161], s.DefaultBlock.Signature[:])
	copy(data[161:169], s.DefaultBlock.PoW[:])

	return data
}

func (s *OpenBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = OpenBlock{}
	}

	blk, err := checkAndCopy(OpenSize, data)
	if err != nil {
		return err
	}

	*s = OpenBlock{
		Source:         NewBlockHash(blk[0:32]),
		Representative: Wallet.NewPublicKey(blk[32:64]),
		Account:        Wallet.NewPublicKey(blk[64:96]),
		DefaultBlock: DefaultBlock{
			Signature: Wallet.NewSignature(blk[96:160]),
			PoW:       NewWork(Util.ReverseBytes(blk[160:168])),
		},
	}

	return nil
}

func (s *OpenBlock) SwitchToUniversalBlock(_ *UniversalBlock, amm *Numbers.RawAmount) *UniversalBlock {
	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			PoW:       s.DefaultBlock.PoW,
			Signature: s.DefaultBlock.Signature,
		},
		Account:        s.Account,
		Representative: s.Representative,
		Balance:        amm,
		Link:           s.Source,
	}
}

func (s *ChangeBlock) MarshalBinary() ([]byte, error) {
	return s.Encode(), nil
}

func (s *ChangeBlock) UnmarshalBinary(b []byte) (error) {
	return s.Decode(b)
}

func (s *ChangeBlock) Encode() (data []byte) {
	data = make([]byte, ChangeExtendedSize)

	data[0] = uint8(s.GetType())
	copy(data[1:33], s.Previous[:])
	copy(data[33:65], s.Representative[:])
	copy(data[65:129], s.DefaultBlock.Signature[:])
	copy(data[129:137], s.DefaultBlock.PoW[:])

	return data
}

func (s *ChangeBlock) Decode(data []byte) (err error) {
	if s == nil {
		*s = ChangeBlock{}
	}

	blk, err := checkAndCopy(ChangeSize, data)
	if err != nil {
		return err
	}

	*s = ChangeBlock{
		Previous:       NewBlockHash(blk[0:32]),
		Representative: Wallet.NewPublicKey(blk[32:64]),
		DefaultBlock: DefaultBlock{
			Signature: Wallet.NewSignature(blk[64:128]),
			PoW:       NewWork(Util.ReverseBytes(blk[128:136])),
		},
	}

	return nil
}

func (s *ChangeBlock) SwitchToUniversalBlock(prevBlock *UniversalBlock, _ *Numbers.RawAmount) *UniversalBlock {
	if prevBlock == nil {
		prevBlock = &UniversalBlock{}
	}

	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			PoW:       s.DefaultBlock.PoW,
			Signature: s.DefaultBlock.Signature,
		},
		Account:        prevBlock.Account,
		Previous:       s.Previous,
		Representative: s.Representative,
		Balance:        prevBlock.Balance,
	}
}

func checkAndCopy(expectedSize int, data []byte) ([]byte, error) {
	switch len(data) {
	case expectedSize:
		// no-op
	case expectedSize + 1:
		data = data[1:]
	default:
		return nil, ErrInvalidSize
	}

	return data, nil
}
