package Util

// StringIsNumeric return false if the input is not between 0 to 9.
func StringIsNumeric(s string) (ok bool) {
	var err int32

	for _, b := range s {
		err |= (((b - 48) >> 8) | (57-b)>>8) & 1
	}

	return err == 0
}

type Order int8

const (
	LittleEndian Order = iota
	BigEndian
)

func BytesToUint(b []byte, order Order) uint64 {
	if order == LittleEndian {
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	} else {
		return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 | uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	}
}

func UintToBytes(i uint64, order Order) []byte {
	if order == LittleEndian {
		return []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24), byte(i >> 32), byte(i >> 40), byte(i >> 48), byte(i >> 56)}
	} else {
		return []byte{byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32), byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
	}
}
