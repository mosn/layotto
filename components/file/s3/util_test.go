package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBucketName(t *testing.T) {
	_, err := GetBucketName("/")
	assert.Equal(t, "invalid fileName format", err.Error())
	_, err = GetBucketName("")
	assert.Equal(t, "invalid fileName format", err.Error())
	_, err = GetBucketName("bucketName")
	assert.Equal(t, "invalid fileName format", err.Error())
	name, err := GetBucketName("bucketName/")
	assert.Nil(t, err)
	assert.Equal(t, name, "bucketName")
}

func TestGetFileName(t *testing.T) {
	name, err := GetFileName("/")
	assert.Equal(t, err.Error(), "file name is empty")
	assert.Equal(t, "", name)
	name, err = GetFileName("aaa")
	assert.Equal(t, err.Error(), "invalid fileName format")
	assert.Equal(t, "", name)
	name, err = GetFileName("/aaa")
	assert.Nil(t, err)
	assert.Equal(t, "aaa", name)
}
