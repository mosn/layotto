//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package oss

import (
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

func TestInit(t *testing.T) {
	data := `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": ["bucket1", "bucket2"]
				}
			]`
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	err := oss.Init(&fc)
	assert.Equal(t, err.Error(), "wrong config for alicloudOss")
	fc.Metadata = []byte(data)
	err = oss.Init(&fc)
	assert.Nil(t, err)
}

func TestSelectClient(t *testing.T) {
	ossObject := &AliCloudOSS{metadata: make(map[string]*OssMetadata), client: make(map[string]*oss.Client)}

	client, err := ossObject.selectClient()
	assert.Equal(t, err.Error(), "should specific endpoint in metadata")
	assert.Nil(t, client)

	client1 := &oss.Client{}
	ossObject.client["127.0.0.1"] = client1
	client, err = ossObject.selectClient()
	assert.Equal(t, client, client1)
	assert.Nil(t, err)

	client2 := &oss.Client{}
	ossObject.client["0.0.0.0"] = client2
	client, err = ossObject.selectClient()
	assert.Equal(t, err.Error(), "should specific endpoint in metadata")
	assert.Nil(t, client)
}

func TestSelectBucket(t *testing.T) {
	ossObject := &AliCloudOSS{metadata: make(map[string]*OssMetadata), client: make(map[string]*oss.Client)}

	bucketName, err := ossObject.selectBucket()
	assert.Equal(t, "", bucketName)
	assert.Equal(t, err.Error(), "no bucket configuration")

	metaData1 := &OssMetadata{Bucket: []string{"test", "test2"}}
	ossObject.metadata["0.0.0.0"] = metaData1
	bucketName, err = ossObject.selectBucket()
	assert.Equal(t, bucketName, "")
	assert.Equal(t, err.Error(), "should specific bucketKey in metadata")

	metaData2 := &OssMetadata{Bucket: []string{"test"}}
	ossObject.metadata["0.0.0.0"] = metaData2
	bucketName, err = ossObject.selectBucket()
	assert.Equal(t, bucketName, "test")
	assert.Nil(t, err)
}

func TestSelectClientAndBucket(t *testing.T) {
	ossObject := &AliCloudOSS{metadata: make(map[string]*OssMetadata), client: make(map[string]*oss.Client)}

	bucket, err := ossObject.selectClientAndBucket(nil)
	assert.Equal(t, err.Error(), "should specific endpoint in metadata")
	assert.Nil(t, bucket)

	client1 := &oss.Client{}
	ossObject.client["127.0.0.1"] = client1
	bucket, err = ossObject.selectClientAndBucket(nil)
	assert.Equal(t, err.Error(), "no bucket configuration")
	assert.Nil(t, bucket)

	metaData1 := &OssMetadata{Bucket: []string{"test", "test2"}}
	ossObject.metadata["127.0.0.1"] = metaData1
	bucket, err = ossObject.selectClientAndBucket(nil)
	assert.Equal(t, err.Error(), "should specific bucketKey in metadata")
	assert.Nil(t, bucket)

	metaData2 := &OssMetadata{Bucket: []string{"test"}}
	ossObject.metadata["127.0.0.1"] = metaData2
	bucket, err = ossObject.selectClientAndBucket(nil)
	assert.Nil(t, err)
	assert.NotNil(t, bucket)
}
