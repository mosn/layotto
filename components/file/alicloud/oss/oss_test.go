package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/file"
)

func TestInit(t *testing.T) {
	fc := file.FileConfig{Metadata: make(map[string]string)}
	fc.Metadata[endpointKey] = " "
	fc.Metadata[accessKeyIDKey] = " "
	fc.Metadata[accessKeySecretKey] = " "
	fc.Metadata[bucketKey] = " "
	fc.Metadata[storageTypeKey] = ""
	oss := NewAliCloudOSS()
	err := oss.Init(&fc)
	assert.NotNil(t, err)
}
