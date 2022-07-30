//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"encoding/json"
	"testing"

	"mosn.io/layotto/components/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	data := `{	"hellos": {
					"helloworld": {
					"hello": "greeting"
					}
				},
				"file": {
					"aliyun.oss": {
                          "metadata":[
                            {
                              "endpoint": "endpoint_address",
                              "accessKeyID": "accessKey",
                              "accessKeySecret": "secret"
                            }
                          ]
					}
				}
			}`
	mscf, err := ParseRuntimeConfig([]byte(data))
	assert.Nil(t, err)
	v := mscf.Files["aliyun.oss"]
	m := make([]*utils.OssMetadata, 0)
	err = json.Unmarshal(v.Metadata, &m)
	assert.Nil(t, err)
	for _, x := range m {
		assert.Equal(t, "endpoint_address", x.Endpoint)
		assert.Equal(t, "accessKey", x.AccessKeyID)
		assert.Equal(t, "secret", x.AccessKeySecret)
	}
}
