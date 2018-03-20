package Util

// StringIsNumeric return false if the input is not between 0 to 9.
func StringIsNumeric(s string) (ok bool) {
	var err int32

	for _, b := range s {
		err |= (((b-48) >> 8) | (57-b)>>8) & 1
	}

	return err == 0
}
