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
	"testing"

	"mosn.io/layotto/components/file"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file/factory"
)

const (
	conf = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
)

func TestInitAliyunOss(t *testing.T) {
	NewAliyunOss()
	f := factory.GetInitFunc(DefaultClientInitFunc)
	assert.Equal(t, f, AliyunDefaultInitFunc)
	clients, err := f([]byte("hello"), map[string]string{})
	assert.Equal(t, err, file.ErrInvalid)
	assert.Equal(t, nil, clients)
	//oss.InitConfig(context.TODO(), file.FileConfig{})
}
