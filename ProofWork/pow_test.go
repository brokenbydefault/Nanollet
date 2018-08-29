package ProofWork

import (
	"fmt"
	"github.com/brokenbydefault/Nanollet/Util"
	"testing"
)

var Genesis, _ = Util.UnsafeHexDecode("991CF190094C00F0B68E2E5F75F6BEE95A2E0BD93CEAA4A6734DB9F19B728948")

func BenchmarkGenerateWork(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Util.SecureHexEncode(GenerateProof(Genesis)[:])
	}
}

func TestGenerateProof(t *testing.T) {
	fmt.Println(Util.SecureHexEncode(GenerateProof(Genesis)[:]))
}

func TestIsValidProof(t *testing.T) {
	if !GenerateProof(Genesis).IsValid(Genesis) {
		t.Error("valid proof reported as invalid")
	}
}

//func TestReferenceGenerateProof(t *testing.T) {
//	fmt.Println(Util.SecureHexEncode(ReferenceGenerateProof(Genesis)))
//}
