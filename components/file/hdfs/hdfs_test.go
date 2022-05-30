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

package hdfs

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"

	"go.beyondstorage.io/v5/types"
)

const (
	// config is the raw json data of component's Metadata configuration
	config = `[
				{
					"endpoint": "tcp:127.0.0.1:9000"
				}
			]`
	endpoint    = "127.0.0.1:9000"
	tcpEndpoint = "tcp:127.0.0.1:9000"
)

func TestHdfs_Init(t *testing.T) {
	hdfs := NewHdfs()

	c := &file.FileConfig{}
	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInvalidConfig)

	c.Metadata = json.RawMessage(config)

	err = hdfs.Init(context.TODO(), c)

	assert.Equal(t, err, ErrInitFailed)
}

func TestHdfs_selectClient(t *testing.T) {
	hdfs := &hdfs{
		client: make(map[string]types.Storager),
		meta:   make(map[string]*HdfsMetaData),
	}

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}
	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	meta := make(map[string]string)
	meta["endpoint"] = tcpEndpoint
	_, err = hdfs.selectClient(meta)
	assert.NotNil(t, err)

	meta["endpoint"] = "tcp:127.0.0.1:9001"
	_, err = hdfs.selectClient(meta)
	assert.NotNil(t, err)

	meta["endpoint"] = "tcp:899.45.7.2:0000"
	_, err = hdfs.selectClient(meta)
	assert.NotNil(t, err)

	meta["endpoint"] = "tcp:333.12.1.5:1222"
	_, err = hdfs.selectClient(meta)
	assert.NotNil(t, err)
}

func TestHdfs_Put(t *testing.T) {
	hdfs := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}
	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	f, _ := os.Open("hdfs.go")

	req := &file.PutFileStu{
		DataStream: f,
		FileName:   "test_put",
		Metadata:   map[string]string{"": ""},
	}
	// missing endpoint
	err = hdfs.Put(context.TODO(), req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = "tcp:127.0.0.1:9001"
	err = hdfs.Put(context.TODO(), req)
	assert.Equal(t, ErrClientNotExist, err)

	// convert from string to int64 failed
	req.Metadata["endpoint"] = "endpoint"
	req.Metadata["fileSize"] = "hdfs"
	err = hdfs.Put(context.TODO(), req)
	assert.NotNil(t, err)

	// convert from string to int64 success
	req.Metadata["endpoint"] = tcpEndpoint
	req.Metadata["fileSize"] = "123"
	err = hdfs.Put(context.TODO(), req)
	assert.NotNil(t, err)

}

func TestHdfs_Get(t *testing.T) {
	hdfs := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	req := &file.GetFileStu{
		FileName: "test_put",
		Metadata: map[string]string{"": ""},
	}

	_, err = hdfs.Get(context.TODO(), req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = endpoint
	_, err = hdfs.Get(context.TODO(), req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = tcpEndpoint
	_, err = hdfs.Get(context.TODO(), req)
	assert.NotNil(t, err)

	//TODO
	//Test checksum content with Get file
}

func TestHdfs_Del(t *testing.T) {
	hdfs := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	req := &file.DelRequest{
		FileName: "test_put",
		Metadata: map[string]string{"": ""},
	}
	err = hdfs.Del(context.TODO(), req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = endpoint
	err = hdfs.Del(context.TODO(), req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = tcpEndpoint
	err = hdfs.Del(context.TODO(), req)
	assert.NotNil(t, err)
}

func TestHdfs_List(t *testing.T) {
	hdfs := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	req := &file.ListRequest{
		DirectoryName: "/",
		PageSize:      1,
		Metadata:      map[string]string{"": ""},
	}

	var resp *file.ListResp
	resp, err = hdfs.List(context.TODO(), req)
	assert.Equal(t, ErrMissingEndPoint, err)
	assert.Nil(t, resp)

	req.Metadata["endpoint"] = endpoint
	_, err = hdfs.List(context.TODO(), req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = tcpEndpoint
	_, err = hdfs.List(context.TODO(), req)
	assert.NotNil(t, err)
}

func TestHdfs_Stat(t *testing.T) {
	hdfs := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := hdfs.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	req := &file.FileMetaRequest{
		FileName: "a.txt",
		Metadata: map[string]string{"": ""},
	}

	var resp *file.FileMetaResp
	resp, err = hdfs.Stat(context.TODO(), req)
	assert.Equal(t, ErrNotSpecifyEndpoint, err)
	assert.Nil(t, resp)

	req.Metadata["endpoint"] = endpoint
	_, err = hdfs.Stat(context.TODO(), req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = tcpEndpoint
	_, err = hdfs.Stat(context.TODO(), req)
	assert.NotNil(t, err)
}

func TestHdfs_CreateHdfsClient(t *testing.T) {
	oss := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := oss.Init(context.TODO(), c)
	assert.Equal(t, err, ErrInitFailed)

	mt := &HdfsMetaData{
		EndPoint: "a",
	}
	store, err := oss.(*hdfs).createHdfsClient(mt)
	assert.Nil(t, store)
	assert.Error(t, err)

	mt.EndPoint = tcpEndpoint
	store, err = oss.(*hdfs).createHdfsClient(mt)
	assert.Nil(t, store)
	assert.Error(t, err)
}

func TestHdfs_IsHdfsMetaValid(t *testing.T) {
	mt := &HdfsMetaData{
		EndPoint: "",
	}

	assert.False(t, mt.isHdfsMetaValid())

	mt.EndPoint = "a"
	assert.True(t, mt.isHdfsMetaValid())
}
