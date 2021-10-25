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

package HdfsOss

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"mosn.io/layotto/components/file"

	"go.beyondstorage.io/services/hdfs"
	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

const (
	endpointKey = "endpoint"
	fileSize    = "filesize"
)

var (
	ErrMissingEndPoint    error = errors.New("missing endpoint info in metadata")
	ErrClientNotExist     error = errors.New("specific client not exist")
	ErrEndPointNotExist   error = errors.New("specific endpoing key not exist")
	ErrInvalidConfig      error = errors.New("invalid hdfs oss config")
	ErrNotSpecifyEndpoint error = errors.New("other error happend in metadata")
	ErrHdfsListFail       error = errors.New("hdfs List opt failed")
)

type HdfsOss struct {
	client map[string]types.Storager
	meta   map[string]*HdfsMetaData
}

type HdfsMetaData struct {
	EndPoint string `json:"endpoint"`
}

func NewHdfsOss() file.File {
	return &HdfsOss{
		client: make(map[string]types.Storager),
		meta:   make(map[string]*HdfsMetaData),
	}
}

func (h *HdfsOss) Init(config *file.FileConfig) error {
	hd := make([]*HdfsMetaData, 0)
	err := json.Unmarshal(config.Metadata, &hd)

	if err != nil {
		return ErrInvalidConfig
	}
	for _, data := range hd {
		if !data.isHdfsMetaValid() {
			return ErrInvalidConfig
		}
		client, err := h.createOssClient(data)

		if err != nil {
			continue
		}

		h.client[data.EndPoint] = client
		h.meta[data.EndPoint] = data
	}

	return nil
}

func (h *HdfsOss) Put(stu *file.PutFileStu) error {
	endpoint := stu.Metadata[endpointKey]

	//It depends on OS HDFS XML ???
	if endpoint == "" {
		return ErrMissingEndPoint
	}

	client, err := h.selectClient(stu.Metadata)
	if err != nil {
		return err
	}

	var size int64
	if filesize, ok := stu.Metadata[fileSize]; ok {
		size, err = strconv.ParseInt(filesize, 10, 64)
		if err != nil {
			return err
		}
	}

	_, err = client.Write(stu.FileName, stu.DataStream, size)

	if err != nil {
		return err
	}

	return nil
}

func (h *HdfsOss) Get(stu *file.GetFileStu) (io.ReadCloser, error) {
	if _, ok := stu.Metadata[endpointKey]; !ok {
		return nil, ErrMissingEndPoint
	}
	client, err := h.selectClient(stu.Metadata)
	if err != nil {
		return nil, err
	}

	var w bytes.Buffer
	_, err = client.Read(stu.FileName, &w)
	if err != nil {
		return nil, err
	}
	r := ioutil.NopCloser(bytes.NewReader(w.Bytes()))

	return r, nil
}

func (h *HdfsOss) List(request *file.ListRequest) (*file.ListResp, error) {
	if _, ok := request.Metadata[endpointKey]; !ok {
		return nil, ErrMissingEndPoint
	}

	client, err := h.selectClient(request.Metadata)
	if err != nil {
		return nil, err
	}
	starter := ""
	resp := &file.ListResp{}

	it, err := client.List(starter)
	if err != nil {
		err = ErrHdfsListFail
	}

	for {
		o, err := it.Next()
		if err != nil && !errors.Is(err, types.IterateDone) {
			fmt.Errorf("fail to List All Next: %v", err)
		}

		if err != nil {
			fmt.Println("list completed")
			break
		}
		resp.FilesName = append(resp.FilesName, o.Path)
	}
	return resp, nil
}

func (h *HdfsOss) Del(request *file.DelRequest) error {

	if _, ok := request.Metadata[endpointKey]; !ok {
		return ErrMissingEndPoint
	}
	client, err := h.selectClient(request.Metadata)
	if err != nil {
		return err
	}
	return client.Delete(request.FileName)
}

func (h *HdfsOss) selectClient(meta map[string]string) (client types.Storager, err error) {
	var endpoint string
	var ok bool

	//endpoint not invaild
	if endpoint, ok = meta[endpointKey]; !ok {

		if len(h.client) == 1 {
			for _, client := range h.client {
				return client, nil
			}
		}
		//May be not use?
		//Because BeyondStorage implemented storage type cannot be assigned a value
		return nil, ErrNotSpecifyEndpoint
	}

	if client, ok = h.client[endpoint]; !ok {
		err = ErrClientNotExist
		return
	}
	return client, err
}

func (h *HdfsOss) createOssClient(meta *HdfsMetaData) (types.Storager, error) {
	client, err := hdfs.NewStorager(pairs.WithEndpoint(meta.EndPoint))

	if err != nil {
		return nil, err
	}
	return client, nil
}

// ishdfsMetaValid check if the metadata is valid
func (hm *HdfsMetaData) isHdfsMetaValid() bool {
	if hm.EndPoint == "" {
		return false
	}
	return true
}
