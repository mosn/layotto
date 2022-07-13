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
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Prefix is an option to set prefix parameter
func Prefix(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.Prefix(value)
}

// KeyMarker is an option to set key-marker parameter
func KeyMarker(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.KeyMarker(value)
}

// MaxUploads is an option to set max-uploads parameter
func MaxUploads(value int) oss.Option {
	if value <= 0 {
		return nil
	}
	return oss.MaxUploads(value)
}

// Delimiter is an option to set delimiler parameter
func Delimiter(value string) oss.Option {
	if value == "" {
		return nil
	}

	return oss.Delimiter(value)
}

// UploadIDMarker is an option to set upload-id-marker parameter
func UploadIDMarker(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.UploadIDMarker(value)
}

// VersionId is an option to set versionId parameter
func VersionId(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.VersionId(value)
}

// ObjectACL is an option to set X-Oss-Object-Acl header
func ObjectACL(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ObjectACL(oss.ACLType(value))
}

// CacheControl is an option to set Cache-Control header
func CacheControl(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.CacheControl(value)
}

// ContentEncoding is an option to set Content-Encoding header
func ContentEncoding(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ContentEncoding(value)
}

// ACL is an option to set X-Oss-Acl header
func ACL(acl string) oss.Option {
	if acl == "" {
		return nil
	}
	return oss.ACL(oss.ACLType(acl))
}

// ContentType is an option to set Content-Type header
func ContentType(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ContentType(value)
}

// ContentLength is an option to set Content-Length header
func ContentLength(length int64) oss.Option {
	if length == 0 {
		return nil
	}
	return oss.ContentLength(length)
}

// ContentDisposition is an option to set Content-Disposition header
func ContentDisposition(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ContentDisposition(value)
}

// SetTagging is an option to set object tagging
func SetTagging(value map[string]string) oss.Option {
	if value == nil {
		return nil
	}
	tagging := oss.Tagging{}
	for k, v := range value {
		tag := oss.Tag{Key: k, Value: v}
		tagging.Tags = append(tagging.Tags, tag)
	}
	return oss.SetTagging(tagging)
}

// ContentLanguage is an option to set Content-Language header
func ContentLanguage(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ContentLanguage(value)
}

// ContentMD5 is an option to set Content-MD5 header
func ContentMD5(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ContentMD5(value)
}

// Expires is an option to set Expires header
func Expires(t int64) oss.Option {
	if t == 0 {
		return nil
	}
	ti := time.Unix(0, t)
	return oss.Expires(ti)
}

// AcceptEncoding is an option to set Accept-Encoding header
func AcceptEncoding(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.AcceptEncoding(value)
}

// IfModifiedSince is an option to set If-Modified-Since header
func IfModifiedSince(t int64) oss.Option {
	if t == 0 {
		return nil
	}
	ti := time.Unix(0, t)
	return oss.IfModifiedSince(ti)
}

// IfUnmodifiedSince is an option to set If-Unmodified-Since header
func IfUnmodifiedSince(t int64) oss.Option {
	if t == 0 {
		return nil
	}
	ti := time.Unix(0, t)
	return oss.IfUnmodifiedSince(ti)
}

// IfMatch is an option to set If-Match header
func IfMatch(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.IfNoneMatch(value)
}

// IfNoneMatch is an option to set IfNoneMatch header
func IfNoneMatch(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.IfNoneMatch(value)
}

// Range is an option to set Range header, [start, end]
func Range(start, end int64) oss.Option {
	if start == 0 && end == 0 {
		return nil
	}
	return oss.Range(start, end)
}

// CopySourceIfMatch is an option to set X-Oss-Copy-Source-If-Match header
func CopySourceIfMatch(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.CopySourceIfMatch(value)
}

// CopySourceIfNoneMatch is an option to set X-Oss-Copy-Source-If-None-Match header
func CopySourceIfNoneMatch(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.CopySourceIfNoneMatch(value)
}

// CopySourceIfModifiedSince is an option to set X-Oss-CopySource-If-Modified-Since header
func CopySourceIfModifiedSince(t int64) oss.Option {
	if t == 0 {
		return nil
	}
	tm := time.Unix(0, t)
	return oss.CopySourceIfModifiedSince(tm)
}

// CopySourceIfUnmodifiedSince is an option to set X-Oss-Copy-Source-If-Unmodified-Since header
func CopySourceIfUnmodifiedSince(t int64) oss.Option {
	if t == 0 {
		return nil
	}
	tm := time.Unix(0, t)
	return oss.CopySourceIfUnmodifiedSince(tm)
}

// MetadataDirective is an option to set X-Oss-Metadata-Directive header
func MetadataDirective(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.MetadataDirective(oss.MetadataDirectiveType(value))
}

// Meta is an option to set Meta header
func Meta(key, value string) oss.Option {
	return oss.Meta(key, value)
}

// ServerSideEncryption is an option to set X-Oss-Server-Side-Encryption header
func ServerSideEncryption(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ServerSideEncryption(value)
}

// ServerSideEncryptionKeyID is an option to set X-Oss-Server-Side-Encryption-Key-Id header
func ServerSideEncryptionKeyID(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ServerSideEncryptionKeyID(value)
}

// ServerSideDataEncryption is an option to set X-Oss-Server-Side-Data-Encryption header
func ServerSideDataEncryption(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ServerSideDataEncryption(value)
}

// SSECAlgorithm is an option to set X-Oss-Server-Side-Encryption-Customer-Algorithm header
func SSECAlgorithm(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.SSECAlgorithm(value)
}

// SSECKey is an option to set X-Oss-Server-Side-Encryption-Customer-Key header
func SSECKey(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.SSECKey(value)
}

// SSECKeyMd5 is an option to set X-Oss-Server-Side-Encryption-Customer-Key-Md5 header
func SSECKeyMd5(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.SSECKeyMd5(value)
}

// Origin is an option to set Origin header
func Origin(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.Origin(value)
}

// RangeBehavior  is an option to set Range value, such as "standard"
func RangeBehavior(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.RangeBehavior(value)
}

func PartHashCtxHeader(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.PartHashCtxHeader(value)
}

func PartMd5CtxHeader(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.PartMd5CtxHeader(value)
}

func PartHashCtxParam(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.PartHashCtxParam(value)
}

func PartMd5CtxParam(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.PartMd5CtxParam(value)
}

// Marker is an option to set marker parameter
func Marker(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.Marker(value)
}

// MaxKeys is an option to set maxkeys parameter
func MaxKeys(value int) oss.Option {
	if value == 0 {
		return nil
	}
	return oss.MaxKeys(value)
}

// EncodingType is an option to set encoding-type parameter
func EncodingType(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.EncodingType(value)
}

// VersionIdMarker is an option to set version-id-marker parameter
func VersionIdMarker(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.VersionIdMarker(value)
}

// MaxParts is an option to set max-parts parameter
func MaxParts(value int) oss.Option {
	if value == 0 {
		return nil
	}
	return oss.MaxParts(value)
}

// PartNumberMarker is an option to set part-number-marker parameter
func PartNumberMarker(value int) oss.Option {
	if value == 0 {
		return nil
	}
	return oss.PartNumberMarker(value)
}

// StorageClass bucket storage class
func StorageClass(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.StorageClass(oss.StorageClassType(value))
}

// ResponseContentType is an option to set response-content-type param
func ResponseContentType(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ResponseContentType(value)
}

// ResponseContentLanguage is an option to set response-content-language param
func ResponseContentLanguage(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ResponseContentLanguage(value)
}

// ResponseExpires is an option to set response-expires param
func ResponseExpires(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ResponseExpires(value)
}

// ResponseCacheControl is an option to set response-cache-control param
func ResponseCacheControl(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ResponseCacheControl(value)
}

// ResponseContentDisposition is an option to set response-content-disposition param
func ResponseContentDisposition(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ResponseContentDisposition(value)
}

// ResponseContentEncoding is an option to set response-content-encoding param
func ResponseContentEncoding(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.ResponseContentEncoding(value)
}

// Process is an option to set x-oss-process param
func Process(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.Process(value)
}

// TrafficLimitParam is a option to set x-oss-traffic-limit
func TrafficLimitParam(value int64) oss.Option {
	if value == 0 {
		return nil
	}
	return oss.TrafficLimitParam(value)
}

// SetHeader Allow users to set personalized http headers
func SetHeader(key string, value interface{}) oss.Option {
	return oss.SetHeader(key, value)
}

// AddParam Allow users to set personalized http params
func AddParam(key string, value interface{}) oss.Option {
	return oss.AddParam(key, value)
}

// RequestPayer is an option to set payer who pay for the request
func RequestPayer(value string) oss.Option {
	if value == "" {
		return nil
	}
	return oss.RequestPayer(oss.PayerType(value))
}
