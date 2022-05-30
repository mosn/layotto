/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"mosn.io/layotto/components/file"
)

const (
	FileMode  = "FileMode"
	FileFlag  = "FileFlag"
	FileIsDir = "IsDir"
)

type LocalStore struct {
}

func NewLocalStore() file.File {
	return &LocalStore{}
}

func (lf *LocalStore) Init(ctx context.Context, f *file.FileConfig) error {
	return nil
}
func (lf *LocalStore) Put(ctx context.Context, f *file.PutFileStu) error {
	mode, ok := f.Metadata[FileMode]
	if !ok {
		return fmt.Errorf("fileMode is required for put file")
	}
	flag, ok := f.Metadata[FileFlag]
	if !ok {
		return fmt.Errorf("fileFlag is required for put file")
	}
	m, err := strconv.ParseUint(mode, 10, 32)
	if err != nil {
		return fmt.Errorf("wrong fileMode value:%+v in metadata", err)
	}
	fl, err := strconv.Atoi(flag)
	if err != nil {
		return fmt.Errorf("wrong fileFlag value:%+v in metadata", err)
	}

	fileObj, err := os.OpenFile(f.FileName, fl, os.FileMode(m))
	if err != nil {
		return err
	}
	defer fileObj.Close()
	data := make([]byte, 512)
	for {
		n, err := f.DataStream.Read(data)
		if err != nil {
			if err == io.EOF {
				if n > 0 {
					_, err = fileObj.Write(data[:n])
					if err != nil {
						return err
					}
				}
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
func (lf *LocalStore) Get(ctx context.Context, f *file.GetFileStu) (io.ReadCloser, error) {
	fileObj, err := os.Open(f.FileName)
	if err != nil {
		return nil, err
	}
	return fileObj, nil
}
func (lf *LocalStore) List(ctx context.Context, f *file.ListRequest) (*file.ListResp, error) {
	res := &file.ListResp{}
	files, err := ioutil.ReadDir(f.DirectoryName)
	if err != nil {
		return nil, err
	}
	marker := ""
	fileNumber := 0
	start := false
	isTruncated := true
	if len(files) == 0 {
		isTruncated = false
	}
	if f.Marker == "" {
		start = true
	}

	for index, fileObj := range files {
		if index == len(files)-1 {
			isTruncated = false
		}
		if fileObj.Name() == f.Marker && !start {
			start = true
			continue
		}
		if start {
			fileNumber++
			marker = fileObj.Name()
			info := &file.FilesInfo{}
			info.Size = fileObj.Size()
			info.LastModified = fileObj.ModTime().String()
			info.FileName = fileObj.Name()
			res.Files = append(res.Files, info)
			if fileNumber == int(f.PageSize) {
				break
			}
		}
	}
	res.Marker = marker
	res.IsTruncated = isTruncated
	return res, nil
}
func (lf *LocalStore) Del(ctx context.Context, f *file.DelRequest) error {
	err := os.Remove(f.FileName)
	return err
}

func (lf *LocalStore) Stat(ctx context.Context, f *file.FileMetaRequest) (*file.FileMetaResp, error) {
	fileInfo, err := os.Stat(f.FileName)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, file.ErrNotExist
		}
		return nil, err
	}
	resp := &file.FileMetaResp{}
	resp.Metadata = make(map[string][]string)
	resp.Size = fileInfo.Size()
	resp.LastModified = fileInfo.ModTime().String()
	mode := int(fileInfo.Mode())
	m := strconv.Itoa(mode)
	t := fileInfo.IsDir()
	isDir := "false"
	if t {
		isDir = "true"
	}
	resp.Metadata[FileMode] = append(resp.Metadata[FileMode], m)
	resp.Metadata[FileIsDir] = append(resp.Metadata[FileIsDir], isDir)
	return resp, nil
}
