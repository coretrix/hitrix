package helper

import (
	"io/ioutil"
)

func ExistsInDir(filename, dir string) bool {
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, file := range fileInfo {
		if filename == file.Name() {
			return true
		}
	}

	return false
}
