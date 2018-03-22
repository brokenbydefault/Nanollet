package Wallet

import (
	"github.com/brokenbydefault/Nanollet/Util"
	"encoding/binary"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"runtime"
	"runtime/debug"
	"errors"
)

type Currency uint32

//@TODO Support Banano
const (
	Base Currency = iota
	Nano
	//Banano
)

var SupportedVersions = [...]uint8{0}

type SeedFY struct {
	Version uint8
	Type    uint8
	Time    uint8
	Memory  uint8
	Thread  uint8
	Salt    []byte
}

// NewSeedFY generate the SeedFY, which is the random salt and the default computational cost parameters
// used in the Argon2id derivation in combination with the password.
func NewSeedFY() (sf SeedFY, err error) {
	sf = SeedFY{
		Version: 0,
		Type:    0,
		Time:    15,
		Memory:  21,
		Thread:  uint8(runtime.NumCPU()),
		Salt:    make([]byte, 32),
	}

	_, err = rand.Read(sf.Salt)
	return
}

var ErrImpossibleDecode = errors.New("impossible to decode the seed")

//ReadSeedFY act like to NewSeedFY, however it creates the struct based on the given hex-encoded SeedFY.
func ReadSeedFY(s string) (sf SeedFY, err error) {
	sb, ok := Util.SecureHexDecode(s)
	if !ok {
		return sf, ErrImpossibleDecode
	}

	sf = SeedFY{
		Version: uint8(sb[0]),
		Type:    uint8(sb[1]),
		Time:    uint8(sb[2]),
		Memory:  uint8(sb[3]) & 0x1F,
		Thread:  sb[4],
		Salt:    []byte(sb[5:]),
	}

	if !sf.IsValid() {
		return sf, ErrImpossibleDecode
	}

	return
}

// Encode will return the hexadecimal representation of the given SeedFY.
func (sf *SeedFY) Encode() (string, error) {
	var seed = make([]byte, 37)

	if !sf.IsValid() {
		return "", ErrImpossibleDecode
	}

	copy(seed, []byte{
		sf.Version,
		sf.Type,
		sf.Time,
		sf.Memory,
		sf.Thread,
	})
	copy(seed[5:], sf.Salt)

	return Util.SecureHexEncode(seed), nil
}

// IsValid will return false if the SeedFY is not supported or don't have
// enough seed-length
func (sf *SeedFY) IsValid() (r bool) {
	for _, v := range SupportedVersions {
		if v == sf.Type {
			r = true
			break
		}
	}

	if len(sf.Salt) != 32 {
		return false
	}

	if sf.Memory > 31 {
		return false
	}

	return
}



// @TODO

type Seed []byte

// RecoverSeedFromSeedFY returns the Seed based on given password and hex-encoded SeedFY.
// SEEDFY: [version][type][time][memory][thread][salt]
func (sf *SeedFY) RecoverSeed(password string, additionaldata []byte) Seed {
	salt := Util.CreateKeyedHash(32, sf.Salt, additionaldata)

	kdf := argon2.IDKey([]byte(password), salt, uint32(sf.Time), uint32(1<<sf.Memory), sf.Thread, 32)
	debug.FreeOSMemory()

	return kdf
}

// CreateKeyPair creates the public-key and secret-key using the given currency and index.
func (s *Seed) CreateKeyPair(coin Currency, index uint32) (PublicKey, SecretKey, error) {
	return RecoverKeyPairFromCoinSeed(RecoverCoinSeed(*s, coin), index)
}

type CoinSeed []byte

// RecoverKeyPairFromSeed will return an seed for given currency.
func RecoverCoinSeed(seed Seed, coin Currency) CoinSeed {
	indexbyte := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexbyte, uint32(coin))

	return Util.CreateKeyedXOFHash(32, seed, indexbyte)
}

// RecoverKeyPairFromSeed will return the Ed25519 key-pair based on given COINSEED and INDEX, it also returns
// non-nil error in case of failure. It uses the Blake2X instead of Blake2b and can support up to 256 keys.
// The Nanollet Seed uses BLAKE2bXOF(size = COINSIZE, key = COINSEED, message = INDEX)
func RecoverKeyPairFromCoinSeed(seed CoinSeed, i uint32) (PublicKey, SecretKey, error) {
	indexbyte := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexbyte, i)

	return CreateKeyPair(Util.CreateKeyedXOFHash(32, seed, indexbyte))
}
