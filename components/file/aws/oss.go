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

package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"io"
	"mosn.io/layotto/components/file/factory"
	"mosn.io/layotto/components/file/util"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"mosn.io/layotto/components/file"
)

const (
	endpointKey              = "endpoint"
	bucketKey                = "bucket"
	defaultCredentialsSource = "provider"
	fileLength               = "length"
)

var (
	ErrNotSpecifyEndpoint error = errors.New("should specific endpoint in metadata")
)

func init() {
	factory.RegisterInitFunc("", AwsDefaultInitFunc)
}
func AwsDefaultInitFunc(staticConf json.RawMessage, DynConf map[string]string) (map[string]interface{}, error) {
	m := make([]*AwsOssMetaData, 0)
	err := json.Unmarshal(staticConf, &m)
	clients := make(map[string]interface{})
	if err != nil {
		return nil, errors.New("invalid config for aws oss")
	}
	for _, data := range m {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if region == data.Region {
				return aws.Endpoint{
					PartitionID:       "aliyun",
					URL:               "https://" + data.EndPoint,
					SigningRegion:     data.Region,
					HostnameImmutable: true,
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		optFunc := []func(options *aws_config.LoadOptions) error{
			aws_config.WithRegion(data.Region),
			aws_config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: data.AccessKeyID, SecretAccessKey: data.AccessKeySecret,
					Source: defaultCredentialsSource,
				},
			}),
			aws_config.WithEndpointResolverWithOptions(customResolver),
		}

		cfg, err := aws_config.LoadDefaultConfig(context.TODO(), optFunc...)
		if err != nil {
			return nil, err
		}
		clients[data.EndPoint] = s3.NewFromConfig(cfg)
	}
	return clients, nil
}

// AwsOss is a binding for aws oss storage.
type AwsOss struct {
	client  map[string]*s3.Client
	meta    map[string]*AwsOssMetaData
	method  string
	rawData json.RawMessage
}

// AwsOssMetaData describe a aws-oss instance.
type AwsOssMetaData struct {
	Region          string `json:"region"`   // eg. us-west-2
	EndPoint        string `json:"endpoint"` // eg. protocol://service-code.region-code.amazonaws.com
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
}

func NewAwsFile() file.File {
	return &AwsOss{
		client: make(map[string]*s3.Client),
		meta:   make(map[string]*AwsOssMetaData),
	}
}

// Init instance by config.
func (a *AwsOss) Init(ctx context.Context, config *file.FileConfig) error {
	m := make([]*AwsOssMetaData, 0)
	err := json.Unmarshal(config.Metadata, &m)
	if err != nil {
		return errors.New("invalid config for aws oss")
	}
	for _, data := range m {
		if !data.isAwsMetaValid() {
			return errors.New("invalid config for aws oss")
		}
		client, err := a.createOssClient(data)
		if err != nil {
			continue
		}
		a.client[data.EndPoint] = client
		a.meta[data.EndPoint] = data
	}
	return nil
}

// isAwsMetaValid check if the metadata valid.
func (am *AwsOssMetaData) isAwsMetaValid() bool {
	if am.AccessKeySecret == "" || am.EndPoint == "" || am.AccessKeyID == "" {
		return false
	}
	return true
}

// createOssClient by input meta info.
func (a *AwsOss) createOssClient(meta *AwsOssMetaData) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if region == meta.Region {
			return aws.Endpoint{
				PartitionID:       "awsoss",
				URL:               "https://" + meta.EndPoint,
				SigningRegion:     meta.Region,
				HostnameImmutable: true,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	optFunc := []func(options *aws_config.LoadOptions) error{
		aws_config.WithRegion(meta.Region),
		aws_config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: meta.AccessKeyID, SecretAccessKey: meta.AccessKeySecret,
				Source: defaultCredentialsSource,
			},
		}),
		aws_config.WithEndpointResolverWithOptions(customResolver),
	}

	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), optFunc...)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

// Put file to aws oss.
func (a *AwsOss) Put(ctx context.Context, st *file.PutFileStu) error {
	//var bodySize int64
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return fmt.Errorf("awsoss put file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return fmt.Errorf("awsoss put file[%s] fail,err: %s", st.FileName, err.Error())
	}
	//if fileLen, ok := st.Metadata[fileLength]; !ok {
	//	return errors.New("please specific the file lenth in metadata with length key")
	//} else {
	//	bodySize, err = strconv.ParseInt(fileLen, 10, 64)
	//	if err != nil {
	//		return errors.New("wrong value of length")
	//	}
	//}

	//input := &s3.PutObjectInput{
	//	Bucket:        &bucket,
	//	Key:           &key,
	//	Body:          st.DataStream,
	//	ContentLength: bodySize,
	//}
	client, err := a.selectClient(st.Metadata, endpointKey)
	if err != nil {
		return err
	}
	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   st.DataStream,
	})

	//_, err = client.PutObject(context.TODO(), input, s3.WithAPIOptions(
	//	v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
	//))
	if err != nil {
		return err
	}
	return nil
}

// selectClient choose aws client from exist client-map, key is endpoint, value is client instance.
func (a *AwsOss) selectClient(meta map[string]string, key string) (*s3.Client, error) {
	// exist specific client with key endpoint
	if ep, ok := meta[key]; ok {
		if client, ok := a.client[ep]; ok {
			return client, nil
		}
	}
	// if not specify endpoint, select default one
	if len(a.client) == 1 {
		for _, client := range a.client {
			return client, nil
		}
	}
	return nil, ErrNotSpecifyEndpoint
}

// Get object from aws oss.
func (a *AwsOss) Get(ctx context.Context, st *file.GetFileStu) (io.ReadCloser, error) {
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("awsoss get file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("awsoss get file[%s] fail,err: %s", st.FileName, err.Error())
	}
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	client, err := a.selectClient(st.Metadata, endpointKey)
	if err != nil {
		return nil, err
	}
	ob, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return ob.Body, nil
}

// List objects from aws oss.
func (a *AwsOss) List(ctx context.Context, st *file.ListRequest) (*file.ListResp, error) {
	bucket, err := util.GetBucketName(st.DirectoryName)
	if err != nil {
		return nil, fmt.Errorf("list bucket[%s] fail, err: %s", st.DirectoryName, err.Error())
	}
	prefix := util.GetFilePrefixName(st.DirectoryName)
	input := &s3.ListObjectsInput{
		Bucket:  &bucket,
		MaxKeys: st.PageSize,
		Marker:  &st.Marker,
		Prefix:  &prefix,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return nil, fmt.Errorf("list bucket[%s] fail, err: %s", st.DirectoryName, err.Error())
	}
	out, err := client.ListObjects(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("list bucket[%s] fail, err: %s", st.DirectoryName, err.Error())
	}
	resp := &file.ListResp{}
	resp.IsTruncated = out.IsTruncated
	marker := ""
	for _, v := range out.Contents {
		file := &file.FilesInfo{}
		file.FileName = *v.Key
		file.Size = v.Size
		file.LastModified = v.LastModified.String()
		resp.Files = append(resp.Files, file)
		marker = *v.Key
	}
	resp.Marker = marker
	return resp, nil
}

// Del object in aws oss.
func (a *AwsOss) Del(ctx context.Context, st *file.DelRequest) error {
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return fmt.Errorf("awsoss put file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return fmt.Errorf("awsoss put file[%s] fail,err: %s", st.FileName, err.Error())
	}
	input := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return err
	}
	_, err = client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
func (a *AwsOss) Stat(ctx context.Context, st *file.FileMetaRequest) (*file.FileMetaResp, error) {
	bucket, err := util.GetBucketName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("awsoss stat file[%s] fail,err: %s", st.FileName, err.Error())
	}
	key, err := util.GetFileName(st.FileName)
	if err != nil {
		return nil, fmt.Errorf("awsoss stat file[%s] fail,err: %s", st.FileName, err.Error())
	}
	input := &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	out, err := client.HeadObject(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "no such key") {
			return nil, file.ErrNotExist
		}
		return nil, fmt.Errorf("awsoss stat file[%s] fail,err: %s", st.FileName, err.Error())
	}
	resp := &file.FileMetaResp{}
	resp.Size = out.ContentLength
	resp.LastModified = out.LastModified.String()
	resp.Metadata = make(map[string][]string)
	resp.Metadata[util.ETag] = append(resp.Metadata[util.ETag], *out.ETag)
	for k, v := range out.Metadata {
		resp.Metadata[k] = append(resp.Metadata[k], v)
	}
	return resp, nil
}

func NewAwsOss() file.Oss {
	return &AwsOss{
		client: make(map[string]*s3.Client),
		meta:   make(map[string]*AwsOssMetaData),
	}
}

func (a *AwsOss) InitConfig(ctx context.Context, config *file.FileConfig) error {
	a.method = config.Method
	a.rawData = config.Metadata
	return nil
}

func (a *AwsOss) InitClient(ctx context.Context, req *file.InitRequest) error {
	initFunc := factory.GetInitFunc(a.method)
	clients, err := initFunc(a.rawData, req.Metadata)
	if err != nil {
		return err
	}
	for k, v := range clients {
		a.client[k] = v.(*s3.Client)
	}
	return nil
}

func (a *AwsOss) GetObject(ctx context.Context, req *file.GetObjectInput) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	client, err := a.selectClient(map[string]string{}, req.ClientName)
	if err != nil {
		return nil, err
	}
	ob, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return ob.Body, nil
}

func (a *AwsOss) PutObject(context.Context) error {
	return nil
}
