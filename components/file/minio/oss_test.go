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

package minio

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

// cfg is the raw json data of component's Metadata configuration
const cfg = `[
				{
					"endpoint": "endpoint",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"region": "us-west-2",
					"SSL":true
				}
			]`

func TestMinioOss_Init(t *testing.T) {
	oss := NewMinioOss()

	initCfg := &file.FileConfig{}
	err := oss.Init(context.TODO(), initCfg)
	assert.Equal(t, err, ErrInvalidConfig)

	initCfg.Metadata = json.RawMessage(cfg)

	err = oss.Init(context.TODO(), initCfg)
	assert.Nil(t, err)
}

func TestMinioOss_selectClient(t *testing.T) {
	minioOss := &MinioOss{
		client: make(map[string]*minio.Core),
		meta:   make(map[string]*MinioMetaData),
	}
	initCfg := &file.FileConfig{
		Metadata: json.RawMessage(cfg),
	}
	err := minioOss.Init(context.TODO(), initCfg)
	assert.Nil(t, err)

	meta := make(map[string]string)
	_, err = minioOss.selectClient(meta)
	assert.Nil(t, err)

	minioOss.client["extra"] = &minio.Core{}
	_, err = minioOss.selectClient(meta)
	assert.Equal(t, ErrNotSpecifyEndPoint, err)

	delete(minioOss.client, "extra")

	meta["endpoint"] = "endpoint1"
	_, err = minioOss.selectClient(meta)
	assert.Equal(t, ErrClientNotExist, err)

	meta["endpoint"] = "endpoint"
	_, err = minioOss.selectClient(meta)
	assert.Nil(t, err)
}

func TestMinioOss_Put(t *testing.T) {
	oss := NewMinioOss()

	initCfg := &file.FileConfig{}
	initCfg.Metadata = json.RawMessage(cfg)
	err := oss.Init(context.TODO(), initCfg)
	assert.Nil(t, err)

	f, _ := os.Open("oss.go")

	putReq := &file.PutFileStu{
		DataStream: f,
		FileName:   "bucket/file",
		Metadata:   map[string]string{"": ""},
	}
	// client not exist
	putReq.Metadata["endpoint"] = "demo-endpoint"
	err = oss.Put(context.TODO(), putReq)
	assert.Equal(t, ErrClientNotExist, err)

	// convert from string to int64 failed
	putReq.Metadata["endpoint"] = "endpoint"
	putReq.Metadata["fileSize"] = "a2"
	err = oss.Put(context.TODO(), putReq)
	assert.NotNil(t, err)

	err = oss.Put(context.TODO(), putReq)
	assert.NotNil(t, err)
}

func TestMinioOss_Get(t *testing.T) {
	oss := NewMinioOss()

	initCfg := &file.FileConfig{}
	initCfg.Metadata = json.RawMessage(cfg)
	err := oss.Init(context.TODO(), initCfg)
	assert.Nil(t, err)

	getReq := &file.GetFileStu{
		FileName: "file",
		Metadata: map[string]string{"": ""},
	}

	_, err = oss.Get(context.TODO(), getReq)
	assert.Equal(t, "minio get file[file] fail,err: invalid fileName format", err.Error())

	// client not exist
	getReq.FileName = "bucketName/file"
	getReq.Metadata["endpoint"] = "demo-endpoint"
	_, err = oss.Get(context.TODO(), getReq)
	assert.Equal(t, ErrClientNotExist, err)

	getReq.Metadata["endpoint"] = "endpoint"
	_, err = oss.Get(context.TODO(), getReq)
	assert.Nil(t, err)
}
