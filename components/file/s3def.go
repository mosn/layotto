package file

import (
	"context"
	"io"
)

type Oss interface {
	InitConfig(context.Context, *FileConfig) error
	InitClient(context.Context, *InitRequest) error
	GetObject(context.Context, *GetObjectInput) (io.ReadCloser, error)
	PutObject(context.Context) error
}

type BaseConfig struct {
}
type InitRequest struct {
	App      string
	Metadata map[string]string
}

type GetObjectInput struct {
	Bucket                     string
	ExpectedBucketOwner        string
	IfMatch                    string
	IfModifiedSince            string
	IfNoneMatch                string
	IfUnmodifiedSince          string
	Key                        string
	PartNumber                 int64
	Range                      string
	RequestPayer               string
	ResponseCacheControl       string
	ResponseContentDisposition string
	ResponseContentEncoding    string
	ResponseContentLanguage    string
	ResponseContentType        string
	ResponseExpires            string
	SseCustomerAlgorithm       string
	SseCustomerKey             string
	SseCustomerKeyMd5          string
	VersionId                  string
}

type PutObjectInput struct {
	Bucket           string
	Key              string
	DataStream       io.Reader
	Acl              string
	BucketKeyEnabled bool
	CacheControl     string
}
