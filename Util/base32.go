package Util

import "encoding/base32"

var unsafeb32 = base32.NewEncoding("13456789abcdefghijkmnopqrstuwxyz").WithPadding(base32.NoPadding)

// UnsafeBase32Encode encodes an byte in non-constant time. It's results
// the Base32 string.
func UnsafeBase32Encode(src []byte) string {
	return unsafeb32.EncodeToString(src)
}

// UnsafeBase32Decode decodes an string in non-constant time. It's results
// the bytes of decoded Base32.
func UnsafeBase32Decode(src string) ([]byte, error) {
	return unsafeb32.DecodeString(src)
}

//@TODO SecureBase32Encode (LowPriority: base32 is used with public-data)
//@TODO SecureBase32Decode (LowPriority: base32 is used with public-data)
