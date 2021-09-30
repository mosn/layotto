package oss

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/file"
	"testing"
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

	meta := map[string]string{}
	_, err = oss.selectClient(meta)
	assert.Equal(t, err.Error(), "specific client not exist")

	meta["endpoint"] = "protocol://service-code.region-code.amazonaws.com"
	client, err := oss.selectClient(meta)
	assert.NotNil(t, client)

	meta["endpoint"] = "protocol://cn-northwest-1.region-code.amazonaws.com"
	client, err = oss.selectClient(meta)
	assert.Equal(t, err.Error(), "specific client not exist")

	oss.client["protocol://cn-northwest-1.region-code.amazonaws.com"] = &s3.Client{}
	client, err = oss.selectClient(meta)
	assert.NotNil(t, client)
}
