package runtime

import (
	"encoding/json"
	"testing"

	"mosn.io/layotto/components/file/alicloud/oss"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	data := `{	"hellos": {
					"helloworld": {
					"hello": "greeting"
					}
				},
				"files": {
					"aliOSS": {
                          "metadata":[
                            {
                              "endpoint": "endpoint_address",
                              "accessKeyID": "accessKey",
                              "accessKeySecret": "secret",
                              "bucket": ["bucket1", "bucket2"]
                            }
                          ]
					}
				}
			}`
	mscf, err := ParseRuntimeConfig([]byte(data))
	assert.Nil(t, err)
	v := mscf.Files["aliOSS"]
	m := make([]*oss.OssMetadata, 0, 0)
	err = json.Unmarshal(v.Metadata, &m)
	assert.Nil(t, err)
	for _, x := range m {
		assert.Equal(t, "endpoint_address", x.Endpoint)
		assert.Equal(t, "accessKey", x.AccessKeyID)
		assert.Equal(t, "secret", x.AccessKeySecret)
		assert.Equal(t, []string{"bucket1", "bucket2"}, x.Bucket)
	}
}
