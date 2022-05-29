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

package hdfs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"mosn.io/layotto/components/file"

	store "go.beyondstorage.io/services/hdfs"
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
	ErrInvalidConfig      error = errors.New("invalid hdfs config")
	ErrNotSpecifyEndpoint error = errors.New("other error happend in metadata")
	ErrHdfsListFail       error = errors.New("hdfs list opt failed")
	ErrInitFailed         error = errors.New("hdfs client init failed")
)

type hdfs struct {
	client map[string]types.Storager
	meta   map[string]*HdfsMetaData
}

type HdfsMetaData struct {
	EndPoint string `json:"endpoint"`
}

func NewHdfs() file.File {
	return &hdfs{
		client: make(map[string]types.Storager),
		meta:   make(map[string]*HdfsMetaData),
	}
}

func (h *hdfs) Init(ctx context.Context, config *file.FileConfig) error {
	hd := make([]*HdfsMetaData, 0)
	err := json.Unmarshal(config.Metadata, &hd)

	if err != nil {
		return ErrInvalidConfig
	}
	for _, data := range hd {
		if !data.isHdfsMetaValid() {
			return ErrInvalidConfig
		}
		client, err := h.createHdfsClient(data)

		if err != nil {
			return ErrInitFailed
		}

		h.client[data.EndPoint] = client
		h.meta[data.EndPoint] = data
	}

	return nil
}

func (h *hdfs) Put(ctx context.Context, stu *file.PutFileStu) error {
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
	return err
}

func (h *hdfs) Get(ctx context.Context, stu *file.GetFileStu) (io.ReadCloser, error) {
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

func (h *hdfs) List(ctx context.Context, request *file.ListRequest) (*file.ListResp, error) {
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
		return nil, ErrHdfsListFail
	}

	marker := ""
	for {
		o, err := it.Next()
		if err != nil && !errors.Is(err, types.IterateDone) {
			return nil, err
		}

		if err != nil {
			fmt.Println("list completed")
			break
		}
		file := &file.FilesInfo{}
		file.FileName = o.Path

		size, ok := o.GetContentLength()
		if !ok {
			return nil, fmt.Errorf("Hdfs list path[%s] size fail, err: %s", o.Path, err.Error())
		}

		file.Size = size

		time, ok := o.GetLastModified()
		if !ok {
			return nil, fmt.Errorf("Hdfs list path[%s] lastModified fail, err: %s", o.Path, err.Error())
		}
		file.LastModified = time.String()

		resp.Files = append(resp.Files, file)

		marker = o.Path
	}

	resp.Marker = marker

	return resp, nil
}

func (h *hdfs) Del(ctx context.Context, request *file.DelRequest) error {
	if _, ok := request.Metadata[endpointKey]; !ok {
		return ErrMissingEndPoint
	}

	client, err := h.selectClient(request.Metadata)
	if err != nil {
		return err
	}
	return client.Delete(request.FileName)
}

func (h *hdfs) Stat(ctx context.Context, request *file.FileMetaRequest) (*file.FileMetaResp, error) {

	clinet, err := h.selectClient(request.Metadata)
	if err != nil {
		return nil, err
	}

	stat, err := clinet.Stat(request.FileName)
	if err != nil {
		return nil, err
	}

	resp := &file.FileMetaResp{}

	size, ok := stat.GetContentLength()
	if !ok {
		return nil, fmt.Errorf("Hdfs stat file[%s] size fail, err: %s", stat.Path, err.Error())
	}

	resp.Size = size

	time, ok := stat.GetLastModified()
	if !ok {
		return nil, fmt.Errorf("Hdfs stat file[%s] lastModified fail, err: %s", stat.Path, err.Error())
	}

	resp.LastModified = time.String()

	return resp, nil
}

func (h *hdfs) selectClient(meta map[string]string) (client types.Storager, err error) {
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

func (h *hdfs) createHdfsClient(meta *HdfsMetaData) (types.Storager, error) {
	return store.NewStorager(pairs.WithEndpoint(meta.EndPoint))
}

// ishdfsMetaValid check if the metadata is valid
func (hm *HdfsMetaData) isHdfsMetaValid() bool {
	return hm.EndPoint != ""
}
