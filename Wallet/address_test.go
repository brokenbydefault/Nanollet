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

func TestAddress_IsValid_Invalid(t *testing.T) {
	addrs := []Address{
		Address("_3tz9pdfskx934ce36cf6h17uspp4hzsamr5hk7u1wd6em1gfsnb618hfsafc"),
		Address("3tz9pdfskx934ce36cf6h17uspp4hzsamr5hk7u1wd6em1gfsnb618hfsafc"),
		Address("nano_3tz9pdfskx934ce36cf6h17uspp4hzsamr5hk7u1wd6em1gfsnb618hfsqfc"),
		Address(""),
		Address("nano$ogdolo.com"),
		Address("nano_$ogdolo.com"),
	}

	for _, addr := range addrs {
		if addr.IsValid() {
			t.Errorf("return valid a invalid address: %s", addr)
			return
		}
	}

}
