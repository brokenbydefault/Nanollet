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

type (
	SecretKey [blakEd25519.PrivateKeySize]byte
	PublicKey [blakEd25519.PublicKeySize]byte
	Signature [blakEd25519.SignatureSize]byte
)

func NewSecretKey(b []byte) (sk SecretKey) {
	copy(sk[:], b)
	return sk
}

func NewPublicKey(b []byte) (pk PublicKey) {
	copy(pk[:], b)
	return pk
}

func NewSignature(b []byte) (sig Signature) {
	copy(sig[:], b)
	return sig
}

var (
	ErrInvalidSecretKeySize = errors.New("wrong size of secret-key")
	ErrBadSigning           = errors.New("impossible to sign with this key")
)

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

func createKeyPair(r io.Reader) (pk PublicKey, sk SecretKey, err error) {
	pkb, skb, err := blakEd25519.GenerateKey(r)
	if err != nil {
		return pk, sk, err
	}

	return NewPublicKey(pkb), NewSecretKey(skb), nil
}

// PublicKeyFromSecretKey extract the Ed25519 public-key
// from the secret key and return the public-key.
func (sk *SecretKey) PublicKey() (pk PublicKey) {
	copy(pk[:], sk[32:])
	return pk
}

// Checksum creates the checksum for given public-key, it returns the checksum
// in byte format.
func (pk *PublicKey) Checksum() []byte {
	return Util.ReverseBytes(Util.CreateHash(5, pk[:]))
}

// IsValidChecksum check the publick-key with arbitrary given checksum, it will return
// true if the checksum matches and false otherwise.
func (pk *PublicKey) IsValidChecksum(checksum []byte) bool {
	return subtle.ConstantTimeCompare(pk.Checksum(), checksum) == 1
	// It's not need to be constant-time since both inputs are public. But this code can be recycled in future, been used in other circumstances.
}

// Sign signs the message with the private-key. It return the signature and error, if any.
func (sk *SecretKey) Sign(message []byte) (sig Signature, err error) {
	sigb := blakEd25519.Sign(sk[:], message)
	if sigb == nil {
		return sig, ErrBadSigning
	}

	pk := sk.PublicKey()
	if !blakEd25519.Verify(pk[:], message, sigb) {
		return sig, ErrBadSigning
	}

	return NewSignature(sigb), nil
}

// MustSign signs the message with the private-key. Its panic if some error occur.
func (sk *SecretKey) MustSign(message []byte) Signature {
	sig, err := sk.Sign(message)
	if err != nil {
		panic(err)
	}

	return sig
}

// IsValidSignature checks the authenticity of the signature based on public-key, it returns
// false if wrong.
func (pk *PublicKey) IsValidSignature(message []byte, sig *Signature) bool {
	return blakEd25519.Verify(pk[:], message, sig[:])
}
