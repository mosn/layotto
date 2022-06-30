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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionNil(t *testing.T) {
	assert.Nil(t, Prefix(""))
	assert.Nil(t, KeyMarker(""))
	assert.Nil(t, MaxUploads(0))
	assert.Nil(t, Delimiter(""))
	assert.Nil(t, UploadIDMarker(""))
	assert.Nil(t, VersionId(""))
	assert.Nil(t, ObjectACL(""))
	assert.Nil(t, CacheControl(""))
	assert.Nil(t, ContentEncoding(""))
	assert.Nil(t, ACL(""))
	assert.Nil(t, ContentType(""))
	assert.Nil(t, ContentLength(0))
	assert.Nil(t, ContentDisposition(""))
	assert.Nil(t, ContentLanguage(""))
	assert.Nil(t, ContentMD5(""))
	assert.Nil(t, Expires(0))
	assert.Nil(t, AcceptEncoding(""))
	assert.Nil(t, IfModifiedSince(0))
	assert.Nil(t, IfMatch(""))
	assert.Nil(t, IfNoneMatch(""))
	assert.Nil(t, Range(0, 0))
	assert.Nil(t, CopySourceIfMatch(""))
	assert.Nil(t, CopySourceIfNoneMatch(""))
	assert.Nil(t, CopySourceIfModifiedSince(0))
	assert.Nil(t, CopySourceIfUnmodifiedSince(0))
	assert.Nil(t, IfUnmodifiedSince(0))
	assert.Nil(t, MetadataDirective(""))
	assert.Nil(t, ServerSideEncryption(""))
	assert.Nil(t, ServerSideEncryptionKeyID(""))
	assert.Nil(t, ServerSideDataEncryption(""))
	assert.Nil(t, SSECAlgorithm(""))
	assert.Nil(t, SSECKey(""))
	assert.Nil(t, SSECKeyMd5(""))
	assert.Nil(t, Origin(""))
	assert.Nil(t, RangeBehavior(""))
	assert.Nil(t, PartHashCtxHeader(""))
	assert.Nil(t, PartNumberMarker(0))
	assert.Nil(t, PartHashCtxParam(""))
	assert.Nil(t, PartMd5CtxHeader(""))
	assert.Nil(t, PartMd5CtxParam(""))
	assert.Nil(t, Marker(""))
	assert.Nil(t, MaxKeys(0))
	assert.Nil(t, EncodingType(""))
	assert.Nil(t, VersionId(""))
	assert.Nil(t, VersionIdMarker(""))
	assert.Nil(t, MaxParts(0))
	assert.Nil(t, StorageClass(""))
	assert.Nil(t, ResponseContentDisposition(""))
	assert.Nil(t, ResponseCacheControl(""))
	assert.Nil(t, ResponseContentEncoding(""))
	assert.Nil(t, ResponseContentLanguage(""))
	assert.Nil(t, ResponseContentType(""))
	assert.Nil(t, ResponseExpires(""))
	assert.Nil(t, Process(""))
	assert.Nil(t, TrafficLimitParam(0))
	assert.Nil(t, RequestPayer(""))
}

func TestOptionNotNil(t *testing.T) {
	assert.NotNil(t, Prefix(" "))
	assert.NotNil(t, KeyMarker(" "))
	assert.NotNil(t, MaxUploads(1))
	assert.NotNil(t, Delimiter(" "))
	assert.NotNil(t, UploadIDMarker(" "))
	assert.NotNil(t, VersionId(" "))
	assert.NotNil(t, ObjectACL(" "))
	assert.NotNil(t, CacheControl(" "))
	assert.NotNil(t, ContentEncoding(" "))
	assert.NotNil(t, ACL(" "))
	assert.NotNil(t, ContentType(" "))
	assert.NotNil(t, ContentLength(1))
	assert.NotNil(t, ContentDisposition(" "))
	assert.NotNil(t, ContentLanguage(" "))
	assert.NotNil(t, ContentMD5(" "))
	assert.NotNil(t, Expires(1))
	assert.NotNil(t, AcceptEncoding(" "))
	assert.NotNil(t, IfModifiedSince(1))
	assert.NotNil(t, IfMatch(" "))
	assert.NotNil(t, IfNoneMatch(" "))
	assert.NotNil(t, Range(1, 1))
	assert.NotNil(t, CopySourceIfMatch(" "))
	assert.NotNil(t, CopySourceIfNoneMatch(" "))
	assert.NotNil(t, CopySourceIfModifiedSince(1))
	assert.NotNil(t, CopySourceIfUnmodifiedSince(1))
	assert.NotNil(t, IfUnmodifiedSince(1))
	assert.NotNil(t, MetadataDirective(" "))
	assert.NotNil(t, Meta(" ", " "))
	assert.NotNil(t, ServerSideEncryption(" "))
	assert.NotNil(t, ServerSideEncryptionKeyID(" "))
	assert.NotNil(t, ServerSideDataEncryption(" "))
	assert.NotNil(t, SSECAlgorithm(" "))
	assert.NotNil(t, SSECKey(" "))
	assert.NotNil(t, SSECKeyMd5(" "))
	assert.NotNil(t, Origin(" "))
	assert.NotNil(t, RangeBehavior(" "))
	assert.NotNil(t, PartHashCtxHeader(" "))
	assert.NotNil(t, PartNumberMarker(1))
	assert.NotNil(t, PartHashCtxParam(" "))
	assert.NotNil(t, PartMd5CtxHeader(" "))
	assert.NotNil(t, PartMd5CtxParam(" "))
	assert.NotNil(t, Marker(" "))
	assert.NotNil(t, MaxKeys(1))
	assert.NotNil(t, EncodingType(" "))
	assert.NotNil(t, VersionId(" "))
	assert.NotNil(t, VersionIdMarker(" "))
	assert.NotNil(t, MaxParts(1))
	assert.NotNil(t, StorageClass(" "))
	assert.NotNil(t, ResponseContentDisposition(" "))
	assert.NotNil(t, ResponseCacheControl(" "))
	assert.NotNil(t, ResponseContentEncoding(" "))
	assert.NotNil(t, ResponseContentLanguage(" "))
	assert.NotNil(t, ResponseContentType(" "))
	assert.NotNil(t, ResponseExpires(" "))
	assert.NotNil(t, Process(" "))
	assert.NotNil(t, TrafficLimitParam(1))
	assert.NotNil(t, SetHeader(" ", " "))
	assert.NotNil(t, AddParam(" ", " "))
	assert.NotNil(t, RequestPayer(" "))
}