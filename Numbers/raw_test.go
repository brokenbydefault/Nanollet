package Numbers

import (
	"testing"
)

func TestNewRawFromString(t *testing.T) {
	n, err := NewRawFromString("340282366920937463463374607431768211456")
	if err != nil {
		t.Error(err)
	}
	x, err := NewRawFromString("340282366920938463463374607431768211455")
	if err != nil {
		t.Error(err)
	}
	n = n.Add(x)

}

func TestRawAmount_IsValid(t *testing.T) {

	if _, err := NewRawFromString("-1"); err == nil {
		t.Error("is valid a invalid amount")
	}

	amm, err := NewRawFromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	if err != nil || amm.IsValid() {
		t.Error("is valid a invalid amount")
	}

}
