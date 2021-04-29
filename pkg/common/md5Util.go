package common

import (
	"crypto/md5"
	"encoding/hex"
)

func CalculateMd5(str string) string {
	md5Calculator := md5.New()
	md5Calculator.Write([]byte(str))
	data := md5Calculator.Sum(nil)
	value := hex.EncodeToString(data)
	return value
}

func CalculateMd5ForBytes(bytes []byte) string {
	md5Calculator := md5.New()
	md5Calculator.Write(bytes)
	data := md5Calculator.Sum(nil)
	value := hex.EncodeToString(data)
	return value
}
