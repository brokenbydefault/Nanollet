package Numbers

import (
	"errors"
	"github.com/brokenbydefault/Nanollet/Util"
	"math/big"
	"io"
)

type RawAmount struct {
	bigint *big.Int
}

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidInput  = errors.New("impossible to convert the input to RawNumbers")
)

var (
	max = new(big.Int).SetBytes([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	min = new(big.Int).SetInt64(0)
)

// NewRaw creates an RawAmount with lowest valid value.
func NewRaw() *RawAmount {
	return NewMin()
}

// NewMax is a wrapper to create a new amount with maximum possible amount.
func NewMax() *RawAmount {
	return &RawAmount{
		bigint: max,
	}
}

// NewMax is a wrapper to create a new amount with minimum possible amount.
func NewMin() *RawAmount {
	return &RawAmount{
		bigint: min,
	}
}

// NewRawFromString creates an RawAmount from numeric string. It returns an error if
// the string is invalid and nil successful.
func NewRawFromString(s string) (*RawAmount, error) {
	i, ok := new(big.Int).SetString(s, 10)
	r := &RawAmount{i}

	if !ok || !r.IsValid() {
		return nil, ErrInvalidInput
	}

	return r, nil
}

// NewRawFromHex creates an RawAmount from hexadecimal string. It returns an non-nil error if
// the value is invalid.
func NewRawFromHex(h string) (*RawAmount, error) {
	b, err := Util.UnsafeHexDecode(h)
	if err != nil {
		return nil, ErrInvalidInput
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
	b := make([]byte, 16)

	if a == nil {
		return b
	}

	bi := a.bigint.Bytes()

	offset := 16 - len(bi)
	if offset < 0 {
		offset = 0
	}

	copy(b[offset:], bi)

	return b
}

func (a *RawAmount) MarshalBinary() (data []byte, err error) {
	return a.ToBytes(), nil
}

func (a *RawAmount) UnmarshalBinary(data []byte) error {
	*a = *NewRawFromBytes(data)

	if !a.IsValid() {
		return ErrInvalidAmount
	}

	return nil
}

func (a *RawAmount) Read(reader io.Reader) (err error) {
	b := make([]byte, 16)
	if n, err := reader.Read(b); n != 16 || err != nil {
		return ErrInvalidInput
	}

	a.bigint = new(big.Int).SetBytes(b)

	return
}

func (a *RawAmount) Write(writer io.Writer) (err error) {
	if n, err := writer.Write(a.ToBytes()); n != 16 || err != nil {
		return ErrInvalidInput
	}

	return
}
