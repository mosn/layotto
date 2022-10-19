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

package aws

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/jinzhu/copier"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"mosn.io/layotto/components/oss"

	"mosn.io/pkg/buffer"

	"github.com/stretchr/testify/assert"
)

const (
	confWithoutUidAndBucket = `
				{
					"endpoint": "endpoint_address",
					"accessKeyID": "accessKey",
					"accessKeySecret": "secret"
				}
			`
)

func TestAwsDefaultInitFunc(t *testing.T) {
	a := &AwsOss{}
	err := a.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte("hello")}})
	assert.Equal(t, err, oss.ErrInvalid)
	assert.Nil(t, a.client)

}

func TestAwsOss(t *testing.T) {
	instance := &AwsOss{}
	err := instance.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	appendObjectResp, err := instance.AppendObject(context.TODO(), &oss.AppendObjectInput{})
	assert.Equal(t, errors.New("AppendObject method not supported on AWS"), err)
	assert.Nil(t, appendObjectResp)

	_, err = instance.AbortMultipartUpload(context.TODO(), &oss.AbortMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CompleteMultipartUpload(context.TODO(), &oss.CompleteMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.CopyObject(context.TODO(), &oss.CopyObjectInput{})
	assert.Equal(t, errors.New("must specific copy_source"), err)

	_, err = instance.CopyObject(context.TODO(), &oss.CopyObjectInput{
		CopySource: &oss.CopySource{CopySourceBucket: "bucket", CopySourceKey: "key"},
	})
	assert.NotEqual(t, errors.New("must specific copy_source"), err)
	_, err = instance.CreateMultipartUpload(context.TODO(), &oss.CreateMultipartUploadInput{})
	assert.NotNil(t, err)

	_, err = instance.DeleteObject(context.TODO(), &oss.DeleteObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjects(context.TODO(), &oss.DeleteObjectsInput{
		Delete: &oss.Delete{Objects: []*oss.ObjectIdentifier{{Key: "object", VersionId: "version"}}},
	})
	assert.NotNil(t, err)
	_, err = instance.DeleteObjectTagging(context.TODO(), &oss.DeleteObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.GetObject(context.TODO(), &oss.GetObjectInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectCannedAcl(context.TODO(), &oss.GetObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.GetObjectTagging(context.TODO(), &oss.GetObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.HeadObject(context.TODO(), &oss.HeadObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.IsObjectExist(context.TODO(), &oss.IsObjectExistInput{})
	assert.NotNil(t, err)

	_, err = instance.ListParts(context.TODO(), &oss.ListPartsInput{})
	assert.NotNil(t, err)

	_, err = instance.ListMultipartUploads(context.TODO(), &oss.ListMultipartUploadsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjects(context.TODO(), &oss.ListObjectsInput{})
	assert.NotNil(t, err)
	_, err = instance.ListObjectVersions(context.TODO(), &oss.ListObjectVersionsInput{})
	assert.NotNil(t, err)

	stream := buffer.NewIoBufferString("hello")
	_, err = instance.PutObject(context.TODO(), &oss.PutObjectInput{DataStream: stream})
	assert.NotNil(t, err)
	_, err = instance.PutObjectCannedAcl(context.TODO(), &oss.PutObjectCannedAclInput{})
	assert.NotNil(t, err)
	_, err = instance.PutObjectTagging(context.TODO(), &oss.PutObjectTaggingInput{})
	assert.NotNil(t, err)

	_, err = instance.RestoreObject(context.TODO(), &oss.RestoreObjectInput{})
	assert.NotNil(t, err)

	_, err = instance.SignURL(context.TODO(), &oss.SignURLInput{})
	assert.NotNil(t, err)

	_, err = instance.UploadPartCopy(context.TODO(), &oss.UploadPartCopyInput{
		CopySource: &oss.CopySource{CopySourceBucket: "bucket", CopySourceKey: "key"},
	})
	assert.NotNil(t, err)

	_, err = instance.UploadPart(context.TODO(), &oss.UploadPartInput{})
	assert.NotNil(t, err)

	err = instance.UpdateDownloadBandwidthRateLimit(context.TODO(), &oss.UpdateBandwidthRateLimitInput{})
	assert.NotNil(t, err)

	err = instance.UpdateUploadBandwidthRateLimit(context.TODO(), &oss.UpdateBandwidthRateLimitInput{})
	assert.NotNil(t, err)

	_, err = oss.GetGetObjectOutput(&s3.GetObjectOutput{})
	assert.Nil(t, err)
	_, err = oss.GetPutObjectOutput(&manager.UploadOutput{})
	assert.Nil(t, err)
	_, err = oss.GetDeleteObjectOutput(&s3.DeleteObjectOutput{})
	assert.Nil(t, err)
	_, err = oss.GetDeleteObjectOutput(&s3.DeleteObjectOutput{})
	assert.Nil(t, err)
	_, err = oss.GetGetObjectTaggingOutput(&s3.GetObjectTaggingOutput{})
	assert.Nil(t, err)
	_, err = oss.GetListObjectsOutput(&s3.ListObjectsOutput{})
	assert.Nil(t, err)
	_, err = oss.GetGetObjectCannedAclOutput(&s3.GetObjectAclOutput{})
	assert.Nil(t, err)
	_, err = oss.GetUploadPartOutput(&s3.UploadPartOutput{})
	assert.Nil(t, err)
	_, err = oss.GetUploadPartCopyOutput(&s3.UploadPartCopyOutput{})
	assert.Nil(t, err)
	_, err = oss.GetListPartsOutput(&s3.ListPartsOutput{})
	assert.Nil(t, err)
	_, err = oss.GetListMultipartUploadsOutput(&s3.ListMultipartUploadsOutput{})
	assert.Nil(t, err)
	_, err = oss.GetListObjectVersionsOutput(&s3.ListObjectVersionsOutput{})
	assert.Nil(t, err)
}

func TestDeepCopy(t *testing.T) {
	value := "hello"
	t1 := time.Now()
	fromValue := &types.ObjectVersion{
		ETag:         &value,
		IsLatest:     true,
		Key:          &value,
		LastModified: &t1,
		Owner:        &types.Owner{DisplayName: &value, ID: &value},
		Size:         10,
		StorageClass: "hello",
		VersionId:    &value,
	}
	tovalue := &oss.ObjectVersion{}
	err := copier.CopyWithOption(tovalue, fromValue, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: []copier.TypeConverter{oss.TimeToInt64}})
	assert.Nil(t, err)
	assert.Equal(t, tovalue.Owner.DisplayName, value)
}
