package Util

import (
	"encoding/hex"
	"strings"
)

// UnsafeBase32Encode encodes an byte in non-constant time. It's results
// the Base32 string.
func UnsafeHexEncode(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

// UnsafeHexDecode decodes an byte without table-lookup. It's results
// the byte-array of decoded hex.
func UnsafeHexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// SecureHexEncode decodes an byte without table-lookup. It's results
// the hexadecimal string representation.
func SecureHexEncode(b []byte) string {
	length := len(b)
	r := make([]byte, length*2)

	for i := 0; i < length; i++ {

		// Splits a single byte into two bytes, resulting in 4 bits for each byte.
		p1, p2 := int(b[i]>>4), int(b[i]&0x0F)

		b1 := 48 + p1
		b2 := 48 + p2

		// The ASCII have a gap between "9" and "A", we need to jump it if the byte is between 10 and 15.
		// if(byte > 9) { b += 7 }
		b1 += ((9 - p1) >> 8) & 7
		b2 += ((9 - p2) >> 8) & 7

		copy(r[i*2:], []byte{byte(b1), byte(b2)})
	}

	return string(r)
}

// SecureHexDecode decodes an byte without table-lookup. It's results
// the byte-array of decoded hex.
func SecureHexDecode(s string) (r []byte, ok bool) {
	length := len(s)

	if (length & 1) == 1 {
		return nil, false
	}

	r = make([]byte, length/2)
	e := 0

	for i := 0; i < length; i++ {

		// "0" - (48) = 0
		// "A" - 48 - (7) = 10
		// "a" - 48 - 7 - (32) = 10
		b := int(s[i] - 48)
		b -= ((9 - b) >> 8) & 7
		b -= ((41 - b) >> 8) & 32

		// if (b < 0 | b > 15) { err = 1 }
		e |= ((b >> 8) | (15-b)>>8) & 1

		st := ^((i & 1) << 8) & (i / 2)
		r[st] <<= 4
		r[st] |= byte(b)

	}

	return r, e == 0
}
