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
	"testing"

	"mosn.io/pkg/buffer"

	"mosn.io/layotto/components/file"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file/factory"
)

const (
	confWithoutUidAndBucket = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
	confWithUid = `[
				{	
					"uid": "123",
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
	confWithUidAndBucket = `[
				{	
					"uid": "123",
					"buckets": ["bucket1","bucket2"],
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
)

func TestInitAliyunOss(t *testing.T) {
	NewAliyunOss()
	f := factory.GetInitFunc(DefaultClientInitFunc)
	clients, err := f([]byte("hello"), map[string]string{})
	assert.Equal(t, err, file.ErrInvalid)
	assert.Nil(t, clients)
	clients, err = f([]byte(confWithoutUidAndBucket), map[string]string{})
	assert.NotEqual(t, file.ErrInvalid, err)
	assert.NotNil(t, clients)
	client, ok := clients[""]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	clients, err = f([]byte(confWithUid), map[string]string{})
	assert.NotEqual(t, file.ErrInvalid, err)
	assert.NotNil(t, clients)
	client, ok = clients[""]
	assert.Equal(t, false, ok)
	assert.Nil(t, client)
	client, ok = clients["123"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	clients, err = f([]byte(confWithUidAndBucket), map[string]string{})
	assert.NotEqual(t, file.ErrInvalid, err)
	assert.NotNil(t, clients)
	client, ok = clients[""]
	assert.Equal(t, false, ok)
	assert.Nil(t, client)

	client, ok = clients["123"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	client, ok = clients["bucket1"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	client, ok = clients["bucket2"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

}

func TestAliyunOss(t *testing.T) {
	instance := NewAliyunOss()
	err := instance.InitConfig(context.TODO(), &file.FileConfig{Method: "", Metadata: []byte(confWithUidAndBucket)})
	assert.Nil(t, err)
	err = instance.InitClient(context.TODO(), &file.InitRequest{})
	assert.Nil(t, err)

	aliyun := instance.(*AliyunOSS)
	clientUid, _ := aliyun.selectClient("123", "")
	assert.Equal(t, clientUid, aliyun.client["123"])

	clientBucket1, _ := aliyun.selectClient("123", "bucket1")
	assert.Equal(t, clientBucket1, aliyun.client["bucket1"])

	clientBucket2, _ := aliyun.selectClient("123", "bucket2")
	assert.Equal(t, clientBucket2, aliyun.client["bucket2"])

	appendObjectResp, err := instance.AppendObject(context.TODO(), &file.AppendObjectInput{})
	assert.NotNil(t, err)
	assert.Nil(t, appendObjectResp)

	_, err = instance.AbortMultipartUpload(context.TODO(), &file.AbortMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CompleteMultipartUpload(context.TODO(), &file.CompleteMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CopyObject(context.TODO(), &file.CopyObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.CreateMultipartUpload(context.TODO(), &file.CreateMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.DeleteObject(context.TODO(), &file.DeleteObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjects(context.TODO(), &file.DeleteObjectsInput{
		Delete: &file.Delete{},
	})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjectTagging(context.TODO(), &file.DeleteObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.GetObject(context.TODO(), &file.GetObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectCannedAcl(context.TODO(), &file.GetObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectTagging(context.TODO(), &file.GetObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.HeadObject(context.TODO(), &file.HeadObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.IsObjectExist(context.TODO(), &file.IsObjectExistInput{})
	assert.NotNil(t, err)

	_, err = instance.ListParts(context.TODO(), &file.ListPartsInput{})
	assert.NotNil(t, err)

	_, err = instance.ListMultipartUploads(context.TODO(), &file.ListMultipartUploadsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjects(context.TODO(), &file.ListObjectsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjectVersions(context.TODO(), &file.ListObjectVersionsInput{})
	assert.NotNil(t, err)

	stream := buffer.NewIoBufferString("hello")
	_, err = instance.PutObject(context.TODO(), &file.PutObjectInput{DataStream: stream})
	assert.NotNil(t, err)
	_, err = instance.PutObjectCannedAcl(context.TODO(), &file.PutObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.PutObjectTagging(context.TODO(), &file.PutObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.RestoreObject(context.TODO(), &file.RestoreObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.SignURL(context.TODO(), &file.SignURLInput{})
	assert.NotNil(t, err)

	_, err = instance.UploadPartCopy(context.TODO(), &file.UploadPartCopyInput{
		CopySource: &file.CopySource{CopySourceBucket: "bucket", CopySourceKey: "key"},
	})
	assert.NotNil(t, err)

	_, err = instance.UploadPart(context.TODO(), &file.UploadPartInput{})
	assert.NotNil(t, err)

	err = instance.UpdateDownLoadBandwidthRateLimit(context.TODO(), &file.UpdateBandwidthRateLimitInput{})
	assert.Nil(t, err)

	err = instance.UpdateUpLoadBandwidthRateLimit(context.TODO(), &file.UpdateBandwidthRateLimitInput{})
	assert.Nil(t, err)

}
