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
	"strconv"

	"mosn.io/layotto/components/pkg/utils"

	l8oss "mosn.io/layotto/components/oss"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	connectTimeoutSec   = "connectTimeoutSec"
	readWriteTimeoutSec = "readWriteTimeout"
)

type AliyunOSS struct {
	client    *oss.Client
	basicConf json.RawMessage
}

func NewAliyunOss() l8oss.Oss {
	return &AliyunOSS{}
}

func (a *AliyunOSS) Init(ctx context.Context, config *l8oss.Config) error {
	connectTimeout, readWriteTimeout := l8oss.DefaultConnectTimeout, l8oss.DefaultReadWriteTimeout
	a.basicConf = config.Metadata[l8oss.BasicConfiguration]
	m := utils.OssMetadata{}
	if err := json.Unmarshal(a.basicConf, &m); err != nil {
		return l8oss.ErrInvalid
	}
	if t, ok := config.Metadata[connectTimeoutSec]; ok {
		if v, err := strconv.Atoi(string(t)); err == nil {
			connectTimeout = v
		}
	}
	if t, ok := config.Metadata[readWriteTimeoutSec]; ok {
		if v, err := strconv.Atoi(string(t)); err == nil {
			readWriteTimeout = v
		}
	}

	client, err := oss.New(m.Endpoint, m.AccessKeyID, m.AccessKeySecret, oss.Timeout(int64(connectTimeout), int64(readWriteTimeout)))
	if err != nil {
		return err
	}
	a.client = client
	return nil
}

func (a *AliyunOSS) GetObject(ctx context.Context, req *l8oss.GetObjectInput) (*l8oss.GetObjectOutput, error) {
	client, err := a.getClient()
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
		return &l8oss.GetObjectOutput{DataStream: body}, err
	}
	body, err := bucket.GetObject(req.Key,
		IfUnmodifiedSince(req.IfUnmodifiedSince),
		IfModifiedSince(req.IfModifiedSince),
		IfMatch(req.IfMatch),
		IfNoneMatch(req.IfNoneMatch),
		Range(req.Start, req.End),
		AcceptEncoding(req.AcceptEncoding),
	)

	return &l8oss.GetObjectOutput{DataStream: body}, err
}

func (a *AliyunOSS) PutObject(ctx context.Context, req *l8oss.PutObjectInput) (*l8oss.PutObjectOutput, error) {
	cli, err := a.getClient()
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
	return &l8oss.PutObjectOutput{}, err
}

func (a *AliyunOSS) DeleteObject(ctx context.Context, req *l8oss.DeleteObjectInput) (*l8oss.DeleteObjectOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.DeleteObject(req.Key, RequestPayer(req.RequestPayer), VersionId(req.VersionId))
	return &l8oss.DeleteObjectOutput{}, err
}
func (a *AliyunOSS) DeleteObjects(ctx context.Context, req *l8oss.DeleteObjectsInput) (*l8oss.DeleteObjectsOutput, error) {
	cli, err := a.getClient()
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
	out := &l8oss.DeleteObjectsOutput{}
	for _, v := range resp.DeletedObjectsDetail {
		object := &l8oss.DeletedObject{Key: v.Key, VersionId: v.VersionId, DeleteMarker: v.DeleteMarker, DeleteMarkerVersionId: v.DeleteMarkerVersionId}
		out.Deleted = append(out.Deleted, object)
	}
	return out, err
}

func (a *AliyunOSS) PutObjectTagging(ctx context.Context, req *l8oss.PutObjectTaggingInput) (*l8oss.PutObjectTaggingOutput, error) {
	cli, err := a.getClient()
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

func (a *AliyunOSS) DeleteObjectTagging(ctx context.Context, req *l8oss.DeleteObjectTaggingInput) (*l8oss.DeleteObjectTaggingOutput, error) {
	cli, err := a.getClient()
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

func (a *AliyunOSS) GetObjectTagging(ctx context.Context, req *l8oss.GetObjectTaggingInput) (*l8oss.GetObjectTaggingOutput, error) {
	cli, err := a.getClient()
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
	out := &l8oss.GetObjectTaggingOutput{Tags: map[string]string{}}
	for _, v := range resp.Tags {
		out.Tags[v.Key] = v.Value
	}
	return out, err
}

func (a *AliyunOSS) GetObjectCannedAcl(ctx context.Context, req *l8oss.GetObjectCannedAclInput) (*l8oss.GetObjectCannedAclOutput, error) {
	cli, err := a.getClient()
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
	output := &l8oss.GetObjectCannedAclOutput{CannedAcl: resp.ACL, Owner: &l8oss.Owner{DisplayName: resp.Owner.DisplayName, ID: resp.Owner.ID}}
	return output, err
}
func (a *AliyunOSS) PutObjectCannedAcl(ctx context.Context, req *l8oss.PutObjectCannedAclInput) (*l8oss.PutObjectCannedAclOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.SetObjectACL(req.Key, oss.ACLType(req.Acl))
	output := &l8oss.PutObjectCannedAclOutput{}
	return output, err
}
func (a *AliyunOSS) ListObjects(ctx context.Context, req *l8oss.ListObjectsInput) (*l8oss.ListObjectsOutput, error) {
	cli, err := a.getClient()
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
	out := &l8oss.ListObjectsOutput{
		CommonPrefixes: resp.CommonPrefixes,
		Delimiter:      resp.Delimiter,
		IsTruncated:    resp.IsTruncated,
		Marker:         resp.Marker,
		MaxKeys:        int32(resp.MaxKeys),
		NextMarker:     resp.NextMarker,
		Prefix:         resp.Prefix,
	}
	for _, v := range resp.Objects {
		object := &l8oss.Object{
			ETag:         v.ETag,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner:        &l8oss.Owner{ID: v.Owner.ID, DisplayName: v.Owner.DisplayName},
			Size:         v.Size,
			StorageClass: v.StorageClass,
		}
		out.Contents = append(out.Contents, object)
	}
	return out, nil
}
func (a *AliyunOSS) CopyObject(ctx context.Context, req *l8oss.CopyObjectInput) (*l8oss.CopyObjectOutput, error) {
	cli, err := a.getClient()
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
	out := &l8oss.CopyObjectOutput{CopyObjectResult: &l8oss.CopyObjectResult{ETag: resp.ETag, LastModified: resp.LastModified.Unix()}}
	return out, err
}

func (a *AliyunOSS) CreateMultipartUpload(ctx context.Context, req *l8oss.CreateMultipartUploadInput) (*l8oss.CreateMultipartUploadOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.InitiateMultipartUpload(req.Key)
	output := &l8oss.CreateMultipartUploadOutput{Bucket: resp.Bucket, Key: resp.Key, UploadId: resp.UploadID}
	return output, err
}
func (a *AliyunOSS) UploadPart(ctx context.Context, req *l8oss.UploadPartInput) (*l8oss.UploadPartOutput, error) {
	cli, err := a.getClient()
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
	output := &l8oss.UploadPartOutput{ETag: resp.ETag}
	return output, err
}
func (a *AliyunOSS) UploadPartCopy(ctx context.Context, req *l8oss.UploadPartCopyInput) (*l8oss.UploadPartCopyOutput, error) {
	cli, err := a.getClient()
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
	output := &l8oss.UploadPartCopyOutput{CopyPartResult: &l8oss.CopyPartResult{ETag: resp.ETag}}
	return output, err
}
func (a *AliyunOSS) CompleteMultipartUpload(ctx context.Context, req *l8oss.CompleteMultipartUploadInput) (*l8oss.CompleteMultipartUploadOutput, error) {
	cli, err := a.getClient()
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
	output := &l8oss.CompleteMultipartUploadOutput{Location: resp.Location, Bucket: resp.Bucket, Key: resp.Key, ETag: resp.ETag}
	return output, err
}
func (a *AliyunOSS) AbortMultipartUpload(ctx context.Context, req *l8oss.AbortMultipartUploadInput) (*l8oss.AbortMultipartUploadOutput, error) {
	cli, err := a.getClient()
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
	output := &l8oss.AbortMultipartUploadOutput{}
	return output, err
}
func (a *AliyunOSS) ListMultipartUploads(ctx context.Context, req *l8oss.ListMultipartUploadsInput) (*l8oss.ListMultipartUploadsOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListMultipartUploads(Prefix(req.Prefix), KeyMarker(req.KeyMarker), MaxUploads(int(req.MaxUploads)), Delimiter(req.Delimiter), UploadIDMarker(req.UploadIdMarker))
	output := &l8oss.ListMultipartUploadsOutput{
		Bucket:             resp.Bucket,
		Delimiter:          resp.Delimiter,
		Prefix:             resp.Prefix,
		KeyMarker:          resp.KeyMarker,
		UploadIDMarker:     resp.UploadIDMarker,
		NextKeyMarker:      resp.NextKeyMarker,
		NextUploadIDMarker: resp.NextUploadIDMarker,
		MaxUploads:         int32(resp.MaxUploads),
		IsTruncated:        resp.IsTruncated,
		CommonPrefixes:     resp.CommonPrefixes,
	}
	for _, v := range resp.Uploads {
		upload := &l8oss.MultipartUpload{Initiated: v.Initiated.Unix(), UploadId: v.UploadID, Key: v.Key}
		output.Uploads = append(output.Uploads, upload)
	}
	return output, err
}

func (a *AliyunOSS) RestoreObject(ctx context.Context, req *l8oss.RestoreObjectInput) (*l8oss.RestoreObjectOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	err = bucket.RestoreObject(req.Key)
	output := &l8oss.RestoreObjectOutput{}
	return output, err
}

func (a *AliyunOSS) ListObjectVersions(ctx context.Context, req *l8oss.ListObjectVersionsInput) (*l8oss.ListObjectVersionsOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.ListObjectVersions()
	output := &l8oss.ListObjectVersionsOutput{
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
		marker := &l8oss.DeleteMarkerEntry{
			IsLatest:     v.IsLatest,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &l8oss.Owner{
				ID:          v.Owner.ID,
				DisplayName: v.Owner.DisplayName,
			},
			VersionId: v.VersionId,
		}
		output.DeleteMarkers = append(output.DeleteMarkers, marker)
	}

	for _, v := range resp.ObjectVersions {
		version := &l8oss.ObjectVersion{
			ETag:         v.ETag,
			IsLatest:     v.IsLatest,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &l8oss.Owner{
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

func (a *AliyunOSS) HeadObject(ctx context.Context, req *l8oss.HeadObjectInput) (*l8oss.HeadObjectOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	output := &l8oss.HeadObjectOutput{ResultMetadata: map[string]string{}}
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

func (a *AliyunOSS) IsObjectExist(ctx context.Context, req *l8oss.IsObjectExistInput) (*l8oss.IsObjectExistOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.IsObjectExist(req.Key)
	return &l8oss.IsObjectExistOutput{FileExist: resp}, err
}

func (a *AliyunOSS) SignURL(ctx context.Context, req *l8oss.SignURLInput) (*l8oss.SignURLOutput, error) {
	cli, err := a.getClient()
	if err != nil {
		return nil, err
	}
	bucket, err := cli.Bucket(req.Bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bucket.SignURL(req.Key, oss.HTTPMethod(req.Method), req.ExpiredInSec)
	return &l8oss.SignURLOutput{SignedUrl: resp}, err
}

// UpdateDownloadBandwidthRateLimit update all client rate
func (a *AliyunOSS) UpdateDownloadBandwidthRateLimit(ctx context.Context, req *l8oss.UpdateBandwidthRateLimitInput) error {
	cli, err := a.getClient()
	if err != nil {
		return err
	}
	err = cli.LimitDownloadSpeed(int(req.AverageRateLimitInBitsPerSec))
	return err
}

// UpdateUploadBandwidthRateLimit update all client rate
func (a *AliyunOSS) UpdateUploadBandwidthRateLimit(ctx context.Context, req *l8oss.UpdateBandwidthRateLimitInput) error {
	cli, err := a.getClient()
	if err != nil {
		return err
	}
	err = cli.LimitUploadSpeed(int(req.AverageRateLimitInBitsPerSec))
	return err
}

func (a *AliyunOSS) AppendObject(ctx context.Context, req *l8oss.AppendObjectInput) (*l8oss.AppendObjectOutput, error) {
	cli, err := a.getClient()
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
	return &l8oss.AppendObjectOutput{AppendPosition: resp}, err
}

func (a *AliyunOSS) ListParts(ctx context.Context, req *l8oss.ListPartsInput) (*l8oss.ListPartsOutput, error) {
	cli, err := a.getClient()
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
	out := &l8oss.ListPartsOutput{
		Bucket:               resp.Bucket,
		Key:                  resp.Key,
		UploadId:             resp.UploadID,
		NextPartNumberMarker: resp.NextPartNumberMarker,
		MaxParts:             int64(resp.MaxParts),
		IsTruncated:          resp.IsTruncated,
	}
	for _, v := range resp.UploadedParts {
		part := &l8oss.Part{Etag: v.ETag, LastModified: v.LastModified.Unix(), PartNumber: int64(v.PartNumber), Size: int64(v.Size)}
		out.Parts = append(out.Parts, part)
	}
	return out, err
}

func (a *AliyunOSS) getClient() (*oss.Client, error) {
	if a.client == nil {
		return nil, utils.ErrNotInitClient
	}
	return a.client, nil
}
