package Wallet

import (
	"encoding/binary"
	"errors"
	"github.com/brokenbydefault/Nanollet/Util"
)

// RecoverKeyPairFromClassicalSeed will return the Ed25519 key-pair based on the hex-encoded HEX and one INDEX
// it uses the same process used in official wallet of Nano. In case of failure the error will be non-nil.
// The Nano Wallet Seed uses BLAKE2b(size = 32, key = nil, message = SEED+INDEX) to construct the key.
func RecoverKeyPairFromClassicalSeed(seed string, i uint32) (pk PublicKey, sk SecretKey, err error) {
	seedbyte, ok := Util.SecureHexDecode(seed)
	if !ok {
		return pk, sk, errors.New("impossible decode the seed")
	}

	indexbyte := make([]byte, 4)
	binary.BigEndian.PutUint32(indexbyte, i)

	return CreateKeyPair(Util.CreateHash(32, seedbyte, indexbyte))
}
