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
	"testing"

	"mosn.io/layotto/components/pkg/utils"

	"mosn.io/layotto/components/oss"

	"mosn.io/pkg/buffer"

	"github.com/stretchr/testify/assert"

	l8oss "mosn.io/layotto/components/oss"
)

const (
	confWithoutUidAndBucket = `
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			`
)

func TestInitAliyunOss(t *testing.T) {
	a := &AliyunOSS{}
	client, err := a.getClient()
	assert.Equal(t, err, utils.ErrNotInitClient)
	assert.Nil(t, client)
	err = a.Init(context.TODO(), &l8oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte("hello")}})
	assert.Equal(t, err, l8oss.ErrInvalid)
	err = a.Init(context.TODO(), &l8oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.NotEqual(t, l8oss.ErrInvalid, err)
	assert.NotNil(t, a.client)

}

func TestAliyunOss(t *testing.T) {
	instance := NewAliyunOss()
	instance.Init(context.TODO(), &l8oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	appendObjectResp, err := instance.AppendObject(context.TODO(), &oss.AppendObjectInput{})
	assert.NotNil(t, err)
	assert.Nil(t, appendObjectResp)

	_, err = instance.AbortMultipartUpload(context.TODO(), &oss.AbortMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CompleteMultipartUpload(context.TODO(), &oss.CompleteMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CopyObject(context.TODO(), &oss.CopyObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.CreateMultipartUpload(context.TODO(), &oss.CreateMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.DeleteObject(context.TODO(), &oss.DeleteObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjects(context.TODO(), &oss.DeleteObjectsInput{
		Delete: &oss.Delete{},
	})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjectTagging(context.TODO(), &oss.DeleteObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.GetObject(context.TODO(), &oss.GetObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectCannedAcl(context.TODO(), &oss.GetObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectTagging(context.TODO(), &oss.GetObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.HeadObject(context.TODO(), &oss.HeadObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.IsObjectExist(context.TODO(), &oss.IsObjectExistInput{})
	assert.NotNil(t, err)

	_, err = instance.ListParts(context.TODO(), &oss.ListPartsInput{})
	assert.NotNil(t, err)

	_, err = instance.ListMultipartUploads(context.TODO(), &oss.ListMultipartUploadsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjects(context.TODO(), &oss.ListObjectsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjectVersions(context.TODO(), &oss.ListObjectVersionsInput{})
	assert.NotNil(t, err)

	stream := buffer.NewIoBufferString("hello")
	_, err = instance.PutObject(context.TODO(), &oss.PutObjectInput{DataStream: stream})
	assert.NotNil(t, err)
	_, err = instance.PutObjectCannedAcl(context.TODO(), &oss.PutObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.PutObjectTagging(context.TODO(), &oss.PutObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.RestoreObject(context.TODO(), &oss.RestoreObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.SignURL(context.TODO(), &oss.SignURLInput{})
	assert.NotNil(t, err)

	_, err = instance.UploadPartCopy(context.TODO(), &oss.UploadPartCopyInput{
		CopySource: &oss.CopySource{CopySourceBucket: "bucket", CopySourceKey: "key"},
	})
	assert.NotNil(t, err)

	_, err = instance.UploadPart(context.TODO(), &oss.UploadPartInput{})
	assert.NotNil(t, err)

	err = instance.UpdateDownloadBandwidthRateLimit(context.TODO(), &oss.UpdateBandwidthRateLimitInput{})
	assert.Nil(t, err)

	err = instance.UpdateUploadBandwidthRateLimit(context.TODO(), &oss.UpdateBandwidthRateLimitInput{})
	assert.Nil(t, err)

}
