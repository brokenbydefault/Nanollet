package Block

import (
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/json"
)


//@TODO UniversalBlock is not finished!
/**
func (s *UniversalBlock) Serialize() ([]byte, error) {

	sb := SerializedUniversalBlock{
		Previous:               DOM.SecureHexEncode(s.Previous),
		Representative:			string(s.Representative),
		Balance:				s.Balance.ToHex(),
		Amount:					s.Amount.ToHex(),
		Link:            		s.Link,
		SerializedDefaultBlock: newSerializeDefault("send", s.DefaultBlock),
	}

	return json.Marshal(sb), err
}
**/

func (s *SendBlock) Serialize() ([]byte, error) {
	sb := SerializedSendBlock{
		Previous:               Util.SecureHexEncode(s.Previous),
		Destination:            string(s.Destination.CreateAddress()),
		Balance:                s.Balance.ToHex(),
		SerializedDefaultBlock: newSerializeDefault("send", s.DefaultBlock),
	}

	return json.Marshal(sb)
}

func (s *ReceiveBlock) Serialize() ([]byte, error) {
	sb := SerializedReceiveBlock{
		Previous:               Util.SecureHexEncode(s.Previous),
		Source:                 Util.SecureHexEncode(s.Source),
		SerializedDefaultBlock: newSerializeDefault("receive", s.DefaultBlock),
	}

	return json.Marshal(sb)
}

func (s *OpenBlock) Serialize() ([]byte, error) {
	sb := SerializedOpenBlock{
		Source:                 Util.SecureHexEncode(s.Source),
		Representative:         string(s.Representative.CreateAddress()),
		Account:                string(s.Account.CreateAddress()),
		SerializedDefaultBlock: newSerializeDefault("open", s.DefaultBlock),
	}

	return json.Marshal(sb)
}

func (s *ChangeBlock) Serialize() ([]byte, error) {
	sb := SerializedChangeBlock{
		Previous:               Util.SecureHexEncode(s.Previous),
		Representative:         string(s.Representative.CreateAddress()),
		SerializedDefaultBlock: newSerializeDefault("change", s.DefaultBlock),
	}

	return json.Marshal(sb)
}

func newSerializeDefault(typ string, v DefaultBlock) SerializedDefaultBlock {
	return SerializedDefaultBlock{
		Type:      typ,
		Work:      Util.SecureHexEncode(v.Work),
		Signature: Util.SecureHexEncode(v.Signature),
	}
}
