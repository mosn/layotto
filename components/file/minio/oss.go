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

package minio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"mosn.io/layotto/components/file/util"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/minio/minio-go/v7"

	"mosn.io/layotto/components/file"
)

const (
	endpointKey = "endpoint"
	fileSize    = "fileSize"
)

var (
	ErrMissingEndPoint    error = errors.New("missing endpoint info in metadata")
	ErrClientNotExist     error = errors.New("specific client not exist")
	ErrInvalidConfig      error = errors.New("invalid minio oss config")
	ErrNotSpecifyEndPoint error = errors.New("not specify endpoint in metadata")
)

type MinioOss struct {
	client map[string]*minio.Core
	meta   map[string]*MinioMetaData
}

type MinioMetaData struct {
	Region          string `json:"region"`
	EndPoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	SSL             bool   `json:"SSL"`
}

func NewMinioOss() file.File {
	return &MinioOss{
		client: make(map[string]*minio.Core),
		meta:   make(map[string]*MinioMetaData),
	}
}

func (m *MinioOss) Init(ctx context.Context, config *file.FileConfig) error {
	md := make([]*MinioMetaData, 0)
	err := json.Unmarshal(config.Metadata, &md)
	if err != nil {
		return ErrInvalidConfig
	}
	for _, data := range md {
		if !data.isMinioMetaValid() {
			return ErrInvalidConfig
		}
		client, err := m.createOssClient(data)
		if err != nil {
			continue
		}
		m.client[data.EndPoint] = client
		m.meta[data.EndPoint] = data
	}
	return nil
}

func (m *MinioOss) Put(ctx context.Context, st *file.PutFileStu) error {
	var (
		size int64 = -1
	)
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return fmt.Errorf("minio put file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return fmt.Errorf("minio put file[%s] fail,err: %s", st.FileName, err.Error())
	}
	core, err := m.selectClient(st.Metadata)
	if err != nil {
		return err
	}
	// specify file size from metadata, default unknown size is -1
	if info, ok := st.Metadata[fileSize]; ok {
		size, err = strconv.ParseInt(info, 10, 64)
		if err != nil {
			return err
		}
	}
	_, err = core.Client.PutObject(ctx, bucket, key, st.DataStream, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioOss) Get(ctx context.Context, st *file.GetFileStu) (io.ReadCloser, error) {
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("minio get file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("minio get file[%s] fail,err: %s", st.FileName, err.Error())
	}
	core, err := m.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	obj, err := core.Client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *MinioOss) List(ctx context.Context, st *file.ListRequest) (*file.ListResp, error) {
	bucket, err := util.GetBucketName(st.DirectoryName)
	marker := ""
	if err != nil {
		return nil, fmt.Errorf("minio list bucket[%s] fail, err: %s", st.DirectoryName, err.Error())
	}
	prefix := util.GetFilePrefixName(st.DirectoryName)

	core, err := m.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	resp := &file.ListResp{}
	out, err := core.ListObjects(bucket, prefix, st.Marker, "", int(st.PageSize))
	if err != nil {
		return nil, err
	}
	resp.IsTruncated = out.IsTruncated
	for _, object := range out.Contents {
		file := &file.FilesInfo{}
		file.FileName = object.Key
		file.Size = object.Size
		file.LastModified = object.LastModified.String()
		resp.Files = append(resp.Files, file)
		marker = object.Key
	}
	resp.Marker = marker
	return resp, nil
}

func (m *MinioOss) Del(ctx context.Context, st *file.DelRequest) error {
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return fmt.Errorf("minio del file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return fmt.Errorf("minio del file[%s] fail,err: %s", st.FileName, err.Error())
	}
	core, err := m.selectClient(st.Metadata)
	if err != nil {
		return fmt.Errorf("minio del file[%s] fail,err: %s", st.FileName, err.Error())
	}
	return core.Client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (m *MinioOss) Stat(ctx context.Context, st *file.FileMetaRequest) (*file.FileMetaResp, error) {
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("minio stat file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("minio stat file[%s] fail,err: %s", st.FileName, err.Error())
	}
	core, err := m.selectClient(st.Metadata)
	if err != nil {
		return nil, fmt.Errorf("minio stat file[%s] fail,err: %s", st.FileName, err.Error())
	}
	info, err := core.Client.StatObject(ctx, bucket, key, minio.GetObjectOptions{})

	if err != nil {
		if err.(minio.ErrorResponse).StatusCode == 404 {
			return nil, file.ErrNotExist
		}
		return nil, err
	}

	resp := &file.FileMetaResp{}
	resp.Metadata = make(map[string][]string)
	resp.LastModified = info.LastModified.String()
	resp.Size = info.Size
	resp.Metadata[util.ETag] = append(resp.Metadata[util.ETag], info.ETag)
	for k, v := range info.Metadata {
		resp.Metadata[k] = v
	}
	return resp, nil
}

func (m *MinioOss) createOssClient(meta *MinioMetaData) (*minio.Core, error) {
	client, err := minio.New(
		meta.EndPoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(meta.AccessKeyID, meta.AccessKeySecret, ""),
			Secure: meta.SSL,
			Region: meta.Region,
		},
	)
	if err != nil {
		return nil, err
	}
	return &minio.Core{Client: client}, nil
}

func (m *MinioOss) selectClient(meta map[string]string) (client *minio.Core, err error) {
	var (
		endpoint string
		ok       bool
	)
	if endpoint, ok = meta[endpointKey]; !ok {
		if len(m.client) == 1 {
			for _, client := range m.client {
				return client, nil
			}
		}
		return nil, ErrNotSpecifyEndPoint
	}
	if client, ok = m.client[endpoint]; !ok {
		err = ErrClientNotExist
		return
	}
	return
}

// isMinioMetaValid check if the metadata is valid
func (mm *MinioMetaData) isMinioMetaValid() bool {
	if mm.AccessKeySecret == "" || mm.EndPoint == "" || mm.AccessKeyID == "" {
		return false
	}
	return true
}
