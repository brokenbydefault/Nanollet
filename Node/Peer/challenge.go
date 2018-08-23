package Peer

import (
	"net"
	"golang.org/x/crypto/blake2b"
	"crypto/rand"
)

type Challenge []byte

func NewChallenge() Challenge {
	c := make([]byte, 32)
	rand.Read(c)

	return Challenge(c)
}

func (c Challenge) Derivative(ip net.IP) []byte {
	if c == nil {
		return make([]byte, 32)
	}

	hash, _ := blake2b.New256(c)
	hash.Write(ip)

	return hash.Sum(nil)
}
