package Storage

import "bytes"

var Frontier []byte

func SetFrontier(b []byte) {
	Frontier = b
}

func IsAccountOpened() bool {
	return !bytes.Equal(Frontier, []byte{0x00})
}
