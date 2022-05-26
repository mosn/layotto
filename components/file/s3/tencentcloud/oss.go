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

package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"

	"mosn.io/layotto/components/file"
)

const (
	endpointKey    = "endpoint"
	aclKey         = "ACL"
	contentTypeKey = "content-type"
)

var (
	ErrClientNotExist     = errors.New("specific client not exist")
	ErrNotSpecifyEndPoint = errors.New("not specify endpoint in metadata")
)

type TencentCloudOSS struct {
	metadata map[string]*OssMetadata
	client   map[string]*cos.Client
}

type OssMetadata struct {
	Endpoint        string `json:"endpoint"`        // bucket url https://console.cloud.tencent.com/cos/bucket
	AccessKeyID     string `json:"accessKeyID"`     // SecretID https://console.cloud.tencent.com/cam/capi
	AccessKeySecret string `json:"accessKeySecret"` // SecretKey
	Timeout         int    `json:"timeout"`         // timeout in milliseconds
	bucketUrl       *url.URL
}

func NewTencentCloudOSS() file.File {
	oss := &TencentCloudOSS{metadata: make(map[string]*OssMetadata), client: make(map[string]*cos.Client)}
	return oss
}

// Init does metadata parsing and connection creation
func (t *TencentCloudOSS) Init(ctx context.Context, metadata *file.FileConfig) error {
	m := make([]*OssMetadata, 0)
	err := json.Unmarshal(metadata.Metadata, &m)
	if err != nil {
		return file.ErrInvalid
	}

	for _, v := range m {
		if !t.checkMetadata(v) {
			return file.ErrInvalid
		}
		client, err := t.getClient(v)
		if err != nil {
			return err
		}
		t.metadata[v.Endpoint] = v
		t.client[v.Endpoint] = client
	}
	return nil
}

func (t *TencentCloudOSS) checkMetadata(m *OssMetadata) bool {
	if m.AccessKeySecret == "" || m.Endpoint == "" || m.AccessKeyID == "" {
		return false
	}

	var endpoint = m.Endpoint
	if !strings.HasPrefix(m.Endpoint, "http") {
		endpoint = "https://" + m.Endpoint
	}
	bucketUrl, err := url.Parse(endpoint)
	if err != nil {
		return false
	}
	m.bucketUrl = bucketUrl

	if m.Timeout <= 0 {
		m.Timeout = 100 * 1000 //100s
	}

	return true
}

func (t *TencentCloudOSS) getClient(metadata *OssMetadata) (*cos.Client, error) {
	b := &cos.BaseURL{BucketURL: metadata.bucketUrl}
	client := cos.NewClient(b, &http.Client{
		//set timeout
		Timeout: time.Duration(metadata.Timeout) * time.Millisecond,
		Transport: &cos.AuthorizationTransport{
			SecretID:  metadata.AccessKeyID,
			SecretKey: metadata.AccessKeySecret,
		},
	})

	return client, nil
}

func (t *TencentCloudOSS) Put(ctx context.Context, st *file.PutFileStu) error {
	if err := t.checkFileName(st.FileName); err != nil {
		return err
	}

	client, err := t.selectClient(st.Metadata)
	if err != nil {
		return err
	}

	opt := &cos.ObjectPutOptions{}
	if v, ok := st.Metadata[contentTypeKey]; ok {
		opt.ObjectPutHeaderOptions = &cos.ObjectPutHeaderOptions{
			ContentType: v,
		}
	}

	if v, ok := st.Metadata[aclKey]; ok {
		opt.ACLHeaderOptions = &cos.ACLHeaderOptions{
			XCosACL: v,
		}
	} else {
		opt.ACLHeaderOptions = &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		}
	}

	_, err = client.Object.Put(ctx, st.FileName, st.DataStream, opt)
	return err
}

func (t *TencentCloudOSS) Get(ctx context.Context, st *file.GetFileStu) (io.ReadCloser, error) {
	if err := t.checkFileName(st.FileName); err != nil {
		return nil, err
	}

	client, err := t.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}

	clientResp, err := client.Object.Get(ctx, st.FileName, nil)
	if err != nil {
		return nil, err
	}

	return clientResp.Body, nil
}

func (t *TencentCloudOSS) List(ctx context.Context, st *file.ListRequest) (*file.ListResp, error) {
	client, err := t.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}

	if st.PageSize > 1000 {
		return nil, errors.New("PageSize must be <=1000")
	}

	if st.PageSize <= 0 {
		return nil, errors.New("PageSize must be >0")
	}

	opt := &cos.BucketGetOptions{
		Prefix:    st.DirectoryName,
		Delimiter: "/",
		MaxKeys:   int(st.PageSize),
		Marker:    st.Marker,
	}

	result, _, err := client.Bucket.Get(ctx, opt)
	if err != nil {
		return nil, err
	}

	resp := &file.ListResp{}
	resp.IsTruncated = result.IsTruncated
	resp.Marker = result.NextMarker

	for _, v := range result.Contents {
		file := &file.FilesInfo{}
		file.FileName = v.Key
		file.Size = v.Size
		file.LastModified = v.LastModified
		resp.Files = append(resp.Files, file)
	}

	return resp, nil
}

func (t *TencentCloudOSS) Del(ctx context.Context, st *file.DelRequest) error {
	if err := t.checkFileName(st.FileName); err != nil {
		return err
	}

	client, err := t.selectClient(st.Metadata)
	if err != nil {
		return err
	}

	_, err = client.Object.Delete(ctx, st.FileName)
	return err
}

func (t *TencentCloudOSS) Stat(ctx context.Context, st *file.FileMetaRequest) (*file.FileMetaResp, error) {
	if err := t.checkFileName(st.FileName); err != nil {
		return nil, err
	}

	client, err := t.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}

	var clientResp *cos.Response
	clientResp, err = client.Object.Head(ctx, st.FileName, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(*clientResp)
	resp := &file.FileMetaResp{}
	resp.Metadata = make(map[string][]string)
	for k, v := range clientResp.Header {
		if k == "Content-Length" {
			if len(v) > 0 {
				l, err := strconv.Atoi(v[0])
				if err == nil {
					resp.Size = int64(l)
				}
			}
			continue
		}
		if k == "Last-Modified" {
			if len(v) > 0 {
				resp.LastModified = v[0]
			}
			continue
		}
		resp.Metadata[k] = append(resp.Metadata[k], v...)
	}

	return resp, nil
}

func (t *TencentCloudOSS) selectClient(meta map[string]string) (*cos.Client, error) {
	var (
		endpoint string
		ok       bool
		client   *cos.Client
	)
	if endpoint, ok = meta[endpointKey]; !ok {
		if len(t.client) == 1 {
			for _, client := range t.client {
				return client, nil
			}
		}
		return nil, ErrNotSpecifyEndPoint
	}

	if client, ok = t.client[endpoint]; !ok {
		return nil, ErrClientNotExist
	}

	return client, nil
}

func (t *TencentCloudOSS) checkFileName(fileName string) error {
	index := strings.Index(fileName, "/")
	if index == 0 {
		return fmt.Errorf("invalid fileName format")
	}

	return nil
}
