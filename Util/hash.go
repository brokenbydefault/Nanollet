package Util

import (
	"golang.org/x/crypto/blake2b"
)

// CreateHash returns the Blake2b hash of message with specified size.
func CreateHash(size int, messagebyte... []byte) []byte {
	hash, err := blake2b.New(size, nil)
	if err != nil {
		panic("hashing failed")
	}

	for _, mb := range messagebyte {
		hash.Write(mb)
	}

	return hash.Sum(nil)
}

// CreateHash returns the Blake2b hash of message with specified size.
func CreateKeyedHash(size int, key []byte, messagebyte... []byte) []byte {
	hash, err := blake2b.New(size, key)
	if err != nil {
		panic("hashing failed")
	}

	for _, mb := range messagebyte {
		hash.Write(mb)
	}

	return hash.Sum(nil)
}

// CreateKeyedXOFHash returns the Blake2b-XOF hash of message using an secret-key, with specified size.
func CreateKeyedXOFHash(size uint32, key []byte, messagebyte... []byte) []byte {
	hash, err := blake2b.NewXOF(size, key)
	if err != nil {
		panic("hashing failed")
	}

	for _, mb := range messagebyte {
		hash.Write(mb)
	}

	r := make([]byte, size)
	hash.Read(r)

	return r
}