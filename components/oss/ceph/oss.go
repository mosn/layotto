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

package ceph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jinzhu/copier"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/pkg/utils"
)

type CephOSS struct {
	client    *s3.Client
	basicConf json.RawMessage
}

func NewCephOss() oss.Oss {
	return &CephOSS{}
}

func (c *CephOSS) Init(ctx context.Context, config *oss.Config) error {
	c.basicConf = config.Metadata[oss.BasicConfiguration]
	m := &utils.OssMetadata{}
	err := json.Unmarshal(c.basicConf, &m)
	if err != nil {
		return oss.ErrInvalid
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: m.Endpoint,
		}, nil
	})
	optFunc := []func(options *aws_config.LoadOptions) error{
		aws_config.WithRegion(m.Region),
		aws_config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: m.AccessKeyID, SecretAccessKey: m.AccessKeySecret,
				Source: "provider",
			},
		}),
		aws_config.WithEndpointResolverWithOptions(customResolver),
	}
	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), optFunc...)
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	c.client = client
	return nil
}

func (c *CephOSS) GetObject(ctx context.Context, req *oss.GetObjectInput) (*oss.GetObjectOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.GetObjectInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	ob, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return oss.GetGetObjectOutput(ob)
}

func (c *CephOSS) PutObject(ctx context.Context, req *oss.PutObjectInput) (*oss.PutObjectOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.PutObjectInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	input.Body = req.DataStream
	uploader := manager.NewUploader(client)
	resp, err := uploader.Upload(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return oss.GetPutObjectOutput(resp)
}

func (c *CephOSS) DeleteObject(ctx context.Context, req *oss.DeleteObjectInput) (*oss.DeleteObjectOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.DeleteObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	resp, err := client.DeleteObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetDeleteObjectOutput(resp)
}

func (c *CephOSS) PutObjectTagging(ctx context.Context, req *oss.PutObjectTaggingInput) (*oss.PutObjectTaggingOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.PutObjectTaggingInput{Tagging: &types.Tagging{}}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	for k, v := range req.Tags {
		k, v := k, v
		input.Tagging.TagSet = append(input.Tagging.TagSet, types.Tag{Key: &k, Value: &v})
	}
	_, err = client.PutObjectTagging(ctx, input)

	return &oss.PutObjectTaggingOutput{}, err
}

func (c *CephOSS) DeleteObjectTagging(ctx context.Context, req *oss.DeleteObjectTaggingInput) (*oss.DeleteObjectTaggingOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.DeleteObjectTaggingInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.DeleteObjectTagging(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetDeleteObjectTaggingOutput(resp)
}

func (c *CephOSS) GetObjectTagging(ctx context.Context, req *oss.GetObjectTaggingInput) (*oss.GetObjectTaggingOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.GetObjectTaggingInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.GetObjectTagging(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetGetObjectTaggingOutput(resp)
}

func (c *CephOSS) CopyObject(ctx context.Context, req *oss.CopyObjectInput) (*oss.CopyObjectOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	if req.CopySource == nil {
		return nil, errors.New("must specific copy_source")
	}

	input := &s3.CopyObjectInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.Int64ToTime}})
	if err != nil {
		return nil, err
	}
	copySource := req.CopySource.CopySourceBucket + "/" + req.CopySource.CopySourceKey
	if req.CopySource.CopySourceVersionId != "" {
		copySource += "?versionId=" + req.CopySource.CopySourceVersionId
	}
	input.CopySource = &copySource
	resp, err := client.CopyObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return oss.GetCopyObjectOutput(resp)
}

func (c *CephOSS) DeleteObjects(ctx context.Context, req *oss.DeleteObjectsInput) (*oss.DeleteObjectsOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.DeleteObjectsInput{
		Bucket: &req.Bucket,
		Delete: &types.Delete{},
	}
	if req.Delete != nil {
		for _, v := range req.Delete.Objects {
			object := &types.ObjectIdentifier{}
			err = copier.CopyWithOption(object, v, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
			if err != nil {
				return nil, err
			}
			input.Delete.Objects = append(input.Delete.Objects, *object)
		}
	}
	resp, err := client.DeleteObjects(ctx, input)
	if err != nil {
		return nil, err
	}

	output := &oss.DeleteObjectsOutput{}
	copier.Copy(output, resp)
	return output, err
}

func (c *CephOSS) ListObjects(ctx context.Context, req *oss.ListObjectsInput) (*oss.ListObjectsOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectsInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.ListObjects(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetListObjectsOutput(resp)
}

func (c *CephOSS) GetObjectCannedAcl(ctx context.Context, req *oss.GetObjectCannedAclInput) (*oss.GetObjectCannedAclOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.GetObjectAclInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.GetObjectAcl(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return oss.GetGetObjectCannedAclOutput(resp)
}

func (c *CephOSS) PutObjectCannedAcl(ctx context.Context, req *oss.PutObjectCannedAclInput) (*oss.PutObjectCannedAclOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	input := &s3.PutObjectAclInput{Bucket: &req.Bucket, Key: &req.Key, ACL: types.ObjectCannedACL(req.Acl)}
	resp, err := client.PutObjectAcl(ctx, input)
	if err != nil {
		return nil, err
	}
	return &oss.PutObjectCannedAclOutput{RequestCharged: string(resp.RequestCharged)}, err
}

func (c *CephOSS) CreateMultipartUpload(ctx context.Context, req *oss.CreateMultipartUploadInput) (*oss.CreateMultipartUploadOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.CreateMultipartUploadInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.Int64ToTime}})
	if err != nil {
		log.DefaultLogger.Errorf("copy CreateMultipartUploadInput fail, err: %+v", err)
		return nil, err
	}
	resp, err := client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}

	output := &oss.CreateMultipartUploadOutput{}
	copier.CopyWithOption(output, resp, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}})
	return output, err
}

func (c *CephOSS) UploadPart(ctx context.Context, req *oss.UploadPartInput) (*oss.UploadPartOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.UploadPartInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	input.Body = req.DataStream
	resp, err := client.UploadPart(ctx, input, s3.WithAPIOptions(v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware))
	if err != nil {
		return nil, err
	}

	return oss.GetUploadPartOutput(resp)
}

func (c *CephOSS) UploadPartCopy(ctx context.Context, req *oss.UploadPartCopyInput) (*oss.UploadPartCopyOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.UploadPartCopyInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	copySource := req.CopySource.CopySourceBucket + "/" + req.CopySource.CopySourceKey
	if req.CopySource.CopySourceVersionId != "" {
		copySource += "?versionId=" + req.CopySource.CopySourceVersionId
	}
	input.CopySource = &copySource
	resp, err := client.UploadPartCopy(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return oss.GetUploadPartCopyOutput(resp)
}

func (c *CephOSS) CompleteMultipartUpload(ctx context.Context, req *oss.CompleteMultipartUploadInput) (*oss.CompleteMultipartUploadOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.CompleteMultipartUploadInput{MultipartUpload: &types.CompletedMultipartUpload{}}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.CompleteMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}

	output := &oss.CompleteMultipartUploadOutput{}
	err = copier.Copy(output, resp)
	return output, err
}

func (c *CephOSS) AbortMultipartUpload(ctx context.Context, req *oss.AbortMultipartUploadInput) (*oss.AbortMultipartUploadOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.AbortMultipartUploadInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.AbortMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}

	output := &oss.AbortMultipartUploadOutput{
		RequestCharged: string(resp.RequestCharged),
	}
	return output, err
}

func (c *CephOSS) ListParts(ctx context.Context, req *oss.ListPartsInput) (*oss.ListPartsOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.ListPartsInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.ListParts(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetListPartsOutput(resp)
}

func (c *CephOSS) ListMultipartUploads(ctx context.Context, req *oss.ListMultipartUploadsInput) (*oss.ListMultipartUploadsOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.ListMultipartUploadsInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.ListMultipartUploads(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetListMultipartUploadsOutput(resp)
}

func (c *CephOSS) ListObjectVersions(ctx context.Context, req *oss.ListObjectVersionsInput) (*oss.ListObjectVersionsOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectVersionsInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.ListObjectVersions(ctx, input)
	if err != nil {
		return nil, err
	}

	return oss.GetListObjectVersionsOutput(resp)
}

func (c *CephOSS) HeadObject(ctx context.Context, req *oss.HeadObjectInput) (*oss.HeadObjectOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	input := &s3.HeadObjectInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	resp, err := client.HeadObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return &oss.HeadObjectOutput{ResultMetadata: resp.Metadata}, nil
}

func (c *CephOSS) IsObjectExist(ctx context.Context, req *oss.IsObjectExistInput) (*oss.IsObjectExistOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	input := &s3.HeadObjectInput{Bucket: &req.Bucket, Key: &req.Key}
	_, err = client.HeadObject(ctx, input)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "StatusCode: 404") {
			return &oss.IsObjectExistOutput{FileExist: false}, nil
		}
		return nil, err
	}
	return &oss.IsObjectExistOutput{FileExist: true}, nil
}

func (c *CephOSS) SignURL(ctx context.Context, req *oss.SignURLInput) (*oss.SignURLOutput, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}
	resignClient := s3.NewPresignClient(client)
	switch strings.ToUpper(req.Method) {
	case "GET":
		input := &s3.GetObjectInput{Bucket: &req.Bucket, Key: &req.Key}
		resp, err := resignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(time.Duration((req.ExpiredInSec)*int64(time.Second))))
		if err != nil {
			return nil, err
		}
		return &oss.SignURLOutput{SignedUrl: resp.URL}, nil
	case "PUT":
		input := &s3.PutObjectInput{Bucket: &req.Bucket, Key: &req.Key}
		resp, err := resignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(time.Duration(req.ExpiredInSec*int64(time.Second))))
		if err != nil {
			return nil, err
		}
		return &oss.SignURLOutput{SignedUrl: resp.URL}, nil
	default:
		return nil, fmt.Errorf("not supported method %+v now", req.Method)
	}
}

func (c *CephOSS) RestoreObject(ctx context.Context, req *oss.RestoreObjectInput) (*oss.RestoreObjectOutput, error) {
	return nil, errors.New("RestoreObject method not supported on CEPH")
}

func (c *CephOSS) UpdateDownloadBandwidthRateLimit(ctx context.Context, req *oss.UpdateBandwidthRateLimitInput) error {
	return errors.New("UpdateDownloadBandwidthRateLimit method not supported now")
}

func (c *CephOSS) UpdateUploadBandwidthRateLimit(ctx context.Context, req *oss.UpdateBandwidthRateLimitInput) error {
	return errors.New("UpdateUploadBandwidthRateLimit method not supported now")
}
func (c *CephOSS) AppendObject(ctx context.Context, req *oss.AppendObjectInput) (*oss.AppendObjectOutput, error) {
	return nil, errors.New("AppendObject method not supported on CEPH")
}

func (c *CephOSS) getClient() (*s3.Client, error) {
	if c.client == nil {
		return nil, utils.ErrNotInitClient
	}
	return c.client, nil
}
