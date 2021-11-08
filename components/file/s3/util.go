package s3

import (
	"fmt"
	"strings"
)

const (
	ETag = "ETag"
)

func GetBucketName(fileName string) (string, error) {
	index := strings.Index(fileName, "/")
	if index == -1 || index == 0 {
		return "", fmt.Errorf("invalid fileName format")
	}
	return fileName[:index], nil
}

func GetFilePrefixName(fileName string) string {
	index := strings.Index(fileName, "/")
	if index == -1 {
		return ""
	}
	return fileName[index+1:]
}

func GetFileName(fileName string) (string, error) {
	index := strings.Index(fileName, "/")
	if index == -1 {
		return "", fmt.Errorf("invalid fileName format")
	}
	name := fileName[index+1:]
	if name == "" {
		return "", fmt.Errorf("file name is empty")
	}
	return name, nil
}
