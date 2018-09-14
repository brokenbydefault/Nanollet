package Block

import (
	"github.com/brokenbydefault/Nanollet/Util"
	"testing"
)

var Genesis = NewBlockHash(Util.UnsafeHexMustDecode("991CF190094C00F0B68E2E5F75F6BEE95A2E0BD93CEAA4A6734DB9F19B728948"))

func BenchmarkGenerateWork(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateProof(&Genesis)
	}
}

func TestIsValidProof(t *testing.T) {
	gen := GenerateProof(&Genesis)
	if !gen.IsValid(&Genesis) {
		t.Error("valid proof reported as invalid")
	}
}
