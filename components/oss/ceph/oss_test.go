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

package ceph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/buffer"

	"mosn.io/layotto/components/oss"
)

const (
	confWithoutUidAndBucket = `
				{
					"endpoint": "http://10.211.55.13:7480",
					"accessKeyID": "QRF1XGIPZ9TB094ETTWU",
					"accessKeySecret": "6gF61QLVduFIFDKPBzc6gKkOsQY2HpAt7vM5mRAA"
				}
			`
)

func TestCephDefaultInitFunc(t *testing.T) {
	a := &CephOss{}
	err := a.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte("hello")}})
	assert.Equal(t, err, oss.ErrInvalid)
	assert.Nil(t, a.client)
}

func TestCephOss_GetObject(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.GetObjectInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
	}
	out, err := cephOss.GetObject(context.Background(), req)
	assert.Nil(t, err)

	data, err := ioutil.ReadAll(out.DataStream)
	assert.Nil(t, err)
	fmt.Println(string(data))
}

func TestCephOss_PutObject(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	reader, err := os.Open("/Users/apple/Desktop/untitled 3.txt")
	assert.Nil(t, err)
	req := &oss.PutObjectInput{
		DataStream: reader,
		Bucket:     "test.bucket",
		Key:        "TestPut.txt",
	}

	out, err := cephOss.PutObject(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_DeleteObject(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.DeleteObjectInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
	}

	out, err := cephOss.DeleteObject(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_DeleteObjects(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	d := &oss.Delete{
		Objects: []*oss.ObjectIdentifier{
			{Key: "TestPut.txt"},
			{Key: "a.txt"},
		},
	}
	req := &oss.DeleteObjectsInput{
		Bucket: "test.bucket",
		Delete: d,
	}
	out, err := cephOss.DeleteObjects(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_ListObjects(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.ListObjectsInput{
		Bucket: "test.bucket",
	}

	out, err := cephOss.ListObjects(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_PutObjectTagging(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	tags := map[string]string{
		"Test": "True",
		"HAHA": "haha",
	}
	req := &oss.PutObjectTaggingInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
		Tags:   tags,
	}

	out, err := cephOss.PutObjectTagging(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_GetObjectTagging(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.GetObjectTaggingInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
	}

	out, err := cephOss.GetObjectTagging(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_DeleteObjectTagging(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.DeleteObjectTaggingInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
	}

	out, err := cephOss.DeleteObjectTagging(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_CopyObject(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	source := oss.CopySource{
		CopySourceBucket: "test.bucket",
		CopySourceKey:    "haha.txt",
	}
	req := &oss.CopyObjectInput{
		Bucket:     "test.bucket",
		Key:        "b.txt",
		CopySource: &source,
	}

	out, err := cephOss.CopyObject(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_PutObjectCannedAcl(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.PutObjectCannedAclInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
		Acl:    "private",
	}

	out, err := cephOss.PutObjectCannedAcl(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_GetObjectCannedAcl(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.GetObjectCannedAclInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
	}

	out, err := cephOss.GetObjectCannedAcl(context.Background(), req)
	assert.Nil(t, err)
	fmt.Printf("%+v\n", out)
}

func TestCephOss_MultipartUploadWithAbort(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	key := "TestMultipartUploadWithAbort"
	uploadId := createMultipartUpload(t, cephOss, key)

	f, err := os.Open("/Users/apple/Downloads/TestMultipartUpload.zip")
	assert.Nil(t, err)
	uploadPart(t, cephOss, key, uploadId, 1, f)

	abortMultipartUpload(t, cephOss, key, uploadId)
}

func TestCephOss_MultipartUploadWithComplete(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	key := "TestMultipartUploadWithComplete"
	uploadId := createMultipartUpload(t, cephOss, key)
	f, err := os.Open("/Users/apple/Downloads/TestMultipartUpload.zip")
	assert.Nil(t, err)
	go listMultipartUploads(t, cephOss)
	eTag := uploadPart(t, cephOss, key, uploadId, 1, f)
	listParts(t, cephOss, key, uploadId)

	completeMultipartUpload(t, cephOss, key, uploadId, eTag, 1)
}

func TestCephOss_UploadPartCopy(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	key := "TestUploadPartCopy"
	uploadId := createMultipartUpload(t, cephOss, key)
	assert.Nil(t, err)
	eTag := uploadPartCopy(t, cephOss, uploadId, 1)

	completeMultipartUpload(t, cephOss, key, uploadId, eTag, 1)
}

func TestCephOss_ListObjectVersions(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.ListObjectVersionsInput{
		Bucket: "test.bucket",
	}

	out, err := cephOss.ListObjectVersions(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_HeadObject(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.HeadObjectInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
	}

	out, err := cephOss.HeadObject(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_IsObjectExist(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.IsObjectExistInput{
		Bucket: "test.bucket",
		Key:    "a.txt",
	}

	out, err := cephOss.IsObjectExist(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func TestCephOss_SignURL(t *testing.T) {
	cephOss := NewCephOss()
	err := cephOss.Init(context.TODO(), &oss.Config{Metadata: map[string]json.RawMessage{oss.BasicConfiguration: []byte(confWithoutUidAndBucket)}})
	assert.Nil(t, err)

	req := &oss.SignURLInput{
		Bucket: "test.bucket",
		Key:    "TestPut.txt",
		Method: "Get",
	}

	out, err := cephOss.SignURL(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func createMultipartUpload(t *testing.T, cephOss oss.Oss, key string) (uploadId string) {
	fmt.Println("=====[CreateMultipartUpload]")
	req := &oss.CreateMultipartUploadInput{
		Bucket: "test.bucket",
		Key:    key,
	}
	out, err := cephOss.CreateMultipartUpload(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
	return out.UploadId
}
func uploadPart(t *testing.T, cephOss oss.Oss, key string, uploadId string, partNumber int32, dataStream io.Reader) (etag string) {
	assert.True(t, (1 <= partNumber) && (partNumber <= 10000))
	req := &oss.UploadPartInput{
		Bucket:     "test.bucket",
		Key:        key,
		UploadId:   uploadId,
		PartNumber: partNumber,
		DataStream: dataStream,
	}
	out, err := cephOss.UploadPart(context.Background(), req)
	assert.Nil(t, err)
	fmt.Printf("=====[UploadPart %d]\n", partNumber)
	printOutput(out)
	return out.ETag
}
func uploadPartCopy(t *testing.T, cephOss oss.Oss, uploadId string, partNumber int32) (etag string) {
	assert.True(t, (1 <= partNumber) && (partNumber <= 10000))
	copySource := oss.CopySource{
		CopySourceBucket: "test.bucket",
		CopySourceKey:    "TestMultipartUpload",
	}
	req := &oss.UploadPartCopyInput{
		Bucket:     "test.bucket",
		Key:        "TestUploadPartCopy",
		UploadId:   uploadId,
		PartNumber: partNumber,
		CopySource: &copySource,
	}
	out, err := cephOss.UploadPartCopy(context.Background(), req)
	assert.Nil(t, err)
	fmt.Printf("=====[UploadPartCopy %d]\n", partNumber)
	printOutput(out)
	return out.CopyPartResult.ETag
}
func abortMultipartUpload(t *testing.T, cephOss oss.Oss, key string, uploadId string) {
	fmt.Println("=====[AbortMultipartUpload]")
	req := &oss.AbortMultipartUploadInput{
		Bucket:   "test.bucket",
		Key:      key,
		UploadId: uploadId,
	}
	out, err := cephOss.AbortMultipartUpload(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}
func completeMultipartUpload(t *testing.T, cephOss oss.Oss, key string, uploadId string, eTag string, partNumber int32) {
	fmt.Println("=====[CompleteMultipartUpload]")
	multipartUpload := &oss.CompletedMultipartUpload{
		Parts: []*oss.CompletedPart{
			{ETag: eTag, PartNumber: partNumber},
		},
	}
	req := &oss.CompleteMultipartUploadInput{
		Bucket:          "test.bucket",
		Key:             key,
		UploadId:        uploadId,
		MultipartUpload: multipartUpload,
	}
	out, err := cephOss.CompleteMultipartUpload(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}
func listMultipartUploads(t *testing.T, cephOss oss.Oss) {
	time.Sleep(time.Second)
	fmt.Println("=====[ListMultipartUploads]")
	req := &oss.ListMultipartUploadsInput{
		Bucket: "test.bucket",
	}
	out, err := cephOss.ListMultipartUploads(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}
func listParts(t *testing.T, cephOss oss.Oss, key, uploadId string) {
	time.Sleep(time.Second)
	fmt.Println("=====[ListParts]")
	req := &oss.ListPartsInput{
		Bucket:   "test.bucket",
		Key:      key,
		UploadId: uploadId,
	}
	out, err := cephOss.ListParts(context.Background(), req)
	assert.Nil(t, err)
	printOutput(out)
}

func printOutput(v interface{}) {
	bs, _ := json.Marshal(v)
	var bf bytes.Buffer
	err := json.Indent(&bf, bs, "", "\t")
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}
	fmt.Println(bf.String())
}

func TestCephOss(t *testing.T) {
	instance := &CephOss{}
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
}
