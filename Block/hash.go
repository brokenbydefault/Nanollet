package Block

import (
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/json"
)

type BlockHash []byte

func (s *SendBlock) Hash() BlockHash {
	return Util.CreateHash(32, s.Previous, s.Destination, s.Balance.ToBytes())
}

func (s *ReceiveBlock) Hash() BlockHash {
	return Util.CreateHash(32, s.Previous, s.Source)
}

func (s *OpenBlock) Hash() BlockHash {
	return Util.CreateHash(32, s.Source, s.Representative, s.Account)
}

func (s *ChangeBlock) Hash() BlockHash {
	return Util.CreateHash(32, s.Previous, s.Representative)
}

// It can change in the future, consider "Destination" as "Target".
func (u *UniversalBlock) Hash() BlockHash {
	return Util.CreateHash(32, u.Previous, u.Destination, u.Representative, u.Balance.ToBytes(), u.Account)
}

func (u *UniversalBlock) HashAsSend() BlockHash {
	return Util.CreateHash(32, u.Previous, u.Destination, u.Balance.ToBytes())
}

func (u *UniversalBlock) HashAsReceive() BlockHash {
	return Util.CreateHash(32, u.Previous, u.Source)
}

func (u *UniversalBlock) HashAsOpen() BlockHash {
	return Util.CreateHash(32, u.Source, u.Representative, u.Account)
}

func (u *UniversalBlock) HashAsChange() BlockHash {
	return Util.CreateHash(32, u.Previous, u.Representative)
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
