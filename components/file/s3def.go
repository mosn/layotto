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
	"io"
)

type Oss interface {
	InitConfig(context.Context, *FileConfig) error
	InitClient(context.Context, *InitRequest) error
	GetObject(context.Context, *GetObjectInput) (io.ReadCloser, error)
	PutObject(context.Context, *PutObjectInput) (*PutObjectOutput, error)
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
	IfModifiedSince            string `protobuf:"bytes,5,opt,name=if_modified_since,json=ifModifiedSince,proto3" json:"if_modified_since,omitempty"`
	IfNoneMatch                string `protobuf:"bytes,6,opt,name=if_none_match,json=ifNoneMatch,proto3" json:"if_none_match,omitempty"`
	IfUnmodifiedSince          string `protobuf:"bytes,7,opt,name=if_unmodified_since,json=ifUnmodifiedSince,proto3" json:"if_unmodified_since,omitempty"`
	Key                        string `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty"`
	PartNumber                 int64  `protobuf:"varint,9,opt,name=part_number,json=partNumber,proto3" json:"part_number,omitempty"`
	Range                      string `protobuf:"bytes,10,opt,name=range,proto3" json:"range,omitempty"`
	RequestPayer               string `protobuf:"bytes,11,opt,name=request_payer,json=requestPayer,proto3" json:"request_payer,omitempty"`
	ResponseCacheControl       string `protobuf:"bytes,12,opt,name=response_cache_control,json=responseCacheControl,proto3" json:"response_cache_control,omitempty"`
	ResponseContentDisposition string `protobuf:"bytes,13,opt,name=response_content_disposition,json=responseContentDisposition,proto3" json:"response_content_disposition,omitempty"`
	ResponseContentEncoding    string `protobuf:"bytes,14,opt,name=response_content_encoding,json=responseContentEncoding,proto3" json:"response_content_encoding,omitempty"`
	ResponseContentLanguage    string `protobuf:"bytes,15,opt,name=response_content_language,json=responseContentLanguage,proto3" json:"response_content_language,omitempty"`
	ResponseContentType        string `protobuf:"bytes,16,opt,name=response_content_type,json=responseContentType,proto3" json:"response_content_type,omitempty"`
	ResponseExpires            string `protobuf:"bytes,17,opt,name=response_expires,json=responseExpires,proto3" json:"response_expires,omitempty"`
	SseCustomerAlgorithm       string `protobuf:"bytes,18,opt,name=sse_customer_algorithm,json=sseCustomerAlgorithm,proto3" json:"sse_customer_algorithm,omitempty"`
	SseCustomerKey             string `protobuf:"bytes,19,opt,name=sse_customer_key,json=sseCustomerKey,proto3" json:"sse_customer_key,omitempty"`
	SseCustomerKeyMd5          string `protobuf:"bytes,20,opt,name=sse_customer_key_md5,json=sseCustomerKeyMd5,proto3" json:"sse_customer_key_md5,omitempty"`
	VersionId                  string `protobuf:"bytes,21,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
}

type PutObjectInput struct {
	DataStream       io.Reader
	Acl              string `protobuf:"bytes,2,opt,name=acl,proto3" json:"acl,omitempty"`
	Bucket           string `protobuf:"bytes,4,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Key              string `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	BucketKeyEnabled bool   `protobuf:"varint,6,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	CacheControl     string `protobuf:"bytes,7,opt,name=cache_control,json=cacheControl,proto3" json:"cache_control,omitempty"`
}

type PutObjectOutput struct {
	BucketKeyEnabled bool   `protobuf:"varint,1,opt,name=bucket_key_enabled,json=bucketKeyEnabled,proto3" json:"bucket_key_enabled,omitempty"`
	Etag             string `protobuf:"bytes,2,opt,name=etag,proto3" json:"etag,omitempty"`
}
