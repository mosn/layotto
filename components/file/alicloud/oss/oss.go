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

package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mosn.io/layotto/components/file"
)

const (
	endpointKey    = "endpoint"
	bucketKey      = "bucket"
	storageTypeKey = "storageType"
)

// AliCloudOSS is a binding for an AliCloud OSS storage bucketKey
type AliCloudOSS struct {
	metadata map[string]*OssMetadata
	client   map[string]*oss.Client
	stream   sync.Map
}

type OssMetadata struct {
	Endpoint        string   `json:"endpoint"`
	AccessKeyID     string   `json:"accessKeyID"`
	AccessKeySecret string   `json:"accessKeySecret"`
	Bucket          []string `json:"bucket"`
}

func NewAliCloudOSS() file.File {
	oss := &AliCloudOSS{metadata: make(map[string]*OssMetadata), client: make(map[string]*oss.Client)}
	return oss
}

// Init does metadata parsing and connection creation
func (s *AliCloudOSS) Init(metadata *file.FileConfig) error {
	m := make([]*OssMetadata, 0)
	err := json.Unmarshal(metadata.Metadata, &m)
	if err != nil {
		return file.ErrInvalid
	}

	for _, v := range m {
		if !s.checkMetadata(v) {
			return file.ErrInvalid
		}
		client, err := s.getClient(v)
		if err != nil {
			return err
		}
		s.metadata[v.Endpoint] = v
		s.client[v.Endpoint] = client
	}
	return nil
}

func (s *AliCloudOSS) Put(ctx context.Context, st *file.PutFileStu) error {
	storageType := st.Metadata[storageTypeKey]
	if storageType == "" {
		storageType = "Standard"
	}
	bucket, err := s.getBucket(st.FileName, st.Metadata)
	fileNameWithoutBucket := s.getFile(st.FileName)
	if err != nil {
		return fmt.Errorf("fail to find bucket for %s: %v", st.Metadata, err)
	}
	err = bucket.PutObject(fileNameWithoutBucket, st.DataStream, oss.ObjectStorageClass(oss.StorageClassType(storageType)), oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		return fmt.Errorf("fail to upload object: %+v", err)
	}

	return nil
}

func (s *AliCloudOSS) Get(ctx context.Context, st *file.GetFileStu) (io.ReadCloser, error) {
	bucket, err := s.getBucket(st.FileName, st.Metadata)
	if err != nil {
		return nil, err
	}
	fileNameWithoutBucket := s.getFile(st.FileName)
	return bucket.GetObject(fileNameWithoutBucket)
}

func (s *AliCloudOSS) List(ctx context.Context, request *file.ListRequest) (*file.ListResp, error) {
	if request.DirectoryName == "" {
		return nil, fmt.Errorf("not specifc directory name")
	}
	bucket, err := s.getBucket(request.DirectoryName, request.Metadata)
	if err != nil {
		return nil, err
	}
	resp := &file.ListResp{}
	prefix := s.getFile(request.DirectoryName)
	object, err := bucket.ListObjectsV2(oss.StartAfter(request.Marker), oss.MaxKeys(int(request.PageSize)), oss.Prefix(prefix))
	if err != nil {
		return nil, err
	}
	resp.IsTruncated = object.IsTruncated
	l := len(object.Objects)
	//last object is marker
	if l > 0 {
		resp.Marker = object.Objects[l-1].Key
	}
	for _, v := range object.Objects {
		file := &file.FilesInfo{}
		file.FileName = v.Key
		file.Size = v.Size
		file.LastModified = v.LastModified.String()
		resp.Files = append(resp.Files, file)
	}
	return resp, nil
}

func (s *AliCloudOSS) Del(ctx context.Context, request *file.DelRequest) error {
	bucket, err := s.getBucket(request.FileName, request.Metadata)
	if err != nil {
		return err
	}
	fileNameWithoutBucket := s.getFile(request.FileName)
	err = bucket.DeleteObject(fileNameWithoutBucket)
	if err != nil {
		return err
	}
	return nil
}

func (s *AliCloudOSS) Stat(ctx context.Context, request *file.FileMetaRequest) (*file.FileMetaResp, error) {
	resp := &file.FileMetaResp{}
	resp.Metadata = make(map[string][]string)
	bucket, err := s.getBucket(request.FileName, request.Metadata)
	if err != nil {
		return nil, err
	}
	fileNameWithoutBucket := s.getFile(request.FileName)
	meta, err := bucket.GetObjectMeta(fileNameWithoutBucket)
	if err != nil {
		if err.(oss.ServiceError).StatusCode == 404 {
			return nil, file.ErrNotExist
		}
		return nil, err
	}
	for k, v := range meta {
		resp.Metadata[k] = append(resp.Metadata[k], v...)
	}
	return resp, nil
}

func (s *AliCloudOSS) checkMetadata(m *OssMetadata) bool {
	if m.AccessKeySecret == "" || m.Endpoint == "" || m.AccessKeyID == "" {
		return false
	}
	return true
}

func (s *AliCloudOSS) getClient(metadata *OssMetadata) (*oss.Client, error) {
	client, err := oss.New(metadata.Endpoint, metadata.AccessKeyID, metadata.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *AliCloudOSS) getBucket(fileName string, metaData map[string]string) (*oss.Bucket, error) {
	ossClient := &oss.Client{}
	bucket := &oss.Bucket{}
	var err error
	// get oss client
	if _, ok := metaData[endpointKey]; ok {
		ossClient = s.client[endpointKey]
	} else {
		// if user not specify endpoint, try to use default client
		ossClient, err = s.selectClient()
		if err != nil {
			return nil, err
		}
	}

	// get oss bucket
	bucketName, err := s.getBucketName(fileName)
	if err != nil {
		return nil, err
	}
	bucket, err = ossClient.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (s *AliCloudOSS) selectClient() (*oss.Client, error) {
	if len(s.client) == 1 {
		for _, client := range s.client {
			return client, nil
		}
	} else {
		return nil, fmt.Errorf("should specific endpoint in metadata")
	}
	return nil, nil
}

func (s *AliCloudOSS) getBucketName(fileName string) (string, error) {
	bucketName := strings.Split(fileName, "/")
	if len(bucketName) == 0 || bucketName[0] == "" {
		return "", fmt.Errorf("invalid fileName format")
	}
	return bucketName[0], nil
}

func (s *AliCloudOSS) getFile(fileName string) string {
	index := strings.Index(fileName, "/")
	if index == -1 {
		return ""
	}
	return fileName[index+1:]
}
