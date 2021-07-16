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
	metadata map[string]*ossMetadata
	client   map[string]*oss.Client
	stream   sync.Map
}

type ossMetadata struct {
	Endpoint        string   `json:"endpoint"`
	AccessKeyID     string   `json:"accessKeyIDKey"`
	AccessKeySecret string   `json:"accessKeySecretKey"`
	Bucket          []string `json:"bucketKey"`
}

func NewAliCloudOSS() file.File {
	oss := &AliCloudOSS{metadata: make(map[string]*ossMetadata), client: make(map[string]*oss.Client)}
	return oss
}

// Init does metadata parsing and connection creation
func (s *AliCloudOSS) Init(metadata *file.FileConfig) error {
	m := &ossMetadata{}
	if len(metadata.Metadata) == 0 {
		return fmt.Errorf("no configuration for aliCloudOSS")
	}
	for _, v := range metadata.Metadata {
		m.Endpoint = v[endpointKey].(string)
		m.AccessKeyID = v[accessKeyIDKey].(string)
		m.AccessKeySecret = v[accessKeySecretKey].(string)
		for _, s := range v[bucketKey].([]interface{}) {
			m.Bucket = append(m.Bucket, s.(string))
		}
		if !s.checkMetadata(m) {
			return fmt.Errorf("wrong configurations for aliCloudOSS")
		}
		client, err := s.getClient(m)
		if err != nil {
			return err
		}
		s.metadata[m.Endpoint] = m
		s.client[m.Endpoint] = client
	}
	return nil
}

func (s *AliCloudOSS) CompletePut(streamId int64, success bool) error {
	if !success {
		s.stream.Delete(streamId)
		return nil
	}
	if v, ok := s.stream.Load(streamId); ok {
		pu := v.(*PartUploadStu)
		_, err := pu.bucket.CompleteMultipartUpload(pu.imur, pu.parts)
		s.stream.Delete(streamId)
		return err
	}
	return fmt.Errorf("file is not uploading")
}

func (s *AliCloudOSS) Put(st *file.PutFileStu) error {
	storageType := st.Metadata[storageTypeKey]
	if v, ok := s.stream.Load(st.StreamId); !ok {
		bucket, err := s.selectClientAndBucket(st.Metadata)
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

func (s *AliCloudOSS) checkMetadata(m *ossMetadata) bool {
	if m.AccessKeySecret == "" || m.Endpoint == "" || m.AccessKeyID == "" || len(m.Bucket) == 0 {
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

func (s *AliCloudOSS) selectClientAndBucket(metaData map[string]string) (*oss.Bucket, error) {
	ossClient := &oss.Client{}
	bucket := &oss.Bucket{}
	var err error
	// get oss client
	if _, ok := metaData[endpointKey]; ok {
		ossClient = s.client[endpointKey]
	} else {
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
