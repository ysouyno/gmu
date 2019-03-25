package utils

import (
	"os"
	"io"
	"fmt"
)

func FileExist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		// file exists
		return true
	} else if os.IsNotExist(err) {
		// file does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist.
		// See err for details

		// Therefore, do *NOT* use !os.IsNotExist(err)
		// to test for file existence
		return false
	}
}

func CopyFile(dst, src string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("\"%s\" is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func Contains(arr []string, str string) bool {
	for _, ele := range arr {
		if ele == str {
			return true
		}
	}

	return false
}
