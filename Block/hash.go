package Block

import (
	"github.com/brokenbydefault/Nanollet/Util"
)

type BlockHash []byte

var universalBlockFlag = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06}

func (s *SendBlock) Hash() BlockHash {
	if s.hash == nil {
		return Util.CreateHash(32, s.Encode()[1:SendHashableSize])
	}

	return s.hash
}

func (s *ReceiveBlock) Hash() BlockHash {
	if s.hash == nil {
		return Util.CreateHash(32, s.Encode()[1:ReceiveHashableSize])
	}

	return s.hash
}

func (s *OpenBlock) Hash() BlockHash {
	if s.hash == nil {
		return Util.CreateHash(32, s.Encode()[1:OpenHashableSize])
	}

	return s.hash
}

func (s *ChangeBlock) Hash() BlockHash {
	if s.hash == nil {
		return Util.CreateHash(32, s.Encode()[1:ChangeHashableSize])
	}

	return s.hash
}

func (u *UniversalBlock) Hash() BlockHash {
	if u.hash == nil {
		return Util.CreateHash(32, universalBlockFlag, u.Encode()[1:StateHashableSize])
	}

	return u.hash
}
