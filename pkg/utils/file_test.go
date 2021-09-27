package utils

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sendData(w *WrapFileStream) {
	w.Write([]byte("1234567890"))
	w.Close()
}
func TestNewFileReader(t *testing.T) {
	s := NewFileReader()
	go sendData(s)
	result := make([]byte, 0, 10)
	for {
		data := make([]byte, 2)
		n, err := s.Read(data)
		if err == io.EOF {
			break
		}
		assert.Equal(t, n, 2)
		assert.Nil(t, err)
		result = append(result, data...)
	}
	assert.Equal(t, string(result), "1234567890")
}
