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
	"errors"
	"io"

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
)

var (
	ErrMissingBucket    error = errors.New("missing bucket info in metadata")
	ErrClientNotExist   error = errors.New("specific client not exist")
	ErrEndPointNotExist error = errors.New("specific endpoing key not exist")
)

// AwsOss is a binding for aws oss storage
type AwsOss struct {
	client map[string]*s3.Client
	meta   map[string]*AwsOssMetaData
}

// AwsOssMetaData describe a aws-oss instance
type AwsOssMetaData struct {
	Region          string `json:"region"`   // eg. us-west-2
	EndPoint        string `json:"endpoint"` // eg. protocol://service-code.region-code.amazonaws.com
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
}

func NewAwsOss() file.File {
	return &AwsOss{
		client: make(map[string]*s3.Client),
		meta:   make(map[string]*AwsOssMetaData),
	}
}

// Init instance by config
func (a *AwsOss) Init(config *file.FileConfig) error {
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

// isAwsMetaValid check if the metadata valid
func (am *AwsOssMetaData) isAwsMetaValid() bool {
	if am.AccessKeySecret == "" || am.EndPoint == "" || am.AccessKeyID == "" || am.Region == "" {
		return false
	}
	return true
}

// createOssClient by input meta info
func (a *AwsOss) createOssClient(meta *AwsOssMetaData) (*s3.Client, error) {
	optFunc := []func(options *aws_config.LoadOptions) error{
		aws_config.WithRegion(meta.Region),
		aws_config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: meta.AccessKeyID, SecretAccessKey: meta.AccessKeySecret,
				Source: defaultCredentialsSource,
			},
		}),
	}
	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), optFunc...)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

// Put file to aws oss
func (a *AwsOss) Put(st *file.PutFileStu) error {
	var (
		bucket string
		key    = st.FileName
	)
	if b, ok := st.Metadata[bucketKey]; ok {
		bucket = b
	} else {
		return ErrMissingBucket
	}
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   st.DataStream,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return err
	}
	_, err = client.PutObject(context.TODO(), input, nil)
	if err != nil {
		return err
	}
	return nil
}

// selectClient choose aws client from exist client-map, key is endpoint, value is client instance
func (a *AwsOss) selectClient(meta map[string]string) (*s3.Client, error) {
	if ep, ok := meta[endpointKey]; ok {
		if client, ok := a.client[ep]; ok {
			return client, nil
		}
		return nil, ErrClientNotExist
	}
	return nil, ErrEndPointNotExist
}

// Get object from aws oss
func (a *AwsOss) Get(st *file.GetFileStu) (io.ReadCloser, error) {
	var (
		bucket string
		key    = st.FileName
	)
	if b, ok := st.Metadata[bucketKey]; ok {
		bucket = b
	} else {
		return nil, ErrMissingBucket
	}
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	ob, err := client.GetObject(context.TODO(), input, nil)
	if err != nil {
		return nil, err
	}
	return ob.Body, nil
}

// List objects from aws oss
func (a *AwsOss) List(st *file.ListRequest) (*file.ListResp, error) {
	var bucket string
	if b, ok := st.Metadata[bucketKey]; ok {
		bucket = b
	} else {
		return nil, ErrMissingBucket
	}
	input := &s3.ListObjectsInput{
		Bucket: &bucket,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	out, err := client.ListObjects(context.TODO(), input, nil)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(out.Contents))
	for _, content := range out.Contents {
		ret = append(ret, *content.Key)
	}
	return &file.ListResp{FilesName: ret}, nil
}

// Del object in aws oss
func (a *AwsOss) Del(st *file.DelRequest) error {
	var (
		bucket string
		key    = st.FileName
	)
	if b, ok := st.Metadata[bucketKey]; ok {
		bucket = b
	} else {
		return ErrMissingBucket
	}
	input := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	client, err := a.selectClient(st.Metadata)
	if err != nil {
		return err
	}
	_, err = client.DeleteObject(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
