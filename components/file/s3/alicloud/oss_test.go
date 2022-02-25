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

package alicloud

import (
	"context"
	"io"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

const data = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`

func TestInit(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	err := oss.Init(context.TODO(), &fc)
	assert.Equal(t, err.Error(), "invalid argument")
	fc.Metadata = []byte(data)
	err = oss.Init(context.TODO(), &fc)
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

func TestGetBucket(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	ac := oss.(*AliCloudOSS)
	mt := make(map[string]string)

	bucket, err := ac.getBucket("/", mt)
	assert.Equal(t, err.Error(), "invalid fileName format")
	assert.Nil(t, bucket)

	bucket, err = ac.getBucket("b/a.txt", mt)
	assert.Equal(t, err.Error(), "bucket name b len is between [3-63],now is 1")
	assert.Nil(t, bucket)

	bucket, err = ac.getBucket("bbbb/a.txt", mt)
	assert.NoError(t, err)
	assert.Equal(t, bucket.BucketName, "bbbb")
}

func TestGetClient(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	ac := oss.(*AliCloudOSS)
	mt := &OssMetadata{
		Endpoint:        "endpoint",
		AccessKeyID:     "ak",
		AccessKeySecret: "ak",
	}

	//TODO test empty endpoint/ak/sk , now will get panic

	client, err := ac.getClient(mt)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestCheckMetadata(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	ac := oss.(*AliCloudOSS)
	mt := &OssMetadata{
		Endpoint:        "",
		AccessKeyID:     "",
		AccessKeySecret: "",
	}

	assert.False(t, ac.checkMetadata(mt))
	mt.Endpoint = "e"
	assert.False(t, ac.checkMetadata(mt))
	mt.AccessKeySecret = "sk"
	assert.False(t, ac.checkMetadata(mt))
	mt.AccessKeyID = "ak"
	assert.True(t, ac.checkMetadata(mt))
}

func TestPut(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.PutFileStu{
		FileName: "",
	}
	err = oss.Put(context.Background(), req)
	assert.Error(t, err)

	req.FileName = "b/a.txt"
	err = oss.Put(context.Background(), req)
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.GetFileStu{
		FileName: "",
	}

	var resp io.ReadCloser
	resp, err = oss.Get(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	req.FileName = "b/a.txt"
	resp, err = oss.Get(context.Background(), req)
	assert.Error(t, err)
}

func TestStat(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.FileMetaRequest{
		FileName: "",
	}

	var resp *file.FileMetaResp
	resp, err = oss.Stat(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	req.FileName = "b/a.txt"
	resp, err = oss.Stat(context.Background(), req)
	assert.Error(t, err)
}

func TestList(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.ListRequest{
		DirectoryName: "",
		PageSize:      0,
	}

	var resp *file.ListResp
	resp, err = oss.List(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	req.DirectoryName = "b/"
	resp, err = oss.List(context.Background(), req)
	assert.Error(t, err)
}

func TestDel(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.DelRequest{
		FileName: "",
	}

	err = oss.Del(context.Background(), req)
	assert.Error(t, err)

	req.FileName = "b/a.txt"
	err = oss.Del(context.Background(), req)
	assert.Error(t, err)
}
