package Util

import (
	"testing"
)

func TestUnsafeBase32Encode(t *testing.T) {

	table := []struct {
		Message string
		Result  string
	}{
		{"Testing the base32", "cjkq8x5bfsmk1x5aeni86rdmensm6"},
	}

	for _, v := range table {
		if UnsafeBase32Encode([]byte(v.Message)) != v.Result {
			t.Errorf("UnsafeBase32Encode failed")
		}
	}

}
