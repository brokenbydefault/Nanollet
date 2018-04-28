package Util

func ReverseBytes(b []byte) []byte {
	length := len(b)
	last := length - 1

	for i := 0; i < length/2; i++ {
		b[i], b[last-i] = b[last-i], b[i]
	}

	return b
}

func ConcatBytes(slice ...[]byte) []byte {
	var l int

	for _, b := range slice {
		l += len(b)
	}

	tmp := make([]byte, l)

	var i int
	for _, b := range slice {
		i += copy(tmp[i:], b)
	}

	return tmp

}