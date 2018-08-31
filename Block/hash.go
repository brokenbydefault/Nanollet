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
	if s.hash == NewBlockHash(nil) {
		copy(s.hash[:], Util.CreateHash(32, s.Encode()[1:SendHashableSize]))
	}

	return s.hash
}

func (s *ReceiveBlock) Hash() BlockHash {
	if s.hash == NewBlockHash(nil) {
		copy(s.hash[:], Util.CreateHash(32, s.Encode()[1:ReceiveHashableSize]))
	}

	return s.hash
}

func (s *OpenBlock) Hash() BlockHash {
	if s.hash == NewBlockHash(nil) {
		copy(s.hash[:], Util.CreateHash(32, s.Encode()[1:OpenHashableSize]))
	}

	return s.hash
}

func (s *ChangeBlock) Hash() BlockHash {
	if s.hash == NewBlockHash(nil) {
		copy(s.hash[:], Util.CreateHash(32, s.Encode()[1:ChangeHashableSize]))
	}

	return s.hash
}

func (u *UniversalBlock) Hash() BlockHash {
	if u.hash == NewBlockHash(nil) {
		copy(u.hash[:], Util.CreateHash(32, universalBlockFlag, u.Encode()[1:StateHashableSize]))
	}

	return u.hash
}

func (s *SendBlock) GetPrevious() (hash BlockHash) {
	return s.Previous
}

func (s *ReceiveBlock) GetPrevious() (hash BlockHash) {
	return s.Previous
}

func (s *OpenBlock) GetPrevious() (hash BlockHash) {
	return hash
}

func (s *ChangeBlock) GetPrevious() (hash BlockHash) {
	return s.Previous
}

func (u *UniversalBlock) GetPrevious() (hash BlockHash) {
	return u.Previous
}
