package Block

import (
	"encoding/json"
	"github.com/brokenbydefault/Nanollet/Util"
)

type BlockHash []byte

var universalblockflag = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06}

func (s *SendBlock) Hash() BlockHash {
	destination, _ := s.Destination.GetPublicKey()
	return Util.CreateHash(32, s.Previous, destination, s.Balance.ToBytes())
}

func (s *ReceiveBlock) Hash() BlockHash {
	return Util.CreateHash(32, s.Previous, s.Source)
}

func (s *OpenBlock) Hash() BlockHash {
	account, _ := s.Account.GetPublicKey()
	representative, _ := s.Representative.GetPublicKey()
	return Util.CreateHash(32, s.Source, representative, account)
}

func (s *ChangeBlock) Hash() BlockHash {
	representative, _ := s.Representative.GetPublicKey()
	return Util.CreateHash(32, s.Previous, representative)
}

func (u *UniversalBlock) Hash() BlockHash {
	var link [32]byte
	copy(link[:], u.Link)

	var previous [32]byte
	copy(previous[:], u.Previous)

	account, _ := u.Account.GetPublicKey()
	representative, _ := u.Representative.GetPublicKey()

	return Util.CreateHash(32, universalblockflag, account, previous[:], representative, u.Balance.ToBytes(), link[:])
}

func (d *BlockHash) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return
	}

	v, err := Util.UnsafeHexDecode(str)
	if err != nil {
		return
	}

	*d = v
	return
}

func (d BlockHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(Util.UnsafeHexEncode(d))
}
