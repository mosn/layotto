package oss

import (
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

func TestInit(t *testing.T) {
	fc := file.FileConfig{}
	metaData := make(map[string]interface{})
	metaData[endpointKey] = " "
	metaData[accessKeyIDKey] = " "
	metaData[accessKeySecretKey] = " "
	metaData[bucketKey] = make([]interface{}, 0)
	metaData[storageTypeKey] = ""
	oss := NewAliCloudOSS()
	err := oss.Init(&fc)
	assert.Equal(t, err.Error(), "no configuration for aliCloudOSS")
	fc.Metadata = append(fc.Metadata, metaData)
	err = oss.Init(&fc)
	assert.Equal(t, err.Error(), "wrong configurations for aliCloudOSS")
}

func TestSelectClient(t *testing.T) {
	ossObject := &AliCloudOSS{metadata: make(map[string]*ossMetadata), client: make(map[string]*oss.Client)}

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

func TestSelectBucket(t *testing.T) {
	ossObject := &AliCloudOSS{metadata: make(map[string]*ossMetadata), client: make(map[string]*oss.Client)}

	bucketName, err := ossObject.selectBucket()
	assert.Equal(t, "", bucketName)
	assert.Equal(t, err.Error(), "no bucket configuration")

	metaData1 := &ossMetadata{Bucket: []string{"test", "test2"}}
	ossObject.metadata["0.0.0.0"] = metaData1
	bucketName, err = ossObject.selectBucket()
	assert.Equal(t, bucketName, "")
	assert.Equal(t, err.Error(), "should specific bucketKey in metadata")

	metaData2 := &ossMetadata{Bucket: []string{"test"}}
	ossObject.metadata["0.0.0.0"] = metaData2
	bucketName, err = ossObject.selectBucket()
	assert.Equal(t, bucketName, "test")
	assert.Nil(t, err)
}
