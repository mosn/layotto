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
	"encoding/json"
	"fmt"
	"io"
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
		return fmt.Errorf("wrong config for alicloudOss")
	}

	for _, v := range m {
		if !s.checkMetadata(v) {
			return fmt.Errorf("wrong configurations for aliCloudOSS")
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

func (s *AliCloudOSS) Put(st *file.PutFileStu) error {
	storageType := st.Metadata[storageTypeKey]
	if storageType == "" {
		storageType = "Standard"
	}
	bucket, err := s.selectClientAndBucket(st.Metadata)
	if err != nil {
		return fmt.Errorf("fail to find bucket for %s: %v", st.Metadata, err)
	}
	err = bucket.PutObject(st.FileName, st.DataStream, oss.ObjectStorageClass(oss.StorageClassType(storageType)), oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		return fmt.Errorf("fail to upload object: %+v", err)
	}

	return nil
}

func (s *AliCloudOSS) Get(st *file.GetFileStu) (io.ReadCloser, error) {
	bucket, err := s.selectClientAndBucket(st.Metadata)
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(st.FileName)
}

func (s *AliCloudOSS) List(request *file.ListRequest) (*file.ListResp, error) {
	if request.DirectoryName == "" {
		return nil, fmt.Errorf("not specifc directory name")
	}
	if request.Metadata != nil {
		request.Metadata = make(map[string]string)
	}
	request.Metadata[bucketKey] = request.DirectoryName
	bucket, err := s.selectClientAndBucket(request.Metadata)
	if err != nil {
		return nil, err
	}
	marker := ""
	resp := &file.ListResp{}
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			return nil, err
		}
		//Return 100 records by default each time
		for _, object := range lsRes.Objects {
			resp.FilesName = append(resp.FilesName, object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return resp, err
}

func (s *AliCloudOSS) Del(request *file.DelRequest) error {
	bucket, err := s.selectClientAndBucket(request.Metadata)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(request.FileName)
	if err != nil {
		return err
	}
	return nil
}

func (s *AliCloudOSS) checkMetadata(m *OssMetadata) bool {
	if m.AccessKeySecret == "" || m.Endpoint == "" || m.AccessKeyID == "" || len(m.Bucket) == 0 {
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

func (s *AliCloudOSS) selectClientAndBucket(metaData map[string]string) (*oss.Bucket, error) {
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
	if _, ok := metaData[bucketKey]; ok {
		bucket, err = ossClient.Bucket(metaData[bucketKey])
		if err != nil {
			return nil, err
		}
	} else {
		bucketName, err := s.selectBucket()
		if err != nil {
			return nil, err
		}
		bucket, err = ossClient.Bucket(bucketName)
		if err != nil {
			return nil, err
		}
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

func (s *AliCloudOSS) selectBucket() (string, error) {
	for _, data := range s.metadata {
		if len(data.Bucket) == 1 {
			return data.Bucket[0], nil
		} else {
			return "", fmt.Errorf("should specific bucketKey in metadata")
		}
	}
	// will be never occur
	return "", fmt.Errorf("no bucket configuration")
}
