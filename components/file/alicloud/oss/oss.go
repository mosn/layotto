package oss

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mosn.io/layotto/components/file"
)

const (
	endpointKey        = "endpoint"
	accessKeyIDKey     = "accessKeyID"
	accessKeySecretKey = "accessKeySecret"
	bucketKey          = "bucket"
	storageTypeKey     = "storageType"
)

type PartUploadStu struct {
	imur   oss.InitiateMultipartUploadResult
	parts  []oss.UploadPart
	bucket *oss.Bucket
}

// AliCloudOSS is a binding for an AliCloud OSS storage bucketKey
type AliCloudOSS struct {
	metadata *ossMetadata
	client   *oss.Client
	stream   sync.Map
}

type ossMetadata struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyIDKey"`
	AccessKeySecret string `json:"accessKeySecretKey"`
	Bucket          string `json:"bucketKey"`
}

// Init does metadata parsing and connection creation
func (s *AliCloudOSS) Init(metadata *file.FileConfig) error {
	m := &ossMetadata{}
	m.Endpoint = metadata.Metadata[endpointKey]
	m.Endpoint = metadata.Metadata[accessKeyIDKey]
	m.Endpoint = metadata.Metadata[accessKeySecretKey]
	m.Endpoint = metadata.Metadata[bucketKey]
	if !s.checkMetadata(m) {
		return fmt.Errorf("wrong configuration for oss")
	}
	client, err := s.getClient(m)
	if err != nil {
		return err
	}
	s.metadata = m
	s.client = client

	return nil
}

func (s *AliCloudOSS) CompletePut(streamId int64) error {
	if v, ok := s.stream.Load(streamId); ok {
		pu := v.(*PartUploadStu)
		_, err := pu.bucket.CompleteMultipartUpload(pu.imur, pu.parts)
		return err
	}
	return fmt.Errorf("file is not uploading")
}

func (s *AliCloudOSS) Put(st *file.PutFileStu) error {
	storageType := st.Metadata[storageTypeKey]
	if v, ok := s.stream.Load(st.StreamId); !ok {
		// create bucket ob
		bucket, err := s.client.Bucket(s.metadata.Bucket)
		if err != nil {
			return err
		}
		//initial multi part upload
		imur, err := bucket.InitiateMultipartUpload(st.FileName, oss.ObjectStorageClass(oss.StorageClassType(storageType)))
		if err != nil {
			return err
		}
		//upload part
		part, err := bucket.UploadPart(imur, bytes.NewReader(st.Data), int64(len(st.Data)), st.ChunkNumber)
		if err != nil {
			return err
		}

		pu := &PartUploadStu{imur: imur, bucket: bucket}
		pu.parts = append(pu.parts, part)
		s.stream.Store(st.StreamId, pu)
		return nil
	} else {
		pu := v.(*PartUploadStu)
		//upload part
		part, err := pu.bucket.UploadPart(pu.imur, bytes.NewReader(st.Data), int64(len(st.Data)), st.ChunkNumber)
		if err != nil {
			return err
		}
		pu.parts = append(pu.parts, part)
		s.stream.Store(st.StreamId, pu)
		return nil
	}
}

func (s *AliCloudOSS) Get(stu *file.GetFileStu) (io.ReadCloser, error) {
	// create bucket ob
	bucket, err := s.client.Bucket(s.metadata.Bucket)
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(stu.ObjectName)
}

func (s *AliCloudOSS) checkMetadata(m *ossMetadata) bool {
	if m.AccessKeySecret == "" || m.Endpoint == "" || m.AccessKeyID == "" || m.Bucket == "" {
		return false
	}
	return true
}

func (s *AliCloudOSS) getClient(metadata *ossMetadata) (*oss.Client, error) {
	client, err := oss.New(metadata.Endpoint, metadata.AccessKeyID, metadata.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}
