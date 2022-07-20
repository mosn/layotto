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
	"io"
	"testing"

	"mosn.io/layotto/components/pkg/utils"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

const (
	data = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
	fileName = "b/a.txt"
)

func TestInit(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
	err := oss.Init(context.TODO(), &fc)
	assert.Equal(t, err.Error(), "invalid argument")
	fc.Metadata = []byte(data)
	err = oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)
}

func TestGetBucket(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	ac := oss.(*AliyunFile)
	mt := make(map[string]string)

	bucket, err := ac.getBucket("/", mt)
	assert.Equal(t, err.Error(), "invalid fileName format")
	assert.Nil(t, bucket)

	bucket, err = ac.getBucket(fileName, mt)
	assert.Equal(t, err.Error(), "bucket name b len is between [3-63],now is 1")
	assert.Nil(t, bucket)

	bucket, err = ac.getBucket("bbbb/a.txt", mt)
	assert.NoError(t, err)
	assert.Equal(t, bucket.BucketName, "bbbb")
}

func TestGetClient(t *testing.T) {
	fc := file.FileConfig{}
	oss := &AliyunFile{}
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	assert.NotNil(t, oss.client)
}

func TestCheckMetadata(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	ac := oss.(*AliyunFile)
	mt := &utils.OssMetadata{
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
	oss := NewAliyunFile()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.PutFileStu{
		FileName: "",
	}
	err = oss.Put(context.Background(), req)
	assert.Error(t, err)

	req.FileName = fileName
	err = oss.Put(context.Background(), req)
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
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

	req.FileName = fileName
	_, err = oss.Get(context.Background(), req)
	assert.Error(t, err)
}

func TestStat(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
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

	req.FileName = fileName
	_, err = oss.Stat(context.Background(), req)
	assert.Error(t, err)
}

func TestList(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
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
	_, err = oss.List(context.Background(), req)
	assert.Error(t, err)
}

func TestDel(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewAliyunFile()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	req := &file.DelRequest{
		FileName: "",
	}

	err = oss.Del(context.Background(), req)
	assert.Error(t, err)

	req.FileName = fileName
	err = oss.Del(context.Background(), req)
	assert.Error(t, err)
}
