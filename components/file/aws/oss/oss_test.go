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

package oss

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/file"
)

const cfg = `[
				{
					"endpoint": "protocol://service-code.region-code.amazonaws.com",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret",
					"region": "us-west-2"
				}
			]`

func TestAwsOss_Init(t *testing.T) {
	oss := NewAwsOss()
	err := oss.Init(&file.FileConfig{})
	assert.Equal(t, err.Error(), "invalid config for aws oss")
	err = oss.Init(&file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)
}

func TestAwsOss_SelectClient(t *testing.T) {
	oss := &AwsOss{
		client: make(map[string]*s3.Client),
		meta:   make(map[string]*AwsOssMetaData),
	}
	err := oss.Init(&file.FileConfig{Metadata: []byte(cfg)})
	assert.Equal(t, nil, err)

	// not specify endpoint, select default client
	meta := map[string]string{}
	_, err = oss.selectClient(meta)
	assert.Nil(t, err)

	// specify endpoint equal config
	meta["endpoint"] = "protocol://service-code.region-code.amazonaws.com"
	client, _ := oss.selectClient(meta)
	assert.NotNil(t, client)

	// specicy not exist endpoint, select default one
	meta["endpoint"] = "protocol://cn-northwest-1.region-code.amazonaws.com"
	client, err = oss.selectClient(meta)
	assert.Nil(t, err)

	// new client with endpoint
	oss.client["protocol://cn-northwest-1.region-code.amazonaws.com"] = &s3.Client{}
	client, _ = oss.selectClient(meta)
	assert.NotNil(t, client)
}
