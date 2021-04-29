package common

import (
	"testing"
)

func TestCalculateMd5(t *testing.T) {
	if CalculateMd5("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ") != "cb838c0b6221441c539e190be8b0b21b" {
		t.Error("md5 value returned by CaculateMd5 has changed!")
	}
}
