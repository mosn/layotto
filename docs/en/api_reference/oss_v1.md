

<a name="oss.proto"></a>

# oss.proto
<a name="top"></a>

This document is automaticallly generated from the [`.proto`](https://github.com/mosn/layotto/tree/main/spec/proto/runtime/v1) files.

The file defined base on s3 protocol, to get an in-depth walkthrough of this file, see:
https://docs.aws.amazon.com/s3/index.html
https://github.com/aws/aws-sdk-go-v2


<a name="spec.proto.extension.v1.ObjectStorageService"></a>

## [gRPC Service] ObjectStorageService
ObjectStorageService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| PutObject | [PutObjectInput](#spec.proto.extension.v1.PutObjectInput) stream | [PutObjectOutput](#spec.proto.extension.v1.PutObjectOutput) | Object CRUD API Adds an object to a bucket. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html |
| GetObject | [GetObjectInput](#spec.proto.extension.v1.GetObjectInput) | [GetObjectOutput](#spec.proto.extension.v1.GetObjectOutput) stream | Retrieves objects. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html |
| DeleteObject | [DeleteObjectInput](#spec.proto.extension.v1.DeleteObjectInput) | [DeleteObjectOutput](#spec.proto.extension.v1.DeleteObjectOutput) | Delete objects. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html |
| CopyObject | [CopyObjectInput](#spec.proto.extension.v1.CopyObjectInput) | [CopyObjectOutput](#spec.proto.extension.v1.CopyObjectOutput) | Creates a copy of an object that is already stored in oss server. Refer https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_CopyObject.html |
| DeleteObjects | [DeleteObjectsInput](#spec.proto.extension.v1.DeleteObjectsInput) | [DeleteObjectsOutput](#spec.proto.extension.v1.DeleteObjectsOutput) | Delete multiple objects from a bucket. Refer https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_DeleteObjects.html |
| ListObjects | [ListObjectsInput](#spec.proto.extension.v1.ListObjectsInput) | [ListObjectsOutput](#spec.proto.extension.v1.ListObjectsOutput) | Returns some or all (up to 1,000) of the objects in a bucket. Refer https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_ListObjects.html |
| HeadObject | [HeadObjectInput](#spec.proto.extension.v1.HeadObjectInput) | [HeadObjectOutput](#spec.proto.extension.v1.HeadObjectOutput) | The HEAD action retrieves metadata from an object without returning the object itself. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html |
| IsObjectExist | [IsObjectExistInput](#spec.proto.extension.v1.IsObjectExistInput) | [IsObjectExistOutput](#spec.proto.extension.v1.IsObjectExistOutput) | This action used to check if the file exists. |
| PutObjectTagging | [PutObjectTaggingInput](#spec.proto.extension.v1.PutObjectTaggingInput) | [PutObjectTaggingOutput](#spec.proto.extension.v1.PutObjectTaggingOutput) | Object Tagging API Sets the supplied tag-set to an object that already exists in a bucket. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObjectTagging.html |
| DeleteObjectTagging | [DeleteObjectTaggingInput](#spec.proto.extension.v1.DeleteObjectTaggingInput) | [DeleteObjectTaggingOutput](#spec.proto.extension.v1.DeleteObjectTaggingOutput) | Removes the entire tag set from the specified object. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjectTagging.html |
| GetObjectTagging | [GetObjectTaggingInput](#spec.proto.extension.v1.GetObjectTaggingInput) | [GetObjectTaggingOutput](#spec.proto.extension.v1.GetObjectTaggingOutput) | Returns the tag-set of an object. Refer https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_GetObjectTagging.html |
| GetObjectCannedAcl | [GetObjectCannedAclInput](#spec.proto.extension.v1.GetObjectCannedAclInput) | [GetObjectCannedAclOutput](#spec.proto.extension.v1.GetObjectCannedAclOutput) | Returns object canned acl. Refer https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#CannedACL |
| PutObjectCannedAcl | [PutObjectCannedAclInput](#spec.proto.extension.v1.PutObjectCannedAclInput) | [PutObjectCannedAclOutput](#spec.proto.extension.v1.PutObjectCannedAclOutput) | Set object canned acl. Refer https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#CannedACL |
| CreateMultipartUpload | [CreateMultipartUploadInput](#spec.proto.extension.v1.CreateMultipartUploadInput) | [CreateMultipartUploadOutput](#spec.proto.extension.v1.CreateMultipartUploadOutput) | Object Multipart Operation API Initiates a multipart upload and returns an upload ID. Refer https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_CreateMultipartUpload.html |
| UploadPart | [UploadPartInput](#spec.proto.extension.v1.UploadPartInput) stream | [UploadPartOutput](#spec.proto.extension.v1.UploadPartOutput) | Uploads a part in a multipart upload. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPart.html |
| UploadPartCopy | [UploadPartCopyInput](#spec.proto.extension.v1.UploadPartCopyInput) | [UploadPartCopyOutput](#spec.proto.extension.v1.UploadPartCopyOutput) | Uploads a part by copying data from an existing object as data source. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPartCopy.html |
| CompleteMultipartUpload | [CompleteMultipartUploadInput](#spec.proto.extension.v1.CompleteMultipartUploadInput) | [CompleteMultipartUploadOutput](#spec.proto.extension.v1.CompleteMultipartUploadOutput) | Completes a multipart upload by assembling previously uploaded parts. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_CompleteMultipartUpload.html |
| AbortMultipartUpload | [AbortMultipartUploadInput](#spec.proto.extension.v1.AbortMultipartUploadInput) | [AbortMultipartUploadOutput](#spec.proto.extension.v1.AbortMultipartUploadOutput) | This action aborts a multipart upload. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_AbortMultipartUpload.html |
| ListMultipartUploads | [ListMultipartUploadsInput](#spec.proto.extension.v1.ListMultipartUploadsInput) | [ListMultipartUploadsOutput](#spec.proto.extension.v1.ListMultipartUploadsOutput) | This action lists in-progress multipart uploads. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListMultipartUploads.html |
| ListParts | [ListPartsInput](#spec.proto.extension.v1.ListPartsInput) | [ListPartsOutput](#spec.proto.extension.v1.ListPartsOutput) | Lists the parts that have been uploaded for a specific multipart upload. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListParts.html |
| ListObjectVersions | [ListObjectVersionsInput](#spec.proto.extension.v1.ListObjectVersionsInput) | [ListObjectVersionsOutput](#spec.proto.extension.v1.ListObjectVersionsOutput) | Returns metadata about all versions of the objects in a bucket. Refer https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectVersions.html |
| SignURL | [SignURLInput](#spec.proto.extension.v1.SignURLInput) | [SignURLOutput](#spec.proto.extension.v1.SignURLOutput) | A presigned URL gives you access to the object identified in the URL, provided that the creator of the presigned URL has permissions to access that object. Refer https://docs.aws.amazon.com/AmazonS3/latest/userguide/PresignedUrlUploadObject.html |
| UpdateDownloadBandwidthRateLimit | [UpdateBandwidthRateLimitInput](#spec.proto.extension.v1.UpdateBandwidthRateLimitInput) | [.google.protobuf.Empty](#google.protobuf.Empty) | This action used to set download bandwidth limit speed. Refer https://github.com/aliyun/aliyun-oss-go-sdk/blob/master/oss/client.go#L2106 |
| UpdateUploadBandwidthRateLimit | [UpdateBandwidthRateLimitInput](#spec.proto.extension.v1.UpdateBandwidthRateLimitInput) | [.google.protobuf.Empty](#google.protobuf.Empty) | This action used to set upload bandwidth limit speed. Refer https://github.com/aliyun/aliyun-oss-go-sdk/blob/master/oss/client.go#L2096 |
| AppendObject | [AppendObjectInput](#spec.proto.extension.v1.AppendObjectInput) stream | [AppendObjectOutput](#spec.proto.extension.v1.AppendObjectOutput) | This action is used to append object. Refer https://help.aliyun.com/document_detail/31981.html or https://github.com/minio/minio-java/issues/980 |
| RestoreObject | [RestoreObjectInput](#spec.proto.extension.v1.RestoreObjectInput) | [RestoreObjectOutput](#spec.proto.extension.v1.RestoreObjectOutput) | Restores an archived copy of an object back. Refer https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_RestoreObject.html |

 <!-- end services -->


<a name="spec.proto.extension.v1.AbortMultipartUploadInput"></a>
<p align="right"><a href="#top">Top</a></p>

## AbortMultipartUploadInput
AbortMultipartUploadInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| upload_id | [string](#string) |  | Upload ID that identifies the multipart upload. This member is required. |






<a name="spec.proto.extension.v1.AbortMultipartUploadOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## AbortMultipartUploadOutput
AbortMultipartUploadOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |






<a name="spec.proto.extension.v1.AppendObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## AppendObjectInput
AppendObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| body | [bytes](#bytes) |  | Object content |
| position | [int64](#int64) |  | Append start position |
| acl | [string](#string) |  | Object ACL |
| cache_control | [string](#string) |  | Sets the Cache-Control header of the response. |
| content_disposition | [string](#string) |  | Sets the Content-Disposition header of the response |
| content_encoding | [string](#string) |  | Sets the Content-Encoding header of the response |
| content_md5 | [string](#string) |  | The base64-encoded 128-bit MD5 digest of the part data. |
| expires | [int64](#int64) |  | Sets the Expires header of the response |
| storage_class | [string](#string) |  | Provides storage class information of the object. Amazon S3 returns this header for all objects except for S3 Standard storage class objects. |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms). |
| meta | [string](#string) |  | Object metadata |
| tags | [AppendObjectInput.TagsEntry](#spec.proto.extension.v1.AppendObjectInput.TagsEntry) | repeated | Object tags |






<a name="spec.proto.extension.v1.AppendObjectInput.TagsEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## AppendObjectInput.TagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.AppendObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## AppendObjectOutput
AppendObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| append_position | [int64](#int64) |  | Next append position |






<a name="spec.proto.extension.v1.CompleteMultipartUploadInput"></a>
<p align="right"><a href="#top">Top</a></p>

## CompleteMultipartUploadInput
CompleteMultipartUploadInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| upload_id | [string](#string) |  | ID for the initiated multipart upload. This member is required. |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| expected_bucket_owner | [string](#string) |  | Expected bucket owner |
| multipart_upload | [CompletedMultipartUpload](#spec.proto.extension.v1.CompletedMultipartUpload) |  | The container for the multipart upload request information. |






<a name="spec.proto.extension.v1.CompleteMultipartUploadOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## CompleteMultipartUploadOutput
CompleteMultipartUploadOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| bucket_key_enabled | [bool](#bool) |  | Indicates whether the multipart upload uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS). |
| etag | [string](#string) |  | Entity tag that identifies the newly created object's data |
| expiration | [string](#string) |  | If the object expiration is configured, this will contain the expiration date (expiry-date) and rule ID (rule-id). The value of rule-id is URL-encoded. |
| location | [string](#string) |  | The URI that identifies the newly created object. |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| sse_kms_keyId | [string](#string) |  | If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object. |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms). |
| version_id | [string](#string) |  | Version ID of the newly created object, in case the bucket has versioning turned on. |






<a name="spec.proto.extension.v1.CompletedMultipartUpload"></a>
<p align="right"><a href="#top">Top</a></p>

## CompletedMultipartUpload
CompletedMultipartUpload


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parts | [CompletedPart](#spec.proto.extension.v1.CompletedPart) | repeated | Array of CompletedPart data types. |






<a name="spec.proto.extension.v1.CompletedPart"></a>
<p align="right"><a href="#top">Top</a></p>

## CompletedPart
CompletedPart


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| etag | [string](#string) |  | Entity tag returned when the part was uploaded. |
| part_number | [int32](#int32) |  | Part number that identifies the part. This is a positive integer between 1 and 10,000. |






<a name="spec.proto.extension.v1.CopyObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## CopyObjectInput
CopyObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The name of the destination bucket. When using this action with an access point This member is required. |
| key | [string](#string) |  | The key of the destination object. This member is required. |
| copy_source | [CopySource](#spec.proto.extension.v1.CopySource) |  | CopySource |
| tagging | [CopyObjectInput.TaggingEntry](#spec.proto.extension.v1.CopyObjectInput.TaggingEntry) | repeated | The tag-set for the object destination object this value must be used in conjunction with the TaggingDirective. The tag-set must be encoded as URL Query parameters. |
| expires | [int64](#int64) |  | The date and time at which the object is no longer cacheable. |
| metadata_directive | [string](#string) |  | Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request. |
| metadata | [CopyObjectInput.MetadataEntry](#spec.proto.extension.v1.CopyObjectInput.MetadataEntry) | repeated | A map of metadata to store with the object in S3. |






<a name="spec.proto.extension.v1.CopyObjectInput.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## CopyObjectInput.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.CopyObjectInput.TaggingEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## CopyObjectInput.TaggingEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.CopyObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## CopyObjectOutput
CopyObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| copy_object_result | [CopyObjectResult](#spec.proto.extension.v1.CopyObjectResult) |  | Container for all response elements. |
| version_id | [string](#string) |  | Version ID of the newly created copy. |
| expiration | [string](#string) |  | If the object expiration is configured, the response includes this header. |






<a name="spec.proto.extension.v1.CopyObjectResult"></a>
<p align="right"><a href="#top">Top</a></p>

## CopyObjectResult
CopyObjectResult


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| etag | [string](#string) |  | Returns the ETag of the new object. The ETag reflects only changes to the contents of an object, not its metadata. |
| last_modified | [int64](#int64) |  | Creation date of the object. |






<a name="spec.proto.extension.v1.CopyPartResult"></a>
<p align="right"><a href="#top">Top</a></p>

## CopyPartResult
CopyPartResult


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| etag | [string](#string) |  | Entity tag of the object. |
| last_modified | [int64](#int64) |  | Last modified time |






<a name="spec.proto.extension.v1.CopySource"></a>
<p align="right"><a href="#top">Top</a></p>

## CopySource
CopySource


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| copy_source_bucket | [string](#string) |  | source object bucket name |
| copy_source_key | [string](#string) |  | source object name |
| copy_source_version_id | [string](#string) |  | source object version |






<a name="spec.proto.extension.v1.CreateMultipartUploadInput"></a>
<p align="right"><a href="#top">Top</a></p>

## CreateMultipartUploadInput
CreateMultipartUploadInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| acl | [string](#string) |  | The canned ACL to apply to the object. This action is not supported by Amazon S3 on Outposts. |
| bucket_key_enabled | [bool](#bool) |  | Specifies whether Amazon S3 should use an S3 Bucket Key for object encryption with server-side encryption using AWS KMS (SSE-KMS). Setting this header to true causes Amazon S3 to use an S3 Bucket Key for object encryption with SSE-KMS. Specifying this header with a PUT action doesn’t affect bucket-level settings for S3 Bucket Key. |
| cache_control | [string](#string) |  | Specifies caching behavior along the request/reply chain |
| content_disposition | [string](#string) |  | Specifies presentational information for the object |
| content_encoding | [string](#string) |  | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. |
| content_language | [string](#string) |  | The language the content is in. |
| content_type | [string](#string) |  | A standard MIME type describing the format of the object data. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner. If the bucket is owned by a different account, the request fails with the HTTP status code 403 Forbidden (access denied). |
| expires | [int64](#int64) |  | The date and time at which the object is no longer cacheable. |
| grant_full_control | [string](#string) |  | Gives the grantee READ, READ_ACP, and WRITE_ACP permissions on the object. This action is not supported by Amazon S3 on Outposts. |
| grant_read | [string](#string) |  | Allows grantee to read the object data and its metadata. This action is not supported by Amazon S3 on Outposts. |
| grant_read_acp | [string](#string) |  | Allows grantee to read the object ACL. This action is not supported by Amazon S3 on Outposts. |
| grant_write_acp | [string](#string) |  | Allows grantee to write the ACL for the applicable object. This action is not supported by Amazon S3 on Outposts. |
| meta_data | [CreateMultipartUploadInput.MetaDataEntry](#spec.proto.extension.v1.CreateMultipartUploadInput.MetaDataEntry) | repeated | A map of metadata to store with the object |
| object_lock_legal_hold_status | [string](#string) |  | Specifies whether you want to apply a legal hold to the uploaded object |
| object_lock_mode | [string](#string) |  | Specifies the Object Lock mode that you want to apply to the uploaded object |
| object_lock_retain_until_date | [int64](#int64) |  | Specifies the date and time when you want the Object Lock to expire |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request |
| sse_customer_algorithm | [string](#string) |  | Specifies the algorithm to use to when encrypting the object (for example, AES256). |
| sse_customer_key | [string](#string) |  | Specifies the customer-provided encryption key to use in encrypting data |
| sse_customer_key_md5 | [string](#string) |  | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321 |
| sse_kms_encryption_context | [string](#string) |  | Specifies the Amazon Web Services KMS Encryption Context to use for object encryption |
| sse_kms_key_id | [string](#string) |  | Specifies the ID of the symmetric customer managed key to use for object encryption |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object |
| storage_class | [string](#string) |  | By default, oss store uses the STANDARD Storage Class to store newly created objects |
| tagging | [CreateMultipartUploadInput.TaggingEntry](#spec.proto.extension.v1.CreateMultipartUploadInput.TaggingEntry) | repeated | The tag-set for the object. The tag-set must be encoded as URL Query parameters. |
| website_redirect_location | [string](#string) |  | If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. |






<a name="spec.proto.extension.v1.CreateMultipartUploadInput.MetaDataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## CreateMultipartUploadInput.MetaDataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.CreateMultipartUploadInput.TaggingEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## CreateMultipartUploadInput.TaggingEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.CreateMultipartUploadOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## CreateMultipartUploadOutput
CreateMultipartUploadOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| abort_date | [int64](#int64) |  | If the bucket has a lifecycle rule configured with an action to abort incomplete multipart uploads and the prefix in the lifecycle rule matches the object name in the request, the response includes this header |
| abort_rule_id | [string](#string) |  | It identifies the applicable lifecycle configuration rule that defines the action to abort incomplete multipart uploads. |
| bucket_key_enabled | [bool](#bool) |  | Indicates whether the multipart upload uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS). |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| sse_customer_algorithm | [string](#string) |  | If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used. |
| sse_customer_key_md5 | [string](#string) |  | If server-side encryption with a customer-provided encryption key was requested, the response will include this header to provide round-trip message integrity verification of the customer-provided encryption key. |
| sse_kms_encryption_context | [string](#string) |  | If present, specifies the Amazon Web Services KMS Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs. |
| sse_kms_key_id | [string](#string) |  | If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object. |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms). |
| upload_id | [string](#string) |  | ID for the initiated multipart upload. |






<a name="spec.proto.extension.v1.Delete"></a>
<p align="right"><a href="#top">Top</a></p>

## Delete
Delete


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objects | [ObjectIdentifier](#spec.proto.extension.v1.ObjectIdentifier) | repeated | ObjectIdentifier |
| quiet | [bool](#bool) |  | Element to enable quiet mode for the request. When you add this element, you must set its value to true. |






<a name="spec.proto.extension.v1.DeleteMarkerEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteMarkerEntry
DeleteMarkerEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_latest | [bool](#bool) |  | Specifies whether the object is (true) or is not (false) the latest version of an object. |
| key | [string](#string) |  | Name of the object key. This member is required. |
| last_modified | [int64](#int64) |  | Date and time the object was last modified. |
| owner | [Owner](#spec.proto.extension.v1.Owner) |  | Owner |
| version_id | [string](#string) |  | Version ID of an object. |






<a name="spec.proto.extension.v1.DeleteObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectInput
DeleteObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name to which the DEL action was initiated This member is required. |
| key | [string](#string) |  | Object key for which the DEL action was initiated. This member is required. |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| version_id | [string](#string) |  | VersionId used to reference a specific version of the object. |






<a name="spec.proto.extension.v1.DeleteObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectOutput
DeleteObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delete_marker | [bool](#bool) |  | Specifies whether the versioned object that was permanently deleted was (true) or was not (false) a delete marker. |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| version_id | [string](#string) |  | Returns the version ID of the delete marker created as a result of the DELETE operation. |






<a name="spec.proto.extension.v1.DeleteObjectTaggingInput"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectTaggingInput
DeleteObjectTaggingInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the objects from which to remove the tags. |
| key | [string](#string) |  | The key that identifies the object in the bucket from which to remove all tags. This member is required. |
| version_id | [string](#string) |  | The versionId of the object that the tag-set will be removed from. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner. If the bucket is owned by a different account, the request fails with the HTTP status code 403 Forbidden (access denied). |






<a name="spec.proto.extension.v1.DeleteObjectTaggingOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectTaggingOutput
DeleteObjectTaggingOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version_id | [string](#string) |  | The versionId of the object the tag-set was removed from. |
| result_metadata | [DeleteObjectTaggingOutput.ResultMetadataEntry](#spec.proto.extension.v1.DeleteObjectTaggingOutput.ResultMetadataEntry) | repeated | Metadata pertaining to the operation's result. |






<a name="spec.proto.extension.v1.DeleteObjectTaggingOutput.ResultMetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectTaggingOutput.ResultMetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.DeleteObjectsInput"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectsInput
DeleteObjectsInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| delete | [Delete](#spec.proto.extension.v1.Delete) |  | Delete objects |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |






<a name="spec.proto.extension.v1.DeleteObjectsOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteObjectsOutput
DeleteObjectsOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deleted | [DeletedObject](#spec.proto.extension.v1.DeletedObject) | repeated | DeletedObject |






<a name="spec.proto.extension.v1.DeletedObject"></a>
<p align="right"><a href="#top">Top</a></p>

## DeletedObject
DeletedObject


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delete_marker | [bool](#bool) |  | Specifies whether the versioned object that was permanently deleted was (true) or was not (false) a delete marker. In a simple DELETE, this header indicates whether (true) or not (false) a delete marker was created. |
| delete_marker_version_id | [string](#string) |  | The version ID of the delete marker created as a result of the DELETE operation. If you delete a specific object version, the value returned by this header is the version ID of the object version deleted. |
| key | [string](#string) |  | The name of the deleted object. |
| version_id | [string](#string) |  | The version ID of the deleted object. |






<a name="spec.proto.extension.v1.GetObjectCannedAclInput"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectCannedAclInput
GetObjectCannedAclInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| version_id | [string](#string) |  | VersionId used to reference a specific version of the object |






<a name="spec.proto.extension.v1.GetObjectCannedAclOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectCannedAclOutput
GetObjectCannedAclOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| canned_acl | [string](#string) |  | Object CannedACL |
| owner | [Owner](#spec.proto.extension.v1.Owner) |  | Owner |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |






<a name="spec.proto.extension.v1.GetObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectInput
GetObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Key of the object to get This member is required |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| if_match | [string](#string) |  | Return the object only if its entity tag (ETag) is the same as the one specified |
| if_modified_since | [int64](#int64) |  | Return the object only if it has been modified since the specified time |
| if_none_match | [string](#string) |  | Return the object only if its entity tag (ETag) is different from the one specified |
| if_unmodified_since | [int64](#int64) |  | Return the object only if it has not been modified since the specified time |
| part_number | [int64](#int64) |  | Part number of the object being read. This is a positive integer between 1 and 10,000. Effectively performs a 'ranged' GET request for the part specified. Useful for downloading just a part of an object. |
| start | [int64](#int64) |  | Downloads the specified range bytes of an object start is used to specify the location where the file starts |
| end | [int64](#int64) |  | end is used to specify the location where the file end |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| response_cache_control | [string](#string) |  | Sets the Cache-Control header of the response. |
| response_content_disposition | [string](#string) |  | Sets the Content-Disposition header of the response |
| response_content_encoding | [string](#string) |  | Sets the Content-Encoding header of the response |
| response_content_language | [string](#string) |  | Sets the Content-Language header of the response |
| response_content_type | [string](#string) |  | Sets the Content-Type header of the response |
| response_expires | [string](#string) |  | Sets the Expires header of the response |
| sse_customer_algorithm | [string](#string) |  | Specifies the algorithm to use to when decrypting the object (for example,AES256) |
| sse_customer_key | [string](#string) |  | Specifies the customer-provided encryption key for Amazon S3 used to encrypt the data. This value is used to decrypt the object when recovering it and must match the one used when storing the data. The key must be appropriate for use with the algorithm specified in the x-amz-server-side-encryption-customer-algorithm header |
| sse_customer_key_md5 | [string](#string) |  | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321 Amazon S3 uses this header for a message integrity check to ensure that the encryption key was transmitted without error. |
| version_id | [string](#string) |  | VersionId used to reference a specific version of the object |
| accept_encoding | [string](#string) |  | Specify Accept-Encoding, aws not supported now |
| signed_url | [string](#string) |  | Specify the signed url of object, user can get object with signed url without ak、sk |






<a name="spec.proto.extension.v1.GetObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectOutput
GetObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| body | [bytes](#bytes) |  | Object data. |
| cache_control | [string](#string) |  | Specifies caching behavior along the request/reply chain. |
| content_disposition | [string](#string) |  | Specifies presentational information for the object. |
| content_encoding | [string](#string) |  | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. |
| content_language | [string](#string) |  | The language the content is in. |
| content_length | [int64](#int64) |  | Size of the body in bytes. |
| content_range | [string](#string) |  | The portion of the object returned in the response. |
| content_type | [string](#string) |  | A standard MIME type describing the format of the object data. |
| delete_marker | [bool](#bool) |  | Specifies whether the object retrieved was (true) or was not (false) a Delete Marker. If false, this response header does not appear in the response. |
| etag | [string](#string) |  | An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL. |
| expiration | [string](#string) |  | If the object expiration is configured (see PUT Bucket lifecycle), the response includes this header. It includes the expiry-date and rule-id key-value pairs providing object expiration information. The value of the rule-id is URL-encoded. |
| expires | [string](#string) |  | The date and time at which the object is no longer cacheable. |
| last_modified | [int64](#int64) |  | Creation date of the object. |
| version_id | [string](#string) |  | Version of the object. |
| tag_count | [int64](#int64) |  | The number of tags, if any, on the object. |
| storage_class | [string](#string) |  | Provides storage class information of the object. Amazon S3 returns this header for all objects except for S3 Standard storage class objects. |
| parts_count | [int64](#int64) |  | The count of parts this object has. This value is only returned if you specify partNumber in your request and the object was uploaded as a multipart upload. |
| metadata | [GetObjectOutput.MetadataEntry](#spec.proto.extension.v1.GetObjectOutput.MetadataEntry) | repeated | A map of metadata to store with the object in S3. Map keys will be normalized to lower-case. |






<a name="spec.proto.extension.v1.GetObjectOutput.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectOutput.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.GetObjectTaggingInput"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectTaggingInput
GetObjectTaggingInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object for which to get the tagging information. This member is required. |
| key | [string](#string) |  | Object key for which to get the tagging information. This member is required. |
| version_id | [string](#string) |  | The versionId of the object for which to get the tagging information. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner. If the bucket is owned by a different account, the request fails with the HTTP status code 403 Forbidden (access denied). |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |






<a name="spec.proto.extension.v1.GetObjectTaggingOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectTaggingOutput
GetObjectTaggingOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tags | [GetObjectTaggingOutput.TagsEntry](#spec.proto.extension.v1.GetObjectTaggingOutput.TagsEntry) | repeated | Contains the tag set. This member is required. |
| version_id | [string](#string) |  | The versionId of the object for which you got the tagging information. |
| result_metadata | [GetObjectTaggingOutput.ResultMetadataEntry](#spec.proto.extension.v1.GetObjectTaggingOutput.ResultMetadataEntry) | repeated | Metadata pertaining to the operation's result. |






<a name="spec.proto.extension.v1.GetObjectTaggingOutput.ResultMetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectTaggingOutput.ResultMetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.GetObjectTaggingOutput.TagsEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetObjectTaggingOutput.TagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.HeadObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## HeadObjectInput
HeadObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| checksum_mode | [string](#string) |  | To retrieve the checksum, this parameter must be enabled |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| if_match | [string](#string) |  | Return the object only if its entity tag (ETag) is the same as the one specified; otherwise, return a 412 (precondition failed) error. |
| if_modified_since | [int64](#int64) |  | Return the object only if it has been modified since the specified time; otherwise, return a 304 (not modified) error. |
| if_none_match | [string](#string) |  | Return the object only if its entity tag (ETag) is different from the one specified |
| if_unmodified_since | [int64](#int64) |  | Return the object only if it has not been modified since the specified time; |
| part_number | [int32](#int32) |  | Part number of the object being read. This is a positive integer between 1 and 10,000. Effectively performs a 'ranged' HEAD request for the part specified. Useful querying about the size of the part and the number of parts in this object. |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| sse_customer_algorithm | [string](#string) |  | Specifies the algorithm to use to when encrypting the object (for example, AES256). |
| sse_customer_key | [string](#string) |  | Specifies the customer-provided encryption key for Amazon S3 to use in encrypting data |
| sse_customer_key_md5 | [string](#string) |  | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. |
| version_id | [string](#string) |  | VersionId used to reference a specific version of the object. |
| with_details | [bool](#bool) |  | Return object details meta |






<a name="spec.proto.extension.v1.HeadObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## HeadObjectOutput
HeadObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result_metadata | [HeadObjectOutput.ResultMetadataEntry](#spec.proto.extension.v1.HeadObjectOutput.ResultMetadataEntry) | repeated | Metadata pertaining to the operation's result. |






<a name="spec.proto.extension.v1.HeadObjectOutput.ResultMetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## HeadObjectOutput.ResultMetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.Initiator"></a>
<p align="right"><a href="#top">Top</a></p>

## Initiator
Initiator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| display_name | [string](#string) |  | Initiator name |
| id | [string](#string) |  | Initiator id |






<a name="spec.proto.extension.v1.IsObjectExistInput"></a>
<p align="right"><a href="#top">Top</a></p>

## IsObjectExistInput
IsObjectExistInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| version_id | [string](#string) |  | Object version id |






<a name="spec.proto.extension.v1.IsObjectExistOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## IsObjectExistOutput
IsObjectExistOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file_exist | [bool](#bool) |  | Object exist or not |






<a name="spec.proto.extension.v1.ListMultipartUploadsInput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListMultipartUploadsInput
ListMultipartUploadsInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| delimiter | [string](#string) |  | Character you use to group keys. All keys that contain the same string between the prefix, if specified, and the first occurrence of the delimiter after the prefix are grouped under a single result element, CommonPrefixes. If you don't specify the prefix parameter, then the substring starts at the beginning of the key. The keys that are grouped under CommonPrefixes result element are not returned elsewhere in the response. |
| encoding_type | [string](#string) |  | Requests Amazon S3 to encode the object keys in the response and specifies the encoding method to use. An object key may contain any Unicode character; |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| key_marker | [string](#string) |  | Together with upload-id-marker, this parameter specifies the multipart upload after which listing should begin. If upload-id-marker is not specified, only the keys lexicographically greater than the specified key-marker will be included in the list. If upload-id-marker is specified, any multipart uploads for a key equal to the key-marker might also be included, provided those multipart uploads have upload IDs lexicographically greater than the specified upload-id-marker. |
| max_uploads | [int64](#int64) |  | Sets the maximum number of multipart uploads, from 1 to 1,000, to return in the response body. 1,000 is the maximum number of uploads that can be returned in a response. |
| prefix | [string](#string) |  | Lists in-progress uploads only for those keys that begin with the specified prefix. You can use prefixes to separate a bucket into different grouping of keys. (You can think of using prefix to make groups in the same way you'd use a folder in a file system.) |
| upload_id_marker | [string](#string) |  | Together with key-marker, specifies the multipart upload after which listing should begin. If key-marker is not specified, the upload-id-marker parameter is ignored. Otherwise, any multipart uploads for a key equal to the key-marker might be included in the list only if they have an upload ID lexicographically greater than the specified upload-id-marker. |






<a name="spec.proto.extension.v1.ListMultipartUploadsOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListMultipartUploadsOutput
ListMultipartUploadsOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| common_prefixes | [string](#string) | repeated | If you specify a delimiter in the request, then the result returns each distinct key prefix containing the delimiter in a CommonPrefixes element. |
| delimiter | [string](#string) |  | Contains the delimiter you specified in the request. If you don't specify a delimiter in your request, this element is absent from the response. |
| encoding_type | [string](#string) |  | Encoding type used by Amazon S3 to encode object keys in the response. |
| is_truncated | [bool](#bool) |  | Indicates whether the returned list of multipart uploads is truncated. A value of true indicates that the list was truncated. The list can be truncated if the number of multipart uploads exceeds the limit allowed or specified by max uploads. |
| key_marker | [string](#string) |  | The key at or after which the listing began. |
| max_uploads | [int32](#int32) |  | Maximum number of multipart uploads that could have been included in the response. |
| next_key_marker | [string](#string) |  | When a list is truncated, this element specifies the value that should be used for the key-marker request parameter in a subsequent request. |
| next_upload_id_marker | [string](#string) |  | When a list is truncated, this element specifies the value that should be used for the upload-id-marker request parameter in a subsequent request. |
| prefix | [string](#string) |  | When a prefix is provided in the request, this field contains the specified prefix. The result contains only keys starting with the specified prefix. |
| upload_id_marker | [string](#string) |  | Upload ID after which listing began. |
| uploads | [MultipartUpload](#spec.proto.extension.v1.MultipartUpload) | repeated | Container for elements related to a particular multipart upload. A response can contain zero or more Upload elements. |






<a name="spec.proto.extension.v1.ListObjectVersionsInput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListObjectVersionsInput
ListObjectVersionsInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| delimiter | [string](#string) |  | A delimiter is a character that you specify to group keys. All keys that contain the same string between the prefix and the first occurrence of the delimiter are grouped under a single result element in CommonPrefixes. These groups are counted as one result against the max-keys limitation. These keys are not returned elsewhere in the response. |
| encoding_type | [string](#string) |  | Requests Amazon S3 to encode the object keys in the response and specifies the encoding method to use. An object key may contain any Unicode character; |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| key_marker | [string](#string) |  | Specifies the key to start with when listing objects in a bucket. |
| max_keys | [int64](#int64) |  | Sets the maximum number of keys returned in the response. By default the action returns up to 1,000 key names. The response might contain fewer keys but will never contain more. If additional keys satisfy the search criteria, but were not returned because max-keys was exceeded, the response contains true. To return the additional keys, see key-marker and version-id-marker. |
| prefix | [string](#string) |  | Use this parameter to select only those keys that begin with the specified prefix. You can use prefixes to separate a bucket into different groupings of keys. (You can think of using prefix to make groups in the same way you'd use a folder in a file system.) You can use prefix with delimiter to roll up numerous objects into a single result under CommonPrefixes. |
| version_id_marker | [string](#string) |  | Specifies the object version you want to start listing from. |






<a name="spec.proto.extension.v1.ListObjectVersionsOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListObjectVersionsOutput
ListObjectVersionsOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| common_prefixes | [string](#string) | repeated | All of the keys rolled up into a common prefix count as a single return when calculating the number of returns. |
| delete_markers | [DeleteMarkerEntry](#spec.proto.extension.v1.DeleteMarkerEntry) | repeated | Container for an object that is a delete marker. |
| delimiter | [string](#string) |  | The delimiter grouping the included keys. |
| encoding_type | [string](#string) |  | Encoding type used by Amazon S3 to encode object key names in the XML response. |
| is_truncated | [bool](#bool) |  | A flag that indicates whether Amazon S3 returned all of the results that satisfied the search criteria |
| key_marker | [string](#string) |  | Marks the last key returned in a truncated response. |
| max_keys | [int64](#int64) |  | Specifies the maximum number of objects to return |
| name | [string](#string) |  | The bucket name. |
| next_key_marker | [string](#string) |  | When the number of responses exceeds the value of MaxKeys, NextKeyMarker specifies the first key not returned that satisfies the search criteria |
| next_version_id_marker | [string](#string) |  | When the number of responses exceeds the value of MaxKeys, NextVersionIdMarker specifies the first object version not returned that satisfies the search criteria. |
| prefix | [string](#string) |  | Selects objects that start with the value supplied by this parameter. |
| version_id_marker | [string](#string) |  | Marks the last version of the key returned in a truncated response. |
| versions | [ObjectVersion](#spec.proto.extension.v1.ObjectVersion) | repeated | Container for version information. |






<a name="spec.proto.extension.v1.ListObjectsInput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListObjectsInput
ListObjectsInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| delimiter | [string](#string) |  | A delimiter is a character you use to group keys. |
| encoding_type | [string](#string) |  | Requests Amazon S3 to encode the object keys in the response and specifies the encoding method to use. An object key may contain any Unicode character; however, XML 1.0 parser cannot parse some characters, such as characters with an ASCII value from 0 to 10. For characters that are not supported in XML 1.0, you can add this parameter to request that Amazon S3 encode the keys in the response. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner. If the bucket is owned by a different account, the request fails with the HTTP status code 403 Forbidden (access denied). |
| marker | [string](#string) |  | Marker is where you want Amazon S3 to start listing from. Amazon S3 starts listing after this specified key. Marker can be any key in the bucket. |
| maxKeys | [int32](#int32) |  | Sets the maximum number of keys returned in the response. By default the action returns up to 1,000 key names. The response might contain fewer keys but will never contain more. |
| prefix | [string](#string) |  | Limits the response to keys that begin with the specified prefix. |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |






<a name="spec.proto.extension.v1.ListObjectsOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListObjectsOutput
ListObjectsOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| common_prefixes | [string](#string) | repeated | CommonPrefixes |
| contents | [Object](#spec.proto.extension.v1.Object) | repeated | Objects contents |
| delimiter | [string](#string) |  | Causes keys that contain the same string between the prefix and the first occurrence of the delimiter to be rolled up into a single result element in the CommonPrefixes collection. These rolled-up keys are not returned elsewhere in the response. Each rolled-up result counts as only one return against the MaxKeys value. |
| encoding_type | [string](#string) |  | Encoding type used by Amazon S3 to encode object keys in the response. |
| is_truncated | [bool](#bool) |  | A flag that indicates whether Amazon S3 returned all of the results that satisfied the search criteria. |
| marker | [string](#string) |  | Indicates where in the bucket listing begins. Marker is included in the response if it was sent with the request. |
| max_keys | [int32](#int32) |  | The maximum number of keys returned in the response body. |
| name | [string](#string) |  | The bucket name. |
| next_marker | [string](#string) |  | When response is truncated (the IsTruncated element value in the response is true), you can use the key name in this field as marker in the subsequent request to get next set of objects. |
| prefix | [string](#string) |  | Keys that begin with the indicated prefix. |






<a name="spec.proto.extension.v1.ListPartsInput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListPartsInput
ListPartsInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| max_parts | [int64](#int64) |  | Sets the maximum number of parts to return |
| part_number_marker | [int64](#int64) |  | Specifies the part after which listing should begin. Only parts with higher part numbers will be listed. |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| upload_id | [string](#string) |  | Upload ID identifying the multipart upload whose parts are being listed. |






<a name="spec.proto.extension.v1.ListPartsOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## ListPartsOutput
ListPartsOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| upload_id | [string](#string) |  | Upload ID identifying the multipart upload whose parts are being listed. |
| next_part_number_marker | [string](#string) |  | When a list is truncated, this element specifies the last part in the list, as well as the value to use for the part-number-marker request parameter in a subsequent request. |
| max_parts | [int64](#int64) |  | Maximum number of parts that were allowed in the response. |
| is_truncated | [bool](#bool) |  | Indicates whether the returned list of parts is truncated. A true value indicates that the list was truncated. A list can be truncated if the number of parts exceeds the limit returned in the MaxParts element. |
| parts | [Part](#spec.proto.extension.v1.Part) | repeated | Container for elements related to a particular part. A response can contain zero or more Part elements. |






<a name="spec.proto.extension.v1.MultipartUpload"></a>
<p align="right"><a href="#top">Top</a></p>

## MultipartUpload
MultipartUpload


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| initiated | [int64](#int64) |  | Date and time at which the multipart upload was initiated. |
| initiator | [Initiator](#spec.proto.extension.v1.Initiator) |  | Identifies who initiated the multipart upload. |
| key | [string](#string) |  | Name of the object key. This member is required. |
| owner | [Owner](#spec.proto.extension.v1.Owner) |  | Specifies the owner of the object that is part of the multipart upload. |
| storage_class | [string](#string) |  | The class of storage used to store the object. |
| upload_id | [string](#string) |  | Upload ID that identifies the multipart upload. |






<a name="spec.proto.extension.v1.Object"></a>
<p align="right"><a href="#top">Top</a></p>

## Object
Object


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| etag | [string](#string) |  | The entity tag is a hash of the object |
| key | [string](#string) |  | The name that you assign to an object. You use the object key to retrieve the object. |
| last_modified | [int64](#int64) |  | Creation date of the object. |
| owner | [Owner](#spec.proto.extension.v1.Owner) |  | The owner of the object |
| size | [int64](#int64) |  | Size in bytes of the object |
| storage_class | [string](#string) |  | The class of storage used to store the object. |






<a name="spec.proto.extension.v1.ObjectIdentifier"></a>
<p align="right"><a href="#top">Top</a></p>

## ObjectIdentifier
ObjectIdentifier


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Key name of the object. This member is required. |
| version_id | [string](#string) |  | VersionId for the specific version of the object to delete. |






<a name="spec.proto.extension.v1.ObjectVersion"></a>
<p align="right"><a href="#top">Top</a></p>

## ObjectVersion
ObjectVersion


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| etag | [string](#string) |  | The entity tag is an MD5 hash of that version of the object. |
| is_latest | [bool](#bool) |  | Specifies whether the object is (true) or is not (false) the latest version of an object. |
| key | [string](#string) |  | Name of the object key. This member is required. |
| last_modified | [int64](#int64) |  | Date and time the object was last modified. |
| owner | [Owner](#spec.proto.extension.v1.Owner) |  | Specifies the owner of the object. |
| size | [int64](#int64) |  | Size in bytes of the object. |
| storage_class | [string](#string) |  | The class of storage used to store the object. |
| version_id | [string](#string) |  | Version ID of an object. |






<a name="spec.proto.extension.v1.Owner"></a>
<p align="right"><a href="#top">Top</a></p>

## Owner
Owner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| display_name | [string](#string) |  | Owner display name |
| id | [string](#string) |  | Owner id |






<a name="spec.proto.extension.v1.Part"></a>
<p align="right"><a href="#top">Top</a></p>

## Part
Part


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| etag | [string](#string) |  | Part Etag |
| last_modified | [int64](#int64) |  | Last modified time |
| part_number | [int64](#int64) |  | Part number |
| size | [int64](#int64) |  | Part size |






<a name="spec.proto.extension.v1.PutObjectCannedAclInput"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectCannedAclInput
PutObjectCannedAclInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| acl | [string](#string) |  | The canned ACL to apply to the object |
| version_id | [string](#string) |  | VersionId used to reference a specific version of the object. |






<a name="spec.proto.extension.v1.PutObjectCannedAclOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectCannedAclOutput
PutObjectCannedAclOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_charged | [string](#string) |  | Request charged |






<a name="spec.proto.extension.v1.PutObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectInput
PutObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| body | [bytes](#bytes) |  | Object data. |
| bucket | [string](#string) |  | The bucket name to which the PUT action was initiated This member is required. |
| key | [string](#string) |  | Object key for which the PUT action was initiated. This member is required. |
| acl | [string](#string) |  | The canned ACL to apply to the object,different oss provider have different acl type |
| bucket_key_enabled | [bool](#bool) |  | Indicates whether the multipart upload uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS). |
| cache_control | [string](#string) |  | Can be used to specify caching behavior along the request/reply chain. |
| content_disposition | [string](#string) |  | Specifies presentational information for the object. For more information, see http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1 (http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1). |
| content_encoding | [string](#string) |  | Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. For more information, see http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11 (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11). |
| expires | [int64](#int64) |  | The date and time at which the object is no longer cacheable. For more information, see http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21 (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21). |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms). |
| signed_url | [string](#string) |  | Specify the signed url of object, user can put object with signed url without ak、sk |
| meta | [PutObjectInput.MetaEntry](#spec.proto.extension.v1.PutObjectInput.MetaEntry) | repeated | A map of metadata to store with the object in S3. |
| tagging | [PutObjectInput.TaggingEntry](#spec.proto.extension.v1.PutObjectInput.TaggingEntry) | repeated | The tag-set for the object. The tag-set must be encoded as URL Query parameters. |






<a name="spec.proto.extension.v1.PutObjectInput.MetaEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectInput.MetaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.PutObjectInput.TaggingEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectInput.TaggingEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.PutObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectOutput
PutObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket_key_enabled | [bool](#bool) |  | Indicates whether the uploaded object uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS). |
| etag | [string](#string) |  | Entity tag for the uploaded object. |
| expiration | [string](#string) |  | If the expiration is configured for the object |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| version_id | [string](#string) |  | Version of the object. |






<a name="spec.proto.extension.v1.PutObjectTaggingInput"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectTaggingInput
PutObjectTaggingInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required. |
| key | [string](#string) |  | Name of the object key. This member is required. |
| tags | [PutObjectTaggingInput.TagsEntry](#spec.proto.extension.v1.PutObjectTaggingInput.TagsEntry) | repeated | Container for the TagSet and Tag elements |
| version_id | [string](#string) |  | The versionId of the object that the tag-set will be added to. |






<a name="spec.proto.extension.v1.PutObjectTaggingInput.TagsEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectTaggingInput.TagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.PutObjectTaggingOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectTaggingOutput
PutObjectTaggingOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version_id | [string](#string) |  | The versionId of the object the tag-set was added to. |
| result_metadata | [PutObjectTaggingOutput.ResultMetadataEntry](#spec.proto.extension.v1.PutObjectTaggingOutput.ResultMetadataEntry) | repeated | Metadata pertaining to the operation's result. |






<a name="spec.proto.extension.v1.PutObjectTaggingOutput.ResultMetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## PutObjectTaggingOutput.ResultMetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.extension.v1.RestoreObjectInput"></a>
<p align="right"><a href="#top">Top</a></p>

## RestoreObjectInput
RestoreObjectInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| version_id | [string](#string) |  | VersionId used to reference a specific version of the object. |






<a name="spec.proto.extension.v1.RestoreObjectOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## RestoreObjectOutput
RestoreObjectOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| restore_output_path | [string](#string) |  | Indicates the path in the provided S3 output location where Select results will be restored to. |






<a name="spec.proto.extension.v1.SignURLInput"></a>
<p align="right"><a href="#top">Top</a></p>

## SignURLInput
SignURLInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| method | [string](#string) |  | the method for sign url, eg. GET、POST |
| expired_in_sec | [int64](#int64) |  | expire time of the sign url |






<a name="spec.proto.extension.v1.SignURLOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## SignURLOutput
SignURLOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signed_url | [string](#string) |  | Object signed url |






<a name="spec.proto.extension.v1.UpdateBandwidthRateLimitInput"></a>
<p align="right"><a href="#top">Top</a></p>

## UpdateBandwidthRateLimitInput
UpdateBandwidthRateLimitInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| average_rate_limit_in_bits_per_sec | [int64](#int64) |  | The average upload/download bandwidth rate limit in bits per second. |
| gateway_resource_name | [string](#string) |  | Resource name of gateway |






<a name="spec.proto.extension.v1.UploadPartCopyInput"></a>
<p align="right"><a href="#top">Top</a></p>

## UploadPartCopyInput
UploadPartCopyInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| copy_source | [CopySource](#spec.proto.extension.v1.CopySource) |  | CopySource |
| part_number | [int32](#int32) |  | Part number of part being copied. This is a positive integer between 1 and 10,000. This member is required. |
| upload_id | [string](#string) |  | Upload ID identifying the multipart upload whose part is being copied. This member is required. |
| start_position | [int64](#int64) |  | The range of bytes to copy from the source object.bytes=start_position-part_size |
| part_size | [int64](#int64) |  | Part size |






<a name="spec.proto.extension.v1.UploadPartCopyOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## UploadPartCopyOutput
UploadPartCopyOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket_key_enabled | [bool](#bool) |  | Indicates whether the multipart upload uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS). |
| copy_part_result | [CopyPartResult](#spec.proto.extension.v1.CopyPartResult) |  | Container for all response elements. |
| copy_source_version_id | [string](#string) |  | The version of the source object that was copied, if you have enabled versioning on the source bucket. |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| sse_customer_algorithm | [string](#string) |  | If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used. |
| sse_customer_key_md5 | [string](#string) |  | If server-side encryption with a customer-provided encryption key was requested, the response will include this header to provide round-trip message integrity verification of the customer-provided encryption key. |
| sse_kms_key_id | [string](#string) |  | If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object. |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms). |






<a name="spec.proto.extension.v1.UploadPartInput"></a>
<p align="right"><a href="#top">Top</a></p>

## UploadPartInput
UploadPartInput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of oss store. |
| bucket | [string](#string) |  | The bucket name containing the object This member is required |
| key | [string](#string) |  | Name of the object key. This member is required. |
| body | [bytes](#bytes) |  | Object data. |
| content_length | [int64](#int64) |  | Size of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically. |
| content_md5 | [string](#string) |  | The base64-encoded 128-bit MD5 digest of the part data. |
| expected_bucket_owner | [string](#string) |  | The account ID of the expected bucket owner |
| part_number | [int32](#int32) |  | Part number of part being uploaded. This is a positive integer between 1 and 10,000. This member is required. |
| request_payer | [string](#string) |  | Confirms that the requester knows that they will be charged for the request. |
| sse_customer_algorithm | [string](#string) |  | Specifies the algorithm to use to when encrypting the object (for example, AES256). |
| sse_customer_key | [string](#string) |  | Specifies the customer-provided encryption key for Amazon S3 to use in encrypting data |
| sse_customer_key_md5 | [string](#string) |  | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. |
| upload_id | [string](#string) |  | Upload ID identifying the multipart upload whose part is being uploaded. This member is required. |






<a name="spec.proto.extension.v1.UploadPartOutput"></a>
<p align="right"><a href="#top">Top</a></p>

## UploadPartOutput
UploadPartOutput


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket_key_enabled | [bool](#bool) |  | Indicates whether the multipart upload uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS). |
| etag | [string](#string) |  | Entity tag for the uploaded object. |
| request_charged | [string](#string) |  | If present, indicates that the requester was successfully charged for the request. |
| sse_customer_algorithm | [string](#string) |  | Specifies the algorithm to use to when encrypting the object (for example, AES256). |
| sse_customer_key_md5 | [string](#string) |  | Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. |
| sse_kms_key_id | [string](#string) |  | Specifies the ID of the symmetric customer managed key to use for object encryption |
| server_side_encryption | [string](#string) |  | The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms). |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

