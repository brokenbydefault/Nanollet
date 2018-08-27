package Wallet

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/subtle"
	"errors"
	"github.com/Inkeliz/blakEd25519"
	"github.com/brokenbydefault/Nanollet/Util"
	"io"
)

type SecretKey []byte
type PublicKey []byte

// CreateKeyPair creates the Ed25519 key-pair from an given input and
// returns the public/private key, and error not nil if something go wrong.
func CreateKeyPair(b []byte) (PublicKey, SecretKey, error) {
	return createKeyPair(bytes.NewReader(b))
}

// GenerateRandomKeyPair creates the Ed25519, using an random input, we
// use the crypto/rand.
func GenerateRandomKeyPair() (PublicKey, SecretKey, error) {
	return createKeyPair(cryptorand.Reader)
}

func createKeyPair(r io.Reader) (PublicKey, SecretKey, error) {
	pk, sk, err := blakEd25519.GenerateKey(r)
	return PublicKey(pk), SecretKey(sk), err
}

// PublicKeyFromSecretKey extract the Ed25519 public-key
// from the secret key and return the public-key.
func (sk SecretKey) PublicKey() (PublicKey, error) {
	if len(sk) != blakEd25519.PrivateKeySize {
		return nil, errors.New("wrong size of secret-key")
	}

	pk := make([]byte, 32)
	copy(pk, sk[32:])

	return PublicKey(pk), nil
}

// Checksum creates the checksum for given public-key, it returns the checksum
// in byte format.
func (pk PublicKey) CreateChecksum() []byte {
	return Util.ReverseBytes(Util.CreateHash(5, pk))
}

// CompareChecksum check the publick-key with arbitrary given checksum, it will return
// true if the checksum matches and false otherwise.
func (pk PublicKey) IsValidChecksum(checksum []byte) bool {
	return subtle.ConstantTimeCompare(pk.CreateChecksum(), checksum) == 1
	// It's not need to be constant-time since both inputs are public. But this code can be recycled in future, been used in other circumstances.
}

type Signature []byte

// CreateSignature signs the message with the private-key. It return
// the signature.
func (sk SecretKey) CreateSignature(message []byte) (Signature, error) {
	if len(sk) != blakEd25519.PrivateKeySize {
		return nil, errors.New("wrong size of secret-key")
	}

	sig := blakEd25519.Sign([]byte(sk), message)

	pk, err := sk.PublicKey()
	if err != nil {
		return nil, err
	}

	if !pk.IsValidSignature(message, sig) {
		return nil, errors.New("signature is not correct")
	}

	return sig, nil
}

// IsValidSignature checks the authenticity of the signature based on public-key, it returns
// false if wrong.
func (pk PublicKey) IsValidSignature(message, sig []byte) bool {
	return blakEd25519.Verify(blakEd25519.PublicKey(pk), message, sig)
}
