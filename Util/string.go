package Util

func HasEmptyString(strs []string) bool {
	for _, str := range strs {
		if str == "" {
			return true
		}
	}

	return false
}

func IsEmptyString(str string) bool {
	return str == ""
}
