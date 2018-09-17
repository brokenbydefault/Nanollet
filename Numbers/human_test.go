package Numbers

import (
	"testing"
)

func TestHumanToRaw(t *testing.T) {

	input := "0.000001"
	amm := NewHumanFromString(input, MegaXRB)
	output, _ := amm.ConvertToBase(MegaXRB, 6)

	if output != input {
		t.Error("wrong value")
	}

}
