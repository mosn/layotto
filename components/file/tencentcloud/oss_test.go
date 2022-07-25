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

package tencentcloud

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5"

	"mosn.io/layotto/components/file"
)

const data = `[
				{
					"endpoint": "https://xxx-1251058690.cos.ap-chengdu.myqcloud.com",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"timeout": 1000
				}
			]`

func TestInit(t *testing.T) {
	data := `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	err := oss.Init(context.TODO(), &fc)
	assert.Equal(t, err.Error(), "invalid argument")
	fc.Metadata = []byte(data)
	err = oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	tcOSS := oss.(*TencentCloudOSS)
	mt := tcOSS.metadata["endpoint_address1"]
	assert.Nil(t, mt)

	mt = tcOSS.metadata["endpoint_address"]
	assert.Equal(t, "endpoint_address", mt.Endpoint)
	assert.Equal(t, "secret", mt.AccessKeySecret)
	assert.Equal(t, "accessKey", mt.AccessKeyID)

	client, ce := tcOSS.getClient(mt)
	assert.NotNil(t, client)
	assert.NoError(t, ce)
}

func TestInitTimeout(t *testing.T) {
	data := `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"timeout": 3000
				}
			]`
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	err := oss.Init(context.TODO(), &fc)
	assert.Equal(t, err.Error(), "invalid argument")
	fc.Metadata = []byte(data)
	err = oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	tcOSS := oss.(*TencentCloudOSS)
	mt := tcOSS.metadata["endpoint_address"]
	assert.NotNil(t, mt)
	assert.Equal(t, 3000, mt.Timeout)
}

func TestInitTimeoutReset(t *testing.T) {
	data := `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"timeout": -1
				}
			]`
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	err := oss.Init(context.TODO(), &fc)
	assert.Equal(t, err.Error(), "invalid argument")
	fc.Metadata = []byte(data)
	err = oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	tcOSS := oss.(*TencentCloudOSS)
	mt := tcOSS.metadata["endpoint_address"]
	assert.NotNil(t, mt)
	assert.Equal(t, 100*1000, mt.Timeout)
}

func TestCheckMetadata(t *testing.T) {
	oss := NewTencentCloudOSS()
	tcOSS := oss.(*TencentCloudOSS)

	mt := &OssMetadata{}
	assert.False(t, tcOSS.checkMetadata(mt))

	mt.Endpoint = "1"
	assert.False(t, tcOSS.checkMetadata(mt))

	mt.AccessKeySecret = "2"
	assert.False(t, tcOSS.checkMetadata(mt))

	mt.AccessKeyID = "3"
	assert.True(t, tcOSS.checkMetadata(mt))

	mt.AccessKeySecret = ""
	assert.False(t, tcOSS.checkMetadata(mt))

	mt.AccessKeySecret = "x"
	mt.Endpoint = "example.com"
	tcOSS.checkMetadata(mt)
	assert.Equal(t, "https", mt.bucketUrl.Scheme)

	mt.Endpoint = "http://example.com"
	tcOSS.checkMetadata(mt)
	assert.Equal(t, "http", mt.bucketUrl.Scheme)
}

func TestSelectClient(t *testing.T) {
	ossObject := &TencentCloudOSS{metadata: make(map[string]*OssMetadata), client: make(map[string]*cos.Client)}

	meta := make(map[string]string)
	client, err := ossObject.selectClient(meta)
	assert.Equal(t, err.Error(), "not specify endpoint in metadata")
	assert.Nil(t, client)

	client1 := &cos.Client{}
	ossObject.client["127.0.0.1"] = client1
	client, err = ossObject.selectClient(meta)
	assert.Equal(t, client, client1)
	assert.Nil(t, err)

	client2 := &cos.Client{}
	ossObject.client["0.0.0.0"] = client2
	client, err = ossObject.selectClient(meta)
	assert.Equal(t, err.Error(), "not specify endpoint in metadata")
	assert.Nil(t, client)

	meta[endpointKey] = "0.0.0.0"
	client, err = ossObject.selectClient(meta)
	assert.Equal(t, err, nil)
	assert.NotNil(t, client)
}

func TestPut(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	st := &file.PutFileStu{
		FileName:   "b/a.txt",
		DataStream: strings.NewReader("aaaa"),
	}

	err = oss.Put(context.Background(), st)
	assert.Error(t, err)

	st.FileName = "/b/a.txt"
	err = oss.Put(context.Background(), st)
	assert.Error(t, err)
}

func TestStat(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	st := &file.FileMetaRequest{
		FileName: "b/a.txt",
	}

	var resp *file.FileMetaResp
	resp, err = oss.Stat(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

}

func TestList(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	st := &file.ListRequest{
		DirectoryName: "/a",
		PageSize:      10,
	}

	var resp *file.ListResp
	resp, err = oss.List(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)

	st.PageSize = 1001
	resp, err = oss.List(context.Background(), st)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "PageSize must be <=1000")
	assert.Nil(t, resp)

	st.PageSize = 0
	resp, err = oss.List(context.Background(), st)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "PageSize must be >0")
	assert.Nil(t, resp)
}

func TestGet(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	st := &file.GetFileStu{
		FileName: "a.txt",
	}

	var resp io.ReadCloser
	resp, err = oss.Get(context.Background(), st)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDel(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	st := &file.DelRequest{
		FileName: "a.txt",
	}

	err = oss.Del(context.Background(), st)
	assert.Error(t, err)
}

func TestCheckFileName(t *testing.T) {
	fc := file.FileConfig{}
	oss := NewTencentCloudOSS()
	fc.Metadata = []byte(data)
	err := oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)

	tcOSS := oss.(*TencentCloudOSS)
	assert.Error(t, tcOSS.checkFileName("/"))
	assert.Error(t, tcOSS.checkFileName("/a"))
	assert.Nil(t, tcOSS.checkFileName("a"))
	assert.Nil(t, tcOSS.checkFileName("a/aa.txt"))
}
