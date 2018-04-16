package Numbers

import (
	"fmt"
	"testing"
)

func TestHumanToRaw(t *testing.T) {

	fmt.Println(MegaXRB)

	input := "0.000001"
	amm := NewHumanFromString(input, MegaXRB)
	output, _ := amm.ConvertToBase(MegaXRB, 6)

	if output != input {
		t.Error("wrong value")
	}

}
