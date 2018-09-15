package Block

import (
	"encoding/binary"
	"golang.org/x/crypto/blake2b"
	"runtime"
	"github.com/brokenbydefault/Nanollet/Util"
)

type Work [8]byte

func NewWork(b []byte) (work Work) {
	copy(work[:], b)
	return work
}

func (s *SendBlock) Work() Work {
	if !s.PoW.IsValid(&s.Previous) {
		s.PoW = GenerateProof(&s.Previous)
	}

	return s.PoW
}

func (s *ReceiveBlock) Work() Work {
	if !s.PoW.IsValid(&s.Previous) {
		s.PoW = GenerateProof(&s.Previous)
	}

	return s.PoW
}

func (s *OpenBlock) Work() Work {
	// Open operation uses the account instead of previous
	previous := NewBlockHash(s.Account[:])

	if !s.PoW.IsValid(&previous) {
		s.PoW = GenerateProof(&previous)
	}

	return s.PoW
}

func (s *ChangeBlock) Work() Work {
	if !s.PoW.IsValid(&s.Previous) {
		s.PoW = GenerateProof(&s.Previous)
	}

	return s.PoW
}

func (u *UniversalBlock) Work() Work {
	var previous BlockHash

	// Open operation uses the account instead of previous
	if Util.IsEmpty(u.Previous[:]) {
		previous = NewBlockHash(u.Account[:])
	} else {
		previous = u.Previous
	}

	if !u.PoW.IsValid(&previous) {
		u.PoW = GenerateProof(&previous)
	}

	return u.PoW
}

func (s *SendBlock) IsValidPOW() bool {
	return s.PoW.IsValid(&s.Previous)
}

func (s *ReceiveBlock) IsValidPOW() bool {
	return s.PoW.IsValid(&s.Previous)
}

func (s *OpenBlock) IsValidPOW() bool {
	hash := BlockHash(s.Account)
	return s.PoW.IsValid(&hash)
}

func (s *ChangeBlock) IsValidPOW() bool {
	return s.PoW.IsValid(&s.Previous)
}

func (u *UniversalBlock) IsValidPOW() bool  {
	var previous BlockHash

	// Open operation uses the account instead of previous
	if Util.IsEmpty(u.Previous[:]) {
		previous = BlockHash(u.Account)
	} else {
		previous = u.Previous
	}

	return u.PoW.IsValid(&previous)
}

var MinimumWork = uint64(0xffffffc000000000)

// GenerateProof will generate the proof of work for given hash, which must be the "previous" or "account" in case
// of open.
func GenerateProof(hash *BlockHash) (nonce Work) {
	limit := uint64(runtime.NumCPU())
	shard := uint64(1<<64-1) / limit

	result := make(chan uint64, 32)
	stop := make(chan bool)

	for i := uint64(0); i < limit; i++ {
		go createProof(hash[:], i*shard, result, stop)
	}

	n := <-result
	close(stop)
	close(result)
	clear(result)

	binary.BigEndian.PutUint64(nonce[:], n)

	return nonce
}

func (w *Work) IsValid(previous *BlockHash) bool {
	nonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonce, binary.BigEndian.Uint64(w[:]))

	return binary.LittleEndian.Uint64(Util.CreateHash(8, nonce, previous[:])) >= MinimumWork
}

func createProof(blockHash []byte, attempt uint64, result chan uint64, stop chan bool) {
	h, _ := blake2b.New(8, nil)
	nonce := make([]byte, 40)
	copy(nonce[8:], blockHash)

	for {
		select {
		default:
			binary.LittleEndian.PutUint64(nonce[:8], attempt)

			h.Reset()
			h.Write(nonce)

			if binary.LittleEndian.Uint64(h.Sum(nil)) >= MinimumWork {
				result <- attempt
			}

			attempt++

		case <-stop:
			return
		}
	}
}

func clear(r chan uint64) {
	for len(r) > 0 {
		<-r
	}
}
