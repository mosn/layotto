package local

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"mosn.io/layotto/components/file"
)

const (
	FileMode = "fileMode"
	FileFlag = "fileFlag"
)

type LocalStore struct {
}

func NewLocalStore() file.File {
	return &LocalStore{}
}

func (lf *LocalStore) Init(f *file.FileConfig) error {
	return nil
}
func (lf *LocalStore) Put(f *file.PutFileStu) error {
	if _, ok := f.Metadata[FileMode]; !ok {
		return fmt.Errorf("fileMode is required for put file")
	}
	if _, ok := f.Metadata[FileFlag]; !ok {
		return fmt.Errorf("fileFlag is required for put file")
	}
	mode := f.Metadata[FileMode]
	m, err := strconv.Atoi(mode)
	if err != nil {
		return fmt.Errorf("wrong fileMode value:%+v in metadata", err)
	}

	flag := f.Metadata[FileFlag]
	fl, err := strconv.Atoi(flag)
	if err != nil {
		return fmt.Errorf("wrong fileFlag value:%+v in metadata", err)
	}

	fileObj, err := os.OpenFile(f.FileName, fl, os.FileMode(m))
	if err != nil {
		return err
	}
	defer fileObj.Close()
	data := make([]byte, 512, 512)
	for {
		n, err := f.DataStream.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = fileObj.Write(data[:n])
		if err != nil {
			return err
		}
	}
	return nil
}
func (lf *LocalStore) Get(f *file.GetFileStu) (io.ReadCloser, error) {
	fileObj, err := os.Open(f.FileName)
	if err != nil {
		return nil, err
	}
	return fileObj, nil
}
func (lf *LocalStore) List(f *file.ListRequest) (*file.ListResp, error) {
	res := &file.ListResp{}
	files, err := ioutil.ReadDir(f.DirectoryName)
	if err != nil {
		return nil, err
	}
	for _, fileObj := range files {
		res.FilesName = append(res.FilesName, fileObj.Name())
	}
	return res, nil
}
func (lf *LocalStore) Del(f *file.DelRequest) error {
	err := os.Remove(f.FileName)
	return err
}
