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

package qiniu

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"mosn.io/layotto/components/file"
)

const (
	endpointKey = "endpoint"
	fileSizeKey = "filesize"
)

var (
	ErrClientNotExist     = errors.New("specific client not exist")
	ErrNotSpecifyEndPoint = errors.New("not specify endpoint in metadata")
)

type QiniuOSS struct {
	metadata map[string]*OssMetadata
	client   map[string]*QiniuOSSClient
}

type OssMetadata struct {
	Endpoint        string `json:"endpoint"`        // bucket url
	AccessKeyID     string `json:"accessKeyID"`     // SecretID
	AccessKeySecret string `json:"accessKeySecret"` // SecretKey
	Bucket          string `json:"bucket"`
	Private         bool   `json:"private"`
	UseHTTPS        bool   `json:"useHTTPS"`
	UseCdnDomains   bool   `json:"useCdnDomains"`
}

func NewQiniuOSS() file.File {
	return &QiniuOSS{
		metadata: make(map[string]*OssMetadata),
		client:   make(map[string]*QiniuOSSClient),
	}
}

func (q *QiniuOSS) Init(ctx context.Context, metadata *file.FileConfig) error {
	m := make([]*OssMetadata, 0)
	err := json.Unmarshal(metadata.Metadata, &m)
	if err != nil {
		return file.ErrInvalid
	}

	for _, v := range m {
		if !v.checkMetadata() {
			return file.ErrInvalid
		}

		var domain string
		if v.UseHTTPS {
			domain = "https://" + v.Endpoint
		} else {
			domain = "http://" + v.Endpoint
		}
		client := newQiniuOSSClient(
			v.AccessKeyID,
			v.AccessKeySecret,
			v.Bucket,
			domain,
			v.Private,
			v.UseHTTPS,
			v.UseCdnDomains,
		)
		q.metadata[v.Endpoint] = v
		q.client[v.Endpoint] = client
	}
	return nil
}

func (m *OssMetadata) checkMetadata() bool {
	if m.AccessKeySecret == "" || m.Endpoint == "" || m.AccessKeyID == "" || m.Bucket == "" {
		return false
	}

	return true
}

func (q *QiniuOSS) selectClient(meta map[string]string) (*QiniuOSSClient, error) {
	var (
		endpoint string
		ok       bool
		client   *QiniuOSSClient
	)
	if endpoint, ok = meta[endpointKey]; !ok {
		if len(q.client) == 1 {
			for _, client := range q.client {
				return client, nil
			}
		}
		return nil, ErrNotSpecifyEndPoint
	}

	if client, ok = q.client[endpoint]; !ok {
		return nil, ErrClientNotExist
	}

	return client, nil
}

func (q *QiniuOSS) Put(ctx context.Context, st *file.PutFileStu) error {
	client, err := q.selectClient(st.Metadata)
	if err != nil {
		return err
	}

	fileSize, ok := st.Metadata[fileSizeKey]
	if !ok {
		return errors.New("filesize must be set")
	}

	val, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil {
		return err
	}

	return client.put(ctx, st.FileName, st.DataStream, val)
}

func (q *QiniuOSS) Get(ctx context.Context, st *file.GetFileStu) (io.ReadCloser, error) {
	client, err := q.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}

	return client.get(ctx, st.FileName)
}

func (q *QiniuOSS) List(ctx context.Context, st *file.ListRequest) (*file.ListResp, error) {
	client, err := q.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}

	items, _, nextMarker, hasNext, err := client.list(ctx, st.DirectoryName, int(st.PageSize), st.Marker)

	if err != nil {
		return nil, err
	}

	resp := &file.ListResp{
		Marker:      nextMarker,
		IsTruncated: hasNext,
		Files:       make([]*file.FilesInfo, len(items)),
	}

	for i := 0; i < len(items); i++ {
		resp.Files[i] = &file.FilesInfo{
			FileName:     items[i].Key,
			Size:         items[i].Fsize,
			LastModified: strconv.FormatInt(items[i].PutTime, 10),
			Meta:         make(map[string]string),
		}

		resp.Files[i].Meta["hash"] = items[i].Hash
		resp.Files[i].Meta["mimeType"] = items[i].MimeType
		resp.Files[i].Meta["type"] = strconv.Itoa(items[i].Type)
		resp.Files[i].Meta["endUser"] = items[i].EndUser
	}

	return resp, nil
}

func (q *QiniuOSS) Del(ctx context.Context, st *file.DelRequest) error {
	client, err := q.selectClient(st.Metadata)
	if err != nil {
		return err
	}

	return client.del(ctx, st.FileName)
}

func (q *QiniuOSS) Stat(ctx context.Context, st *file.FileMetaRequest) (*file.FileMetaResp, error) {
	client, err := q.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}

	fileInfo, err := client.stat(ctx, st.FileName)
	if err != nil {
		return nil, err
	}

	resp := &file.FileMetaResp{
		Size:         fileInfo.Fsize,
		LastModified: strconv.FormatInt(fileInfo.PutTime, 10),
		Metadata:     make(map[string][]string),
	}

	return resp, nil
}
