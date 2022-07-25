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

package qiniu

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"

	"github.com/golang/mock/gomock"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/pkg/mock"
)

const (
	endpointAddress = "endpoint_address"
	metadata        = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": "xc2022"
				}
			]`

	metadata2 = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": "",
					"useHTTPS": true
				}
			]`
)

func TestNew(t *testing.T) {
	s := NewQiniuOSS()
	assert.NotNil(t, s)
}

func TestInit(t *testing.T) {
	oss := NewQiniuOSS()
	fc := file.FileConfig{}
	err := oss.Init(context.Background(), &fc)
	assert.Error(t, err)

	fc.Metadata = []byte(metadata)
	err = oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	fc.Metadata = []byte(metadata2)
	err = oss.Init(context.Background(), &fc)
	assert.Error(t, err)

	fc.Metadata = []byte(metadata + ",")
	err = oss.Init(context.Background(), &fc)
	assert.Error(t, err)

	data3 := `[
				{
					"endpoint": "",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": "",
					"useHTTPS": true
				}
			]`
	fc.Metadata = []byte(data3)
	err = oss.Init(context.Background(), &fc)
	assert.Error(t, err)

	data4 := `[
				{
					"endpoint": "xxxxx",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": "cc",
					"useHTTPS": true
				}
			]`
	fc.Metadata = []byte(data4)
	err = oss.Init(context.Background(), &fc)
	assert.Nil(t, err)
}

func TestCheckMetadata(t *testing.T) {
	m := &OssMetadata{}

	assert.False(t, m.checkMetadata())

	m.Bucket = "1"
	assert.False(t, m.checkMetadata())

	m.AccessKeyID = "1"
	assert.False(t, m.checkMetadata())

	m.AccessKeySecret = "1"
	assert.False(t, m.checkMetadata())

	m.Endpoint = "1"
	assert.True(t, m.checkMetadata())

}

func TestSelectClient(t *testing.T) {
	oss := NewQiniuOSS().(*QiniuOSS)
	fc := file.FileConfig{}
	fc.Metadata = []byte(metadata)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	mt := make(map[string]string)
	var client *QiniuOSSClient
	client, err = oss.selectClient(mt)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	mt[endpointKey] = "1"
	client, err = oss.selectClient(mt)
	assert.Error(t, err)
	assert.Nil(t, client)

	mt[endpointKey] = endpointAddress
	client, err = oss.selectClient(mt)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestSelectClientWithMulti(t *testing.T) {
	oss := NewQiniuOSS().(*QiniuOSS)
	data := `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": "xc2022"
				},
				{
					"endpoint": "endpoint_address2",
					"accessKeyID": "accessKey2",
					"accessKeySecret": "secret2",
					"bucket": "xc20222"
				}
			]`
	fc := file.FileConfig{}
	fc.Metadata = []byte(data)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	mt := make(map[string]string)
	var client *QiniuOSSClient
	client, err = oss.selectClient(mt)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestPut(t *testing.T) {
	oss := NewQiniuOSS()
	fc := file.FileConfig{}
	fc.Metadata = []byte(metadata)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	st := &file.PutFileStu{
		FileName:   "a.txt",
		DataStream: strings.NewReader("aaa"),
		Metadata:   make(map[string]string),
	}

	st.Metadata[endpointKey] = "1"
	err = oss.Put(context.Background(), st)
	assert.Error(t, err)

	st.Metadata[endpointKey] = endpointAddress
	err = oss.Put(context.Background(), st)
	assert.Error(t, err)

	st.Metadata[fileSizeKey] = endpointAddress
	err = oss.Put(context.Background(), st)
	assert.Error(t, err)

	st.Metadata[fileSizeKey] = "10"
	err = oss.Put(context.Background(), st)
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `=~^http(s)?://example.com`,
		httpmock.NewStringResponder(200, ""))

	oss := NewQiniuOSS()
	data := `[
				{
					"endpoint": "example.com",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"bucket": "xc2022"
				}
			]`
	fc := file.FileConfig{}
	fc.Metadata = []byte(data)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	st := &file.GetFileStu{
		FileName: "a.txt",
		Metadata: make(map[string]string),
	}

	var resp io.ReadCloser
	st.Metadata[endpointKey] = "1"
	resp, err = oss.Get(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

	st.Metadata[endpointKey] = "example.com"
	resp, err = oss.Get(context.Background(), st)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDel(t *testing.T) {
	oss := NewQiniuOSS()
	fc := file.FileConfig{}
	fc.Metadata = []byte(metadata)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	st := &file.DelRequest{
		FileName: "a.txt",
		Metadata: make(map[string]string),
	}

	st.Metadata[endpointKey] = "1"
	err = oss.Del(context.Background(), st)
	assert.Error(t, err)

	st.Metadata[endpointKey] = endpointAddress
	err = oss.Del(context.Background(), st)
	assert.Error(t, err)
}

func TestStat(t *testing.T) {
	oss := NewQiniuOSS().(*QiniuOSS)
	fc := file.FileConfig{}
	fc.Metadata = []byte(metadata)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	st := &file.FileMetaRequest{
		FileName: "a.txt",
		Metadata: make(map[string]string),
	}

	var resp *file.FileMetaResp
	st.Metadata[endpointKey] = "1"
	resp, err = oss.Stat(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

	st.Metadata[endpointKey] = endpointAddress
	resp, err = oss.Stat(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fu := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()
	bm.EXPECT().Stat(gomock.Eq("xc2022"), gomock.Eq("a.txt")).Return(storage.FileInfo{}, nil)
	mockOss(oss, bm, fu)

	st.Metadata[endpointKey] = endpointAddress
	resp, err = oss.Stat(context.Background(), st)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

}

func TestList(t *testing.T) {
	oss := NewQiniuOSS().(*QiniuOSS)
	fc := file.FileConfig{}
	fc.Metadata = []byte(metadata)
	err := oss.Init(context.Background(), &fc)
	assert.NoError(t, err)

	st := &file.ListRequest{
		DirectoryName: "b/",
		Metadata:      make(map[string]string),
		PageSize:      1,
	}

	var resp *file.ListResp
	st.Metadata[endpointKey] = "1"
	resp, err = oss.List(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

	st.Metadata[endpointKey] = endpointAddress
	resp, err = oss.List(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fu := mock.NewMockFormUploader(ctrl)
	bm := mock.NewMockBucketManager(ctrl)
	defer ctrl.Finish()

	items := make([]storage.ListItem, 1)
	items[0] = storage.ListItem{
		Key:     "a.txt",
		Fsize:   3,
		PutTime: time.Now().UnixNano() / 1e9,
	}

	bm.EXPECT().ListFiles(gomock.Eq("xc2022"), gomock.Eq("b/"), gomock.Any(), gomock.Any(), gomock.Any()).Return(items, make([]string, 0), "", false, nil)
	mockOss(oss, bm, fu)

	st.Metadata[endpointKey] = endpointAddress
	resp, err = oss.List(context.Background(), st)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

}

func mockOss(oss *QiniuOSS, bm BucketManager, fu FormUploader) {
	for _, v := range oss.client {
		v.bm = bm
		v.fu = fu
	}
}
