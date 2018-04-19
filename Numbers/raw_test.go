package Numbers

import (
	"fmt"
	"testing"
)

func TestNewRawFromString(t *testing.T) {

	n, err := NewRawFromString("340282366920937463463374607431768211456")
	if err != nil {
		panic(err)
	}
	x, err := NewRawFromString("340282366920938463463374607431768211455")
	if err != nil {
		panic(err)
	}
	n = n.Add(x)
	fmt.Println(n.ToString(), n.IsValid())
	fmt.Println(n.ToBytes())

}

func TestRawAmount_IsValid(t *testing.T) {

	if _, err := NewRawFromString("-1"); err == nil {
		t.Error("is valid a invalid amount")
	}

	if _, err := NewRawFromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"); err == nil {
		t.Error("is valid a invalid amount")
	}

}
