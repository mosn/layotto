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

package s3

import (
	"context"
	"errors"
	"testing"

	mockoss "mosn.io/layotto/pkg/mock/components/oss"

	"mosn.io/pkg/buffer"

	mocks3 "mosn.io/layotto/pkg/mock/runtime/oss"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	s3 "mosn.io/layotto/spec/proto/extension/v1"

	"mosn.io/layotto/components/file"
	l8s3 "mosn.io/layotto/components/file"

	"github.com/golang/mock/gomock"

	"mosn.io/layotto/pkg/grpc"
)

type MockDataStream struct {
	buffer.IoBuffer
}

func (m *MockDataStream) Close() error {
	m.CloseWithError(nil)
	return nil
}

//TestInitClient
func TestInitClient(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	// Test InitClient function
	initReq := &s3.InitInput{StoreName: "NoStore", Metadata: map[string]string{"k": "v"}}
	ctx := context.TODO()
	_, err := s3Server.InitClient(ctx, initReq)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	mockossServer.EXPECT().InitClient(ctx, &l8s3.InitRequest{Metadata: initReq.Metadata}).Return(nil)
	initReq.StoreName = "mockossServer"
	_, err = s3Server.InitClient(ctx, initReq)
	assert.Nil(t, err)
	mockossServer.EXPECT().InitClient(ctx, &l8s3.InitRequest{Metadata: initReq.Metadata}).Return(errors.New("init fail"))
	_, err = s3Server.InitClient(ctx, initReq)
	assert.Equal(t, err.Error(), "init fail")

}

// TestGetObject
func TestGetObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	// Test GetObject function
	ctx := context.TODO()
	mockServer := mocks3.NewMockObjectStorageService_GetObjectServer(ctrl)
	getObjectReq := &s3.GetObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object"}
	err := s3Server.GetObject(getObjectReq, mockServer)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	iobuf := buffer.NewIoBufferBytes([]byte("hello"))
	dataStream := &MockDataStream{iobuf}
	output := &file.GetObjectOutput{Etag: "tag"}
	output.DataStream = dataStream
	mockServer.EXPECT().Context().Return(ctx)
	mockossServer.EXPECT().GetObject(ctx, &l8s3.GetObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	getObjectReq.StoreName = "mockossServer"
	mockServer.EXPECT().Send(&s3.GetObjectOutput{Body: []byte("hello"), Etag: "tag"}).Times(1)
	err = s3Server.GetObject(getObjectReq, mockServer)
	assert.Nil(t, err)
}

// TestDeleteObject
func TestDeleteObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	deleteObjectReq := &s3.DeleteObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object"}
	_, err := s3Server.DeleteObject(ctx, deleteObjectReq)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.DeleteObjectOutput{DeleteMarker: false, VersionId: "123"}
	mockossServer.EXPECT().DeleteObject(ctx, &l8s3.DeleteObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	deleteObjectReq.StoreName = "mockossServer"
	resp, err := s3Server.DeleteObject(ctx, deleteObjectReq)
	assert.Nil(t, err)
	assert.Equal(t, false, resp.DeleteMarker)
	assert.Equal(t, "123", resp.VersionId)
}

//TestPutObjectTagging
func TestPutObjectTagging(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.PutObjectTaggingInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", Tags: map[string]string{"key": "value"}, VersionId: "123"}
	_, err := s3Server.PutObjectTagging(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.PutObjectTaggingOutput{}
	mockossServer.EXPECT().PutObjectTagging(ctx, &l8s3.PutObjectTaggingInput{Bucket: "layotto", Key: "object", VersionId: "123", Tags: map[string]string{"key": "value"}}).Return(output, nil)
	req.StoreName = "mockossServer"
	_, err = s3Server.PutObjectTagging(ctx, req)
	assert.Nil(t, err)
}

//TestDeleteObjectTagging
func TestDeleteObjectTagging(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.DeleteObjectTaggingInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", VersionId: "123"}
	_, err := s3Server.DeleteObjectTagging(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.DeleteObjectTaggingOutput{VersionId: "123"}
	mockossServer.EXPECT().DeleteObjectTagging(ctx, &l8s3.DeleteObjectTaggingInput{Bucket: "layotto", Key: "object", VersionId: "123"}).Return(output, nil)
	req.StoreName = "mockossServer"
	_, err = s3Server.DeleteObjectTagging(ctx, req)
	assert.Nil(t, err)
}

//TestGetObjectTagging
func TestGetObjectTagging(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.GetObjectTaggingInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", VersionId: "123"}
	_, err := s3Server.GetObjectTagging(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.GetObjectTaggingOutput{Tags: map[string]string{"key": "value"}, VersionId: "123"}
	mockossServer.EXPECT().GetObjectTagging(ctx, &l8s3.GetObjectTaggingInput{Bucket: "layotto", Key: "object", VersionId: "123"}).Return(output, nil)
	req.StoreName = "mockossServer"
	resp, err := s3Server.GetObjectTagging(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "value", resp.Tags["key"])
	assert.Equal(t, "123", resp.VersionId)
}

//TestCopyObject
func TestCopyObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.CopyObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object"}
	_, err := s3Server.CopyObject(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.CopyObjectOutput{CopyObjectResult: &l8s3.CopyObjectResult{ETag: "etag"}}
	mockossServer.EXPECT().CopyObject(ctx, &l8s3.CopyObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	req.StoreName = "mockossServer"
	resp, err := s3Server.CopyObject(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "etag", resp.CopyObjectResult.Etag)
}

//TestDeleteObjects
func TestDeleteObjects(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.DeleteObjectsInput{StoreName: "NoStore", Bucket: "layotto", Delete: &s3.Delete{Quiet: true, Objects: []*s3.ObjectIdentifier{{Key: "object", VersionId: "version"}}}}
	_, err := s3Server.DeleteObjects(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.DeleteObjectsOutput{Deleted: []*l8s3.DeletedObject{{DeleteMarker: true, VersionId: "version"}}}
	mockossServer.EXPECT().DeleteObjects(ctx, &l8s3.DeleteObjectsInput{Bucket: "layotto", Delete: &l8s3.Delete{Quiet: true, Objects: []*l8s3.ObjectIdentifier{{Key: "object", VersionId: "version"}}}}).Return(output, nil)
	req.StoreName = "mockossServer"
	resp, err := s3Server.DeleteObjects(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, true, resp.Deleted[0].DeleteMarker)
	assert.Equal(t, "version", resp.Deleted[0].VersionId)
}

//TestListObjects
func TestListObjects(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss["mockossServer"] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.ListObjectsInput{
		StoreName:           "NoStore",
		Bucket:              "layotto",
		Delimiter:           "delimiter",
		EncodingType:        "EncodingType",
		ExpectedBucketOwner: "ExpectedBucketOwner",
		Marker:              "Marker",
		MaxKeys:             1,
		Prefix:              "Prefix",
		RequestPayer:        "RequestPayer",
	}
	_, err := s3Server.ListObjects(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.ListObjectsOutput{Delimiter: "delimiter", IsTruncated: true}
	mockossServer.EXPECT().ListObjects(ctx,
		&l8s3.ListObjectsInput{
			Bucket:              "layotto",
			Delimiter:           "delimiter",
			EncodingType:        "EncodingType",
			ExpectedBucketOwner: "ExpectedBucketOwner",
			Marker:              "Marker",
			MaxKeys:             1,
			Prefix:              "Prefix",
			RequestPayer:        "RequestPayer",
		},
	).Return(output, nil)
	req.StoreName = "mockossServer"
	resp, err := s3Server.ListObjects(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, true, resp.IsTruncated)
	assert.Equal(t, "delimiter", resp.Delimiter)
}
