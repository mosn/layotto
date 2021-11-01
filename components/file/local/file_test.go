package local

import (
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

const (
	FileName = "test.txt"
)

func WriteFile(writer *io.PipeWriter) {
	writer.Write([]byte("hello"))
	writer.Close()
}

func CheckFileExist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}

func TestFile(t *testing.T) {
	ls := &LocalStore{}
	f := &file.PutFileStu{}
	f.FileName = FileName
	f.Metadata = make(map[string]string)
	reader, writer := io.Pipe()
	err := ls.Put(f)
	assert.Equal(t, err.Error(), "fileMode is required for put file")
	mode := 0777
	f.Metadata[FileMode] = strconv.Itoa(mode)
	err = ls.Put(f)
	assert.Equal(t, err.Error(), "fileFlag is required for put file")
	f.Metadata[FileFlag] = strconv.Itoa(os.O_RDWR | os.O_CREATE)
	f.DataStream = reader
	go WriteFile(writer)
	err = ls.Put(f)
	exist, err := CheckFileExist(f.FileName)
	assert.Nil(t, err)
	assert.Equal(t, true, exist)

	data := make([]byte, 10, 10)
	fs := &file.GetFileStu{}
	fs.FileName = FileName
	stream, err := ls.Get(fs)
	assert.Nil(t, err)
	n, err := stream.Read(data)
	stream.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(data[:n]), "hello")

	fr := &file.ListRequest{}
	fr.DirectoryName = "."
	res, err := ls.List(fr)
	assert.Nil(t, err)
	t.Log(res.FilesName)

	fd := &file.DelRequest{}
	fd.FileName = FileName
	err = ls.Del(fd)
	assert.Nil(t, err)
	exist, err = CheckFileExist(FileName)
	assert.Nil(t, err)
	assert.Equal(t, exist, false)
}
