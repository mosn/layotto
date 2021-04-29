package common

import (
	"os"
)

func GetFileSize(f string) int64 {
	fi, err := os.Stat(f)
	if err == nil {
		return fi.Size()
	}

	return -1
}
