package tencentcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5"

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
	oss := NewTencentCloudOSS()
	err := oss.Init(context.TODO(), &fc)
	assert.Equal(t, err.Error(), "invalid argument")
	fc.Metadata = []byte(data)
	err = oss.Init(context.TODO(), &fc)
	assert.Nil(t, err)
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
