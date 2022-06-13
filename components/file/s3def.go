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

package file

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNotSpecifyEndpoint = errors.New("should specific endpoint in metadata")
)

type Oss interface {
	InitConfig(context.Context, *FileConfig) error
	InitClient(context.Context, *InitRequest) error
	GetObject(context.Context, *GetObjectInput) (io.ReadCloser, error)
	PutObject(context.Context, *PutObjectInput) (*PutObjectOutput, error)
	DeleteObject(context.Context, *DeleteObjectInput) (*DeleteObjectOutput, error)
	PutObjectTagging(context.Context, *PutBucketTaggingInput) (*PutBucketTaggingOutput, error)
	DeleteObjectTagging(context.Context, *DeleteObjectTaggingInput) (*DeleteObjectTaggingOutput, error)
	GetObjectTagging(context.Context, *GetObjectTaggingInput) (*GetObjectTaggingOutput, error)
	CopyObject(context.Context, *CopyObjectInput) (*CopyObjectOutput, error)
	DeleteObjects(context.Context, *DeleteObjectsInput) (*DeleteObjectsOutput, error)
	ListObjects(context.Context, *ListObjectsInput) (*ListObjectsOutput, error)
	GetObjectAcl(context.Context, *GetObjectAclInput) (*GetObjectAclOutput, error)
	PutObjectAcl(context.Context, *PutObjectAclInput) (*PutObjectAclOutput, error)
	RestoreObject(context.Context, *RestoreObjectInput) (*RestoreObjectOutput, error)
	CreateMultipartUpload(context.Context, *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error)
	UploadPart(context.Context, *UploadPartInput) (*UploadPartOutput, error)
	UploadPartCopy(context.Context, *UploadPartCopyInput) (*UploadPartCopyOutput, error)
	CompleteMultipartUpload(context.Context, *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error)
	AbortMultipartUpload(context.Context, *AbortMultipartUploadInput) (*AbortMultipartUploadOutput, error)
	ListMultipartUploads(context.Context, *ListMultipartUploadsInput) (*ListMultipartUploadsOutput, error)
	ListObjectVersions(context.Context, *ListObjectVersionsInput) (*ListObjectVersionsOutput, error)
	HeadObject(context.Context, *HeadObjectInput) (*HeadObjectOutput, error)
	IsObjectExist(context.Context, *IsObjectExistInput) (*IsObjectExistOutput, error)
}

type BaseConfig struct {
}
type InitRequest struct {
	App      string
	Metadata map[string]string
}

type GetObjectInput struct {
	Bucket                     string `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	ExpectedBucketOwner        string `protobuf:"bytes,3,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	IfMatch                    string `protobuf:"bytes,4,opt,name=if_match,json=ifMatch,proto3" json:"if_match,omitempty"`
	IfModifiedSince            int64  `protobuf:"varint,5,opt,name=if_modified_since,json=ifModifiedSince,proto3" json:"if_modified_since,omitempty"`
	IfNoneMatch                string `protobuf:"bytes,6,opt,name=if_none_match,json=ifNoneMatch,proto3" json:"if_none_match,omitempty"`
	IfUnmodifiedSince          int64  `protobuf:"varint,7,opt,name=if_unmodified_since,json=ifUnmodifiedSince,proto3" json:"if_unmodified_since,omitempty"`
	Key                        string `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty"`
	PartNumber                 int64  `protobuf:"varint,9,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
	Start                      int64  `protobuf:"varint,10,opt,name=start,proto3" json:"start,omitempty"`
	End                        int64  `protobuf:"varint,11,opt,name=end,proto3" json:"end,omitempty"`
	RequestPayer               string `protobuf:"bytes,12,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	ResponseCacheControl       string `protobuf:"bytes,13,opt,name=response_cache_control,json=responseCacheControl,proto3" json:"response_cache_control,omitempty"`
	ResponseContentDisposition string `protobuf:"bytes,14,opt,name=response_content_disposition,json=responseContentDisposition,proto3" json:"response_content_disposition,omitempty"`
	ResponseContentEncoding    string `protobuf:"bytes,15,opt,name=response_content_encoding,json=responseContentEncoding,proto3" json:"response_content_encoding,omitempty"`
	ResponseContentLanguage    string `protobuf:"bytes,16,opt,name=response_content_language,json=responseContentLanguage,proto3" json:"response_content_language,omitempty"`
	ResponseContentType        string `protobuf:"bytes,17,opt,name=response_content_type,json=responseContentType,proto3" json:"response_content_type,omitempty"`
	ResponseExpires            string `protobuf:"bytes,18,opt,name=response_expires,json=responseExpires,proto3" json:"response_expires,omitempty"`
	SseCustomerAlgorithm       string `protobuf:"bytes,19,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SseCustomerKey             string `protobuf:"bytes,20,opt,name=sse_customer_key,json=sseCustomerKey,proto3" json:"sse_customer_key,omitempty"`
	SseCustomerKeyMd5          string `protobuf:"bytes,21,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	VersionId                  string `protobuf:"bytes,22,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
	AcceptEncoding             string `protobuf:"bytes,23,opt,name=accept_encoding,json=acceptEncoding,proto3" json:"accept_encoding,omitempty"`
}

type PutObjectInput struct {
	DataStream           io.Reader
	ACL                  string            `protobuf:"bytes,2,opt,name=acl,proto3" json:"acl,omitempty"`
	Body                 []byte            `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	Bucket               string            `protobuf:"bytes,4,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                  string            `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	BucketKeyEnabled     bool              `protobuf:"varint,6,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	CacheControl         string            `protobuf:"bytes,7,opt,name=cache_control,json=cacheControl,proto3" json:"cache_control,omitempty"`
	ContentDisposition   string            `protobuf:"bytes,8,opt,name=content_disposition,json=contentDisposition,proto3" json:"content_disposition,omitempty"`
	ContentEncoding      string            `protobuf:"bytes,9,opt,name=content_encoding,json=contentEncoding,proto3" json:"content_encoding,omitempty"`
	Expires              int64             `protobuf:"varint,10,opt,name=expires,proto3" json:"expires,omitempty"`
	ServerSideEncryption string            `protobuf:"bytes,11,opt,name=server_side_encryption,json=serverSideEncryption,proto3" json:"server_side_encryption,omitempty"`
	Meta                 map[string]string `protobuf:"bytes,12,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type PutObjectOutput struct {
	BucketKeyEnabled bool   `protobuf:"varint,1,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	ETag             string `protobuf:"bytes,2,opt,name=etag,proto3" json:"etag,omitempty"`
}

type DeleteObjectInput struct {
	Bucket       string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key          string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	RequestPayer string `protobuf:"bytes,4,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
}
type DeleteObjectOutput struct {
	DeleteMarker   bool   `protobuf:"varint,1,opt,name=delete_marker,json=deleteMarker,proto3" json:"delete_marker,omitempty"`
	RequestCharged string `protobuf:"bytes,2,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
	VersionId      string `protobuf:"bytes,3,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type PutBucketTaggingInput struct {
	Bucket    string            `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key       string            `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Tags      map[string]string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	VersionId string            `protobuf:"bytes,5,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}
type PutBucketTaggingOutput struct {
}

type DeleteObjectTaggingInput struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}
type DeleteObjectTaggingOutput struct {
	VersionId string `protobuf:"bytes,1,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type GetObjectTaggingInput struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}
type GetObjectTaggingOutput struct {
	Tags map[string]string `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type CopySource struct {
	CopySourceBucket    string `protobuf:"bytes,1,opt,name=copy_source_bucket,json=copySourceBucket,proto3" json:"copy_source_bucket,omitempty"`
	CopySourceKey       string `protobuf:"bytes,2,opt,name=copy_source_key,json=copySourceKey,proto3" json:"copy_source_key,omitempty"`
	CopySourceVersionId string `protobuf:"bytes,3,opt,name=copy_source_version_id,json=copySourceVersionId,proto3" json:"copy_source_version_id,omitempty"`
}

type CopyObjectInput struct {
	StoreName  string      `protobuf:"bytes,1,opt,name=store_name,json=storeName,proto3" json:"store_name,omitempty"`
	Bucket     string      `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key        string      `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	CopySource *CopySource `protobuf:"bytes,4,opt,name=copy_source,json=copySource,proto3" json:"copy_source,omitempty"`
}
type CopyObjectOutput struct {
	CopyObjectResult *CopyObjectResult `protobuf:"bytes,1,opt,name=CopyObjectResult,proto3" json:"CopyObjectResult,omitempty"`
}
type CopyObjectResult struct {
	ETag         string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	LastModified int64  `protobuf:"varint,2,opt,name=LastModified,proto3" json:"LastModified,omitempty"`
}

type DeleteObjectsInput struct {
	Bucket string  `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Delete *Delete `protobuf:"bytes,2,opt,name=Delete,proto3" json:"Delete,omitempty"`
}
type Delete struct {
	Objects []*ObjectIdentifier `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
	Quiet   bool                `protobuf:"varint,2,opt,name=quiet,proto3" json:"quiet,omitempty"`
}
type ObjectIdentifier struct {
	Key       string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	VersionId string `protobuf:"bytes,2,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type DeleteObjectsOutput struct {
	Deleted []*DeletedObject `protobuf:"bytes,1,rep,name=deleted,proto3" json:"deleted,omitempty"`
}

type DeletedObject struct {
	DeleteMarker          bool   `protobuf:"varint,1,opt,name=delete_marker,json=deleteMarker,proto3" json:"delete_marker,omitempty"`
	DeleteMarkerVersionId string `protobuf:"bytes,2,opt,name=delete_marker_version_id,json=deleteMarkerVersionId,proto3" json:"delete_marker_version_id,omitempty"`
	Key                   string `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	VersionId             string `protobuf:"bytes,4,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type ListObjectsInput struct {
	Bucket              string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Delimiter           string `protobuf:"bytes,2,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	EncodingType        string `protobuf:"bytes,3,opt,name=encoding_type,json=encodingType,proto3" json:"encoding_type,omitempty"`
	ExpectedBucketOwner string `protobuf:"bytes,4,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	Marker              string `protobuf:"bytes,5,opt,name=marker,proto3" json:"marker,omitempty"`
	MaxKeys             int32  `protobuf:"varint,6,opt,name=maxKeys,proto3" json:"maxKeys,omitempty"`
	Prefix              string `protobuf:"bytes,7,opt,name=prefix,proto3" json:"prefix,omitempty"`
	RequestPayer        string `protobuf:"bytes,8,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
}
type ListObjectsOutput struct {
	CommonPrefixes []string  `protobuf:"bytes,1,rep,name=common_prefixes,json=commonPrefixes,proto3" json:"common_prefixes,omitempty"`
	Contents       []*Object `protobuf:"bytes,2,rep,name=contents,proto3" json:"contents,omitempty"`
	Delimiter      string    `protobuf:"bytes,3,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	EncodingType   string    `protobuf:"bytes,4,opt,name=encoding_type,json=encodingType,proto3" json:"encoding_type,omitempty"`
	IsTruncated    bool      `protobuf:"varint,5,opt,name=is_truncated,json=isTruncated,proto3" json:"is_truncated,omitempty"`
	Marker         string    `protobuf:"bytes,6,opt,name=marker,proto3" json:"marker,omitempty"`
	MaxKeys        int32     `protobuf:"varint,7,opt,name=max_keys,json=maxKeys,proto3" json:"max_keys,omitempty"`
	Name           string    `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	NextMarker     string    `protobuf:"bytes,9,opt,name=next_marker,json=nextMarker,proto3" json:"next_marker,omitempty"`
	Prefix         string    `protobuf:"bytes,10,opt,name=prefix,proto3" json:"prefix,omitempty"`
}
type Object struct {
	ETag         string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	Key          string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	LastModified int64  `protobuf:"bytes,3,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
	Owner        *Owner `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	Size         int64  `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	StorageClass string `protobuf:"bytes,6,opt,name=storage_class,json=storageClass,proto3" json:"storage_class,omitempty"`
}
type Owner struct {
	DisplayName string `protobuf:"bytes,1,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	ID          string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

type GetObjectAclInput struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}
type GetObjectAclOutput struct {
	Grants         []*Grant `protobuf:"bytes,1,rep,name=grants,proto3" json:"grants,omitempty"`
	Owner          *Owner   `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	RequestCharged string   `protobuf:"bytes,3,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
}
type Grant struct {
	Grantee    *Grantee `protobuf:"bytes,1,opt,name=grantee,proto3" json:"grantee,omitempty"`
	Permission string   `protobuf:"bytes,2,opt,name=permission,proto3" json:"permission,omitempty"`
}
type Grantee struct {
	DisplayName  string `protobuf:"bytes,1,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	EmailAddress string `protobuf:"bytes,2,opt,name=email_address,json=emailAddress,proto3" json:"email_address,omitempty"`
	ID           string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Type         string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	URI          string `protobuf:"bytes,5,opt,name=uri,proto3" json:"uri,omitempty"`
}

type PutObjectAclInput struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Acl    string `protobuf:"bytes,3,opt,name=acl,proto3" json:"acl,omitempty"`
}
type PutObjectAclOutput struct {
	RequestCharged string `protobuf:"bytes,1,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
}

type RestoreObjectInput struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}
type RestoreObjectOutput struct {
	RequestCharged    string `protobuf:"bytes,1,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
	RestoreOutputPath string `protobuf:"bytes,2,opt,name=restore_output_path,json=restoreOutputPath,proto3" json:"restore_output_path,omitempty"`
}

type CreateMultipartUploadInput struct {
	Bucket                    string            `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                       string            `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	ACL                       string            `protobuf:"bytes,3,opt,name=acl,proto3" json:"acl,omitempty"`
	BucketKeyEnabled          bool              `protobuf:"varint,4,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	CacheControl              string            `protobuf:"bytes,5,opt,name=cache_control,json=cacheControl,proto3" json:"cache_control,omitempty"`
	ContentDisposition        string            `protobuf:"bytes,6,opt,name=content_disposition,json=contentDisposition,proto3" json:"content_disposition,omitempty"`
	ContentEncoding           string            `protobuf:"bytes,7,opt,name=content_encoding,json=contentEncoding,proto3" json:"content_encoding,omitempty"`
	ContentLanguage           string            `protobuf:"bytes,8,opt,name=content_language,json=contentLanguage,proto3" json:"content_language,omitempty"`
	ContentType               string            `protobuf:"bytes,9,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	ExpectedBucketOwner       string            `protobuf:"bytes,10,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	Expires                   int64             `protobuf:"bytes,11,opt,name=expires,proto3" json:"expires,omitempty"`
	GrantFullControl          string            `protobuf:"bytes,12,opt,name=grant_full_control,json=grantFullControl,proto3" json:"grant_full_control,omitempty"`
	GrantRead                 string            `protobuf:"bytes,13,opt,name=grant_read,json=grantRead,proto3" json:"grant_read,omitempty"`
	GrantReadACP              string            `protobuf:"bytes,14,opt,name=grant_read_acp,json=grantReadAcp,proto3" json:"grant_read_acp,omitempty"`
	GrantWriteACP             string            `protobuf:"bytes,15,opt,name=grant_write_acp,json=grantWriteAcp,proto3" json:"grant_write_acp,omitempty"`
	MetaData                  map[string]string `protobuf:"bytes,16,rep,name=meta_data,json=metaData,proto3" json:"meta_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ObjectLockLegalHoldStatus string            `protobuf:"bytes,17,opt,name=object_lock_legal_hold_status,json=objectLockLegalHoldStatus,proto3" json:"object_lock_legal_hold_status,omitempty"`
	ObjectLockMode            string            `protobuf:"bytes,18,opt,name=object_lock_mode,json=objectLockMode,proto3" json:"object_lock_mode,omitempty"`
	ObjectLockRetainUntilDate int64             `protobuf:"bytes,19,opt,name=object_lock_retain_until_date,json=objectLockRetainUntilDate,proto3" json:"object_lock_retain_until_date,omitempty"`
	RequestPayer              string            `protobuf:"bytes,20,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	SSECustomerAlgorithm      string            `protobuf:"bytes,21,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SSECustomerKey            string            `protobuf:"bytes,22,opt,name=sse_customer_key,json=sseCustomerKey,proto3" json:"sse_customer_key,omitempty"`
	SSECustomerKeyMD5         string            `protobuf:"bytes,23,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	SSEKMSEncryptionContext   string            `protobuf:"bytes,24,opt,name=sse_kms_encryption_context,json=sseKmsEncryptionContext,proto3" json:"sse_kms_encryption_context,omitempty"`
	SSEKMSKeyId               string            `protobuf:"bytes,25,opt,name=sse_kms_key_id,json=sseKmsKeyId,proto3" json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption      string            `protobuf:"bytes,26,opt,name=server_side_encryption,json=serverSideEncryption,proto3" json:"server_side_encryption,omitempty"`
	StorageClass              string            `protobuf:"bytes,27,opt,name=storage_class,json=storageClass,proto3" json:"storage_class,omitempty"`
	Tagging                   string            `protobuf:"bytes,28,opt,name=tagging,proto3" json:"tagging,omitempty"`
	WebsiteRedirectLocation   string            `protobuf:"bytes,29,opt,name=website_redirect_location,json=websiteRedirectLocation,proto3" json:"website_redirect_location,omitempty"`
}
type CreateMultipartUploadOutput struct {
	Bucket                  string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                     string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	AbortDate               int64  `protobuf:"bytes,3,opt,name=abort_date,json=abortDate,proto3" json:"abort_date,omitempty"`
	AbortRuleId             string `protobuf:"bytes,4,opt,name=abort_rule_id,json=abortRuleId,proto3" json:"abort_rule_id,omitempty"`
	BucketKeyEnabled        bool   `protobuf:"varint,5,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	RequestCharged          string `protobuf:"bytes,6,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
	SSECustomerAlgorithm    string `protobuf:"bytes,7,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SSECustomerKeyMD5       string `protobuf:"bytes,8,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	SSEKMSEncryptionContext string `protobuf:"bytes,9,opt,name=sse_kms_encryption_context,json=sseKmsEncryptionContext,proto3" json:"sse_kms_encryption_context,omitempty"`
	SSEKMSKeyId             string `protobuf:"bytes,10,opt,name=sse_kms_key_id,json=sseKmsKeyId,proto3" json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption    string `protobuf:"bytes,11,opt,name=server_side_encryption,json=serverSideEncryption,proto3" json:"server_side_encryption,omitempty"`
	UploadId                string `protobuf:"bytes,12,opt,name=upload_id,json=uploadId,proto3" json:"upload_id,omitempty"`
}

type UploadPartInput struct {
	DataStream io.Reader
	Bucket     string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key        string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	//Body                 []byte `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	ContentLength        int64  `protobuf:"varint,4,opt,name=content_length,json=contentLength,proto3" json:"content_length,omitempty"`
	ContentMd5           string `protobuf:"bytes,5,opt,name=content_md5,json=contentMd5,proto3" json:"content_md5,omitempty"`
	ExpectedBucketOwner  string `protobuf:"bytes,6,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	PartNumber           int32  `protobuf:"varint,7,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
	RequestPayer         string `protobuf:"bytes,8,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	SseCustomerAlgorithm string `protobuf:"bytes,9,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SseCustomerKey       string `protobuf:"bytes,10,opt,name=sse_customer_key,json=sseCustomerKey,proto3" json:"sse_customer_key,omitempty"`
	SseCustomerKeyMd5    string `protobuf:"bytes,11,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	UploadId             string `protobuf:"bytes,12,opt,name=upload_id,json=uploadId,proto3" json:"upload_id,omitempty"`
}
type UploadPartOutput struct {
	BucketKeyEnabled     bool   `protobuf:"varint,1,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	ETag                 string `protobuf:"bytes,2,opt,name=etag,proto3" json:"etag,omitempty"`
	RequestCharged       string `protobuf:"bytes,3,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
	SseCustomerAlgorithm string `protobuf:"bytes,4,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SseCustomerKeyMd5    string `protobuf:"bytes,5,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	SseKmsKeyId          string `protobuf:"bytes,6,opt,name=sse_kms_key_id,json=sseKmsKeyId,proto3" json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption string `protobuf:"bytes,7,opt,name=server_side_encryption,json=serverSideEncryption,proto3" json:"server_side_encryption,omitempty"`
}

type UploadPartCopyInput struct {
	StoreName     string      `protobuf:"bytes,1,opt,name=store_name,json=storeName,proto3" json:"store_name,omitempty"`
	Bucket        string      `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key           string      `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	CopySource    *CopySource `protobuf:"bytes,4,opt,name=copy_source,json=copySource,proto3" json:"copy_source,omitempty"`
	PartNumber    int32       `protobuf:"varint,5,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
	UploadId      string      `protobuf:"bytes,6,opt,name=upload_id,json=uploadId,proto3" json:"upload_id,omitempty"`
	StartPosition int64       `protobuf:"varint,7,opt,name=start_position,json=startPosition,proto3" json:"start_position,omitempty"`
	PartSize      int64       `protobuf:"varint,8,opt,name=part_size,json=partSize,proto3" json:"part_size,omitempty"`
}
type UploadPartCopyOutput struct {
	BucketKeyEnabled     bool            `protobuf:"varint,1,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	CopyPartResult       *CopyPartResult `protobuf:"bytes,2,opt,name=copy_part_result,json=copyPartResult,proto3" json:"copy_part_result,omitempty"`
	CopySourceVersionId  string          `protobuf:"bytes,3,opt,name=copy_source_version_id,json=copySourceVersionId,proto3" json:"copy_source_version_id,omitempty"`
	RequestCharged       string          `protobuf:"bytes,4,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
	SseCustomerAlgorithm string          `protobuf:"bytes,5,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SseCustomerKeyMd5    string          `protobuf:"bytes,6,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	SseKmsKeyId          string          `protobuf:"bytes,7,opt,name=sse_kms_key_id,json=sseKmsKeyId,proto3" json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption string          `protobuf:"bytes,8,opt,name=server_side_encryption,json=serverSideEncryption,proto3" json:"server_side_encryption,omitempty"`
}
type CopyPartResult struct {
	ETag         string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	LastModified int64  `protobuf:"bytes,2,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
}

type CompleteMultipartUploadInput struct {
	Bucket              string                    `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                 string                    `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	UploadId            string                    `protobuf:"bytes,3,opt,name=upload_id,json=uploadId,proto3" json:"upload_id,omitempty"`
	RequestPayer        string                    `protobuf:"bytes,4,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	ExpectedBucketOwner string                    `protobuf:"bytes,5,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	MultipartUpload     *CompletedMultipartUpload `protobuf:"bytes,6,opt,name=multipart_upload,json=multipartUpload,proto3" json:"multipart_upload,omitempty"`
}
type CompletedMultipartUpload struct {
	Parts []*CompletedPart `protobuf:"bytes,1,rep,name=parts,proto3" json:"parts,omitempty"`
}
type CompletedPart struct {
	ETag       string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	PartNumber int32  `protobuf:"varint,2,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
}
type CompleteMultipartUploadOutput struct {
	Bucket               string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                  string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	BucketKeyEnabled     bool   `protobuf:"varint,3,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	ETag                 string `protobuf:"bytes,4,opt,name=etag,proto3" json:"etag,omitempty"`
	Expiration           string `protobuf:"bytes,5,opt,name=expiration,proto3" json:"expiration,omitempty"`
	Location             string `protobuf:"bytes,6,opt,name=location,proto3" json:"location,omitempty"`
	RequestCharged       string `protobuf:"bytes,7,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
	SseKmsKeyId          string `protobuf:"bytes,8,opt,name=sse_kms_keyId,json=sseKmsKeyId,proto3" json:"sse_kms_keyId,omitempty"`
	ServerSideEncryption string `protobuf:"bytes,9,opt,name=server_side_encryption,json=serverSideEncryption,proto3" json:"server_side_encryption,omitempty"`
	VersionId            string `protobuf:"bytes,10,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type AbortMultipartUploadInput struct {
	Bucket              string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                 string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	ExpectedBucketOwner string `protobuf:"bytes,3,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	RequestPayer        string `protobuf:"bytes,4,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	UploadId            string `protobuf:"bytes,5,opt,name=upload_id,json=uploadId,proto3" json:"upload_id,omitempty"`
}
type AbortMultipartUploadOutput struct {
	RequestCharged string `protobuf:"bytes,1,opt,name=request_charged,json=requestCharged,proto3" json:"request_charged,omitempty"`
}

type ListMultipartUploadsInput struct {
	Bucket              string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Delimiter           string `protobuf:"bytes,2,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	EncodingType        string `protobuf:"bytes,3,opt,name=encoding_type,json=encodingType,proto3" json:"encoding_type,omitempty"`
	ExpectedBucketOwner string `protobuf:"bytes,4,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	KeyMarker           string `protobuf:"bytes,5,opt,name=key_marker,json=keyMarker,proto3" json:"key_marker,omitempty"`
	MaxUploads          int64  `protobuf:"varint,6,opt,name=max_uploads,json=maxUploads,proto3" json:"max_uploads,omitempty"`
	Prefix              string `protobuf:"bytes,7,opt,name=prefix,proto3" json:"prefix,omitempty"`
	UploadIdMarker      string `protobuf:"bytes,8,opt,name=upload_id_marker,json=uploadIdMarker,proto3" json:"upload_id_marker,omitempty"`
}
type ListMultipartUploadsOutput struct {
	Bucket             string             `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	CommonPrefixes     []string           `protobuf:"bytes,2,rep,name=common_prefixes,json=commonPrefixes,proto3" json:"common_prefixes,omitempty"`
	Delimiter          string             `protobuf:"bytes,3,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	EncodingType       string             `protobuf:"bytes,4,opt,name=encoding_type,json=encodingType,proto3" json:"encoding_type,omitempty"`
	IsTruncated        bool               `protobuf:"varint,5,opt,name=is_truncated,json=isTruncated,proto3" json:"is_truncated,omitempty"`
	KeyMarker          string             `protobuf:"bytes,6,opt,name=key_marker,json=keyMarker,proto3" json:"key_marker,omitempty"`
	MaxUploads         int32              `protobuf:"varint,7,opt,name=max_uploads,json=maxUploads,proto3" json:"max_uploads,omitempty"`
	NextKeyMarker      string             `protobuf:"bytes,8,opt,name=next_key_marker,json=nextKeyMarker,proto3" json:"next_key_marker,omitempty"`
	NextUploadIdMarker string             `protobuf:"bytes,9,opt,name=next_upload_id_marker,json=nextUploadIdMarker,proto3" json:"next_upload_id_marker,omitempty"`
	Prefix             string             `protobuf:"bytes,10,opt,name=prefix,proto3" json:"prefix,omitempty"`
	UploadIdMarker     string             `protobuf:"bytes,11,opt,name=upload_id_marker,json=uploadIdMarker,proto3" json:"upload_id_marker,omitempty"`
	Uploads            []*MultipartUpload `protobuf:"bytes,12,rep,name=uploads,proto3" json:"uploads,omitempty"`
}
type MultipartUpload struct {
	Initiated    int64      `protobuf:"bytes,1,opt,name=initiated,proto3" json:"initiated,omitempty"`
	Initiator    *Initiator `protobuf:"bytes,2,opt,name=initiator,proto3" json:"initiator,omitempty"`
	Key          string     `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Owner        *Owner     `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	StorageClass string     `protobuf:"bytes,5,opt,name=storage_class,json=storageClass,proto3" json:"storage_class,omitempty"`
	UploadId     string     `protobuf:"bytes,6,opt,name=upload_id,json=uploadId,proto3" json:"upload_id,omitempty"`
}
type Initiator struct {
	DisplayName string `protobuf:"bytes,1,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	ID          string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

type ListObjectVersionsInput struct {
	Bucket              string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Delimiter           string `protobuf:"bytes,2,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	EncodingType        string `protobuf:"bytes,3,opt,name=encoding_type,json=encodingType,proto3" json:"encoding_type,omitempty"`
	ExpectedBucketOwner string `protobuf:"bytes,4,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	KeyMarker           string `protobuf:"bytes,5,opt,name=key_marker,json=keyMarker,proto3" json:"key_marker,omitempty"`
	MaxKeys             int32  `protobuf:"varint,6,opt,name=max_keys,json=maxKeys,proto3" json:"max_keys,omitempty"`
	Prefix              string `protobuf:"bytes,7,opt,name=prefix,proto3" json:"prefix,omitempty"`
	VersionIdMarker     string `protobuf:"bytes,8,opt,name=version_id_marker,json=versionIdMarker,proto3" json:"version_id_marker,omitempty"`
}
type ListObjectVersionsOutput struct {
	CommonPrefixes      []string             `protobuf:"bytes,1,rep,name=common_prefixes,json=commonPrefixes,proto3" json:"common_prefixes,omitempty"`
	DeleteMarkers       []*DeleteMarkerEntry `protobuf:"bytes,2,rep,name=delete_markers,json=deleteMarkers,proto3" json:"delete_markers,omitempty"`
	Delimiter           string               `protobuf:"bytes,3,opt,name=delimiter,proto3" json:"delimiter,omitempty"`
	EncodingType        string               `protobuf:"bytes,4,opt,name=encoding_type,json=encodingType,proto3" json:"encoding_type,omitempty"`
	IsTruncated         bool                 `protobuf:"varint,5,opt,name=is_truncated,json=isTruncated,proto3" json:"is_truncated,omitempty"`
	KeyMarker           string               `protobuf:"bytes,6,opt,name=key_marker,json=keyMarker,proto3" json:"key_marker,omitempty"`
	MaxKeys             int32                `protobuf:"varint,7,opt,name=max_keys,json=maxKeys,proto3" json:"max_keys,omitempty"`
	Name                string               `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	NextKeyMarker       string               `protobuf:"bytes,9,opt,name=next_key_marker,json=nextKeyMarker,proto3" json:"next_key_marker,omitempty"`
	NextVersionIdMarker string               `protobuf:"bytes,10,opt,name=next_version_id_marker,json=nextVersionIdMarker,proto3" json:"next_version_id_marker,omitempty"`
	Prefix              string               `protobuf:"bytes,11,opt,name=prefix,proto3" json:"prefix,omitempty"`
	VersionIdMarker     string               `protobuf:"bytes,12,opt,name=version_id_marker,json=versionIdMarker,proto3" json:"version_id_marker,omitempty"`
	Versions            []*ObjectVersion     `protobuf:"bytes,13,rep,name=versions,proto3" json:"versions,omitempty"`
}
type DeleteMarkerEntry struct {
	IsLatest     bool   `protobuf:"varint,1,opt,name=is_latest,json=isLatest,proto3" json:"is_latest,omitempty"`
	Key          string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	LastModified int64  `protobuf:"bytes,3,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
	Owner        *Owner `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	VersionId    string `protobuf:"bytes,5,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}
type ObjectVersion struct {
	ETag         string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	IsLatest     bool   `protobuf:"varint,2,opt,name=is_latest,json=isLatest,proto3" json:"is_latest,omitempty"`
	Key          string `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	LastModified int64  `protobuf:"bytes,4,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
	Owner        *Owner `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	Size         int64  `protobuf:"varint,6,opt,name=size,proto3" json:"size,omitempty"`
	StorageClass string `protobuf:"bytes,7,opt,name=storage_class,json=storageClass,proto3" json:"storage_class,omitempty"`
	VersionId    string `protobuf:"bytes,8,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type HeadObjectInput struct {
	Bucket               string `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key                  string `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	ChecksumMode         string `protobuf:"bytes,4,opt,name=checksum_mode,json=checksumMode,proto3" json:"checksum_mode,omitempty"`
	ExpectedBucketOwner  string `protobuf:"bytes,5,opt,name=expected_bucket_owner,json=expectedBucketOwner,proto3" json:"expected_bucket_owner,omitempty"`
	IfMatch              string `protobuf:"bytes,6,opt,name=if_match,json=ifMatch,proto3" json:"if_match,omitempty"`
	IfModifiedSince      int64  `protobuf:"varint,7,opt,name=if_modified_since,json=ifModifiedSince,proto3" json:"if_modified_since,omitempty"`
	IfNoneMatch          string `protobuf:"bytes,8,opt,name=if_none_match,json=ifNoneMatch,proto3" json:"if_none_match,omitempty"`
	IfUnmodifiedSince    int64  `protobuf:"varint,9,opt,name=if_unmodified_since,json=ifUnmodifiedSince,proto3" json:"if_unmodified_since,omitempty"`
	PartNumber           int32  `protobuf:"varint,10,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
	RequestPayer         string `protobuf:"bytes,11,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	SseCustomerAlgorithm string `protobuf:"bytes,12,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SseCustomerKey       string `protobuf:"bytes,13,opt,name=sse_customer_key,json=sseCustomerKey,proto3" json:"sse_customer_key,omitempty"`
	SseCustomerKeyMd5    string `protobuf:"bytes,14,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	VersionId            string `protobuf:"bytes,15,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}
type HeadObjectOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata map[string]string `protobuf:"bytes,1,rep,name=ResultMetadata,proto3" json:"ResultMetadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type IsObjectExistInput struct {
	Bucket string `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key    string `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
}
type IsObjectExistOutput struct {
	FileExist bool `protobuf:"varint,1,opt,name=file_exist,json=fileExist,proto3" json:"file_exist,omitempty"`
}
