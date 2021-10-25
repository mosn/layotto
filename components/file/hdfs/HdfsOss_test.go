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
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/file"

	"go.beyondstorage.io/v5/types"
)

// config is the raw json data of component's Metadata configuration
const config = `[
				{
					"endpoint": "tcp:127.0.0.1:9000"
				}
			]`

func TestHdfsOss_Init(t *testing.T) {
	oss := NewHdfs()

	c := &file.FileConfig{}
	err := oss.Init(c)
	assert.Equal(t, err, ErrInvalidConfig)

	c.Metadata = json.RawMessage(config)

	err = oss.Init(c)

	assert.Nil(t, err)
}

func TestHdfsOss_selectClient(t *testing.T) {
	HdfsOss := &hdfs{
		client: make(map[string]types.Storager),
		meta:   make(map[string]*HdfsMetaData),
	}

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}
	err := HdfsOss.Init(c)
	assert.Nil(t, err)

	meta := make(map[string]string)
	meta["endpoint"] = "tcp:127.0.0.1:9000"
	_, err = HdfsOss.selectClient(meta)
	assert.Nil(t, err)

	meta["endpoint"] = "tcp:127.0.0.1:9001"
	_, err = HdfsOss.selectClient(meta)
	assert.NotNil(t, err)

	meta["endpoint"] = "tcp:899.45.7.2:0000"
	_, err = HdfsOss.selectClient(meta)
	assert.NotNil(t, err)

	meta["endpoint"] = "tcp:333.12.1.5:1222"
	_, err = HdfsOss.selectClient(meta)
	assert.NotNil(t, err)
}

func TestHdfsOss_Put(t *testing.T) {
	oss := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}
	err := oss.Init(c)
	assert.Nil(t, err)

	f, _ := os.Open("oss.go")

	req := &file.PutFileStu{
		DataStream: f,
		FileName:   "test_put",
		Metadata:   map[string]string{"": ""},
	}
	// missing endpoint
	err = oss.Put(req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = "tcp:127.0.0.1:9001"
	err = oss.Put(req)
	assert.Equal(t, ErrClientNotExist, err)

	// convert from string to int64 failed
	req.Metadata["endpoint"] = "endpoint"
	req.Metadata["fileSize"] = "hdfs"
	err = oss.Put(req)
	assert.NotNil(t, err)

	// convert from string to int64 success
	req.Metadata["endpoint"] = "tcp:127.0.0.1:9000"
	req.Metadata["fileSize"] = "123"
	err = oss.Put(req)
	assert.Nil(t, err)

}

func TestMinioOss_Get(t *testing.T) {
	oss := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := oss.Init(c)
	assert.Nil(t, err)

	req := &file.GetFileStu{
		FileName: "test_put",
		Metadata: map[string]string{"": ""},
	}

	_, err = oss.Get(req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = "127.0.0.1:9000"
	_, err = oss.Get(req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = "tcp:127.0.0.1:9000"
	_, err = oss.Get(req)
	assert.Nil(t, err)

	//TODO
	//Test checksum content with Get file
}

func TestHdfsOss_List(t *testing.T) {
	oss := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := oss.Init(c)
	assert.Nil(t, err)

	req := &file.ListRequest{
		DirectoryName: "test_put",
		Metadata:      map[string]string{"": ""},
	}

	_, err = oss.List(req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = "127.0.0.1:9000"
	_, err = oss.List(req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = "tcp:127.0.0.1:9000"
	_, err = oss.List(req)
	assert.Nil(t, err)

	//TODO
	//Test "" Directory and exist files Directory

}

func TestHdfsOss_Del(t *testing.T) {
	oss := NewHdfs()

	c := &file.FileConfig{
		Metadata: json.RawMessage(config),
	}

	err := oss.Init(c)
	assert.Nil(t, err)

	req := &file.DelRequest{
		FileName: "test_put",
		Metadata: map[string]string{"": ""},
	}
	err = oss.Del(req)
	assert.Equal(t, ErrMissingEndPoint, err)

	// client not exist
	req.Metadata["endpoint"] = "127.0.0.1:9000"
	err = oss.Del(req)
	assert.Equal(t, ErrClientNotExist, err)

	req.Metadata["endpoint"] = "tcp:127.0.0.1:9000"
	err = oss.Del(req)
	assert.Nil(t, err)
}
