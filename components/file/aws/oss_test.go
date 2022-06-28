package aws

import (
	"context"
	"errors"
	"testing"

	"mosn.io/pkg/buffer"

	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/file/factory"

	"github.com/stretchr/testify/assert"
)

const (
	confWithoutUidAndBucket = `[
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
	confWithUid = `[
				{	
					"uid": "123",
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
	confWithUidAndBucket = `[
				{	
					"uid": "123",
					"buckets": ["bucket1","bucket2"],
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			]`
)

func TestAwsDefaultInitFunc(t *testing.T) {
	NewAwsOss()
	f := factory.GetInitFunc(DefaultClientInitFunc)
	clients, err := f([]byte("hello"), map[string]string{})
	assert.Equal(t, err, file.ErrInvalid)
	assert.Nil(t, clients)
	clients, err = f([]byte(confWithoutUidAndBucket), map[string]string{})
	assert.NotEqual(t, file.ErrInvalid, err)
	assert.NotNil(t, clients)
	client, ok := clients[""]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	clients, err = f([]byte(confWithUid), map[string]string{})
	assert.NotEqual(t, file.ErrInvalid, err)
	assert.NotNil(t, clients)
	client, ok = clients[""]
	assert.Equal(t, false, ok)
	assert.Nil(t, client)
	client, ok = clients["123"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	clients, err = f([]byte(confWithUidAndBucket), map[string]string{})
	assert.NotEqual(t, file.ErrInvalid, err)
	assert.NotNil(t, clients)
	client, ok = clients[""]
	assert.Equal(t, false, ok)
	assert.Nil(t, client)

	client, ok = clients["123"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	client, ok = clients["bucket1"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)

	client, ok = clients["bucket2"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, client)
}

func TestAwsOss(t *testing.T) {
	instance := NewAwsOss()
	err := instance.InitConfig(context.TODO(), &file.FileConfig{Method: "", Metadata: []byte(confWithoutUidAndBucket)})
	assert.Nil(t, err)
	err = instance.InitClient(context.TODO(), &file.InitRequest{})
	assert.Nil(t, err)

	appendObjectResp, err := instance.AppendObject(context.TODO(), &file.AppendObjectInput{})
	assert.Equal(t, errors.New("AppendObject method not supported on AWS"), err)
	assert.Nil(t, appendObjectResp)

	_, err = instance.AbortMultipartUpload(context.TODO(), &file.AbortMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CompleteMultipartUpload(context.TODO(), &file.CompleteMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CopyObject(context.TODO(), &file.CopyObjectInput{})
	assert.Equal(t, errors.New("must specific copy_source"), err)

	_, err = instance.CopyObject(context.TODO(), &file.CopyObjectInput{
		CopySource: &file.CopySource{CopySourceBucket: "bucket", CopySourceKey: "key"},
	})
	assert.NotEqual(t, errors.New("must specific copy_source"), err)
	_, err = instance.CreateMultipartUpload(context.TODO(), &file.CreateMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.DeleteObject(context.TODO(), &file.DeleteObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjects(context.TODO(), &file.DeleteObjectsInput{
		Delete: &file.Delete{},
	})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjectTagging(context.TODO(), &file.DeleteObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.GetObject(context.TODO(), &file.GetObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectCannedAcl(context.TODO(), &file.GetObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectTagging(context.TODO(), &file.GetObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.HeadObject(context.TODO(), &file.HeadObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.IsObjectExist(context.TODO(), &file.IsObjectExistInput{})
	assert.NotNil(t, err)

	_, err = instance.ListParts(context.TODO(), &file.ListPartsInput{})
	assert.NotNil(t, err)

	_, err = instance.ListMultipartUploads(context.TODO(), &file.ListMultipartUploadsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjects(context.TODO(), &file.ListObjectsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjectVersions(context.TODO(), &file.ListObjectVersionsInput{})
	assert.NotNil(t, err)

	stream := buffer.NewIoBufferString("hello")
	_, err = instance.PutObject(context.TODO(), &file.PutObjectInput{DataStream: stream})
	assert.NotNil(t, err)
	_, err = instance.PutObjectCannedAcl(context.TODO(), &file.PutObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.PutObjectTagging(context.TODO(), &file.PutObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.RestoreObject(context.TODO(), &file.RestoreObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.SignURL(context.TODO(), &file.SignURLInput{})
	assert.NotNil(t, err)

	_, err = instance.UploadPartCopy(context.TODO(), &file.UploadPartCopyInput{
		CopySource: &file.CopySource{CopySourceBucket: "bucket", CopySourceKey: "key"},
	})
	assert.NotNil(t, err)

	_, err = instance.UploadPart(context.TODO(), &file.UploadPartInput{})
	assert.NotNil(t, err)

	err = instance.UpdateDownLoadBandwidthRateLimit(context.TODO(), &file.UpdateBandwidthRateLimitInput{})
	assert.NotNil(t, err)

	err = instance.UpdateUpLoadBandwidthRateLimit(context.TODO(), &file.UpdateBandwidthRateLimitInput{})
	assert.NotNil(t, err)

}
