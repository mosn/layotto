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
package huaweicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/pkg/utils"
)

const (
	config = `{
		"endpoint": "your endpoint",
		"accessKeyID": "your accessKeyID",
		"accessKeySecret": "your accessKeySecret",
		"region": "your region"
	}`
	bucket = "your bucket"

	// The AK cannot be escrowed at present, skip the checking logic when false
	haveAK = false
)

var h *HuaweicloudOSS

func init() {
	h = &HuaweicloudOSS{}
	h.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(config)}})
}

func Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if haveAK {
		return assert.Equal(t, expected, actual, msgAndArgs...)
	}
	return true
}

func Nil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {
	if haveAK {
		return assert.Nil(t, object, msgAndArgs...)
	}
	return true
}

func NotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {
	if haveAK {
		return assert.NotNil(t, object, msgAndArgs...)
	}
	return true
}

func True(t *testing.T, value bool, msgAndArgs ...interface{}) bool {
	if haveAK {
		return assert.True(t, value, msgAndArgs...)
	}
	return true
}

func NotZero(t *testing.T, i interface{}, msgAndArgs ...interface{}) bool {
	if haveAK {
		return assert.NotZero(t, i, msgAndArgs...)
	}
	return true
}

func testRecover() {
	if !haveAK {
		_ = recover()
	}
}

func TestInitHuaweicloudOBS(t *testing.T) {
	defer testRecover()
	h := &HuaweicloudOSS{}
	_, err := h.getClient()
	Equal(t, err, utils.ErrNotInitClient)

	err = h.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte("hello")}})
	Equal(t, oss.ErrInvalid, err)
	err = h.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(config)}})
	Nil(t, err)

	cli, err := h.getClient()
	Nil(t, err)
	NotNil(t, cli)
}

func TestGetObject(t *testing.T) {
	defer testRecover()
	key := "create_multipart_upload_completed_test"
	putObject(t, key)

	input := &oss.GetObjectInput{Bucket: bucket, Key: key}
	output, err := h.GetObject(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	defer output.DataStream.Close()
	body, err := io.ReadAll(output.DataStream)
	Nil(t, err)
	content := string(body)
	println(content)
}

func TestPutObject(t *testing.T) {
	defer testRecover()
	key := "put_object_test"
	putObject(t, key)
}

func TestDeleteObject(t *testing.T) {
	defer testRecover()
	key := "delete_object_test"
	putObject(t, key)

	input := &oss.DeleteObjectInput{Bucket: bucket, Key: key}
	output, err := h.DeleteObject(context.TODO(), input)
	Nil(t, err)
	True(t, output.DeleteMarker)
}

func TestCopyObject(t *testing.T) {
	defer testRecover()
	key := "copy_obejct_test"
	tempKey := "copy_obejct_test_temp"
	putObject(t, key)

	input := &oss.CopyObjectInput{
		Bucket: bucket,
		Key:    tempKey,
		CopySource: &oss.CopySource{
			CopySourceBucket: bucket,
			CopySourceKey:    key,
		},
	}
	output, err := h.CopyObject(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
}

func TestDeleteObjects(t *testing.T) {
	defer testRecover()
	key1 := "delete_objects_test_1"
	key2 := "delete_objects_test_2"
	putObject(t, key1)
	putObject(t, key2)

	input := &oss.DeleteObjectsInput{
		Bucket: bucket,
		Delete: &oss.Delete{
			Objects: nil,
			Quiet:   false,
		},
	}
	oList := make([]*oss.ObjectIdentifier, 0, 2)
	o1 := &oss.ObjectIdentifier{
		Key: key1,
	}
	o2 := &oss.ObjectIdentifier{
		Key: key2,
	}
	oList = append(oList, o1, o2)
	input.Delete.Objects = oList

	output, err := h.DeleteObjects(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	fmt.Printf("output:%+v", output)
}

func TestListObjects(t *testing.T) {
	defer testRecover()
	input := &oss.ListObjectsInput{Bucket: "layotto", MaxKeys: 100}
	output, err := h.ListObjects(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
}

func TestPutObjectCannedAcl(t *testing.T) {
	defer testRecover()
	key := "put_object_canned_acl_test"
	putObject(t, key)

	input := &oss.PutObjectCannedAclInput{
		Bucket: bucket,
		Key:    key,
		Acl:    "public-read-write",
	}
	output, err := h.PutObjectCannedAcl(context.TODO(), input)
	Nil(t, err)
	fmt.Printf("output:%+v", output)
}

func TestRestoreObject(t *testing.T) {
	defer testRecover()
	key := "cold_file"

	input := &oss.PutObjectInput{
		DataStream:    strings.NewReader("cold file test"),
		Bucket:        bucket,
		Key:           key,
		StorageClass:  "COLD",
		ContentLength: 0,
	}
	output, err := h.PutObject(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)

	input1 := &oss.RestoreObjectInput{
		Bucket: bucket,
		Key:    key,
		RestoreRequest: oss.RestoreRequest{
			Days: 1,
			Tier: "Expedited",
		},
	}
	output1, err := h.RestoreObject(context.TODO(), input1)
	Nil(t, err)
	NotNil(t, output1)
}

func TestMultipartCreateUploadCompleted(t *testing.T) {
	defer testRecover()
	key := "create_multipart_upload_completed_test"

	input := &oss.CreateMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
	}
	output, err := h.CreateMultipartUpload(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)

	file, err := os.Open("/Users/mucan/Downloads/multipart_test")
	Nil(t, err)

	input1 := &oss.UploadPartInput{
		DataStream: file,
		Bucket:     bucket,
		Key:        key,
		UploadId:   output.UploadId,
		PartNumber: 1,
	}

	output1, err := h.UploadPart(context.TODO(), input1)
	Nil(t, err)
	NotNil(t, output1)
	file.Close()

	file, err = os.Open("/Users/mucan/Downloads/multipart_test")
	Nil(t, err)
	input1.PartNumber = 2
	input1.DataStream = file
	output1, err = h.UploadPart(context.TODO(), input1)
	Nil(t, err)
	NotNil(t, output1)
	file.Close()

	input2 := &oss.CompleteMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: output.UploadId,
	}
	completedMultipartUpload := &oss.CompletedMultipartUpload{}
	parts := make([]*oss.CompletedPart, 0, 2)
	part := &oss.CompletedPart{
		ETag:       output1.ETag,
		PartNumber: 1,
	}
	parts = append(parts, part)
	completedMultipartUpload.Parts = parts
	input2.MultipartUpload = completedMultipartUpload
	output4, err := h.CompleteMultipartUpload(context.TODO(), input2)
	Nil(t, err)
	NotNil(t, output4)
}

func TestMultipartCreateCopyCompleted(t *testing.T) {
	defer testRecover()
	key := "create_multipart_copy_completed_test"
	sourceKey := "copy_part_source"
	putObject(t, sourceKey)

	input := &oss.CreateMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
	}
	output, err := h.CreateMultipartUpload(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)

	input1 := &oss.UploadPartCopyInput{
		Bucket: bucket,
		Key:    key,
		CopySource: &oss.CopySource{
			CopySourceBucket: bucket,
			CopySourceKey:    sourceKey,
		},
		PartNumber: 1,
		UploadId:   output.UploadId,
	}

	output1, err := h.UploadPartCopy(context.TODO(), input1)
	Nil(t, err)
	NotNil(t, output1)

	input2 := &oss.CompleteMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: output.UploadId,
	}
	completedMultipartUpload := &oss.CompletedMultipartUpload{}
	parts := make([]*oss.CompletedPart, 0, 2)
	part := &oss.CompletedPart{
		ETag:       output1.CopyPartResult.ETag,
		PartNumber: 1,
	}
	parts = append(parts, part)
	completedMultipartUpload.Parts = parts
	input2.MultipartUpload = completedMultipartUpload
	output4, err := h.CompleteMultipartUpload(context.TODO(), input2)
	Nil(t, err)
	NotNil(t, output4)
}

func TestAbortMultipartUpload(t *testing.T) {
	defer testRecover()
	key := "multi_part_upload_6"
	uploadId := initiateMultipartUpload(t, key)
	NotZero(t, uploadId)
	input := &oss.AbortMultipartUploadInput{Bucket: bucket, Key: key, UploadId: uploadId}
	output, err := h.AbortMultipartUpload(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
}

func TestListMultipartUploads(t *testing.T) {
	defer testRecover()
	key := "multi_part_upload_6"
	uploadId := initiateMultipartUpload(t, key)
	input := &oss.ListMultipartUploadsInput{Bucket: bucket, UploadIdMarker: uploadId, KeyMarker: key}
	output, err := h.ListMultipartUploads(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
}

func TestListObjectVersions(t *testing.T) {
	defer testRecover()
	input := &oss.ListObjectVersionsInput{Bucket: "layotto"}
	output, err := h.ListObjectVersions(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	fmt.Printf("output:%+v", output)
}

func TestHeadObject(t *testing.T) {
	defer testRecover()
	input := &oss.HeadObjectInput{
		Bucket:      "layotto",
		Key:         "cold_file",
		WithDetails: true,
	}
	output, err := h.HeadObject(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	fmt.Printf("output:%+v", output)
}

func TestIsObjectExist(t *testing.T) {
	defer testRecover()
	key := "is_object_exist_test"
	putObject(t, key)

	input := &oss.IsObjectExistInput{
		Bucket: bucket,
		Key:    key,
	}
	output, err := h.IsObjectExist(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	True(t, output.FileExist)
}

func TestSignUrl(t *testing.T) {
	defer testRecover()
	key := "sign_url_test"
	putObject(t, key)

	input := &oss.SignURLInput{
		Bucket: bucket,
		Key:    key,
		Method: "GET",
	}
	input.ExpiredInSec = 1 * 60 * 60

	output, err := h.SignURL(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	fmt.Printf("url:%s", output.SignedUrl)
}

func TestAppendObject(t *testing.T) {
	defer testRecover()
	key := "append_object_test_2"

	input := &oss.AppendObjectInput{
		DataStream: strings.NewReader(" append part1"),
		Bucket:     bucket,
		Key:        key,
		// Position:   ,
	}

	output, err := h.AppendObject(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)
	fmt.Printf("pos:%d", output.AppendPosition)
}

func TestListParts(t *testing.T) {
	defer testRecover()
	key := "create_multipart_upload_completed_test"

	input := &oss.CreateMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
	}
	output, err := h.CreateMultipartUpload(context.TODO(), input)
	Nil(t, err)
	NotNil(t, output)

	file, err := os.Open("/Users/mucan/Downloads/multipart_test")
	Nil(t, err)

	input1 := &oss.UploadPartInput{
		DataStream: file,
		Bucket:     bucket,
		Key:        key,
		UploadId:   output.UploadId,
		PartNumber: 1,
	}

	output1, err := h.UploadPart(context.TODO(), input1)
	Nil(t, err)
	NotNil(t, output1)
	file.Close()

	file, err = os.Open("/Users/mucan/Downloads/multipart_test")
	Nil(t, err)
	input1.PartNumber = 2
	input1.DataStream = file
	output1, err = h.UploadPart(context.TODO(), input1)
	Nil(t, err)
	NotNil(t, output1)
	file.Close()

	input3 := &oss.ListPartsInput{Bucket: bucket, Key: key, UploadId: output.UploadId}
	output3, err := h.ListParts(context.TODO(), input3)
	Nil(t, err)
	NotNil(t, output3)

	input2 := &oss.CompleteMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: output.UploadId,
	}
	completedMultipartUpload := &oss.CompletedMultipartUpload{}
	parts := make([]*oss.CompletedPart, 0, 2)
	part := &oss.CompletedPart{
		ETag:       output1.ETag,
		PartNumber: 1,
	}
	parts = append(parts, part)
	completedMultipartUpload.Parts = parts
	input2.MultipartUpload = completedMultipartUpload
	output4, err := h.CompleteMultipartUpload(context.TODO(), input2)
	Nil(t, err)
	NotNil(t, output4)
}

func initiateMultipartUpload(t *testing.T, key string) (uploadId string) {
	defer testRecover()
	input := &oss.CreateMultipartUploadInput{Bucket: bucket, Key: key}
	output, err := h.CreateMultipartUpload(context.TODO(), input)
	Nil(t, err)
	uploadId = output.UploadId
	return
}

func putObject(t *testing.T, key string) {
	defer testRecover()
	input := &oss.PutObjectInput{
		DataStream: strings.NewReader("hello obs"),
		Bucket:     bucket,
		Key:        key,
	}
	_, err := h.PutObject(context.TODO(), input)
	Nil(t, err)
}
