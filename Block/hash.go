package Block

import (
	"github.com/brokenbydefault/Nanollet/Util"
	"golang.org/x/crypto/blake2b"
)

type BlockHash [blake2b.Size256]byte

func NewBlockHash(b []byte) (hash BlockHash) {
	copy(hash[:], b)
	return hash
}

var universalBlockFlag = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06}

func (s *SendBlock) Hash() BlockHash {
	return NewBlockHash(Util.CreateHash(32, s.Encode()[1:SendHashableSize]))
}

func (s *ReceiveBlock) Hash() BlockHash {
	return NewBlockHash(Util.CreateHash(32, s.Encode()[1:ReceiveHashableSize]))
}

func (s *OpenBlock) Hash() BlockHash {
	return NewBlockHash(Util.CreateHash(32, s.Encode()[1:OpenHashableSize]))
}

func (s *ChangeBlock) Hash() BlockHash {
	return NewBlockHash(Util.CreateHash(32, s.Encode()[1:ChangeHashableSize]))
}

func (u *UniversalBlock) Hash() BlockHash {
	return NewBlockHash(Util.CreateHash(32, universalBlockFlag, u.Encode()[1:StateHashableSize]))
}

func (s *SendBlock) GetPrevious() (hash BlockHash) {
	return s.Previous
}

func (s *ReceiveBlock) GetPrevious() (hash BlockHash) {
	return s.Previous
}

func (s *OpenBlock) GetPrevious() (hash BlockHash) {
	return NewBlockHash(nil)
}

func (s *ChangeBlock) GetPrevious() (hash BlockHash) {
	return s.Previous
}

func (u *UniversalBlock) GetPrevious() (hash BlockHash) {
	if Util.IsEmpty(u.Previous[:]) {
		return NewBlockHash(nil)
	}

	return u.Previous
}
