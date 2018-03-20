package Util

func ReverseBytes(b []byte) []byte {
	length := len(b)
	last := length - 1

	for i := 0; i < length/2; i++ {
		b[i], b[last-i] = b[last-i], b[i]
	}

	return b
}