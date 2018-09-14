package Ephemeral

import (
	"golang.org/x/crypto/curve25519"
	"crypto/rand"
	"golang.org/x/crypto/blake2b"
)

const (
	X25519Size = 32
)

type PublicKey [X25519Size]byte
type SecretKey [X25519Size]byte

func NewEphemeral() (sk SecretKey) {
	if _, err := rand.Read(sk[:]); err != nil {
		panic("impossible create key")
	}

	return sk
}

func (e SecretKey) PublicKey() (pk PublicKey) {
	pkb, sk := [32]byte{}, [32]byte(e)
	curve25519.ScalarBaseMult(&pkb, &sk)

	return PublicKey(pkb)
}

func (e SecretKey) Exchange(partner PublicKey) (key [32]byte) {
	exkey, pk, sk := [32]byte{}, [32]byte(partner), [32]byte(e)
	curve25519.ScalarMult(&exkey, &sk, &pk)

	return blake2b.Sum256(exkey[:])
}
