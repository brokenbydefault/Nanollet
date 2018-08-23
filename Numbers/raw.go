package Numbers

import (
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
	i, ok := new(big.Int).SetString(s, 10)
	r := &RawAmount{i}

	if !ok {
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

	return NewRawFromBytes(b), nil
}

// NewRawFromBytes creates an RawAmount from byte-array. It returns an non-nil error if
// the string is invalid and nil successful.
func NewRawFromBytes(b []byte) *RawAmount {
	return &RawAmount{
		bigint: new(big.Int).SetBytes(b),
	}
}

func (a *RawAmount) Copy(src []byte) (i int) {
	if a == nil {
		*a = RawAmount{}
	}

	*a = *NewRawFromBytes(src)
	return len(src)
}

// ToString transforms the RawAmount to string, which can be printable.
func (a *RawAmount) ToString() string {
	if a == nil {
		return ""
	}

	return new(big.Int).Set(a.bigint).String()
}

// ToPaddedHex transforms the RawAmount to hexadecimal string with 16 byte,
// left zero-padded. It can be used in RPC.
func (a *RawAmount) ToHex() string {
	if a == nil {
		return ""
	}

	return Util.UnsafeHexEncode(a.ToBytes())
}

// ToBytes transforms the RawAmount to 16 byte, left zero-padded. It can be used in
// block signature and in RPC.
func (a *RawAmount) ToBytes() []byte {
	if a == nil {
		return nil
	}

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
