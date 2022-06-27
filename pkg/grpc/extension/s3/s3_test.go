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

	"mosn.io/pkg/buffer"

	mock_s3 "mosn.io/layotto/pkg/mock/runtime/oss"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	s3 "mosn.io/layotto/spec/proto/extension/v1"

	"mosn.io/layotto/components/file"
	l8s3 "mosn.io/layotto/components/file"

	"github.com/golang/mock/gomock"

	mock_oss "mosn.io/layotto/pkg/mock/components/oss"

	"mosn.io/layotto/pkg/grpc"
)

type MockDataStream struct {
	buffer.IoBuffer
}

func (m *MockDataStream) Close() error {
	m.CloseWithError(nil)
	return nil
}

func TestS3Server(t *testing.T) {
	ac := &grpc.ApplicationContext{AppId: "test", Oss: map[string]file.Oss{}}
	ctrl := gomock.NewController(t)
	mockoss := mock_oss.NewMockOss(ctrl)
	ac.Oss["mockoss"] = mockoss
	NewS3Server(ac)
	s3Server := &S3Server{appId: ac.AppId, ossInstance: ac.Oss}

	// Test InitClient function
	initReq := &s3.InitInput{StoreName: "NoStore", Metadata: map[string]string{"k": "v"}}
	ctx := context.TODO()
	_, err := s3Server.InitClient(ctx, initReq)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	mockoss.EXPECT().InitClient(ctx, &l8s3.InitRequest{Metadata: initReq.Metadata}).Return(nil)
	initReq.StoreName = "mockoss"
	_, err = s3Server.InitClient(ctx, initReq)
	assert.Nil(t, err)

	mockoss.EXPECT().InitClient(ctx, &l8s3.InitRequest{Metadata: initReq.Metadata}).Return(errors.New("init fail"))
	_, err = s3Server.InitClient(ctx, initReq)
	assert.Equal(t, err.Error(), "init fail")

	// Test GetObject function
	mockServer := mock_s3.NewMockObjectStorageService_GetObjectServer(ctrl)
	getObjectReq := &s3.GetObjectInput{StoreName: "NoStore", Bucket: "layotto", Key: "object"}
	err = s3Server.GetObject(getObjectReq, mockServer)
	assert.Equal(t, status.Errorf(codes.InvalidArgument, NotSupportStoreName, "NoStore"), err)
	iobuf := buffer.NewIoBufferBytes([]byte("hello"))
	dataStream := &MockDataStream{iobuf}
	output := &file.GetObjectOutput{Etag: "tag"}
	output.DataStream = dataStream
	mockServer.EXPECT().Context().Return(ctx)
	mockoss.EXPECT().GetObject(ctx, &l8s3.GetObjectInput{Bucket: "layotto", Key: "object"}).Return(output, nil)
	getObjectReq.StoreName = "mockoss"
	mockServer.EXPECT().Send(&s3.GetObjectOutput{Body: []byte("hello"), Etag: "tag"}).Times(1)
	err = s3Server.GetObject(getObjectReq, mockServer)
	assert.Nil(t, err)

}
