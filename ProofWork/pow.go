package ProofWork

import (
	"encoding/binary"
	"golang.org/x/crypto/blake2b"
	"runtime"
	"encoding/json"
	"github.com/brokenbydefault/Nanollet/Util"
)

type Work []byte

var MinimumWork = uint64(0xffffffc000000000)

// GenerateProof will generate the proof of work, the nonce itself is one uint64 in BigEndian, it's generate as follows:
// Pick one unique Nonce and concatenate with the Blockhash:
// [LittleEndian UINT64 Nonce][BlockHash]
// Now computes the hash of the previous concatenation:
// Blake2(size = 8, message = [LittleEndian UINT64 Nonce][BlockHash])
// Now you need to use this value as one UINT64 and compare against the minimum work:
// LitleEndian(Blake2(...)) > MinimumWork
// If it's correct then you have in hand one correct nonce/pow, you need to reverse it so use the BigEndian.
func GenerateProof(blockHash []byte) Work {
	limit := uint64(runtime.NumCPU())
	shard := uint64(1<<64-1) / limit

	result := make(chan uint64, 32)
	stop := make(chan bool)

	for i := uint64(0); i < limit; i++ {
		go createProof(blockHash, i*shard, result, stop)
	}

	nonce := <-result
	close(stop)
	close(result)
	clear(result)

	n := make([]byte, 8)
	binary.BigEndian.PutUint64(n, nonce)

	return Work(n)
}

func (w Work) IsValid(blockHash []byte) bool {
	n := make([]byte, 8)
	copy(n, w)

	binary.LittleEndian.PutUint64(n, binary.BigEndian.Uint64(n))

	h, _ := blake2b.New(8, nil)
	h.Write(n)
	h.Write(blockHash)

	return binary.LittleEndian.Uint64(h.Sum(nil)) >= MinimumWork
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

func (d *Work) MarshalJSON() ([]byte, error) {
	return json.Marshal(Util.UnsafeHexEncode(*d))
}

func (d *Work) UnmarshalJSON(data []byte) (err error) {
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
