package huaweiyun

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/pkg/utils"
	"strconv"
	"time"
)

const connectTimeoutSec = "connectTimeoutSec"

type HuaweiyunOSS struct {
	client   *obs.ObsClient
	metadata utils.OssMetadata
}

func (h *HuaweiyunOSS) Init(ctx context.Context, config *oss.Config) error {
	connectTimeout := 30
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

func (h *HuaweiyunOSS) GetObject(ctx context.Context, input *oss.GetObjectInput) (*oss.GetObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}

	obsInput, err := h.convertGetObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.GetObject(obsInput)
	if err != nil {
		return nil, err
	}

	output, err := h.convertGetObjectOutput(obsOutput)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (h *HuaweiyunOSS) convertGetObjectInput(input *oss.GetObjectInput) (*obs.GetObjectInput, error) {
	obsInput := &obs.GetObjectInput{
		GetObjectMetadataInput: obs.GetObjectMetadataInput{
			Bucket:        input.Bucket,
			Key:           input.Key,
			VersionId:     input.VersionId,
			Origin:        "",
			RequestHeader: "",
			SseHeader:     nil,
		},
		IfMatch:                    input.IfMatch,
		IfNoneMatch:                input.IfNoneMatch,
		IfUnmodifiedSince:          time.Unix(input.IfUnmodifiedSince, 0),
		IfModifiedSince:            time.Unix(input.IfModifiedSince, 0),
		RangeStart:                 input.Start,
		RangeEnd:                   input.End,
		ImageProcess:               "",
		ResponseCacheControl:       input.ResponseCacheControl,
		ResponseContentDisposition: input.ResponseContentDisposition,
		ResponseContentEncoding:    input.ResponseContentEncoding,
		ResponseContentLanguage:    input.ResponseContentLanguage,
		ResponseContentType:        input.ResponseContentType,
		ResponseExpires:            input.ResponseExpires,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertGetObjectOutput(obsOutput *obs.GetObjectOutput) (*oss.GetObjectOutput, error) {
	output := &oss.GetObjectOutput{
		DataStream:         obsOutput.Body,
		CacheControl:       obsOutput.CacheControl,
		ContentDisposition: obsOutput.ContentDisposition,
		ContentEncoding:    obsOutput.ContentEncoding,
		ContentLanguage:    obsOutput.ContentLanguage,
		ContentLength:      obsOutput.ContentLength,
		ContentRange:       "",
		ContentType:        obsOutput.ContentType,
		DeleteMarker:       obsOutput.DeleteMarker,
		Etag:               obsOutput.ETag,
		Expiration:         obsOutput.Expiration,
		Expires:            obsOutput.Expires,
		LastModified:       obsOutput.LastModified.Unix(),
		VersionId:          obsOutput.VersionId,
		TagCount:           0,
		StorageClass:       string(obsOutput.StorageClass),
		PartsCount:         0,
		Metadata:           obsOutput.Metadata,
	}
	return output, nil
}

func (h *HuaweiyunOSS) PutObject(ctx context.Context, input *oss.PutObjectInput) (*oss.PutObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertPutObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.PutObject(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertPutObjectOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertPutObjectInput(input *oss.PutObjectInput) (*obs.PutObjectInput, error) {
	obsInput := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket:                  input.Bucket,
				Key:                     input.Key,
				ACL:                     obs.AclType(input.ACL),
				GrantReadId:             "",
				GrantReadAcpId:          "",
				GrantWriteAcpId:         "",
				GrantFullControlId:      "",
				StorageClass:            obs.StorageClassType(input.StorageClass),
				WebsiteRedirectLocation: "",
				Expires:                 input.Expires,
				SseHeader:               nil,
				Metadata:                input.Meta,
			},
			ContentType:     "",
			ContentMD5:      "",
			ContentLength:   input.ContentLength,
			ContentEncoding: input.ContentEncoding,
		},
		Body: input.DataStream,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertPutObjectOutput(obsOutput *obs.PutObjectOutput) (*oss.PutObjectOutput, error) {
	output := &oss.PutObjectOutput{
		BucketKeyEnabled: false,
		ETag:             obsOutput.ETag,
	}
	return output, nil
}

func (h *HuaweiyunOSS) DeleteObject(ctx context.Context, input *oss.DeleteObjectInput) (*oss.DeleteObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertDeleteObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.DeleteObject(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertDeleteObjectOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertDeleteObjectInput(input *oss.DeleteObjectInput) (*obs.DeleteObjectInput, error) {
	obsInput := &obs.DeleteObjectInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionId: input.VersionId,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertDeleteObjectOutput(obsOutput *obs.DeleteObjectOutput) (*oss.DeleteObjectOutput, error) {
	output := &oss.DeleteObjectOutput{
		DeleteMarker:   obsOutput.DeleteMarker,
		RequestCharged: "",
		VersionId:      obsOutput.VersionId,
	}
	return output, nil
}

func (h *HuaweiyunOSS) PutObjectTagging(ctx context.Context, input *oss.PutObjectTaggingInput) (*oss.PutObjectTaggingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) DeleteObjectTagging(ctx context.Context, input *oss.DeleteObjectTaggingInput) (*oss.DeleteObjectTaggingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) GetObjectTagging(ctx context.Context, input *oss.GetObjectTaggingInput) (*oss.GetObjectTaggingOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) CopyObject(ctx context.Context, input *oss.CopyObjectInput) (*oss.CopyObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertCopyObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.CopyObject(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertCopyObjectOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertCopyObjectInput(input *oss.CopyObjectInput) (*obs.CopyObjectInput, error) {
	obsInput := &obs.CopyObjectInput{
		ObjectOperationInput: obs.ObjectOperationInput{
			Bucket:                  input.Bucket,
			Key:                     input.Key,
			ACL:                     "",
			GrantReadId:             "",
			GrantReadAcpId:          "",
			GrantWriteAcpId:         "",
			GrantFullControlId:      "",
			StorageClass:            "",
			WebsiteRedirectLocation: "",
			Expires:                 input.Expires,
			SseHeader:               nil,
			Metadata:                input.Metadata,
		},
		CopySourceBucket:            input.CopySource.CopySourceBucket,
		CopySourceKey:               input.CopySource.CopySourceKey,
		CopySourceVersionId:         input.CopySource.CopySourceVersionId,
		CopySourceIfMatch:           "",
		CopySourceIfNoneMatch:       "",
		CopySourceIfUnmodifiedSince: time.Time{},
		CopySourceIfModifiedSince:   time.Time{},
		SourceSseHeader:             nil,
		CacheControl:                "",
		ContentDisposition:          "",
		ContentEncoding:             "",
		ContentLanguage:             "",
		ContentType:                 "",
		Expires:                     "",
		MetadataDirective:           "",
		SuccessActionRedirect:       "",
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertCopyObjectOutput(obsOutput *obs.CopyObjectOutput) (*oss.CopyObjectOutput, error) {
	output := &oss.CopyObjectOutput{CopyObjectResult: &oss.CopyObjectResult{
		ETag:         obsOutput.ETag,
		LastModified: obsOutput.LastModified.Unix(),
	}}
	return output, nil
}

func (h *HuaweiyunOSS) DeleteObjects(ctx context.Context, input *oss.DeleteObjectsInput) (*oss.DeleteObjectsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertDeleteObjectsInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.DeleteObjects(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertDeleteObjectsOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertDeleteObjectsInput(input *oss.DeleteObjectsInput) (*obs.DeleteObjectsInput, error) {
	obsInput := &obs.DeleteObjectsInput{
		Bucket:       input.Bucket,
		XMLName:      xml.Name{},
		Quiet:        input.Delete.Quiet,
		Objects:      nil,
		EncodingType: "",
	}
	objects := make([]obs.ObjectToDelete, len(input.Delete.Objects))
	for _, v := range input.Delete.Objects {
		object := obs.ObjectToDelete{
			XMLName:   xml.Name{},
			Key:       v.Key,
			VersionId: v.VersionId,
		}
		objects = append(objects, object)
	}
	obsInput.Objects = objects
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertDeleteObjectsOutput(obsOutput *obs.DeleteObjectsOutput) (*oss.DeleteObjectsOutput, error) {
	output := &oss.DeleteObjectsOutput{}
	deleteObjects := make([]*oss.DeletedObject, len(obsOutput.Deleteds))
	for _, v := range obsOutput.Deleteds {
		object := &oss.DeletedObject{
			DeleteMarker:          v.DeleteMarker,
			DeleteMarkerVersionId: v.DeleteMarkerVersionId,
			Key:                   v.Key,
			VersionId:             v.VersionId,
		}
		deleteObjects = append(deleteObjects, object)
	}
	output.Deleted = deleteObjects
	return output, nil
}

func (h *HuaweiyunOSS) ListObjects(ctx context.Context, input *oss.ListObjectsInput) (*oss.ListObjectsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertListObjectsInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.ListObjects(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertListObjectsOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertListObjectsInput(input *oss.ListObjectsInput) (*obs.ListObjectsInput, error) {
	obsInput := &obs.ListObjectsInput{
		ListObjsInput: obs.ListObjsInput{
			Prefix:        input.Prefix,
			MaxKeys:       int(input.MaxKeys),
			Delimiter:     input.Delimiter,
			Origin:        "",
			RequestHeader: "",
			EncodingType:  input.EncodingType,
		},
		Bucket: input.Bucket,
		Marker: input.Marker,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertListObjectsOutput(obsOutput *obs.ListObjectsOutput) (*oss.ListObjectsOutput, error) {
	output := &oss.ListObjectsOutput{
		CommonPrefixes: obsOutput.CommonPrefixes,
		Contents:       nil,
		Delimiter:      obsOutput.Delimiter,
		EncodingType:   obsOutput.EncodingType,
		IsTruncated:    obsOutput.IsTruncated,
		Marker:         obsOutput.Marker,
		MaxKeys:        int32(obsOutput.MaxKeys),
		Name:           obsOutput.Name,
		NextMarker:     obsOutput.NextMarker,
		Prefix:         obsOutput.Prefix,
	}
	contexts := make([]*oss.Object, len(obsOutput.Contents))
	for _, v := range obsOutput.Contents {
		context := &oss.Object{
			ETag:         v.ETag,
			Key:          v.Key,
			LastModified: v.LastModified.Unix(),
			Owner: &oss.Owner{
				DisplayName: v.Owner.DisplayName,
				ID:          v.Owner.ID,
			},
			Size:         v.Size,
			StorageClass: string(v.StorageClass),
		}
		contexts = append(contexts, context)
	}
	output.Contents = contexts
	return output, nil
}

func (h *HuaweiyunOSS) GetObjectCannedAcl(ctx context.Context, input *oss.GetObjectCannedAclInput) (*oss.GetObjectCannedAclOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertGetObjectCannedAclInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.GetObjectAcl(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertGetObjectCannedAclOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertGetObjectCannedAclInput(input *oss.GetObjectCannedAclInput) (*obs.GetObjectAclInput, error) {
	obsInput := &obs.GetObjectAclInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionId: input.VersionId,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertGetObjectCannedAclOutput(obsOutput *obs.GetObjectAclOutput) (*oss.GetObjectCannedAclOutput, error) {
	output := &oss.GetObjectCannedAclOutput{
		CannedAcl: "",
		Owner: &oss.Owner{
			DisplayName: obsOutput.Owner.DisplayName,
			ID:          obsOutput.Owner.ID,
		},
		RequestCharged: "",
	}
	return output, nil
}

func (h *HuaweiyunOSS) PutObjectCannedAcl(ctx context.Context, input *oss.PutObjectCannedAclInput) (*oss.PutObjectCannedAclOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertPutObjectCannedAclInput(input)
	if err != nil {
		return nil, err
	}
	obsBaseModel, err := client.SetObjectAcl(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertPutObjectCannedAclOutput(obsBaseModel)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertPutObjectCannedAclInput(input *oss.PutObjectCannedAclInput) (*obs.SetObjectAclInput, error) {
	obsInput := &obs.SetObjectAclInput{
		Bucket:              input.Bucket,
		Key:                 input.Key,
		VersionId:           input.VersionId,
		ACL:                 obs.AclType(input.Acl),
		AccessControlPolicy: obs.AccessControlPolicy{},
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertPutObjectCannedAclOutput(baseModel *obs.BaseModel) (*oss.PutObjectCannedAclOutput, error) {
	output := &oss.PutObjectCannedAclOutput{RequestCharged: ""}
	return output, nil
}

func (h *HuaweiyunOSS) RestoreObject(ctx context.Context, input *oss.RestoreObjectInput) (*oss.RestoreObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertRestoreObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsBaseModel, err := client.RestoreObject(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertRestoreObjectOutput(obsBaseModel)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertRestoreObjectInput(input *oss.RestoreObjectInput) (*obs.RestoreObjectInput, error) {
	obsInput := &obs.RestoreObjectInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionId: input.VersionId,
		XMLName:   xml.Name{},
		Days:      int(input.RestoreRequest.Days),
		Tier:      obs.RestoreTierType(input.RestoreRequest.Tier),
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertRestoreObjectOutput(baseModel *obs.BaseModel) (*oss.RestoreObjectOutput, error) {
	output := &oss.RestoreObjectOutput{
		RequestCharged:    "",
		RestoreOutputPath: "",
	}
	return output, nil
}

func (h *HuaweiyunOSS) CreateMultipartUpload(ctx context.Context, input *oss.CreateMultipartUploadInput) (*oss.CreateMultipartUploadOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UploadPart(ctx context.Context, input *oss.UploadPartInput) (*oss.UploadPartOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertUploadPartInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.UploadPart(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertUploadPartOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertUploadPartInput(input *oss.UploadPartInput) (*obs.UploadPartInput, error) {
	obsInput := &obs.UploadPartInput{
		Bucket:     input.Bucket,
		Key:        input.Key,
		PartNumber: int(input.PartNumber),
		UploadId:   input.UploadId,
		ContentMD5: input.ContentMd5,
		SseHeader:  nil,
		Body:       input.DataStream,
		SourceFile: "",
		Offset:     0,
		PartSize:   int64(input.PartNumber),
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertUploadPartOutput(baseModel *obs.UploadPartOutput) (*oss.UploadPartOutput, error) {
	output := &oss.UploadPartOutput{
		BucketKeyEnabled:     false,
		ETag:                 "",
		RequestCharged:       "",
		SSECustomerAlgorithm: "",
		SSECustomerKeyMD5:    "",
		SSEKMSKeyId:          "",
		ServerSideEncryption: "",
	}
	return output, nil
}

func (h *HuaweiyunOSS) UploadPartCopy(ctx context.Context, input *oss.UploadPartCopyInput) (*oss.UploadPartCopyOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) CompleteMultipartUpload(ctx context.Context, input *oss.CompleteMultipartUploadInput) (*oss.CompleteMultipartUploadOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertCompleteMultipartUploadInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.CompleteMultipartUpload(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertCompleteMultipartUploadOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertCompleteMultipartUploadInput(input *oss.CompleteMultipartUploadInput) (*obs.CompleteMultipartUploadInput, error) {
	obsInput := &obs.CompleteMultipartUploadInput{
		Bucket:       input.Bucket,
		Key:          input.Key,
		UploadId:     input.UploadId,
		XMLName:      xml.Name{},
		Parts:        nil,
		EncodingType: "",
	}
	parts := make([]obs.Part, len(input.MultipartUpload.Parts))
	for _, v := range input.MultipartUpload.Parts {
		part := obs.Part{
			XMLName:      xml.Name{},
			PartNumber:   int(v.PartNumber),
			ETag:         v.ETag,
			LastModified: time.Time{},
			Size:         0,
		}
		parts = append(parts, part)
	}
	obsInput.Parts = parts
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertCompleteMultipartUploadOutput(baseModel *obs.CompleteMultipartUploadOutput) (*oss.CompleteMultipartUploadOutput, error) {
	output := &oss.CompleteMultipartUploadOutput{
		Bucket:               "",
		Key:                  "",
		BucketKeyEnabled:     false,
		ETag:                 "",
		Expiration:           "",
		Location:             "",
		RequestCharged:       "",
		SSEKMSKeyId:          "",
		ServerSideEncryption: "",
		VersionId:            "",
	}
	return output, nil
}

func (h *HuaweiyunOSS) AbortMultipartUpload(ctx context.Context, input *oss.AbortMultipartUploadInput) (*oss.AbortMultipartUploadOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertAbortMultipartUploadInput(input)
	if err != nil {
		return nil, err
	}
	obsBaseModel, err := client.AbortMultipartUpload(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertAbortMultipartUploadOutput(obsBaseModel)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertAbortMultipartUploadInput(input *oss.AbortMultipartUploadInput) (*obs.AbortMultipartUploadInput, error) {
	obsInput := &obs.AbortMultipartUploadInput{
		Bucket:   input.Bucket,
		Key:      input.Key,
		UploadId: input.UploadId,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertAbortMultipartUploadOutput(baseModel *obs.BaseModel) (*oss.AbortMultipartUploadOutput, error) {
	output := &oss.AbortMultipartUploadOutput{RequestCharged: ""}
	return output, nil
}

func (h *HuaweiyunOSS) ListMultipartUploads(ctx context.Context, input *oss.ListMultipartUploadsInput) (*oss.ListMultipartUploadsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertListMultipartUploadsInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.ListMultipartUploads(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertListMultipartUploadsOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertListMultipartUploadsInput(input *oss.ListMultipartUploadsInput) (*obs.ListMultipartUploadsInput, error) {
	obsInput := &obs.ListMultipartUploadsInput{
		Bucket:         input.Bucket,
		Prefix:         input.Prefix,
		MaxUploads:     int(input.MaxUploads),
		Delimiter:      input.Delimiter,
		KeyMarker:      input.KeyMarker,
		UploadIdMarker: input.UploadIdMarker,
		EncodingType:   input.EncodingType,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertListMultipartUploadsOutput(obsOutput *obs.ListMultipartUploadsOutput) (*oss.ListMultipartUploadsOutput, error) {
	output := &oss.ListMultipartUploadsOutput{
		Bucket:             obsOutput.Bucket,
		CommonPrefixes:     obsOutput.CommonPrefixes,
		Delimiter:          obsOutput.Delimiter,
		EncodingType:       obsOutput.EncodingType,
		IsTruncated:        obsOutput.IsTruncated,
		KeyMarker:          obsOutput.KeyMarker,
		MaxUploads:         int32(obsOutput.MaxUploads),
		NextKeyMarker:      obsOutput.NextKeyMarker,
		NextUploadIDMarker: obsOutput.NextUploadIdMarker,
		Prefix:             obsOutput.Prefix,
		UploadIDMarker:     obsOutput.UploadIdMarker,
		Uploads:            nil,
	}
	uploads := make([]*oss.MultipartUpload, len(obsOutput.Uploads))
	for _, v := range obsOutput.Uploads {
		upload := &oss.MultipartUpload{
			Initiated: v.Initiated.Unix(),
			Initiator: &oss.Initiator{
				DisplayName: v.Initiator.DisplayName,
				ID:          v.Initiator.ID,
			},
			Key: v.Key,
			Owner: &oss.Owner{
				DisplayName: v.Owner.DisplayName,
				ID:          v.Owner.ID,
			},
			StorageClass: string(v.StorageClass),
			UploadId:     v.UploadId,
		}
		uploads = append(uploads, upload)
	}
	output.Uploads = uploads
	return output, nil
}

func (h *HuaweiyunOSS) ListObjectVersions(ctx context.Context, input *oss.ListObjectVersionsInput) (*oss.ListObjectVersionsOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) HeadObject(ctx context.Context, input *oss.HeadObjectInput) (*oss.HeadObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertHeadObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsBaseModel, err := client.HeadObject(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertHeadObjectOutput(obsBaseModel)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertHeadObjectInput(input *oss.HeadObjectInput) (*obs.HeadObjectInput, error) {
	obsInput := &obs.HeadObjectInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionId: input.VersionId,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertHeadObjectOutput(baseModel *obs.BaseModel) (*oss.HeadObjectOutput, error) {
	output := &oss.HeadObjectOutput{ResultMetadata: nil}
	return output, nil
}

func (h *HuaweiyunOSS) IsObjectExist(ctx context.Context, input *oss.IsObjectExistInput) (*oss.IsObjectExistOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) SignURL(ctx context.Context, input *oss.SignURLInput) (*oss.SignURLOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UpdateDownloadBandwidthRateLimit(ctx context.Context, input *oss.UpdateBandwidthRateLimitInput) error {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) UpdateUploadBandwidthRateLimit(ctx context.Context, input *oss.UpdateBandwidthRateLimitInput) error {
	//TODO implement me
	panic("implement me")
}

func (h *HuaweiyunOSS) AppendObject(ctx context.Context, input *oss.AppendObjectInput) (*oss.AppendObjectOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertAppendObjectInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.AppendObject(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertAppendObjectOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertAppendObjectInput(input *oss.AppendObjectInput) (*obs.AppendObjectInput, error) {
	obsInput := &obs.AppendObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket:                  input.Bucket,
				Key:                     input.Key,
				ACL:                     obs.AclType(input.ACL),
				GrantReadId:             "",
				GrantReadAcpId:          "",
				GrantWriteAcpId:         "",
				GrantFullControlId:      "",
				StorageClass:            obs.StorageClassType(input.StorageClass),
				WebsiteRedirectLocation: "",
				Expires:                 0,
				SseHeader:               nil,
				Metadata:                nil,
			},
			ContentType:     "",
			ContentMD5:      input.ContentMd5,
			ContentLength:   0,
			ContentEncoding: input.ContentEncoding,
		},
		Body:     input.DataStream,
		Position: input.Position,
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertAppendObjectOutput(obsOutput *obs.AppendObjectOutput) (*oss.AppendObjectOutput, error) {
	output := &oss.AppendObjectOutput{AppendPosition: obsOutput.NextAppendPosition}
	return output, nil
}

func (h *HuaweiyunOSS) ListParts(ctx context.Context, input *oss.ListPartsInput) (*oss.ListPartsOutput, error) {
	client, err := h.getClient()
	if err != nil {
		return nil, err
	}
	obsInput, err := h.convertListPartsInput(input)
	if err != nil {
		return nil, err
	}
	obsOutput, err := client.ListParts(obsInput)
	if err != nil {
		return nil, err
	}
	output, err := h.convertListPartsOutput(obsOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (h *HuaweiyunOSS) convertListPartsInput(input *oss.ListPartsInput) (*obs.ListPartsInput, error) {
	obsInput := &obs.ListPartsInput{
		Bucket:           input.Bucket,
		Key:              input.Key,
		UploadId:         input.UploadId,
		MaxParts:         int(input.MaxParts),
		PartNumberMarker: int(input.PartNumberMarker),
		EncodingType:     "",
	}
	return obsInput, nil
}

func (h *HuaweiyunOSS) convertListPartsOutput(obsOutput *obs.ListPartsOutput) (*oss.ListPartsOutput, error) {
	output := &oss.ListPartsOutput{
		Bucket:               obsOutput.Bucket,
		Key:                  obsOutput.Key,
		UploadId:             obsOutput.UploadId,
		NextPartNumberMarker: string(obsOutput.NextPartNumberMarker),
		MaxParts:             int64(obsOutput.MaxParts),
		IsTruncated:          obsOutput.IsTruncated,
		Parts:                nil,
	}
	parts := make([]*oss.Part, len(output.Parts))
	for _, v := range output.Parts {
		part := &oss.Part{
			Etag:         v.Etag,
			LastModified: v.LastModified,
			PartNumber:   v.PartNumber,
			Size:         v.Size,
		}
		parts = append(parts, part)
	}
	output.Parts = parts
	return output, nil
}

func (h *HuaweiyunOSS) getClient() (*obs.ObsClient, error) {
	if h.client == nil {
		return nil, utils.ErrNotInitClient
	}
	return h.client, nil
}
