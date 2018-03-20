package Wallet

import (
	"github.com/brokenbydefault/Nanollet/Util"
	"errors"
	"strings"
	"encoding/json"
)

const ADDRESS_PREFIX = "xrb"

var ALLOWED_PREFIX = [...]string{"xrb"}

type Address string

// CreateAddress creates the encoded address using the public-key. It returns
// the address (with identifier, public-key and checksum) as string, encoded
// with base32.
func (pk *PublicKey) CreateAddress() Address {
	addr := ADDRESS_PREFIX
	addr += "_"
	addr += Util.UnsafeBase32Encode(append([]byte{0, 0, 0}, []byte(*pk)...))[4:]
	addr += Util.UnsafeBase32Encode(pk.CreateChecksum())

	return Address(addr)
}

// PublicKeyFromFromAddress gets the Ed25519 public-key from the encoded address,
// returning the public-key. It's return an non-nil error if something bad happens.
func (addr Address) GetPublicKey() (PublicKey, error) {
	if addr.IsCorrectlyFormatted() == false {
		return nil, errors.New("invalid address")
	}

	addr = "1111" + addr.RemovePrefix()

	pkBytes, err := Util.UnsafeBase32Decode(string(addr[:56]))
	if err != nil {
		return nil, err
	}

	return PublicKey(pkBytes[3:]), nil
}

// GetChecksum extract the existing checksum of the address, returns the checksum
// as byte-array.
func (addr Address) GetChecksum() ([]byte, error) {
	if addr.IsCorrectlyFormatted() == false {
		return nil, errors.New("invalid address")
	}

	addr = "1111" + addr.RemovePrefix()

	checksum, err := Util.UnsafeBase32Decode(string(addr[len(addr)-8:]))
	if err != nil {
		return nil, err
	}

	return checksum, nil
}

// GetPrefix extract the existing prefix of the address, everything before
// the first underscore.
func (addr Address) GetPrefix() string {
	return strings.SplitN(string(addr), "_", 2)[0]
}

// UpdateAddress modify the prefix of the address returning the address with new
// prefix identifier. (Can be used if "xrb_" be replaced by "nano_" in future)
func (addr Address) UpdatePrefix() Address {
	return Address(ADDRESS_PREFIX+"_") + addr.RemovePrefix()
}

// RemovePrefix remove the prefix of the address, returns an address
// without the prefix.
func (addr Address) RemovePrefix() Address {
	return Address(strings.SplitN(string(addr), "_", 2)[1])
}

// IsValid returns true if the given encoded address have an correct formatting and
// also the checksum is correct.
func (addr Address) IsValid() bool {
	pk, err := addr.GetPublicKey()
	if err != nil {
		return false
	}

	checksum, err := addr.GetChecksum()
	if err != nil {
		return false
	}

	return pk.CompareChecksum(checksum)
}

// IsCorrectlyFormatted returns true if the given encoded address have an correct
// format. It return true if had an valid prefix and length, but checksum doesn't matter.
func (addr Address) IsCorrectlyFormatted() bool {
	if len(addr) == 0 || string(addr) == addr.GetPrefix() || len(addr.RemovePrefix()) != 60 {
		return false
	}

	prefix := addr.GetPrefix()
	for _, allowed := range ALLOWED_PREFIX {
		if prefix == allowed {
			return true
		}
	}

	return false
}

func (d *Address) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return
	}

	*d = Address(str)
	return
}

func (d *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*d))
}

func (d *PublicKey) UnmarshalJSON(data []byte) (err error) {
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

func (d PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(Util.UnsafeHexEncode(d))
}
