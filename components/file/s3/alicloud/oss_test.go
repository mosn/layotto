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
					"accessKeySecret": "secret"
				}
			]`
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
