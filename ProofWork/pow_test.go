package ProofWork

import "testing"

var Genesis = []byte("991CF190094C00F0B68E2E5F75F6BEE95A2E0BD93CEAA4A6734DB9F19B728948")

func BenchmarkGenerateWork(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateProof(Genesis)
	}
}