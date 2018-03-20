package Numbers

import (
	"testing"
	"fmt"
)

func TestNewRawFromString(t *testing.T) {

	n, _ := NewRawFromString("1")
	fmt.Println(n.ToHex())

}
