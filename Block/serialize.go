package Block

import (
	"encoding/json"
	"strings"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"bytes"
)

type BlockType string

const (
	None    = ""
	Open    = "open"
	Receive = "receive"
	Send    = "send"
	Change  = "change"
	State   = "state"
)

func (u *UniversalBlock) Serialize() ([]byte, error) {
	subtype := u.SubType
	u.SubType = None
	blk, err := json.Marshal(u)
	u.SubType = subtype

	return blk, err
}

func (u *UniversalBlock) SwitchTo(t BlockType) BlockTransaction {
	destination, source := u.GetTarget()

	switch t {
	case Open:
		return &OpenBlock{
			Representative: u.Representative,
			Account:        u.Account,
			Source:         source,
		}
	case Receive:
		return &ReceiveBlock{
			Previous: u.Previous,
			Source:   source,
		}
	case Send:
		return &SendBlock{
			Previous:    u.Previous,
			Destination: destination,
			Balance:     u.Balance,
		}
	case Change:
		return &ChangeBlock{
			Previous:       u.Previous,
			Representative: u.Representative,
		}
	}

	return u
}

func (u *UniversalBlock) SwitchToUniversalBlock() *UniversalBlock {
	return u
}

func (s *SendBlock) Serialize() ([]byte, error) {
	jsn, err := json.Marshal(s)

	// RPC INCONSISTENCY: All calls uses `balance` as numeric string, but the
	// `Block` uses hex, except state-block, which uses numeric-string too.
	jsn = bytes.Replace(jsn, []byte(`"balance":"`+s.Balance.ToString()+`"`), []byte(`"balance":"`+s.Balance.ToHex()+`"`), 1)
	//////////////////////////////////////////////////////////////////////

	return jsn, err
}

func (s *SendBlock) SwitchToUniversalBlock() *UniversalBlock {
	destpk, _ := s.Destination.GetPublicKey()
	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    Send,
			SubType: "",
			PoW:     s.PoW,
		},
		Account:        "",
		Previous:       s.Previous,
		Representative: "",
		Balance:        s.Balance,
		Link:           BlockHash(destpk),
	}
}

func (s *ReceiveBlock) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ReceiveBlock) SwitchToUniversalBlock() *UniversalBlock {
	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    Receive,
			SubType: "",
			PoW:     s.PoW,
		},
		Account:        "",
		Previous:       s.Previous,
		Representative: "",
		Balance:        Numbers.NewRaw(),
		Link:           s.Source,
	}
}

func (s *OpenBlock) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s *OpenBlock) SwitchToUniversalBlock() *UniversalBlock {
	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    Open,
			SubType: None,
			PoW:     s.PoW,
		},
		Account:        s.Account,
		Previous:       make([]byte, 32),
		Representative: s.Representative,
		Balance:        Numbers.NewRaw(),
		Link:           s.Source,
	}
}

func (s *ChangeBlock) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ChangeBlock) SwitchToUniversalBlock() *UniversalBlock {
	return &UniversalBlock{
		DefaultBlock: DefaultBlock{
			Type:    Change,
			SubType: None,
			PoW:     s.PoW,
		},
		Account:        "",
		Previous:       s.Previous,
		Representative: s.Representative,
		Balance:        Numbers.NewRaw(),
		Link:           make([]byte, 32),
	}
}

func (d *BlockType) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return
	}

	*d = BlockType(strings.ToLower(str))
	return
}

func (d BlockType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}
