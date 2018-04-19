package Numbers

import (
	"encoding/json"
	"errors"
	"github.com/brokenbydefault/Nanollet/Util"
	"math/big"
)

type RawAmount struct {
	bigint *big.Int
}

var max = new(big.Int).SetBytes([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
var min = new(big.Int).SetInt64(0)

func NewRaw() *RawAmount {
	return &RawAmount{new(big.Int).SetInt64(0)}
}

// NewRawFromString creates an RawAmount from numeric string. It returns an error if
// the string is invalid and nil successful.
func NewRawFromString(s string) (*RawAmount, error) {
	i, success := new(big.Int).SetString(s, 10)
	r := &RawAmount{i}

	if !success || !r.IsValid() {
		return nil, errors.New("invalid string")
	}

	return r, nil
}

// NewRawFromHex creates an RawAmount from hexadecimal string. It returns an non-nil error if
// the value is invalid.
func NewRawFromHex(h string) (*RawAmount, error) {
	b, err := Util.UnsafeHexDecode(h)
	if err != nil {
		return nil, errors.New("invalid hex")
	}

	return NewRawFromBytes(b)
}

// NewRawFromBytes creates an RawAmount from byte-array. It returns an non-nil error if
// the string is invalid and nil successful.
func NewRawFromBytes(b []byte) (*RawAmount, error) {
	i := new(big.Int).SetBytes(b)

	r := &RawAmount{i}
	if !r.IsValid() {
		return nil, errors.New("invalid string")
	}

	return r, nil
}

// ToString transforms the RawAmount to string, which can be printable.
func (a *RawAmount) ToString() string {
	return new(big.Int).Set(a.bigint).String()
}

// ToPaddedHex transforms the RawAmount to hexadecimal string with 16 byte,
// left zero-padded. It can be used in RPC.
func (a *RawAmount) ToHex() string {
	return Util.UnsafeHexEncode(a.ToBytes())
}

// ToBytes transforms the RawAmount to 16 byte, left zero-padded. It can be used in
// block signature and in RPC.
func (a *RawAmount) ToBytes() []byte {
	bi := a.bigint.Bytes()

	// If the value is larger than Uint128 we return,
	// because it already have more than 16 bytes
	// however it's a invalid value for the Nano.
	if l := len(bi); l >= 16 {
		return bi
	}

	b := make([]byte, 16)
	copy(b[16-len(bi):], bi)

	return b
}

func (d *RawAmount) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.ToString())
}

func (d *RawAmount) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return
	}

	v, err := NewRawFromString(str)
	if err != nil {
		return
	}

	*d = *v
	return
}
