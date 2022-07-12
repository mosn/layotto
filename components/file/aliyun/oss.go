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

package aliyun

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/file/factory"
)

const (
	DefaultClientInitFunc = "aliyun"
)

func NewAliyunOss() file.Oss {
	return &AliyunOSS{
		client: make(map[string]*oss.Client),
	}
}

func init() {
	factory.RegisterInitFunc(DefaultClientInitFunc, AliyunDefaultInitFunc)
}

func AliyunDefaultInitFunc(staticConf json.RawMessage, DynConf map[string]string) (map[string]interface{}, error) {
	m := make([]*file.OssMetadata, 0)
	clients := make(map[string]interface{})
	err := json.Unmarshal(staticConf, &m)
	if err != nil {
		return nil, file.ErrInvalid
	}
	for _, v := range m {
		client, err := oss.New(v.Endpoint, v.AccessKeyID, v.AccessKeySecret)
		if err != nil {
			return nil, err
		}
		clients[v.Uid] = client
		for _, bucketName := range v.Buckets {
			if _, ok := clients[bucketName]; ok {
				continue
			}
			clients[bucketName] = client
		}
	}
	return clients, nil
}

func (a *AliyunOSS) InitConfig(ctx context.Context, config *file.OssConfig) error {
	a.method = config.Method
	a.rawData = config.Metadata
	return nil
}

func (a *AliyunOSS) InitClient(ctx context.Context, req *file.InitRequest) error {
	if a.method == "" {
		a.method = DefaultClientInitFunc
	}
	initFunc := factory.GetInitFunc(a.method)
	clients, err := initFunc(a.rawData, req.Metadata)
	if err != nil {
		return err
	}
	for k, v := range clients {
		a.client[k] = v.(*oss.Client)
	}
	return nil
}

func (a *AliyunOSS) GetObject(ctx context.Context, req *file.GetObjectInput) (*file.GetObjectOutput, error) {
	client, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	//user can use SignedUrl to get file without ak、sk
	if req.SignedUrl != "" {
		body, err := bucket.GetObjectWithURL(req.SignedUrl)
		return &file.GetObjectOutput{DataStream: body}, err
	}
	body, err := bucket.GetObject(req.Key,
		IfUnmodifiedSince(req.IfUnmodifiedSince),
		IfModifiedSince(req.IfModifiedSince),
		IfMatch(req.IfMatch),
		IfNoneMatch(req.IfNoneMatch),
		Range(req.Start, req.End),
		AcceptEncoding(req.AcceptEncoding),
	)

	return &file.GetObjectOutput{DataStream: body}, err
}

func (a *AliyunOSS) PutObject(ctx context.Context, req *file.PutObjectInput) (*file.PutObjectOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	metaOption := []oss.Option{
		CacheControl(req.CacheControl),
		ContentDisposition(req.ContentDisposition),
		ContentEncoding(req.ContentEncoding),
		Expires(req.Expires),
		ServerSideEncryption(req.ServerSideEncryption),
		ObjectACL(req.ACL),
		SetTagging(req.Tagging),
	}
	for k, v := range req.Meta {
		o := oss.Meta(k, v)
		metaOption = append(metaOption, o)
	}
	//user can use SignedUrl to put file without ak、sk
	if req.SignedUrl != "" {
		err = bucket.PutObjectWithURL(req.SignedUrl, req.DataStream,
			metaOption...,
		)
	} else {
		err = bucket.PutObject(req.Key, req.DataStream,
			metaOption...,
		)
	}
	return &file.PutObjectOutput{}, err
}

func (a *AliyunOSS) DeleteObject(ctx context.Context, req *file.DeleteObjectInput) (*file.DeleteObjectOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.DeleteObject(req.Key, RequestPayer(req.RequestPayer), VersionId(req.VersionId))
	return &file.DeleteObjectOutput{}, err
}
func (a *AliyunOSS) DeleteObjects(ctx context.Context, req *file.DeleteObjectsInput) (*file.DeleteObjectsOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	var objects []oss.DeleteObject
	for _, v := range req.Delete.Objects {
		object := oss.DeleteObject{Key: v.Key, VersionId: v.VersionId}
		objects = append(objects, object)
	}
	resp, err := bucket.DeleteObjectVersions(objects, oss.DeleteObjectsQuiet(req.Delete.Quiet))
	if err != nil {
		return nil, err
	}
	out := &file.DeleteObjectsOutput{}
	for _, v := range resp.DeletedObjectsDetail {
		object := &file.DeletedObject{Key: v.Key, VersionId: v.VersionId, DeleteMarker: v.DeleteMarker, DeleteMarkerVersionId: v.DeleteMarkerVersionId}
		out.Deleted = append(out.Deleted, object)
	}
	return out, err
}

func (a *AliyunOSS) PutObjectTagging(ctx context.Context, req *file.PutObjectTaggingInput) (*file.PutObjectTaggingOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	tagging := oss.Tagging{}
	for k, v := range req.Tags {
		tag := oss.Tag{Key: k, Value: v}
		tagging.Tags = append(tagging.Tags, tag)
	}
	err = bucket.PutObjectTagging(req.Key, tagging, VersionId(req.VersionId))
	return nil, err
}

func (a *AliyunOSS) DeleteObjectTagging(ctx context.Context, req *file.DeleteObjectTaggingInput) (*file.DeleteObjectTaggingOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.DeleteObjectTagging(req.Key)
	return nil, err
}

func (a *AliyunOSS) GetObjectTagging(ctx context.Context, req *file.GetObjectTaggingInput) (*file.GetObjectTaggingOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.GetObjectTagging(req.Key)
	if err != nil {
		return nil, err
	}
	out := &file.GetObjectTaggingOutput{Tags: map[string]string{}}
	for _, v := range resp.Tags {
		out.Tags[v.Key] = v.Value
	}
	return out, err
}

func (a *AliyunOSS) GetObjectCannedAcl(ctx context.Context, req *file.GetObjectCannedAclInput) (*file.GetObjectCannedAclOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.GetObjectACL(req.Key)
	if err != nil {
		return nil, err
	}
	output := &file.GetObjectCannedAclOutput{CannedAcl: resp.ACL, Owner: &file.Owner{DisplayName: resp.Owner.DisplayName, ID: resp.Owner.ID}}
	return output, err
}
func (a *AliyunOSS) PutObjectCannedAcl(ctx context.Context, req *file.PutObjectCannedAclInput) (*file.PutObjectCannedAclOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.SetObjectACL(req.Key, oss.ACLType(req.Acl))
	output := &file.PutObjectCannedAclOutput{}
	return output, err
}
func (a *AliyunOSS) ListObjects(ctx context.Context, req *file.ListObjectsInput) (*file.ListObjectsOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListObjects()
	if err != nil {
		return nil, err
	}
	out := &file.ListObjectsOutput{
		CommonPrefixes: resp.CommonPrefixes,
		Delimiter:      resp.Delimiter,
		IsTruncated:    resp.IsTruncated,
		Marker:         resp.Marker,
		MaxKeys:        int32(resp.MaxKeys),
		NextMarker:     resp.NextMarker,
		Prefix:         resp.Prefix,
	}
	for _, v := range resp.Objects {
		object := &file.Object{
			ETag:         v.ETag,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner:        &file.Owner{ID: v.Owner.ID, DisplayName: v.Owner.DisplayName},
			Size:         v.Size,
			StorageClass: v.StorageClass,
		}
		out.Contents = append(out.Contents, object)
	}
	return out, nil
}
func (a *AliyunOSS) CopyObject(ctx context.Context, req *file.CopyObjectInput) (*file.CopyObjectOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	var options []oss.Option
	for k, v := range req.Metadata {
		option := Meta(k, v)
		options = append(options, option)
	}
	options = append(options, MetadataDirective(req.MetadataDirective))
	options = append(options, VersionId(req.CopySource.CopySourceVersionId))
	resp, err := bucket.CopyObject(req.CopySource.CopySourceKey, req.Key, options...)
	if err != nil {
		return nil, err
	}
	out := &file.CopyObjectOutput{CopyObjectResult: &file.CopyObjectResult{ETag: resp.ETag, LastModified: resp.LastModified.Unix()}}
	return out, err
}

func (a *AliyunOSS) CreateMultipartUpload(ctx context.Context, req *file.CreateMultipartUploadInput) (*file.CreateMultipartUploadOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.InitiateMultipartUpload(req.Key)
	output := &file.CreateMultipartUploadOutput{Bucket: resp.Bucket, Key: resp.Key, UploadId: resp.UploadID}
	return output, err
}
func (a *AliyunOSS) UploadPart(ctx context.Context, req *file.UploadPartInput) (*file.UploadPartOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.UploadPart(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		req.DataStream,
		req.ContentLength,
		int(req.PartNumber))
	output := &file.UploadPartOutput{ETag: resp.ETag}
	return output, err
}
func (a *AliyunOSS) UploadPartCopy(ctx context.Context, req *file.UploadPartCopyInput) (*file.UploadPartCopyOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.UploadPartCopy(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		req.CopySource.CopySourceBucket,
		req.CopySource.CopySourceKey,
		req.StartPosition,
		req.PartSize,
		int(req.PartNumber),
		VersionId(req.CopySource.CopySourceVersionId),
	)
	output := &file.UploadPartCopyOutput{CopyPartResult: &file.CopyPartResult{ETag: resp.ETag}}
	return output, err
}
func (a *AliyunOSS) CompleteMultipartUpload(ctx context.Context, req *file.CompleteMultipartUploadInput) (*file.CompleteMultipartUploadOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	parts := make([]oss.UploadPart, 0)
	if req.MultipartUpload != nil {
		for _, v := range req.MultipartUpload.Parts {
			part := oss.UploadPart{PartNumber: int(v.PartNumber), ETag: v.ETag}
			parts = append(parts, part)
		}
	}
	resp, err := bucket.CompleteMultipartUpload(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		parts,
	)
	output := &file.CompleteMultipartUploadOutput{Location: resp.Location, Bucket: resp.Bucket, Key: resp.Key, ETag: resp.ETag}
	return output, err
}
func (a *AliyunOSS) AbortMultipartUpload(ctx context.Context, req *file.AbortMultipartUploadInput) (*file.AbortMultipartUploadOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	err = bucket.AbortMultipartUpload(
		oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
	)
	output := &file.AbortMultipartUploadOutput{}
	return output, err
}
func (a *AliyunOSS) ListMultipartUploads(ctx context.Context, req *file.ListMultipartUploadsInput) (*file.ListMultipartUploadsOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListMultipartUploads(Prefix(req.Prefix), KeyMarker(req.KeyMarker), MaxUploads(int(req.MaxUploads)), Delimiter(req.Delimiter), UploadIDMarker(req.UploadIdMarker))
	output := &file.ListMultipartUploadsOutput{
		Bucket:             resp.Bucket,
		Delimiter:          resp.Delimiter,
		Prefix:             resp.Prefix,
		KeyMarker:          resp.KeyMarker,
		UploadIdMarker:     resp.UploadIDMarker,
		NextKeyMarker:      resp.NextKeyMarker,
		NextUploadIdMarker: resp.NextUploadIDMarker,
		MaxUploads:         int32(resp.MaxUploads),
		IsTruncated:        resp.IsTruncated,
		CommonPrefixes:     resp.CommonPrefixes,
	}
	for _, v := range resp.Uploads {
		upload := &file.MultipartUpload{Initiated: v.Initiated.Unix(), UploadId: v.UploadID, Key: v.Key}
		output.Uploads = append(output.Uploads, upload)
	}
	return output, err
}

func (a *AliyunOSS) RestoreObject(ctx context.Context, req *file.RestoreObjectInput) (*file.RestoreObjectOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.RestoreObject(req.Key)
	output := &file.RestoreObjectOutput{}
	return output, err
}

func (a *AliyunOSS) ListObjectVersions(ctx context.Context, req *file.ListObjectVersionsInput) (*file.ListObjectVersionsOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListObjectVersions()
	output := &file.ListObjectVersionsOutput{
		Name:                resp.Name,
		Prefix:              resp.Prefix,
		KeyMarker:           resp.KeyMarker,
		VersionIdMarker:     resp.VersionIdMarker,
		MaxKeys:             int32(resp.MaxKeys),
		Delimiter:           resp.Delimiter,
		IsTruncated:         resp.IsTruncated,
		NextKeyMarker:       resp.NextKeyMarker,
		NextVersionIdMarker: resp.NextVersionIdMarker,
		CommonPrefixes:      resp.CommonPrefixes,
	}
	for _, v := range resp.ObjectDeleteMarkers {
		marker := &file.DeleteMarkerEntry{
			IsLatest:     v.IsLatest,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &file.Owner{
				ID:          v.Owner.ID,
				DisplayName: v.Owner.DisplayName,
			},
			VersionId: v.VersionId,
		}
		output.DeleteMarkers = append(output.DeleteMarkers, marker)
	}

	for _, v := range resp.ObjectVersions {
		version := &file.ObjectVersion{
			ETag:         v.ETag,
			IsLatest:     v.IsLatest,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &file.Owner{
				ID:          v.Owner.ID,
				DisplayName: v.Owner.DisplayName,
			},
			Size:         v.Size,
			StorageClass: v.StorageClass,
			VersionId:    v.VersionId,
		}
		output.Versions = append(output.Versions, version)
	}

	return output, err
}

func (a *AliyunOSS) HeadObject(ctx context.Context, req *file.HeadObjectInput) (*file.HeadObjectOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	output := &file.HeadObjectOutput{ResultMetadata: map[string]string{}}
	var resp http.Header
	if req.WithDetails {
		resp, err = bucket.GetObjectDetailedMeta(req.Key)
	} else {
		resp, err = bucket.GetObjectMeta(req.Key)
	}
	if err != nil {
		return nil, err
	}
	for k, v := range resp {
		for _, t := range v {
			//if key exist,concatenated with commas
			if _, ok := output.ResultMetadata[k]; ok {
				output.ResultMetadata[k] = output.ResultMetadata[k] + "," + t
			} else {
				output.ResultMetadata[k] = t
			}
		}
	}
	return output, err
}

func (a *AliyunOSS) IsObjectExist(ctx context.Context, req *file.IsObjectExistInput) (*file.IsObjectExistOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.IsObjectExist(req.Key)
	return &file.IsObjectExistOutput{FileExist: resp}, err
}

func (a *AliyunOSS) SignURL(ctx context.Context, req *file.SignURLInput) (*file.SignURLOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.SignURL(req.Key, oss.HTTPMethod(req.Method), req.ExpiredInSec)
	return &file.SignURLOutput{SignedUrl: resp}, err
}

//UpdateDownLoadBandwidthRateLimit update all client rate
func (a *AliyunOSS) UpdateDownLoadBandwidthRateLimit(ctx context.Context, req *file.UpdateBandwidthRateLimitInput) error {
	for _, cli := range a.client {
		err := cli.LimitDownloadSpeed(int(req.AverageRateLimitInBitsPerSec))
		return err
	}
	return nil
}

//UpdateUpLoadBandwidthRateLimit update all client rate
func (a *AliyunOSS) UpdateUpLoadBandwidthRateLimit(ctx context.Context, req *file.UpdateBandwidthRateLimitInput) error {
	for _, cli := range a.client {
		err := cli.LimitUploadSpeed(int(req.AverageRateLimitInBitsPerSec))
		return err
	}
	return nil
}

func (a *AliyunOSS) AppendObject(ctx context.Context, req *file.AppendObjectInput) (*file.AppendObjectOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.AppendObject(req.Key, req.DataStream, req.Position,
		CacheControl(req.CacheControl),
		ContentDisposition(req.ContentDisposition),
		ContentEncoding(req.ContentEncoding),
		Expires(req.Expires),
		ServerSideEncryption(req.ServerSideEncryption),
		ObjectACL(req.ACL),
	)
	if err != nil {
		return nil, err
	}
	return &file.AppendObjectOutput{AppendPosition: resp}, err
}

func (a *AliyunOSS) ListParts(ctx context.Context, req *file.ListPartsInput) (*file.ListPartsOutput, error) {
	cli, err := a.selectClient(req.Uid, req.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListUploadedParts(oss.InitiateMultipartUploadResult{Bucket: req.Bucket, Key: req.Key, UploadID: req.UploadId},
		MaxParts(int(req.MaxParts)),
		PartNumberMarker(int(req.PartNumberMarker)),
		RequestPayer(req.RequestPayer),
	)
	if err != nil {
		return nil, err
	}
	out := &file.ListPartsOutput{
		Bucket:               resp.Bucket,
		Key:                  resp.Key,
		UploadId:             resp.UploadID,
		NextPartNumberMarker: resp.NextPartNumberMarker,
		MaxParts:             int64(resp.MaxParts),
		IsTruncated:          resp.IsTruncated,
	}
	for _, v := range resp.UploadedParts {
		part := &file.Part{Etag: v.ETag, LastModified: v.LastModified.Unix(), PartNumber: int64(v.PartNumber), Size: int64(v.Size)}
		out.Parts = append(out.Parts, part)
	}
	return out, err
}
