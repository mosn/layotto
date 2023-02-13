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
package huaweicloud

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/jinzhu/copier"

	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/pkg/utils"
)

const connectTimeoutSec = "connectTimeoutSec"

type HuaweicloudOSS struct {
	client   *obs.ObsClient
	metadata utils.OssMetadata
}

func NewHuaweicloudOSS() oss.Oss {
	return &HuaweicloudOSS{}
}

func (h *HuaweicloudOSS) Init(ctx context.Context, config *oss.Config) error {
	connectTimeout := oss.DefaultConnectTimeout
	jsonRawMessage := config.Metadata[oss.BasicConfiguration]
	err := json.Unmarshal(jsonRawMessage, &h.metadata)
	if err != nil {
		return oss.ErrInvalid
	}
	if t, ok := config.Metadata[connectTimeoutSec]; ok {
		if v, err := strconv.Atoi(string(t)); err != nil {
			connectTimeout = v
		}
	}

	client, err := obs.New(h.metadata.AccessKeyID, h.metadata.AccessKeySecret, h.metadata.Endpoint, obs.WithConnectTimeout(connectTimeout))
	if err != nil {
		return err
	}
	h.client = client
	return nil
}

func (h *HuaweicloudOSS) GetObject(ctx context.Context, input *oss.GetObjectInput) (*oss.GetObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.GetObjectInput{}
	if err = copier.CopyWithOption(obsInput, input, copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{oss.Int64ToTime}}); err != nil {
		return nil, err
	}
	metadataInput := &obs.GetObjectMetadataInput{}
	copier.Copy(metadataInput, input)
	obsInput.GetObjectMetadataInput = *metadataInput
	obsInput.RangeStart = input.Start
	obsInput.RangeEnd = input.End

	obsOutput, err := client.GetObject(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.GetObjectOutput{}
	if err = copier.CopyWithOption(output, obsOutput, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
		return nil, err
	}
	output.DataStream = obsOutput.Body

	return output, nil
}

func (h *HuaweicloudOSS) PutObject(ctx context.Context, input *oss.PutObjectInput) (*oss.PutObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.PutObjectInput{}
	obsInput.Body = input.DataStream
	basicInput := &obs.PutObjectBasicInput{}
	if err = copier.Copy(basicInput, input); err != nil {
		return nil, err
	}
	obsInput.PutObjectBasicInput = *basicInput
	operationInput := &obs.ObjectOperationInput{}
	if err = copier.CopyWithOption(operationInput, input, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	basicInput.ObjectOperationInput = *operationInput

	obsOutput, err := client.PutObject(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.PutObjectOutput{}
	if err = copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}

	return output, nil
}

func (h *HuaweicloudOSS) DeleteObject(ctx context.Context, input *oss.DeleteObjectInput) (*oss.DeleteObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.DeleteObjectInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}

	obsOutput, err := client.DeleteObject(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.DeleteObjectOutput{}
	if err = copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}

	return output, nil
}

func (h *HuaweicloudOSS) PutObjectTagging(ctx context.Context, input *oss.PutObjectTaggingInput) (*oss.PutObjectTaggingOutput, error) {
	return nil, ErrHaveNotTag
}

func (h *HuaweicloudOSS) DeleteObjectTagging(ctx context.Context, input *oss.DeleteObjectTaggingInput) (*oss.DeleteObjectTaggingOutput, error) {
	return nil, ErrHaveNotTag
}

func (h *HuaweicloudOSS) GetObjectTagging(ctx context.Context, input *oss.GetObjectTaggingInput) (*oss.GetObjectTaggingOutput, error) {
	return nil, ErrHaveNotTag
}

func (h *HuaweicloudOSS) CopyObject(ctx context.Context, input *oss.CopyObjectInput) (*oss.CopyObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.CopyObjectInput{}
	if err = copier.CopyWithOption(obsInput, input.CopySource, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	operationInput := &obs.ObjectOperationInput{}
	if err = copier.CopyWithOption(operationInput, input, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	obsInput.ObjectOperationInput = *operationInput

	obsOutput, err := client.CopyObject(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.CopyObjectOutput{}
	if err = copier.CopyWithOption(output, obsOutput, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
		return nil, err
	}

	return output, nil
}

func (h *HuaweicloudOSS) DeleteObjects(ctx context.Context, input *oss.DeleteObjectsInput) (*oss.DeleteObjectsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.DeleteObjectsInput{}
	obsInput.Bucket = input.Bucket
	obsInput.Quiet = input.Delete.Quiet
	objects := make([]obs.ObjectToDelete, 0, len(input.Delete.Objects))
	for _, v := range input.Delete.Objects {
		object := &obs.ObjectToDelete{}
		if err = copier.Copy(object, v); err != nil {
			return nil, err
		}
		objects = append(objects, *object)
	}
	obsInput.Objects = objects

	obsOutput, err := client.DeleteObjects(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.DeleteObjectsOutput{Deleted: make([]*oss.DeletedObject, 0, len(obsOutput.Deleteds))}
	for _, v := range obsOutput.Deleteds {
		deleteObject := &oss.DeletedObject{}
		if err = copier.Copy(deleteObject, v); err != nil {
			return nil, err
		}
		output.Deleted = append(output.Deleted, deleteObject)
	}

	return output, nil
}

func (h *HuaweicloudOSS) ListObjects(ctx context.Context, input *oss.ListObjectsInput) (*oss.ListObjectsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.ListObjectsInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	objsInput := &obs.ListObjsInput{}
	if err = copier.Copy(objsInput, input); err != nil {
		return nil, err
	}
	obsInput.ListObjsInput = *objsInput

	obsOutput, err := client.ListObjects(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.ListObjectsOutput{}
	if err = copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}
	contents := make([]*oss.Object, 0, len(obsOutput.Contents))
	for _, v := range obsOutput.Contents {
		content := &oss.Object{}
		if err = copier.CopyWithOption(content, v, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
			return nil, err
		}
		owner := &oss.Owner{}
		if err = copier.Copy(owner, v.Owner); err != nil {
			return nil, err
		}
		content.Owner = owner
		contents = append(contents, content)
	}
	output.Contents = contents

	return output, nil
}

func (h *HuaweicloudOSS) GetObjectCannedAcl(ctx context.Context, input *oss.GetObjectCannedAclInput) (*oss.GetObjectCannedAclOutput, error) {
	return nil, ErrNotSupportAclGet
}

func (h *HuaweicloudOSS) PutObjectCannedAcl(ctx context.Context, input *oss.PutObjectCannedAclInput) (*oss.PutObjectCannedAclOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.SetObjectAclInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	obsInput.ACL = obs.AclType(input.Acl)

	_, err = client.SetObjectAcl(obsInput)
	if err != nil {
		return nil, err
	}
	output := &oss.PutObjectCannedAclOutput{}
	return output, nil
}

func (h *HuaweicloudOSS) RestoreObject(ctx context.Context, input *oss.RestoreObjectInput) (*oss.RestoreObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.RestoreObjectInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	obsInput.Days = int(input.RestoreRequest.Days)
	obsInput.Tier = obs.RestoreTierType(input.RestoreRequest.Tier)

	_, err = client.RestoreObject(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.RestoreObjectOutput{}

	return output, nil
}

func (h *HuaweicloudOSS) CreateMultipartUpload(ctx context.Context, input *oss.CreateMultipartUploadInput) (*oss.CreateMultipartUploadOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.InitiateMultipartUploadInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	operationInput := &obs.ObjectOperationInput{}
	if err = copier.Copy(operationInput, input); err != nil {
		return nil, err
	}
	operationInput.GrantReadId = input.GrantRead
	operationInput.GrantReadAcpId = input.GrantReadACP
	operationInput.GrantWriteAcpId = input.GrantWriteACP
	operationInput.GrantFullControlId = input.GrantFullControl

	obsOutput, err := client.InitiateMultipartUpload(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.CreateMultipartUploadOutput{}
	if err = copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}

	return output, nil
}

func (h *HuaweicloudOSS) UploadPart(ctx context.Context, input *oss.UploadPartInput) (*oss.UploadPartOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput := &obs.UploadPartInput{}
	if err = copier.CopyWithOption(obsInput, input, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	obsInput.Body = input.DataStream
	obsOutput, err := client.UploadPart(obsInput)
	if err != nil {
		return nil, err
	}
	output := &oss.UploadPartOutput{}
	if err = copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweicloudOSS) UploadPartCopy(ctx context.Context, input *oss.UploadPartCopyInput) (*oss.UploadPartCopyOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.CopyPartInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	if err = copier.Copy(obsInput, input.CopySource); err != nil {
		return nil, err
	}

	obsOutput, err := client.CopyPart(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.UploadPartCopyOutput{}
	partResult := &oss.CopyPartResult{}
	if err = copier.CopyWithOption(partResult, obsOutput, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
		return nil, err
	}
	output.CopyPartResult = partResult
	return output, nil
}

func (h *HuaweicloudOSS) CompleteMultipartUpload(ctx context.Context, input *oss.CompleteMultipartUploadInput) (*oss.CompleteMultipartUploadOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput := &obs.CompleteMultipartUploadInput{}
	copier.CopyWithOption(obsInput, input, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if input.MultipartUpload != nil && len(input.MultipartUpload.Parts) != 0 {
		obsInput.Parts = make([]obs.Part, 0, len(input.MultipartUpload.Parts))
		for _, v := range input.MultipartUpload.Parts {
			part := &obs.Part{}
			copier.CopyWithOption(part, v, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			obsInput.Parts = append(obsInput.Parts, *part)
		}
	}
	obsOutput, err := client.CompleteMultipartUpload(obsInput)
	if err != nil {
		return nil, err
	}
	output := &oss.CompleteMultipartUploadOutput{}
	copier.CopyWithOption(output, obsOutput, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return output, nil
}

func (h *HuaweicloudOSS) AbortMultipartUpload(ctx context.Context, input *oss.AbortMultipartUploadInput) (*oss.AbortMultipartUploadOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput := &obs.AbortMultipartUploadInput{}
	copier.Copy(obsInput, input)
	_, err = client.AbortMultipartUpload(obsInput)
	if err != nil {
		return nil, err
	}
	output := &oss.AbortMultipartUploadOutput{}
	return output, nil
}

func (h *HuaweicloudOSS) ListMultipartUploads(ctx context.Context, input *oss.ListMultipartUploadsInput) (*oss.ListMultipartUploadsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.ListMultipartUploadsInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}

	obsOutput, err := client.ListMultipartUploads(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.ListMultipartUploadsOutput{}
	if err = copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}
	uploads := make([]*oss.MultipartUpload, 0, len(obsOutput.Uploads))
	for _, v := range obsOutput.Uploads {
		upload := &oss.MultipartUpload{}
		if err = copier.CopyWithOption(upload, v, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
			return nil, err
		}
		initiator := &oss.Initiator{}
		copier.Copy(initiator, v.Initiator)
		upload.Initiator = initiator
		owner := &oss.Owner{}
		if err = copier.Copy(owner, v.Owner); err != nil {
			return nil, err
		}
		upload.Owner = owner
		uploads = append(uploads, upload)
	}
	output.Uploads = uploads

	return output, nil
}

func (h *HuaweicloudOSS) ListObjectVersions(ctx context.Context, input *oss.ListObjectVersionsInput) (*oss.ListObjectVersionsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.ListVersionsInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	objsInput := &obs.ListObjsInput{}
	if err = copier.Copy(objsInput, input); err != nil {
		return nil, err
	}
	obsInput.ListObjsInput = *objsInput

	obsOutput, err := client.ListVersions(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.ListObjectVersionsOutput{}
	if err = copier.CopyWithOption(output, obsOutput, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	deleteMarkers := make([]*oss.DeleteMarkerEntry, 0, len(obsOutput.DeleteMarkers))
	for _, v := range obsOutput.DeleteMarkers {
		entry := &oss.DeleteMarkerEntry{}
		if err = copier.CopyWithOption(entry, v, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
			return nil, err
		}
		deleteMarkers = append(deleteMarkers, entry)
	}
	output.DeleteMarkers = deleteMarkers
	versions := make([]*oss.ObjectVersion, 0, len(obsOutput.Versions))
	for _, v := range obsOutput.Versions {
		version := &oss.ObjectVersion{}
		if err = copier.CopyWithOption(version, v, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}}); err != nil {
			return nil, err
		}
		owner := &oss.Owner{}
		if err = copier.Copy(owner, v.Owner); err != nil {
			return nil, err
		}
		version.Owner = owner
		versions = append(versions, version)
	}
	output.Versions = versions

	return output, nil
}

func (h *HuaweicloudOSS) HeadObject(ctx context.Context, input *oss.HeadObjectInput) (*oss.HeadObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	if !input.WithDetails {
		obsInput := &obs.HeadObjectInput{}
		if err = copier.Copy(obsInput, input); err != nil {
			return nil, err
		}
		obsBaseModel, err := client.HeadObject(obsInput)
		if err != nil {
			return nil, err
		}
		output := &oss.HeadObjectOutput{ResultMetadata: make(map[string]string, len(obsBaseModel.ResponseHeaders))}
		for k, v := range obsBaseModel.ResponseHeaders {
			for _, t := range v {
				if _, ok := output.ResultMetadata[k]; ok {
					output.ResultMetadata[k] = output.ResultMetadata[t] + "," + t
				} else {
					output.ResultMetadata[k] = t
				}
			}
		}
		return output, nil
	}

	obsInput := &obs.GetObjectMetadataInput{}
	copier.CopyWithOption(obsInput, input, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	obsOutput, err := client.GetObjectMetadata(obsInput)
	if err != nil {
		return nil, err
	}
	output := &oss.HeadObjectOutput{ResultMetadata: map[string]string{}}
	copier.CopyWithOption(output.ResultMetadata, obsOutput.Metadata, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	jsonByte, err := json.Marshal(obsOutput)
	if err != nil {
		return nil, err
	}
	valueMap := make(map[string]interface{})
	if err := json.Unmarshal(jsonByte, &valueMap); err != nil {
		return nil, err
	}
	for k, v := range valueMap {
		if value, ok := v.(string); ok {
			output.ResultMetadata[k] = value
		}
	}
	return output, nil
}

func (h *HuaweicloudOSS) IsObjectExist(ctx context.Context, input *oss.IsObjectExistInput) (*oss.IsObjectExistOutput, error) {
	listObjectsInput := &oss.ListObjectsInput{Bucket: input.Bucket}
	listObjectsOutput, err := h.ListObjects(ctx, listObjectsInput)
	if err != nil {
		return nil, err
	}
	isExist := false
	for _, v := range listObjectsOutput.Contents {
		if v == nil {
			continue
		}
		if v.Key == input.Key {
			isExist = true
			break
		}
	}
	output := &oss.IsObjectExistOutput{FileExist: isExist}
	return output, nil
}

func (h *HuaweicloudOSS) SignURL(ctx context.Context, input *oss.SignURLInput) (*oss.SignURLOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.CreateSignedUrlInput{}
	if err = copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	obsInput.Expires = int(input.ExpiredInSec)

	obsOutput, err := client.CreateSignedUrl(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.SignURLOutput{SignedUrl: obsOutput.SignedUrl}

	return output, nil
}

func (h *HuaweicloudOSS) UpdateDownloadBandwidthRateLimit(ctx context.Context, input *oss.UpdateBandwidthRateLimitInput) error {
	return ErrDownloadNotBandwidthLimit
}

func (h *HuaweicloudOSS) UpdateUploadBandwidthRateLimit(ctx context.Context, input *oss.UpdateBandwidthRateLimitInput) error {
	return ErrUploadNotBandwidthLimit
}

func (h *HuaweicloudOSS) AppendObject(ctx context.Context, input *oss.AppendObjectInput) (*oss.AppendObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput := &obs.AppendObjectInput{}
	obsInput.Position = input.Position
	obsInput.Body = input.DataStream
	basicInput := &obs.PutObjectBasicInput{}
	basicInput.ContentMD5 = input.ContentMd5
	basicInput.ContentEncoding = input.ContentEncoding
	operationInput := &obs.ObjectOperationInput{}
	if err = copier.Copy(operationInput, input); err != nil {
		return nil, err
	}
	basicInput.ObjectOperationInput = *operationInput
	obsInput.PutObjectBasicInput = *basicInput

	if err != nil {
		return nil, err
	}
	obsOutput, err := client.AppendObject(obsInput)
	if err != nil {
		return nil, err
	}

	output := &oss.AppendObjectOutput{AppendPosition: obsOutput.NextAppendPosition}

	return output, nil
}

// todo 测试异常
func (h *HuaweicloudOSS) ListParts(ctx context.Context, input *oss.ListPartsInput) (*oss.ListPartsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput := &obs.ListPartsInput{}
	if err := copier.Copy(obsInput, input); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.ListParts(obsInput)
	if err != nil {
		return nil, err
	}
	output := &oss.ListPartsOutput{}
	if err := copier.Copy(output, obsOutput); err != nil {
		return nil, err
	}
	parts := make([]*oss.Part, 0, len(obsOutput.Parts))
	for _, v := range obsOutput.Parts {
		part := &oss.Part{}
		copier.Copy(part, v)
		parts = append(parts, part)
	}
	output.Parts = parts
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweicloudOSS) getClient() (*obs.ObsClient, error) {
	if h.client == nil {
		return nil, utils.ErrNotInitClient
	}
	return h.client, nil
}
