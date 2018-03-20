package Wallet

import (
	"testing"
)

func TestAddress_IsValid(t *testing.T) {
	addr := Address("xrb_1kmp7wp1muz4ct3pzyceorxjm7pzeuszqrm86ky8ysytmfjxe9awfe8ngc8i")
	if !addr.IsValid() {
		t.Error("valid address is invalid")
	}
}
