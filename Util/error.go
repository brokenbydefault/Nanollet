package Util

// CheckError will return error if one of the input is non-nil.
func CheckError(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// ExistEmpty will return true if one of the input is empty.
func ExistEmpty(s ...string) bool {
	for _, ss := range s {
		if ss == "" {
			return true
		}
	}
	return false
}
