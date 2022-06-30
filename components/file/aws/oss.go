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
	"net/url"
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

	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/file/factory"
)

const (
	DefaultClientInitFunc = "aws"
)

func init() {
	factory.RegisterInitFunc(DefaultClientInitFunc, AwsDefaultInitFunc)
}

func AwsDefaultInitFunc(staticConf json.RawMessage, DynConf map[string]string) (map[string]interface{}, error) {
	m := make([]*file.OssMetadata, 0)
	err := json.Unmarshal(staticConf, &m)
	clients := make(map[string]interface{})
	if err != nil {
		return nil, file.ErrInvalid
	}
	for _, data := range m {
		optFunc := []func(options *aws_config.LoadOptions) error{
			aws_config.WithRegion(data.Region),
			aws_config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: data.AccessKeyID, SecretAccessKey: data.AccessKeySecret,
					Source: defaultCredentialsSource,
				},
			}),
		}

		cfg, err := aws_config.LoadDefaultConfig(context.TODO(), optFunc...)
		if err != nil {
			return nil, err
		}
		client := s3.NewFromConfig(cfg)
		clients[data.Uid] = client
		for _, bucketName := range data.Buckets {
			if _, ok := clients[bucketName]; ok {
				continue
			}
			clients[bucketName] = client
		}
	}
	return clients, nil
}

func NewAwsOss() file.Oss {
	return &AwsOss{
		client: make(map[string]*s3.Client),
		meta:   make(map[string]*file.OssMetadata),
	}
}
func (a *AwsOss) InitConfig(ctx context.Context, config *file.FileConfig) error {
	a.method = config.Method
	a.rawData = config.Metadata
	return nil
}

func (a *AwsOss) InitClient(ctx context.Context, req *file.InitRequest) error {
	if a.method == "" {
		a.method = DefaultClientInitFunc
	}
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

func (a *AwsOss) GetObject(ctx context.Context, req *file.GetObjectInput) (*file.GetObjectOutput, error) {
	input := &s3.GetObjectInput{}
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	ob, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	out := &file.GetObjectOutput{}
	err = copier.Copy(out, ob)
	if err != nil {
		return nil, err
	}
	out.DataStream = ob.Body
	return out, nil
}

func (a *AwsOss) PutObject(ctx context.Context, req *file.PutObjectInput) (*file.PutObjectOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	out := &file.PutObjectOutput{}
	err = copier.Copy(out, resp)
	if err != nil {
		return nil, err
	}
	return out, err
}

func (a *AwsOss) DeleteObject(ctx context.Context, req *file.DeleteObjectInput) (*file.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := client.DeleteObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.DeleteObjectOutput{DeleteMarker: resp.DeleteMarker, RequestCharged: string(resp.RequestCharged), VersionId: *resp.VersionId}, err
}

func (a *AwsOss) PutObjectTagging(ctx context.Context, req *file.PutObjectTaggingInput) (*file.PutObjectTaggingOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	return &file.PutObjectTaggingOutput{}, err
}
func (a *AwsOss) DeleteObjectTagging(ctx context.Context, req *file.DeleteObjectTaggingInput) (*file.DeleteObjectTaggingOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	return &file.DeleteObjectTaggingOutput{VersionId: *resp.VersionId}, err
}

func (a *AwsOss) GetObjectTagging(ctx context.Context, req *file.GetObjectTaggingInput) (*file.GetObjectTaggingOutput, error) {
	client, err := a.selectClient(req.Bucket)
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

	output := &file.GetObjectTaggingOutput{Tags: map[string]string{}}
	for _, tags := range resp.TagSet {
		output.Tags[*tags.Key] = *tags.Value
	}
	return output, err
}

func (a *AwsOss) CopyObject(ctx context.Context, req *file.CopyObjectInput) (*file.CopyObjectOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}

	if req.CopySource == nil {
		return nil, errors.New("must specific copy_source")
	}

	//TODO: should support objects accessed through access points
	copySource := req.CopySource.CopySourceBucket + "/" + req.CopySource.CopySourceKey
	if req.CopySource.CopySourceVersionId != "" {
		copySource += "?versionId=" + req.CopySource.CopySourceVersionId
	}
	copySourceUrlEncode := url.QueryEscape(copySource)
	input := &s3.CopyObjectInput{
		Bucket:     &req.Bucket,
		Key:        &req.Key,
		CopySource: &copySourceUrlEncode,
	}
	resp, err := client.CopyObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.CopyObjectOutput{CopyObjectResult: &file.CopyObjectResult{ETag: *resp.CopyObjectResult.ETag, LastModified: resp.CopyObjectResult.LastModified.Unix()}}, err
}
func (a *AwsOss) DeleteObjects(ctx context.Context, req *file.DeleteObjectsInput) (*file.DeleteObjectsOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	output := &file.DeleteObjectsOutput{}
	copier.Copy(output, resp)
	return output, err
}
func (a *AwsOss) ListObjects(ctx context.Context, req *file.ListObjectsInput) (*file.ListObjectsOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectsInput{}
	err = copier.CopyWithOption(
		input,
		req,
		copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters:  []copier.TypeConverter{},
		},
	)
	if err != nil {
		return nil, err
	}
	resp, err := client.ListObjects(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.ListObjectsOutput{}
	err = copier.CopyWithOption(output, resp, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: &time.Time{},
				DstType: int64(0),
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(*time.Time)
					if !ok {
						return nil, errors.New("src type not matching")
					}
					return s.Unix(), nil
				},
			},
		}})
	// if not return NextMarker, use the value of the last Key in the response as the marker
	if output.IsTruncated && output.NextMarker == "" {
		index := len(output.Contents) - 1
		output.NextMarker = output.Contents[index].Key
	}
	return output, err
}
func (a *AwsOss) GetObjectCannedAcl(ctx context.Context, req *file.GetObjectCannedAclInput) (*file.GetObjectCannedAclOutput, error) {
	return nil, errors.New("GetObjectCannedAcl method not supported on AWS")
}
func (a *AwsOss) PutObjectCannedAcl(ctx context.Context, req *file.PutObjectCannedAclInput) (*file.PutObjectCannedAclOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	input := &s3.PutObjectAclInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
		ACL:    types.ObjectCannedACL(req.Acl),
	}
	resp, err := client.PutObjectAcl(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.PutObjectCannedAclOutput{RequestCharged: string(resp.RequestCharged)}, err
}
func (a *AwsOss) RestoreObject(ctx context.Context, req *file.RestoreObjectInput) (*file.RestoreObjectOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	input := &s3.RestoreObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	resp, err := client.RestoreObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.RestoreObjectOutput{RequestCharged: string(resp.RequestCharged), RestoreOutputPath: *resp.RestoreOutputPath}, err
}
func (a *AwsOss) CreateMultipartUpload(ctx context.Context, req *file.CreateMultipartUploadInput) (*file.CreateMultipartUploadOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	input := &s3.CreateMultipartUploadInput{}
	err = copier.CopyWithOption(
		input,
		req,
		copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters:  []copier.TypeConverter{int642time},
		},
	)
	if err != nil {
		log.DefaultLogger.Errorf("copy CreateMultipartUploadInput fail, err: %+v", err)
		return nil, err
	}
	resp, err := client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.CreateMultipartUploadOutput{}
	copier.CopyWithOption(
		output,
		resp,
		copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters:  []copier.TypeConverter{time2int64},
		},
	)
	return output, err
}
func (a *AwsOss) UploadPart(ctx context.Context, req *file.UploadPartInput) (*file.UploadPartOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	input := &s3.UploadPartInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	input.Body = req.DataStream
	resp, err := client.UploadPart(
		ctx,
		input,
		s3.WithAPIOptions(
			v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
		),
	)
	if err != nil {
		return nil, err
	}
	output := &file.UploadPartOutput{}
	err = copier.Copy(output, resp)
	if err != nil {
		return nil, err
	}
	return output, err
}
func (a *AwsOss) UploadPartCopy(ctx context.Context, req *file.UploadPartCopyInput) (*file.UploadPartCopyOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}

	//TODO: should support objects accessed through access points
	copySource := req.CopySource.CopySourceBucket + "/" + req.CopySource.CopySourceKey
	if req.CopySource.CopySourceVersionId != "" {
		copySource += "?versionId=" + req.CopySource.CopySourceVersionId
	}
	input := &s3.UploadPartCopyInput{}
	err = copier.CopyWithOption(input, req, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{}})
	if err != nil {
		return nil, err
	}
	input.CopySource = &copySource
	resp, err := client.UploadPartCopy(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.UploadPartCopyOutput{}
	err = copier.Copy(output, resp)
	return output, err
}
func (a *AwsOss) CompleteMultipartUpload(ctx context.Context, req *file.CompleteMultipartUploadInput) (*file.CompleteMultipartUploadOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	input := &s3.CompleteMultipartUploadInput{MultipartUpload: &types.CompletedMultipartUpload{}}
	err = copier.CopyWithOption(
		input,
		req,
		copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
			Converters:  []copier.TypeConverter{},
		},
	)
	if err != nil {
		return nil, err
	}
	resp, err := client.CompleteMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.CompleteMultipartUploadOutput{}
	err = copier.Copy(output, resp)
	return output, err
}
func (a *AwsOss) AbortMultipartUpload(ctx context.Context, req *file.AbortMultipartUploadInput) (*file.AbortMultipartUploadOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	output := &file.AbortMultipartUploadOutput{
		RequestCharged: string(resp.RequestCharged),
	}
	return output, err
}
func (a *AwsOss) ListMultipartUploads(ctx context.Context, req *file.ListMultipartUploadsInput) (*file.ListMultipartUploadsOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	output := &file.ListMultipartUploadsOutput{CommonPrefixes: []string{}, Uploads: []*file.MultipartUpload{}}
	err = copier.Copy(output, resp)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.CommonPrefixes {
		output.CommonPrefixes = append(output.CommonPrefixes, *v.Prefix)
	}
	for _, v := range resp.Uploads {
		upload := &file.MultipartUpload{
			//Initiated:    timestamppb.New(*v.Initiated),
			Initiator:    &file.Initiator{DisplayName: *v.Initiator.DisplayName, ID: *v.Initiator.ID},
			Key:          *v.Key,
			Owner:        &file.Owner{ID: *v.Owner.ID, DisplayName: *v.Owner.DisplayName},
			StorageClass: string(v.StorageClass),
			UploadId:     *v.UploadId,
		}
		output.Uploads = append(output.Uploads, upload)
	}
	return output, err
}
func (a *AwsOss) ListObjectVersions(ctx context.Context, req *file.ListObjectVersionsInput) (*file.ListObjectVersionsOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	output := &file.ListObjectVersionsOutput{}
	err = copier.Copy(output, resp)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.CommonPrefixes {
		output.CommonPrefixes = append(output.CommonPrefixes, *v.Prefix)
	}
	for _, v := range resp.DeleteMarkers {
		entry := &file.DeleteMarkerEntry{
			IsLatest: v.IsLatest,
			Key:      *v.Key,
			//LastModified: timestamppb.New(*v.LastModified),
			Owner:     &file.Owner{DisplayName: *v.Owner.DisplayName, ID: *v.Owner.ID},
			VersionId: *v.VersionId,
		}
		output.DeleteMarkers = append(output.DeleteMarkers, entry)
	}
	for _, v := range resp.Versions {
		version := &file.ObjectVersion{
			ETag:         *v.ETag,
			IsLatest:     v.IsLatest,
			Key:          *v.Key,
			LastModified: v.LastModified.Unix(),
			Owner:        &file.Owner{DisplayName: *v.Owner.DisplayName, ID: *v.Owner.ID},
			Size:         v.Size,
			StorageClass: string(v.StorageClass),
			VersionId:    *v.VersionId,
		}
		output.Versions = append(output.Versions, version)
	}
	return output, err
}

func (a *AwsOss) HeadObject(ctx context.Context, req *file.HeadObjectInput) (*file.HeadObjectOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
	return &file.HeadObjectOutput{ResultMetadata: resp.Metadata}, nil
}

func (a *AwsOss) IsObjectExist(ctx context.Context, req *file.IsObjectExistInput) (*file.IsObjectExistOutput, error) {
	client, err := a.selectClient(req.Bucket)
	if err != nil {
		return nil, err
	}
	input := &s3.HeadObjectInput{Bucket: &req.Bucket, Key: &req.Key}
	_, err = client.HeadObject(ctx, input)
	if err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "StatusCode: 404") {
			return &file.IsObjectExistOutput{FileExist: false}, nil
		}
		return nil, err
	}
	return &file.IsObjectExistOutput{FileExist: true}, nil
}

func (a *AwsOss) SignURL(ctx context.Context, req *file.SignURLInput) (*file.SignURLOutput, error) {
	client, err := a.selectClient(req.Bucket)
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
		return &file.SignURLOutput{SignedUrl: resp.URL}, nil
	case "PUT":
		input := &s3.PutObjectInput{Bucket: &req.Bucket, Key: &req.Key}
		resp, err := resignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(time.Duration(req.ExpiredInSec*int64(time.Second))))
		if err != nil {
			return nil, err
		}
		return &file.SignURLOutput{SignedUrl: resp.URL}, nil
	default:
		return nil, fmt.Errorf("not supported method %+v now", req.Method)
	}
}

func (a *AwsOss) UpdateDownLoadBandwidthRateLimit(ctx context.Context, req *file.UpdateBandwidthRateLimitInput) error {
	return errors.New("UpdateDownLoadBandwidthRateLimit method not supported now")
}

func (a *AwsOss) UpdateUpLoadBandwidthRateLimit(ctx context.Context, req *file.UpdateBandwidthRateLimitInput) error {
	return errors.New("UpdateUpLoadBandwidthRateLimit method not supported now")
}
func (a *AwsOss) AppendObject(ctx context.Context, req *file.AppendObjectInput) (*file.AppendObjectOutput, error) {
	return nil, errors.New("AppendObject method not supported on AWS")
}

func (a *AwsOss) ListParts(ctx context.Context, req *file.ListPartsInput) (*file.ListPartsOutput, error) {
	return nil, errors.New("ListParts method not supported on AWS")
}
