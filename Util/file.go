package Util

import "io/ioutil"

func FileToString(filepath string) string {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return string(file)
}
