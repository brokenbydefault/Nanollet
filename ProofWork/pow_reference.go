// +build ignore

package ProofWork

import (
	"encoding/binary"
	"github.com/brokenbydefault/Nanollet/Util"
	"fmt"
)

func ReferenceGenerateProof(blockHash []byte) []byte {
	var attempt uint64
	var nonce = make([]byte, 8)

	for attempt = 0; attempt < 1<<64-1; attempt++ {

		binary.LittleEndian.PutUint64(nonce, attempt)
		hash := Util.CreateHash(8, nonce, blockHash)

		if attempt < 100 {
			fmt.Println(hash)
		}

		if binary.LittleEndian.Uint64(hash) >= MinimumWork {
			break
		}

	}

	return Util.ReverseBytes(nonce)
}
