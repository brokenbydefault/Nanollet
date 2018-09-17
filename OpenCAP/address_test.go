package OpenCAP

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func TestAddress_GetPublicKey(t *testing.T) {
	expected := Wallet.Address("nano_3tz9pdfskx934ce36cf6h17uspp4hzsamr5hk7u1wd6em1gfsnb618hfsafc").MustGetPublicKey()

	pk, err := Address("nanollet-gotest$ogdolo.com").GetPublicKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pk != expected {
		t.Error("invalid public-key received")
		return
	}
}

func TestAddress_IsValid(t *testing.T) {
	addrs := []Address{
		Address("nanollet-gotest$ogdolo.com"),
		Address("name$site.com.br"),
		Address("name$subdomain.site.com.br"),
	}

	for _, addr := range addrs {
		if !addr.IsValid() {
			t.Errorf("returns invalid a valid address: %s", addr )
			return
		}
	}
}

func TestAddress_IsValid_Invalid(t *testing.T) {
	addrs := []Address{
		Address("nano_3tz9pdfskx934ce36cf6h17uspp4hzsamr5hk7u1wd6em1gfsnb618hfsafc"),
		Address(""),
		Address("$ogdolo.com"),
		Address("abc$.com"),
		Address("abc$site."),
		Address("$.com"),
		Address("ac$$.com"),
		Address("ac$$."),
		Address("ac$....."),
		Address("ac$.."),
		Address("$b.c.d"),
	}

	for _, addr := range addrs {
		if addr.IsValid() {
			t.Errorf("return valid a invalid address: %s", addr)
			return
		}
	}

}
