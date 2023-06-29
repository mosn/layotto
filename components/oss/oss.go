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
	"io"
)

const (
	ServiceName = "oss"
)

// Oss is the interface for ObjectStorageService components.
// For more details about the fields explanation, please refer to the `.proto` file.
type Oss interface {
	Init(context.Context, *Config) error
	GetObject(context.Context, *GetObjectInput) (*GetObjectOutput, error)
	PutObject(context.Context, *PutObjectInput) (*PutObjectOutput, error)
	DeleteObject(context.Context, *DeleteObjectInput) (*DeleteObjectOutput, error)
	PutObjectTagging(context.Context, *PutObjectTaggingInput) (*PutObjectTaggingOutput, error)
	DeleteObjectTagging(context.Context, *DeleteObjectTaggingInput) (*DeleteObjectTaggingOutput, error)
	GetObjectTagging(context.Context, *GetObjectTaggingInput) (*GetObjectTaggingOutput, error)
	CopyObject(context.Context, *CopyObjectInput) (*CopyObjectOutput, error)
	DeleteObjects(context.Context, *DeleteObjectsInput) (*DeleteObjectsOutput, error)
	ListObjects(context.Context, *ListObjectsInput) (*ListObjectsOutput, error)
	GetObjectCannedAcl(context.Context, *GetObjectCannedAclInput) (*GetObjectCannedAclOutput, error)
	PutObjectCannedAcl(context.Context, *PutObjectCannedAclInput) (*PutObjectCannedAclOutput, error)
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
	SignURL(context.Context, *SignURLInput) (*SignURLOutput, error)
	UpdateDownloadBandwidthRateLimit(context.Context, *UpdateBandwidthRateLimitInput) error
	UpdateUploadBandwidthRateLimit(context.Context, *UpdateBandwidthRateLimitInput) error
	AppendObject(context.Context, *AppendObjectInput) (*AppendObjectOutput, error)
	ListParts(context.Context, *ListPartsInput) (*ListPartsOutput, error)
}

type GetObjectInput struct {
	Bucket                     string `json:"bucket,omitempty"`
	ExpectedBucketOwner        string `json:"expected_bucket_owner,omitempty"`
	IfMatch                    string `json:"if_match,omitempty"`
	IfModifiedSince            int64  `json:"if_modified_since,omitempty"`
	IfNoneMatch                string `json:"if_none_match,omitempty"`
	IfUnmodifiedSince          int64  `json:"if_unmodified_since,omitempty"`
	Key                        string `json:"key,omitempty"`
	PartNumber                 int64  `json:"part_number,omitempty"`
	Start                      int64  `json:"start,omitempty"`
	End                        int64  `json:"end,omitempty"`
	RequestPayer               string `json:"request_payer,omitempty"`
	ResponseCacheControl       string `json:"response_cache_control,omitempty"`
	ResponseContentDisposition string `json:"response_content_disposition,omitempty"`
	ResponseContentEncoding    string `json:"response_content_encoding,omitempty"`
	ResponseContentLanguage    string `json:"response_content_language,omitempty"`
	ResponseContentType        string `json:"response_content_type,omitempty"`
	ResponseExpires            string `json:"response_expires,omitempty"`
	SseCustomerAlgorithm       string `json:"sse_customer_algorithm,omitempty"`
	SseCustomerKey             string `json:"sse_customer_key,omitempty"`
	SseCustomerKeyMd5          string `json:"sse_customer_key_md5,omitempty"`
	VersionId                  string `json:"version_id,omitempty"`
	AcceptEncoding             string `json:"accept_encoding,omitempty"`
	SignedUrl                  string `json:"signed_url,omitempty"`
}

type GetObjectOutput struct {
	DataStream         io.ReadCloser
	CacheControl       string            `json:"cache_control,omitempty"`
	ContentDisposition string            `json:"content_disposition,omitempty"`
	ContentEncoding    string            `json:"content_encoding,omitempty"`
	ContentLanguage    string            `json:"content_language,omitempty"`
	ContentLength      int64             `json:"content_length,omitempty"`
	ContentRange       string            `json:"content_range,omitempty"`
	ContentType        string            `json:"content_type,omitempty"`
	DeleteMarker       bool              `json:"delete_marker,omitempty"`
	Etag               string            `json:"etag,omitempty"`
	Expiration         string            `json:"expiration,omitempty"`
	Expires            string            `json:"expires,omitempty"`
	LastModified       int64             `json:"last_modified,omitempty"`
	VersionId          string            `json:"version_id,omitempty"`
	TagCount           int64             `json:"tag_count,omitempty"`
	StorageClass       string            `json:"storage_class,omitempty"`
	PartsCount         int64             `json:"parts_count,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type PutObjectInput struct {
	DataStream           io.Reader
	ACL                  string            `json:"acl,omitempty"`
	Bucket               string            `json:"bucket,omitempty"`
	Key                  string            `json:"key,omitempty"`
	BucketKeyEnabled     bool              `json:"bucket_key_enabled,omitempty"`
	CacheControl         string            `json:"cache_control,omitempty"`
	ContentDisposition   string            `json:"content_disposition,omitempty"`
	ContentEncoding      string            `json:"content_encoding,omitempty"`
	Expires              int64             `json:"expires,omitempty"`
	ServerSideEncryption string            `json:"server_side_encryption,omitempty"`
	SignedUrl            string            `json:"signed_url,omitempty"`
	Meta                 map[string]string `json:"meta,omitempty"`
	Tagging              map[string]string `json:"tagging,omitempty"`
	StorageClass         string            `json:"storage_class,omitempty"`
	ContentLength        int64             `json:"content_length,omitempty"`
}

type PutObjectOutput struct {
	BucketKeyEnabled bool              `json:"bucket_key_enabled,omitempty"`
	ETag             string            `json:"etag,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

type DeleteObjectInput struct {
	Bucket       string `json:"bucket,omitempty"`
	Key          string `json:"key,omitempty"`
	RequestPayer string `json:"request_payer,omitempty"`
	VersionId    string `json:"version_id,omitempty"`
}
type DeleteObjectOutput struct {
	DeleteMarker   bool              `json:"delete_marker,omitempty"`
	RequestCharged string            `json:"request_charged,omitempty"`
	VersionId      string            `json:"version_id,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type PutObjectTaggingInput struct {
	Bucket    string            `json:"bucket,omitempty"`
	Key       string            `json:"key,omitempty"`
	Tags      map[string]string `json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	VersionId string            `json:"version_id,omitempty"`
}
type PutObjectTaggingOutput struct {
}

type DeleteObjectTaggingInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Key                 string `json:"key,omitempty"`
	VersionId           string `json:"version_id,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
}
type DeleteObjectTaggingOutput struct {
	VersionId string `json:"version_id,omitempty"`
}

type GetObjectTaggingInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Key                 string `json:"key,omitempty"`
	VersionId           string ` json:"version_id,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
	RequestPayer        string ` json:"request_payer,omitempty"`
}
type GetObjectTaggingOutput struct {
	Tags           map[string]string `json:"tags,omitempty"`
	VersionId      string            `json:"version_id,omitempty"`
	ResultMetadata map[string]string `json:"result_metadata,omitempty"`
}

type CopySource struct {
	CopySourceBucket    string `json:"copy_source_bucket,omitempty"`
	CopySourceKey       string `json:"copy_source_key,omitempty"`
	CopySourceVersionId string `json:"copy_source_version_id,omitempty"`
}

type CopyObjectInput struct {
	Bucket     string            `json:"bucket,omitempty"`
	Key        string            `json:"key,omitempty"`
	CopySource *CopySource       `json:"copy_source,omitempty"`
	Tagging    map[string]string `json:"tagging,omitempty"`
	Expires    int64             `json:"expires,omitempty"`
	// Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request.
	MetadataDirective string `json:"metadata_directive,omitempty"`
	// A map of metadata to store with the object in S3.
	Metadata map[string]string `json:"metadata,omitempty"`
}
type CopyObjectOutput struct {
	CopyObjectResult *CopyObjectResult `json:"copy_object_result,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}
type CopyObjectResult struct {
	ETag         string `json:"etag,omitempty"`
	LastModified int64  `json:"LastModified,omitempty"`
}

type DeleteObjectsInput struct {
	Bucket string  `json:"bucket,omitempty"`
	Delete *Delete `json:"delete,omitempty"`
}
type Delete struct {
	Objects []*ObjectIdentifier `json:"objects,omitempty"`
	Quiet   bool                `json:"quiet,omitempty"`
}
type ObjectIdentifier struct {
	Key       string `json:"key,omitempty"`
	VersionId string `json:"version_id,omitempty"`
}

type DeleteObjectsOutput struct {
	Deleted  []*DeletedObject  `json:"deleted,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type DeletedObject struct {
	DeleteMarker          bool   `json:"delete_marker,omitempty"`
	DeleteMarkerVersionId string `json:"delete_marker_version_id,omitempty"`
	Key                   string `json:"key,omitempty"`
	VersionId             string `json:"version_id,omitempty"`
}

type ListObjectsInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Delimiter           string `json:"delimiter,omitempty"`
	EncodingType        string `json:"encoding_type,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
	Marker              string `json:"marker,omitempty"`
	MaxKeys             int32  `json:"maxKeys,omitempty"`
	Prefix              string `json:"prefix,omitempty"`
	RequestPayer        string `json:"request_payer,omitempty"`
}
type ListObjectsOutput struct {
	CommonPrefixes []string          `json:"common_prefixes,omitempty"`
	Contents       []*Object         `json:"contents,omitempty"`
	Delimiter      string            `json:"delimiter,omitempty"`
	EncodingType   string            `json:"encoding_type,omitempty"`
	IsTruncated    bool              `json:"is_truncated,omitempty"`
	Marker         string            `json:"marker,omitempty"`
	MaxKeys        int32             `json:"max_keys,omitempty"`
	Name           string            `json:"name,omitempty"`
	NextMarker     string            `json:"next_marker,omitempty"`
	Prefix         string            `json:"prefix,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}
type Object struct {
	ETag         string `json:"etag,omitempty"`
	Key          string `json:"key,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
	Owner        *Owner `json:"owner,omitempty"`
	Size         int64  `json:"size,omitempty"`
	StorageClass string `json:"storage_class,omitempty"`
}
type Owner struct {
	DisplayName string `json:"display_name,omitempty"`
	ID          string `json:"id,omitempty"`
}

type GetObjectCannedAclInput struct {
	Bucket    string `json:"bucket,omitempty"`
	Key       string `json:"key,omitempty"`
	VersionId string `json:"version_id,omitempty"`
}
type GetObjectCannedAclOutput struct {
	CannedAcl      string            `json:"canned_acl,omitempty"`
	Owner          *Owner            `json:"owner,omitempty"`
	RequestCharged string            `json:"request_charged,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type PutObjectCannedAclInput struct {
	Bucket    string `json:"bucket,omitempty"`
	Key       string `json:"key,omitempty"`
	Acl       string `json:"acl,omitempty"`
	VersionId string `json:"version_id,omitempty"`
}
type PutObjectCannedAclOutput struct {
	RequestCharged string            `json:"request_charged,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type GlacierJobParameters struct {
	Tier string `json:"tier,omitempty"`
}
type OutputLocation struct {
	BucketName string `json:"bucket_name,omitempty"`
	Prefix     string `json:"prefix,omitempty"`
}
type CSVInput struct {
	AllowQuotedRecordDelimiter bool   `json:"allow_quoted_record_delimiter,omitempty"`
	Comments                   string `json:"comments,omitempty"`
	FieldDelimiter             string `json:"field_delimiter,omitempty"`
	FileHeaderInfo             string `json:"file_header_info,omitempty"`
	QuoteCharacter             string `json:"quote_character,omitempty"`
	QuoteEscapeCharacter       string `json:"quote_escape_character,omitempty"`
	RecordDelimiter            string `json:"record_delimiter,omitempty"`
}
type InputSerialization struct {
	CSV             CSVInput `json:"csv,omitempty"`
	CompressionType string   `json:"compression_type,omitempty"`
	JSON            string   `json:"json,omitempty"`
}
type CSVOutput struct {
	FieldDelimiter       string `json:"field_delimiter,omitempty"`
	QuoteCharacter       string `json:"quote_character,omitempty"`
	QuoteEscapeCharacter string `json:"quote_escape_character,omitempty"`
	QuoteFields          string `json:"quote_fields,omitempty"`
	RecordDelimiter      string `json:"record_delimiter,omitempty"`
}
type OutputSerialization struct {
	Csv  CSVOutput `json:"csv,omitempty"`
	Json string    `json:"json,omitempty"`
}
type SelectParameters struct {
	Expression          string              `json:"expression,omitempty"`
	ExpressionType      string              `json:"expression_type,omitempty"`
	InputSerialization  InputSerialization  `json:"input_serialization,omitempty"`
	OutputSerialization OutputSerialization `json:"output_serialization,omitempty"`
}
type RestoreRequest struct {
	Days                 int32                `json:"days,omitempty"`
	Description          string               `json:"description,omitempty"`
	GlacierJobParameters GlacierJobParameters `json:"glacier_job_parameters,omitempty"`
	OutputLocation       OutputLocation       `json:"output_location,omitempty"`
	SelectParameters     SelectParameters     `json:"select_parameters,omitempty"`
	Tier                 string               `json:"tier,omitempty"`
	Type                 string               `json:"type,omitempty"`
}

type RestoreObjectInput struct {
	Bucket         string         `json:"bucket,omitempty"`
	Key            string         `json:"key,omitempty"`
	RestoreRequest RestoreRequest `json:"restore_request,omitempty"`
	VersionId      string         `json:"version_id,omitempty"`
}
type RestoreObjectOutput struct {
	RequestCharged    string `json:"request_charged,omitempty"`
	RestoreOutputPath string `json:"restore_output_path,omitempty"`
}

type CreateMultipartUploadInput struct {
	Bucket                    string            `json:"bucket,omitempty"`
	Key                       string            `json:"key,omitempty"`
	ACL                       string            `json:"acl,omitempty"`
	BucketKeyEnabled          bool              `json:"bucket_key_enabled,omitempty"`
	CacheControl              string            `json:"cache_control,omitempty"`
	ContentDisposition        string            `json:"content_disposition,omitempty"`
	ContentEncoding           string            `json:"content_encoding,omitempty"`
	ContentLanguage           string            `json:"content_language,omitempty"`
	ContentType               string            `json:"content_type,omitempty"`
	ExpectedBucketOwner       string            `json:"expected_bucket_owner,omitempty"`
	Expires                   int64             `json:"expires,omitempty"`
	GrantFullControl          string            `json:"grant_full_control,omitempty"`
	GrantRead                 string            `json:"grant_read,omitempty"`
	GrantReadACP              string            `json:"grant_read_acp,omitempty"`
	GrantWriteACP             string            `json:"grant_write_acp,omitempty"`
	MetaData                  map[string]string `json:"meta_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ObjectLockLegalHoldStatus string            `json:"object_lock_legal_hold_status,omitempty"`
	ObjectLockMode            string            `json:"object_lock_mode,omitempty"`
	ObjectLockRetainUntilDate int64             `json:"object_lock_retain_until_date,omitempty"`
	RequestPayer              string            `json:"request_payer,omitempty"`
	SSECustomerAlgorithm      string            `json:"sse_customer_algorithm,omitempty"`
	SSECustomerKey            string            `json:"sse_customer_key,omitempty"`
	SSECustomerKeyMD5         string            `json:"sse_customer_key_md5,omitempty"`
	SSEKMSEncryptionContext   string            `json:"sse_kms_encryption_context,omitempty"`
	SSEKMSKeyId               string            `json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption      string            `json:"server_side_encryption,omitempty"`
	StorageClass              string            `json:"storage_class,omitempty"`
	Tagging                   map[string]string `json:"tagging,omitempty"`
	WebsiteRedirectLocation   string            `json:"website_redirect_location,omitempty"`
}
type CreateMultipartUploadOutput struct {
	Bucket                  string            `json:"bucket,omitempty"`
	Key                     string            `json:"key,omitempty"`
	AbortDate               int64             `json:"abort_date,omitempty"`
	AbortRuleId             string            `json:"abort_rule_id,omitempty"`
	BucketKeyEnabled        bool              `json:"bucket_key_enabled,omitempty"`
	RequestCharged          string            `json:"request_charged,omitempty"`
	SSECustomerAlgorithm    string            `json:"sse_customer_algorithm,omitempty"`
	SSECustomerKeyMD5       string            `json:"sse_customer_key_md5,omitempty"`
	SSEKMSEncryptionContext string            `json:"sse_kms_encryption_context,omitempty"`
	SSEKMSKeyId             string            `json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption    string            `json:"server_side_encryption,omitempty"`
	UploadId                string            `json:"upload_id,omitempty"`
	Metadata                map[string]string `json:"metadata,omitempty"`
}

type UploadPartInput struct {
	DataStream io.Reader
	Bucket     string `json:"bucket,omitempty"`
	Key        string `json:"key,omitempty"`
	//Body                 []byte `json:"body,omitempty"`
	ContentLength        int64  `json:"content_length,omitempty"`
	ContentMd5           string `json:"content_md5,omitempty"`
	ExpectedBucketOwner  string `json:"expected_bucket_owner,omitempty"`
	PartNumber           int32  `json:"part_number,omitempty"`
	RequestPayer         string `json:"request_payer,omitempty"`
	SseCustomerAlgorithm string `json:"sse_customer_algorithm,omitempty"`
	SseCustomerKey       string `json:"sse_customer_key,omitempty"`
	SseCustomerKeyMd5    string `json:"sse_customer_key_md5,omitempty"`
	UploadId             string `json:"upload_id,omitempty"`
}
type UploadPartOutput struct {
	BucketKeyEnabled     bool              `json:"bucket_key_enabled,omitempty"`
	ETag                 string            `json:"etag,omitempty"`
	RequestCharged       string            `json:"request_charged,omitempty"`
	SSECustomerAlgorithm string            `json:"sse_customer_algorithm,omitempty"`
	SSECustomerKeyMD5    string            `json:"sse_customer_key_md5,omitempty"`
	SSEKMSKeyId          string            `json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption string            `json:"server_side_encryption,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

type UploadPartCopyInput struct {
	Bucket        string      `json:"bucket,omitempty"`
	Key           string      `json:"key,omitempty"`
	CopySource    *CopySource `json:"copy_source,omitempty"`
	PartNumber    int32       `json:"part_number,omitempty"`
	UploadId      string      `json:"upload_id,omitempty"`
	StartPosition int64       `json:"start_position,omitempty"`
	PartSize      int64       `json:"part_size,omitempty"`
}
type UploadPartCopyOutput struct {
	BucketKeyEnabled     bool              `json:"bucket_key_enabled,omitempty"`
	CopyPartResult       *CopyPartResult   `json:"copy_part_result,omitempty"`
	CopySourceVersionId  string            `json:"copy_source_version_id,omitempty"`
	RequestCharged       string            `json:"request_charged,omitempty"`
	SSECustomerAlgorithm string            `json:"sse_customer_algorithm,omitempty"`
	SSECustomerKeyMD5    string            `json:"sse_customer_key_md5,omitempty"`
	SSEKMSKeyId          string            `json:"sse_kms_key_id,omitempty"`
	ServerSideEncryption string            `json:"server_side_encryption,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}
type CopyPartResult struct {
	ETag         string `json:"etag,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
}

type CompleteMultipartUploadInput struct {
	Bucket              string                    `json:"bucket,omitempty"`
	Key                 string                    `json:"key,omitempty"`
	UploadId            string                    `json:"upload_id,omitempty"`
	RequestPayer        string                    `json:"request_payer,omitempty"`
	ExpectedBucketOwner string                    `json:"expected_bucket_owner,omitempty"`
	MultipartUpload     *CompletedMultipartUpload `json:"multipart_upload,omitempty"`
}
type CompletedMultipartUpload struct {
	Parts []*CompletedPart `json:"parts,omitempty"`
}
type CompletedPart struct {
	ETag       string `json:"etag,omitempty"`
	PartNumber int32  `json:"part_number,omitempty"`
}
type CompleteMultipartUploadOutput struct {
	Bucket               string            `json:"bucket,omitempty"`
	Key                  string            `json:"key,omitempty"`
	BucketKeyEnabled     bool              `json:"bucket_key_enabled,omitempty"`
	ETag                 string            `json:"etag,omitempty"`
	Expiration           string            `json:"expiration,omitempty"`
	Location             string            `json:"location,omitempty"`
	RequestCharged       string            `json:"request_charged,omitempty"`
	SSEKMSKeyId          string            `json:"sse_kms_keyId,omitempty"`
	ServerSideEncryption string            `json:"server_side_encryption,omitempty"`
	VersionId            string            `json:"version_id,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

type AbortMultipartUploadInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Key                 string `json:"key,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
	RequestPayer        string `json:"request_payer,omitempty"`
	UploadId            string `json:"upload_id,omitempty"`
}
type AbortMultipartUploadOutput struct {
	RequestCharged string            `json:"request_charged,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type ListMultipartUploadsInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Delimiter           string `json:"delimiter,omitempty"`
	EncodingType        string `json:"encoding_type,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
	KeyMarker           string `json:"key_marker,omitempty"`
	MaxUploads          int64  `json:"max_uploads,omitempty"`
	Prefix              string `json:"prefix,omitempty"`
	UploadIdMarker      string `json:"upload_id_marker,omitempty"`
}
type ListMultipartUploadsOutput struct {
	Bucket             string             `json:"bucket,omitempty"`
	CommonPrefixes     []string           `json:"common_prefixes,omitempty"`
	Delimiter          string             `json:"delimiter,omitempty"`
	EncodingType       string             `json:"encoding_type,omitempty"`
	IsTruncated        bool               `json:"is_truncated,omitempty"`
	KeyMarker          string             `json:"key_marker,omitempty"`
	MaxUploads         int32              `json:"max_uploads,omitempty"`
	NextKeyMarker      string             `json:"next_key_marker,omitempty"`
	NextUploadIDMarker string             `json:"next_upload_id_marker,omitempty"`
	Prefix             string             `json:"prefix,omitempty"`
	UploadIDMarker     string             `json:"upload_id_marker,omitempty"`
	Uploads            []*MultipartUpload `json:"uploads,omitempty"`
	Metadata           map[string]string  `json:"metadata,omitempty"`
}
type MultipartUpload struct {
	Initiated    int64      `json:"initiated,omitempty"`
	Initiator    *Initiator `json:"initiator,omitempty"`
	Key          string     `json:"key,omitempty"`
	Owner        *Owner     `json:"owner,omitempty"`
	StorageClass string     `json:"storage_class,omitempty"`
	UploadId     string     `json:"upload_id,omitempty"`
}
type Initiator struct {
	DisplayName string `json:"display_name,omitempty"`
	ID          string `json:"id,omitempty"`
}

type ListObjectVersionsInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Delimiter           string `json:"delimiter,omitempty"`
	EncodingType        string `json:"encoding_type,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
	KeyMarker           string `json:"key_marker,omitempty"`
	MaxKeys             int32  `json:"max_keys,omitempty"`
	Prefix              string `json:"prefix,omitempty"`
	VersionIdMarker     string `json:"version_id_marker,omitempty"`
}
type ListObjectVersionsOutput struct {
	CommonPrefixes      []string             `json:"common_prefixes,omitempty"`
	DeleteMarkers       []*DeleteMarkerEntry `json:"delete_markers,omitempty"`
	Delimiter           string               `json:"delimiter,omitempty"`
	EncodingType        string               `json:"encoding_type,omitempty"`
	IsTruncated         bool                 `json:"is_truncated,omitempty"`
	KeyMarker           string               `json:"key_marker,omitempty"`
	MaxKeys             int32                `json:"max_keys,omitempty"`
	Name                string               `json:"name,omitempty"`
	NextKeyMarker       string               `json:"next_key_marker,omitempty"`
	NextVersionIdMarker string               `json:"next_version_id_marker,omitempty"`
	Prefix              string               `json:"prefix,omitempty"`
	VersionIdMarker     string               `json:"version_id_marker,omitempty"`
	Versions            []*ObjectVersion     `json:"versions,omitempty"`
	Metadata            map[string]string    `json:"metadata,omitempty"`
}
type DeleteMarkerEntry struct {
	IsLatest     bool   `json:"is_latest,omitempty"`
	Key          string `json:"key,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
	Owner        *Owner `json:"owner,omitempty"`
	VersionId    string `json:"version_id,omitempty"`
}
type ObjectVersion struct {
	ETag         string `json:"etag,omitempty"`
	IsLatest     bool   `json:"is_latest,omitempty"`
	Key          string `json:"key,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
	Owner        *Owner `json:"owner,omitempty"`
	Size         int64  `json:"size,omitempty"`
	StorageClass string `json:"storage_class,omitempty"`
	VersionId    string `json:"version_id,omitempty"`
}

type HeadObjectInput struct {
	Bucket               string `json:"bucket,omitempty"`
	Key                  string `json:"key,omitempty"`
	ChecksumMode         string `json:"checksum_mode,omitempty"`
	ExpectedBucketOwner  string `json:"expected_bucket_owner,omitempty"`
	IfMatch              string `json:"if_match,omitempty"`
	IfModifiedSince      int64  `json:"if_modified_since,omitempty"`
	IfNoneMatch          string `json:"if_none_match,omitempty"`
	IfUnmodifiedSince    int64  `json:"if_unmodified_since,omitempty"`
	PartNumber           int32  `json:"part_number,omitempty"`
	RequestPayer         string `json:"request_payer,omitempty"`
	SSECustomerAlgorithm string `json:"sse_customer_algorithm,omitempty"`
	SSECustomerKey       string `json:"sse_customer_key,omitempty"`
	SSECustomerKeyMD5    string `json:"sse_customer_key_md5,omitempty"`
	VersionId            string `json:"version_id,omitempty"`
	WithDetails          bool   `json:"with_details,omitempty"`
}
type HeadObjectOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata map[string]string `json:"result_metadata,omitempty"`
}

type IsObjectExistInput struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
}
type IsObjectExistOutput struct {
	FileExist bool              `json:"file_exist,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type SignURLInput struct {
	Bucket       string `json:"bucket,omitempty"`
	Key          string `json:"key,omitempty"`
	Method       string `json:"method,omitempty"`
	ExpiredInSec int64  `json:"expired_in_sec,omitempty"`
}
type SignURLOutput struct {
	SignedUrl string `json:"signed_url,omitempty"`
}

type UpdateBandwidthRateLimitInput struct {
	// The average upload/download bandwidth rate limit in bits per second.
	AverageRateLimitInBitsPerSec int64 `json:"average_rate_limit_in_bits_per_sec,omitempty"`
	//Resource name of gateway
	GatewayResourceName string `json:"gateway_resource_name,omitempty"`
}

type AppendObjectInput struct {
	DataStream           io.Reader
	Bucket               string            `json:"bucket,omitempty"`
	Key                  string            `json:"key,omitempty"`
	Position             int64             `json:"position,omitempty"`
	ACL                  string            `json:"acl,omitempty"`
	CacheControl         string            `json:"cache_control,omitempty"`
	ContentDisposition   string            `json:"content_disposition,omitempty"`
	ContentEncoding      string            `json:"content_encoding,omitempty"`
	ContentMd5           string            `json:"content_md5,omitempty"`
	Expires              int64             `json:"expires,omitempty"`
	StorageClass         string            `json:"storage_class,omitempty"`
	ServerSideEncryption string            `json:"server_side_encryption,omitempty"`
	Meta                 string            `json:"meta,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty"`
}
type AppendObjectOutput struct {
	AppendPosition int64             `json:"append_position,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type ListPartsInput struct {
	Bucket              string `json:"bucket,omitempty"`
	Key                 string `json:"key,omitempty"`
	ExpectedBucketOwner string `json:"expected_bucket_owner,omitempty"`
	MaxParts            int64  `json:"max_parts,omitempty"`
	PartNumberMarker    int64  `json:"part_number_marker,omitempty"`
	RequestPayer        string `json:"request_payer,omitempty"`
	UploadId            string `json:"upload_id,omitempty"`
}
type ListPartsOutput struct {
	Bucket               string            `json:"bucket,omitempty"`
	Key                  string            `json:"key,omitempty"`
	UploadId             string            `json:"upload_id,omitempty"`
	NextPartNumberMarker string            `json:"next_part_number_marker,omitempty"`
	MaxParts             int64             `json:"max_parts,omitempty"`
	IsTruncated          bool              `json:"is_truncated,omitempty"`
	Parts                []*Part           `json:"parts,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

type Part struct {
	Etag         string `json:"etag,omitempty"`
	LastModified int64  `json:"last_modified,omitempty"`
	PartNumber   int64  `json:"part_number,omitempty"`
	Size         int64  `json:"size,omitempty"`
}
