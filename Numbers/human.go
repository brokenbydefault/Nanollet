package Numbers

import (
	"errors"
	"math/big"
	"strings"
)

type UnitBase int
type HumanAmount struct {
	whole   string
	decimal string
	base    UnitBase
}

const (
	MicroXRB UnitBase = 18 + (3 * iota)
	MiliXRB
	XRB
	KiloXRB
	MegaXRB
	GigaXRB
	RAW UnitBase = 0
)

func NewHumanFromString(n string, base UnitBase) *HumanAmount {
	values := make([]string, 2)

	split := strings.Split(n, ".")
	if len(split) > 2 {
		panic("value is invalid, having more than one dot")
	}

	copy(values, split)

	r := HumanAmount{
		whole:   values[0],
		decimal: values[1],
		base:    base,
	}

	return &r
}

func NewHumanFromRaw(n *RawAmount) *HumanAmount {
	return NewHumanFromString(n.ToString(), 0)
}

func (n *HumanAmount) ConvertToRawAmount() (*RawAmount, error) {
	exp := int(n.base) - len(n.decimal)
	if exp < 0 {
		return nil, errors.New("lowest decimal already reached")
	}

	amm, err := NewRawFromString(n.whole + n.decimal)
	if err != nil {
		return nil, err
	}

	r := amm.Multiply(pow(10, exp))

	return r, nil
}

func (n *HumanAmount) ConvertToBase(base UnitBase, scale int) (string, error) {
	exp := int(base) - scale
	if exp < 0 {
		return "", errors.New("lowest decimal already reached")
	}

	raw, err := n.ConvertToRawAmount()
	if err != nil {
		return "", err
	}

	ramm := raw.Divide(pow(10, exp)).ToString()
	r := whole(ramm, scale) + decimal(ramm, scale)

	return r, nil
}

func pow(a, b int) *RawAmount {
	return &RawAmount{new(big.Int).Exp(new(big.Int).SetInt64(int64(a)), new(big.Int).SetInt64(int64(b)), nil)}
}

func whole(s string, scale int) string {
	if scale == 0 {
		return s
	}

	if len(s) > scale {
		return s[:len(s)-scale]
	}

	return "0"
}

func decimal(s string, scale int) string {
	if scale == 0 {
		return ""
	}

	if len(s) >= scale {
		return "." + s[len(s)-scale:]
	}

	r := make([]byte, scale)
	for i := range r {
		r[i] = 0x30
	}

	copy(r[scale-len(s):], s)

	return "." + string(r)
}
