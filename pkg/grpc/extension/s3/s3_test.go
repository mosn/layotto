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
	"io"
	"testing"

	"mosn.io/layotto/spec/proto/extension/v1/s3"

	l8s3 "mosn.io/layotto/components/oss"

	mockoss "mosn.io/layotto/pkg/mock/components/oss"

	"mosn.io/pkg/buffer"

	mocks3 "mosn.io/layotto/pkg/mock/runtime/oss"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"

	"mosn.io/layotto/pkg/grpc"
)

const (
	MOCKSERVER = "mockossServer"
	ByteSize   = 5
)

type MockDataStream struct {
	buffer.IoBuffer
}

func (m *MockDataStream) Close() error {
	m.CloseWithError(nil)
	return nil
}

// TestGetObject
func TestGetObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
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
	output := &l8s3.GetObjectOutput{Etag: "tag"}
	output.DataStream = dataStream
	mockServer.EXPECT().Context().Return(ctx)
	mockossServer.EXPECT().GetObject(ctx, &l8s3.GetObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	getObjectReq.StoreName = MOCKSERVER
	mockServer.EXPECT().Send(&s3.GetObjectOutput{Body: []byte("hello"), Etag: "tag"}).Times(1)
	err = s3Server.GetObject(getObjectReq, mockServer)
	assert.Nil(t, err)
}

// TestPutObject
func TestPutObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	// Test GetObject function
	ctx := context.TODO()
	mockStream := mocks3.NewMockObjectStorageService_PutObjectServer(ctrl)
	putObjectReq := &s3.PutObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", Body: []byte("put")}
	mockStream.EXPECT().Recv().Return(putObjectReq, nil)
	err := s3Server.PutObject(mockStream)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)

	putObjectReq.StoreName = MOCKSERVER
	output := &l8s3.PutObjectOutput{ETag: "tag"}
	mockStream.EXPECT().Context().Return(ctx)
	mockStream.EXPECT().Recv().Return(putObjectReq, nil)
	mockStream.EXPECT().SendAndClose(&s3.PutObjectOutput{Etag: "tag"}).Times(1)
	mockossServer.EXPECT().PutObject(ctx, &l8s3.PutObjectInput{DataStream: newPutObjectStreamReader(putObjectReq.Body, mockStream), Bucket: "layotto", Key: "object"}).Return(output, nil)
	err = s3Server.PutObject(mockStream)
	assert.Nil(t, err)

	mockStream.EXPECT().Recv().Return(nil, io.EOF)
	stream := newPutObjectStreamReader(putObjectReq.Body, mockStream)
	data := make([]byte, ByteSize)
	n, err := stream.Read(data)
	assert.Equal(t, 3, n)
	assert.Equal(t, io.EOF, err)
}

// TestUploadPart
func TestUploadPart(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	// Test GetObject function
	ctx := context.TODO()
	mockStream := mocks3.NewMockObjectStorageService_UploadPartServer(ctrl)
	UploadPartReq := &s3.UploadPartInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", Body: []byte("put")}
	mockStream.EXPECT().Recv().Return(UploadPartReq, nil)
	err := s3Server.UploadPart(mockStream)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)

	UploadPartReq.StoreName = MOCKSERVER
	output := &l8s3.UploadPartOutput{ETag: "tag"}
	mockStream.EXPECT().Context().Return(ctx)
	mockStream.EXPECT().Recv().Return(UploadPartReq, nil)
	mockStream.EXPECT().SendAndClose(&s3.UploadPartOutput{Etag: "tag"}).Times(1)
	mockossServer.EXPECT().UploadPart(ctx, &l8s3.UploadPartInput{DataStream: newUploadPartStreamReader(UploadPartReq.Body, mockStream), Bucket: "layotto", Key: "object"}).Return(output, nil)
	err = s3Server.UploadPart(mockStream)
	assert.Nil(t, err)

	mockStream.EXPECT().Recv().Return(nil, io.EOF)
	stream := newUploadPartStreamReader(UploadPartReq.Body, mockStream)
	data := make([]byte, ByteSize)
	n, err := stream.Read(data)
	assert.Equal(t, 3, n)
	assert.Equal(t, io.EOF, err)
}

// TestAppendObject
func TestAppendObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	// Test GetObject function
	ctx := context.TODO()
	mockStream := mocks3.NewMockObjectStorageService_AppendObjectServer(ctrl)
	req := &s3.AppendObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", Body: []byte("put")}
	mockStream.EXPECT().Recv().Return(req, nil)
	err := s3Server.AppendObject(mockStream)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)

	req.StoreName = MOCKSERVER
	output := &l8s3.AppendObjectOutput{AppendPosition: 123}
	mockStream.EXPECT().Context().Return(ctx)
	mockStream.EXPECT().Recv().Return(req, nil)
	mockStream.EXPECT().SendAndClose(&s3.AppendObjectOutput{AppendPosition: 123}).Times(1)
	mockossServer.EXPECT().AppendObject(ctx, &l8s3.AppendObjectInput{DataStream: newAppendObjectStreamReader(req.Body, mockStream), Bucket: "layotto", Key: "object"}).Return(output, nil)
	err = s3Server.AppendObject(mockStream)
	assert.Nil(t, err)

	mockStream.EXPECT().Recv().Return(nil, io.EOF)
	stream := newAppendObjectStreamReader(req.Body, mockStream)
	data := make([]byte, ByteSize)
	n, err := stream.Read(data)
	assert.Equal(t, 3, n)
	assert.Equal(t, io.EOF, err)
}

// TestDeleteObject
func TestDeleteObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	deleteObjectReq := &s3.DeleteObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object"}
	_, err := s3Server.DeleteObject(ctx, deleteObjectReq)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.DeleteObjectOutput{DeleteMarker: false, VersionId: "123"}
	mockossServer.EXPECT().DeleteObject(ctx, &l8s3.DeleteObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	deleteObjectReq.StoreName = MOCKSERVER
	resp, err := s3Server.DeleteObject(ctx, deleteObjectReq)
	assert.Nil(t, err)
	assert.Equal(t, false, resp.DeleteMarker)
	assert.Equal(t, "123", resp.VersionId)
}

// TestPutObjectTagging
func TestPutObjectTagging(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.PutObjectTaggingInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", Tags: map[string]string{"key": "value"}, VersionId: "123"}
	_, err := s3Server.PutObjectTagging(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.PutObjectTaggingOutput{}
	mockossServer.EXPECT().PutObjectTagging(ctx, &l8s3.PutObjectTaggingInput{Bucket: "layotto", Key: "object", VersionId: "123", Tags: map[string]string{"key": "value"}}).Return(output, nil)
	req.StoreName = MOCKSERVER
	_, err = s3Server.PutObjectTagging(ctx, req)
	assert.Nil(t, err)
}

// TestDeleteObjectTagging
func TestDeleteObjectTagging(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.DeleteObjectTaggingInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", VersionId: "123"}
	_, err := s3Server.DeleteObjectTagging(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.DeleteObjectTaggingOutput{VersionId: "123"}
	mockossServer.EXPECT().DeleteObjectTagging(ctx, &l8s3.DeleteObjectTaggingInput{Bucket: "layotto", Key: "object", VersionId: "123"}).Return(output, nil)
	req.StoreName = MOCKSERVER
	_, err = s3Server.DeleteObjectTagging(ctx, req)
	assert.Nil(t, err)
}

// TestGetObjectTagging
func TestGetObjectTagging(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.GetObjectTaggingInput{StoreName: "NoStore", Bucket: "layotto", Key: "object", VersionId: "123"}
	_, err := s3Server.GetObjectTagging(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.GetObjectTaggingOutput{Tags: map[string]string{"key": "value"}, VersionId: "123"}
	mockossServer.EXPECT().GetObjectTagging(ctx, &l8s3.GetObjectTaggingInput{Bucket: "layotto", Key: "object", VersionId: "123"}).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.GetObjectTagging(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "value", resp.Tags["key"])
	assert.Equal(t, "123", resp.VersionId)
}

// TestCopyObject
func TestCopyObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.CopyObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object"}
	_, err := s3Server.CopyObject(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.CopyObjectOutput{CopyObjectResult: &l8s3.CopyObjectResult{ETag: "etag"}}
	mockossServer.EXPECT().CopyObject(ctx, &l8s3.CopyObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.CopyObject(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "etag", resp.CopyObjectResult.Etag)
}

// TestDeleteObjects
func TestDeleteObjects(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.DeleteObjectsInput{StoreName: "NoStore", Bucket: "layotto", Delete: &s3.Delete{Quiet: true, Objects: []*s3.ObjectIdentifier{{Key: "object", VersionId: "version"}}}}
	_, err := s3Server.DeleteObjects(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.DeleteObjectsOutput{Deleted: []*l8s3.DeletedObject{{DeleteMarker: true, VersionId: "version"}}}
	mockossServer.EXPECT().DeleteObjects(ctx, &l8s3.DeleteObjectsInput{Bucket: "layotto", Delete: &l8s3.Delete{Quiet: true, Objects: []*l8s3.ObjectIdentifier{{Key: "object", VersionId: "version"}}}}).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.DeleteObjects(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, true, resp.Deleted[0].DeleteMarker)
	assert.Equal(t, "version", resp.Deleted[0].VersionId)
}

// TestListObjects
func TestListObjects(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
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
	req.StoreName = MOCKSERVER
	resp, err := s3Server.ListObjects(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, true, resp.IsTruncated)
	assert.Equal(t, "delimiter", resp.Delimiter)
}

// TestGetObjectCannedAcl
func TestGetObjectCannedAcl(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.GetObjectCannedAclInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
		VersionId: "versionId",
	}
	_, err := s3Server.GetObjectCannedAcl(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.GetObjectCannedAclOutput{CannedAcl: "public-read-write", RequestCharged: "yes"}
	mockossServer.EXPECT().GetObjectCannedAcl(ctx,
		&l8s3.GetObjectCannedAclInput{
			Bucket:    "layotto",
			Key:       "key",
			VersionId: "versionId",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.GetObjectCannedAcl(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "public-read-write", resp.CannedAcl)
	assert.Equal(t, "yes", resp.RequestCharged)
}

// TestPutObjectCannedAcl
func TestPutObjectCannedAcl(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.PutObjectCannedAclInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
		VersionId: "versionId",
	}
	_, err := s3Server.PutObjectCannedAcl(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.PutObjectCannedAclOutput{RequestCharged: "yes"}
	mockossServer.EXPECT().PutObjectCannedAcl(ctx,
		&l8s3.PutObjectCannedAclInput{
			Bucket:    "layotto",
			Key:       "key",
			VersionId: "versionId",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.PutObjectCannedAcl(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "yes", resp.RequestCharged)
}

// TestRestoreObject
func TestRestoreObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.RestoreObjectInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
		VersionId: "versionId",
	}
	_, err := s3Server.RestoreObject(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.RestoreObjectOutput{RestoreOutputPath: "yes", RequestCharged: "yes"}
	mockossServer.EXPECT().RestoreObject(ctx,
		&l8s3.RestoreObjectInput{
			Bucket:    "layotto",
			Key:       "key",
			VersionId: "versionId",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.RestoreObject(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "yes", resp.RequestCharged)
	assert.Equal(t, "yes", resp.RestoreOutputPath)
}

// TestCreateMultipartUpload
func TestCreateMultipartUpload(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.CreateMultipartUploadInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
	}
	_, err := s3Server.CreateMultipartUpload(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.CreateMultipartUploadOutput{Bucket: "layotto", Key: "object", UploadId: "123"}
	mockossServer.EXPECT().CreateMultipartUpload(ctx,
		&l8s3.CreateMultipartUploadInput{
			Bucket: "layotto",
			Key:    "key",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.CreateMultipartUpload(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "123", resp.UploadId)
	assert.Equal(t, "layotto", resp.Bucket)
	assert.Equal(t, "object", resp.Key)
}

// TestUploadPartCopy
func TestUploadPartCopy(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.UploadPartCopyInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
	}
	_, err := s3Server.UploadPartCopy(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.UploadPartCopyOutput{BucketKeyEnabled: true, CopyPartResult: &l8s3.CopyPartResult{ETag: "123", LastModified: 456}}
	mockossServer.EXPECT().UploadPartCopy(ctx,
		&l8s3.UploadPartCopyInput{
			Bucket: "layotto",
			Key:    "key",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.UploadPartCopy(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "123", resp.CopyPartResult.Etag)
	assert.Equal(t, int64(456), resp.CopyPartResult.LastModified)
	assert.Equal(t, true, resp.BucketKeyEnabled)
}

// TestCompleteMultipartUpload
func TestCompleteMultipartUpload(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.CompleteMultipartUploadInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
		UploadId:  "123",
	}
	_, err := s3Server.CompleteMultipartUpload(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.CompleteMultipartUploadOutput{
		BucketKeyEnabled: true,
		Expiration:       "expiration",
		ETag:             "etag",
	}
	mockossServer.EXPECT().CompleteMultipartUpload(ctx,
		&l8s3.CompleteMultipartUploadInput{
			Bucket:   "layotto",
			Key:      "key",
			UploadId: "123",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.CompleteMultipartUpload(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "expiration", resp.Expiration)
	assert.Equal(t, "etag", resp.Etag)
	assert.Equal(t, true, resp.BucketKeyEnabled)
}

// TestAbortMultipartUpload
func TestAbortMultipartUpload(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.AbortMultipartUploadInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "key",
		UploadId:  "123",
	}
	_, err := s3Server.AbortMultipartUpload(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.AbortMultipartUploadOutput{
		RequestCharged: "true",
	}
	mockossServer.EXPECT().AbortMultipartUpload(ctx,
		&l8s3.AbortMultipartUploadInput{
			Bucket:   "layotto",
			Key:      "key",
			UploadId: "123",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.AbortMultipartUpload(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "true", resp.RequestCharged)
}

// TestListMultipartUploads
func TestListMultipartUploads(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.ListMultipartUploadsInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
	}
	_, err := s3Server.ListMultipartUploads(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.ListMultipartUploadsOutput{
		Bucket: "layotto",
	}
	mockossServer.EXPECT().ListMultipartUploads(ctx,
		&l8s3.ListMultipartUploadsInput{
			Bucket: "layotto",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.ListMultipartUploads(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "layotto", resp.Bucket)
}

// TestListObjectVersions
func TestListObjectVersions(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.ListObjectVersionsInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		KeyMarker: "marker",
	}
	_, err := s3Server.ListObjectVersions(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.ListObjectVersionsOutput{
		Delimiter: "layotto",
	}
	mockossServer.EXPECT().ListObjectVersions(ctx,
		&l8s3.ListObjectVersionsInput{
			Bucket:    "layotto",
			KeyMarker: "marker",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.ListObjectVersions(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "layotto", resp.Delimiter)
}

// TestHeadObject
func TestHeadObject(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.HeadObjectInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "object",
	}
	_, err := s3Server.HeadObject(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.HeadObjectOutput{
		ResultMetadata: map[string]string{"key": "value"},
	}
	mockossServer.EXPECT().HeadObject(ctx,
		&l8s3.HeadObjectInput{
			Bucket: "layotto",
			Key:    "object",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.HeadObject(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"key": "value"}, resp.ResultMetadata)
}

// TestIsObjectExist
func TestIsObjectExist(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.IsObjectExistInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "object",
	}
	_, err := s3Server.IsObjectExist(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.IsObjectExistOutput{
		FileExist: true,
	}
	mockossServer.EXPECT().IsObjectExist(ctx,
		&l8s3.IsObjectExistInput{
			Bucket: "layotto",
			Key:    "object",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.IsObjectExist(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, true, resp.FileExist)
}

// TestSignURL
func TestSignURL(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.SignURLInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
		Key:       "object",
	}
	_, err := s3Server.SignURL(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.SignURLOutput{
		SignedUrl: "http://object",
	}
	mockossServer.EXPECT().SignURL(ctx,
		&l8s3.SignURLInput{
			Bucket: "layotto",
			Key:    "object",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.SignURL(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "http://object", resp.SignedUrl)
}

// TestUpdateDownLoadBandwidthRateLimit
func TestUpdateDownLoadBandwidthRateLimit(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.UpdateBandwidthRateLimitInput{
		StoreName:                    "NoStore",
		AverageRateLimitInBitsPerSec: 1,
	}
	_, err := s3Server.UpdateDownloadBandwidthRateLimit(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	mockossServer.EXPECT().UpdateDownloadBandwidthRateLimit(ctx,
		&l8s3.UpdateBandwidthRateLimitInput{
			AverageRateLimitInBitsPerSec: 1,
		},
	).Return(nil)
	req.StoreName = MOCKSERVER
	_, err = s3Server.UpdateDownloadBandwidthRateLimit(ctx, req)
	assert.Nil(t, err)
}

// TestUpdateUpLoadBandwidthRateLimit
func TestUpdateUpLoadBandwidthRateLimit(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.UpdateBandwidthRateLimitInput{
		StoreName:                    "NoStore",
		AverageRateLimitInBitsPerSec: 1,
	}
	_, err := s3Server.UpdateUploadBandwidthRateLimit(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	mockossServer.EXPECT().UpdateUploadBandwidthRateLimit(ctx,
		&l8s3.UpdateBandwidthRateLimitInput{
			AverageRateLimitInBitsPerSec: 1,
		},
	).Return(nil)
	req.StoreName = MOCKSERVER
	_, err = s3Server.UpdateUploadBandwidthRateLimit(ctx, req)
	assert.Nil(t, err)
}

// TestListParts
func TestListParts(t *testing.T) {
	// prepare oss server
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]l8s3.Oss{}}
	ctrl := gomock.NewController(t)
	mockossServer := mockoss.NewMockOss(ctrl)
	ac.Oss[MOCKSERVER] = mockossServer
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	ctx := context.TODO()
	req := &s3.ListPartsInput{
		StoreName: "NoStore",
		Bucket:    "layotto",
	}
	_, err := s3Server.ListParts(ctx, req)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	output := &l8s3.ListPartsOutput{
		Bucket: "layotto",
		Key:    "object",
	}
	mockossServer.EXPECT().ListParts(ctx,
		&l8s3.ListPartsInput{
			Bucket: "layotto",
		},
	).Return(output, nil)
	req.StoreName = MOCKSERVER
	resp, err := s3Server.ListParts(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "layotto", resp.Bucket)
}
