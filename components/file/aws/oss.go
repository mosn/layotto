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
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	m := make([]*AwsOssMetaData, 0)
	err := json.Unmarshal(staticConf, &m)
	clients := make(map[string]interface{})
	if err != nil {
		return nil, errors.New("invalid config for aws oss")
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
		clients[data.EndPoint] = s3.NewFromConfig(cfg)
	}
	return clients, nil
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

func (a *AwsOss) GetObject(ctx context.Context, req *file.GetObjectInput) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	client, err := a.selectClient(map[string]string{}, "")
	if err != nil {
		return nil, err
	}
	ob, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return ob.Body, nil
}

func (a *AwsOss) PutObject(ctx context.Context, req *file.PutObjectInput) (*file.PutObjectOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	uploader := manager.NewUploader(client)
	resp, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
		Body:   req.DataStream,
	})
	if err != nil {
		return nil, err
	}
	return &file.PutObjectOutput{BucketKeyEnabled: resp.BucketKeyEnabled, ETag: *resp.ETag}, err
}

func (a *AwsOss) DeleteObject(ctx context.Context, req *file.DeleteObjectInput) (*file.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	resp, err := client.DeleteObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.DeleteObjectOutput{DeleteMarker: resp.DeleteMarker, RequestCharged: string(resp.RequestCharged), VersionId: *resp.VersionId}, err
}

func (a *AwsOss) PutObjectTagging(ctx context.Context, req *file.PutBucketTaggingInput) (*file.PutBucketTaggingOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.PutObjectTaggingInput{
		Bucket:    &req.Bucket,
		Key:       &req.Key,
		Tagging:   &types.Tagging{},
		VersionId: &req.VersionId,
	}

	for k, v := range req.Tags {
		k, v := k, v
		input.Tagging.TagSet = append(input.Tagging.TagSet, types.Tag{Key: &k, Value: &v})
	}
	_, err = client.PutObjectTagging(ctx, input)
	return &file.PutBucketTaggingOutput{}, err
}
func (a *AwsOss) DeleteObjectTagging(ctx context.Context, req *file.DeleteObjectTaggingInput) (*file.DeleteObjectTaggingOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.DeleteObjectTaggingInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	resp, err := client.DeleteObjectTagging(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.DeleteObjectTaggingOutput{VersionId: *resp.VersionId}, err
}

func (a *AwsOss) GetObjectTagging(ctx context.Context, req *file.GetObjectTaggingInput) (*file.GetObjectTaggingOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.GetObjectTaggingInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
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
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}

	//TODO: should support objects accessed through access points
	copySource := req.CopySource.CopySourceBucket + "/" + req.CopySource.CopySourceKey + "?versionId=" + req.CopySource.CopySourceVersionId
	input := &s3.CopyObjectInput{
		Bucket:     &req.Bucket,
		Key:        &req.Key,
		CopySource: &copySource,
	}
	_, err = client.CopyObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return &file.CopyObjectOutput{}, err
}
func (a *AwsOss) DeleteObjects(ctx context.Context, req *file.DeleteObjectsInput) (*file.DeleteObjectsOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.DeleteObjectsInput{
		Bucket: &req.Bucket,
		Delete: &types.Delete{},
	}
	for _, v := range req.Delete.Objects {
		object := types.ObjectIdentifier{Key: &v.Key, VersionId: &v.VersionId}
		input.Delete.Objects = append(input.Delete.Objects, object)
	}
	_, err = client.DeleteObjects(ctx, input)
	if err != nil {
		return nil, err
	}
	return &file.DeleteObjectsOutput{}, err
}
func (a *AwsOss) ListObjects(ctx context.Context, req *file.ListObjectsInput) (*file.ListObjectsOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	//input := &s3.ListObjectsInput{
	//	Bucket:              &req.Bucket,
	//	Delimiter:           &req.Delimiter,
	//	EncodingType:        types.EncodingType(req.EncodingType),
	//	ExpectedBucketOwner: &req.ExpectedBucketOwner,
	//	Marker:              &req.Marker,
	//	MaxKeys:             req.MaxKeys,
	//	Prefix:              &req.Prefix,
	//	RequestPayer:        types.RequestPayer(req.RequestPayer),
	//}

	input := &s3.ListObjectsInput{}
	copier.Copy(input, req)
	resp, err := client.ListObjects(ctx, input)
	if err != nil {
		return nil, err
	}

	//output := &file.ListObjectsOutput{Delimiter: *(resp.Delimiter), EncodingType: string(resp.EncodingType), IsTruncated: resp.IsTruncated, Marker: *(resp.Marker),
	//	MaxKeys: resp.MaxKeys, Name: *(resp.Name), NextMarker: *(resp.NextMarker), Prefix: *(resp.Prefix),
	//}
	//
	//for _, v := range resp.CommonPrefixes {
	//	output.CommonPrefixes = append(output.CommonPrefixes, *v.Prefix)
	//}
	//
	//for _, v := range resp.Contents {
	//	object := &file.Object{Etag: *v.ETag, Key: *v.Key, Size: v.Size, StorageClass: string(v.StorageClass), Owner: &file.Owner{DisplayName: *v.Owner.DisplayName, ID: *v.Owner.ID}, LastModified: &timestamppb.Timestamp{Seconds: int64(v.LastModified.Second()), Nanos: int32(v.LastModified.Nanosecond())}}
	//	output.Contents = append(output.Contents, object)
	//}
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
	return output, err
}
func (a *AwsOss) GetObjectAcl(ctx context.Context, req *file.GetObjectAclInput) (*file.GetObjectAclOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.GetObjectAclInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	}
	resp, err := client.GetObjectAcl(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.GetObjectAclOutput{Owner: &file.Owner{}, RequestCharged: string(resp.RequestCharged)}
	for _, v := range resp.Grants {
		output.Grants = append(output.Grants, &file.Grant{Grantee: &file.Grantee{DisplayName: *v.Grantee.DisplayName, EmailAddress: *v.Grantee.DisplayName,
			ID: *v.Grantee.ID, Type: string(v.Grantee.Type), URI: *v.Grantee.URI,
		}, Permission: string(v.Permission)})
	}
	return output, err
}
func (a *AwsOss) PutObjectAcl(ctx context.Context, req *file.PutObjectAclInput) (*file.PutObjectAclOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
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
	return &file.PutObjectAclOutput{RequestCharged: string(resp.RequestCharged)}, err
}
func (a *AwsOss) RestoreObject(ctx context.Context, req *file.RestoreObjectInput) (*file.RestoreObjectOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
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
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.CreateMultipartUploadInput{}
	err = copier.CopyWithOption(input, req, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: int64(0),
				DstType: &time.Time{},
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(int64)
					if !ok {
						return nil, errors.New("src type not matching")
					}
					t := time.Unix(0, s)
					return &t, nil
				},
			},
		}})
	if err != nil {
		log.DefaultLogger.Errorf("copy CreateMultipartUploadInput fail, err: %+v", err)
		return nil, err
	}
	resp, err := client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.CreateMultipartUploadOutput{}
	copier.CopyWithOption(output, resp, copier.Option{
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
		},
	})
	return output, err
}
func (a *AwsOss) UploadPart(ctx context.Context, req *file.UploadPartInput) (*file.UploadPartOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.UploadPartInput{
		Body:                 req.DataStream,
		Bucket:               &req.Bucket,
		Key:                  &req.Key,
		ContentLength:        req.ContentLength,
		ContentMD5:           &req.ContentMd5,
		ExpectedBucketOwner:  &req.ExpectedBucketOwner,
		PartNumber:           req.PartNumber,
		RequestPayer:         types.RequestPayer(req.RequestPayer),
		SSECustomerAlgorithm: &req.SseCustomerAlgorithm,
		SSECustomerKey:       &req.SseCustomerKey,
		SSECustomerKeyMD5:    &req.SseCustomerKeyMd5,
		UploadId:             &req.UploadId,
	}
	resp, err := client.UploadPart(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.UploadPartOutput{BucketKeyEnabled: resp.BucketKeyEnabled, ETag: *resp.ETag, RequestCharged: string(resp.RequestCharged),
		SseCustomerAlgorithm: *resp.SSECustomerAlgorithm, SseCustomerKeyMd5: *resp.SSECustomerKeyMD5, SseKmsKeyId: *resp.SSECustomerKeyMD5,
		ServerSideEncryption: string(resp.ServerSideEncryption),
	}
	return output, err
}
func (a *AwsOss) UploadPartCopy(ctx context.Context, req *file.UploadPartCopyInput) (*file.UploadPartCopyOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}

	//TODO: should support objects accessed through access points
	copySource := req.CopySource.CopySourceBucket + "/" + req.CopySource.CopySourceKey + "?versionId=" + req.CopySource.CopySourceVersionId
	input := &s3.UploadPartCopyInput{
		Bucket:     &req.Bucket,
		Key:        &req.Key,
		CopySource: &copySource,
		PartNumber: req.PartNumber,
		UploadId:   &req.UploadId,
	}
	resp, err := client.UploadPartCopy(ctx, input)
	if err != nil {
		return nil, err
	}
	//LastModified := &timestamppb.Timestamp{Seconds: int64(resp.CopyPartResult.LastModified.Second()), Nanos: int32(resp.CopyPartResult.LastModified.Nanosecond())}
	output := &file.UploadPartCopyOutput{BucketKeyEnabled: resp.BucketKeyEnabled, RequestCharged: string(resp.RequestCharged),
		//CopyPartResult:       &file.CopyPartResult{ETag: *resp.CopyPartResult.ETag, LastModified: LastModified},
		CopySourceVersionId:  *resp.CopySourceVersionId,
		SseCustomerAlgorithm: *resp.SSECustomerAlgorithm, SseCustomerKeyMd5: *resp.SSECustomerKeyMD5, SseKmsKeyId: *resp.SSECustomerKeyMD5,
		ServerSideEncryption: string(resp.ServerSideEncryption),
	}
	return output, err
}
func (a *AwsOss) CompleteMultipartUpload(ctx context.Context, req *file.CompleteMultipartUploadInput) (*file.CompleteMultipartUploadOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:              &req.Bucket,
		Key:                 &req.Key,
		UploadId:            &req.UploadId,
		RequestPayer:        types.RequestPayer(req.RequestPayer),
		ExpectedBucketOwner: &req.ExpectedBucketOwner,
		MultipartUpload:     &types.CompletedMultipartUpload{},
	}
	for _, v := range req.MultipartUpload.Parts {
		input.MultipartUpload.Parts = append(input.MultipartUpload.Parts, types.CompletedPart{ETag: &v.ETag, PartNumber: v.PartNumber})
	}
	resp, err := client.CompleteMultipartUpload(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.CompleteMultipartUploadOutput{
		Bucket:               *resp.Bucket,
		Key:                  *resp.Key,
		BucketKeyEnabled:     resp.BucketKeyEnabled,
		ETag:                 *resp.ETag,
		Expiration:           *resp.Expiration,
		Location:             *resp.Location,
		RequestCharged:       string(resp.RequestCharged),
		SseKmsKeyId:          *resp.SSEKMSKeyId,
		ServerSideEncryption: string(resp.ServerSideEncryption),
		VersionId:            *resp.VersionId,
	}
	return output, err
}
func (a *AwsOss) AbortMultipartUpload(ctx context.Context, req *file.AbortMultipartUploadInput) (*file.AbortMultipartUploadOutput, error) {
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.AbortMultipartUploadInput{
		Bucket:              &req.Bucket,
		Key:                 &req.Key,
		UploadId:            &req.UploadId,
		RequestPayer:        types.RequestPayer(req.RequestPayer),
		ExpectedBucketOwner: &req.ExpectedBucketOwner,
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
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.ListMultipartUploadsInput{
		Bucket:              &req.Bucket,
		ExpectedBucketOwner: &req.ExpectedBucketOwner,
	}
	resp, err := client.ListMultipartUploads(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.ListMultipartUploadsOutput{
		Bucket:             *resp.Bucket,
		Delimiter:          *resp.Delimiter,
		EncodingType:       string(resp.EncodingType),
		IsTruncated:        resp.IsTruncated,
		KeyMarker:          *resp.KeyMarker,
		MaxUploads:         resp.MaxUploads,
		NextKeyMarker:      *resp.NextKeyMarker,
		NextUploadIdMarker: *resp.NextUploadIdMarker,
		Prefix:             *resp.Prefix,
		UploadIdMarker:     *resp.UploadIdMarker,
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
	client, err := a.selectClient(map[string]string{}, endpointKey)
	if err != nil {
		return nil, err
	}
	input := &s3.ListObjectVersionsInput{
		Bucket:              &req.Bucket,
		Delimiter:           &req.Delimiter,
		EncodingType:        types.EncodingType(req.EncodingType),
		ExpectedBucketOwner: &req.ExpectedBucketOwner,
		KeyMarker:           &req.KeyMarker,
		MaxKeys:             req.MaxKeys,
		Prefix:              &req.Prefix,
		VersionIdMarker:     &req.VersionIdMarker,
	}
	resp, err := client.ListObjectVersions(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &file.ListObjectVersionsOutput{
		Delimiter:           *resp.Delimiter,
		EncodingType:        string(resp.EncodingType),
		IsTruncated:         resp.IsTruncated,
		KeyMarker:           *resp.KeyMarker,
		MaxKeys:             resp.MaxKeys,
		Name:                *resp.Name,
		NextKeyMarker:       *resp.NextKeyMarker,
		NextVersionIdMarker: *resp.NextVersionIdMarker,
		Prefix:              *resp.Prefix,
		VersionIdMarker:     *resp.VersionIdMarker,
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
			ETag:     *v.ETag,
			IsLatest: v.IsLatest,
			Key:      *v.Key,
			//LastModified: timestamppb.New(*v.LastModified),
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
	return nil, nil
}

func (a *AwsOss) IsObjectExist(context.Context, *file.IsObjectExistInput) (*file.IsObjectExistOutput, error) {
	return nil, nil
}

func (a *AwsOss) SignURL(ctx context.Context, req *file.SignURLInput) (*file.SignURLOutput, error) {
	return nil, nil
}

func (a *AwsOss) UpdateDownLoadBandwidthRateLimit(ctx context.Context, req *file.UpdateBandwidthRateLimitInput) error {
	return nil
}

func (a *AwsOss) UpdateUpLoadBandwidthRateLimit(ctx context.Context, req *file.UpdateBandwidthRateLimitInput) error {
	return nil
}
func (a *AwsOss) AppendObject(ctx context.Context, req *file.AppendObjectInput) (*file.AppendObjectOutput, error) {
	return nil, nil
}

func (a *AwsOss) ListParts(ctx context.Context, req *file.ListPartsInput) (*file.ListPartsOutput, error) {
	return nil, nil
}
