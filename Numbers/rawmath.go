package Numbers

import (
	"math/big"
)

// Divide divides the a by b
func (a *RawAmount) Divide(b *RawAmount) *RawAmount {
	return &RawAmount{new(big.Int).Div(a.bigint, b.bigint)}
}

// Multiply multiplies the a by b
func (a *RawAmount) Multiply(b *RawAmount) *RawAmount {
	return &RawAmount{new(big.Int).Mul(a.bigint, b.bigint)}
}

// Subtract subtracts the a by b
func (a *RawAmount) Subtract(b *RawAmount) *RawAmount {
	return &RawAmount{new(big.Int).Sub(a.bigint, b.bigint)}
}

// Add adds the a by b
func (a *RawAmount) Add(b *RawAmount) *RawAmount {
	return &RawAmount{new(big.Int).Add(a.bigint, b.bigint)}
}

// Abs returns the absolute value
func (a *RawAmount) Abs() *RawAmount {
	return &RawAmount{new(big.Int).Abs(a.bigint)}
}

// Compare compares a with b, return
// -1 if a <  b
//  0 if a == b
// +1 if a >  b
func (a *RawAmount) Compare(b *RawAmount) int {
	return a.bigint.Cmp(b.bigint)
}

// IsValid returns true if the value is between 0 to 1<<128-1
func (a *RawAmount) IsValid() bool {
	return a.bigint.Cmp(min) >= 0 && a.bigint.Cmp(max) <= 0
}
