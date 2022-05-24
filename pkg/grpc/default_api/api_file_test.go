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

package default_api

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mosn.io/layotto/components/file"
	"mosn.io/layotto/pkg/mock"
	"mosn.io/layotto/pkg/mock/runtime"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func SendData(w net.Conn) {
	w.Write([]byte("testFile"))
	w.Close()
}
func TestGetFile(t *testing.T) {
	r, w := net.Pipe()
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	mockStream := runtime.NewMockRuntime_GetFileServer(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	err := api.GetFile(&runtimev1pb.GetFileRequest{StoreName: "mock1"}, mockStream)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not supported store type: mock1"))
	metadata := make(map[string]string)
	mockFile.EXPECT().Get(context.Background(), &file.GetFileStu{FileName: "", Metadata: metadata}).Return(r, nil).Times(1)
	mockStream.EXPECT().Send(&runtimev1pb.GetFileResponse{Data: []byte("testFile")}).Times(1)
	mockStream.EXPECT().Context().Return(context.Background())
	go SendData(w)
	err = api.GetFile(&runtimev1pb.GetFileRequest{StoreName: "mock"}, mockStream)
	assert.Nil(t, err)
}

func TestPutFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	mockStream := runtime.NewMockRuntime_PutFileServer(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	err := api.PutFile(mockStream)
	assert.Nil(t, err)

	mockStream.EXPECT().Recv().Return(&runtimev1pb.PutFileRequest{StoreName: "mock1"}, nil).Times(1)
	err = api.PutFile(mockStream)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))

	mockStream.EXPECT().Recv().Return(&runtimev1pb.PutFileRequest{StoreName: "mock"}, nil).Times(1)
	stream := newPutObjectStreamReader(nil, mockStream)
	Metadata := make(map[string]string)
	mockStream.EXPECT().Context().Return(context.Background())
	mockFile.EXPECT().Put(context.Background(), &file.PutFileStu{DataStream: stream, FileName: "", Metadata: Metadata}).Return(errors.New("err occur")).Times(1)
	err = api.PutFile(mockStream)
	s, _ := status.FromError(err)
	assert.Equal(t, s.Message(), "err occur")
}

func TestListFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	request := &runtimev1pb.FileRequest{StoreName: "mock1"}
	request.Metadata = make(map[string]string)
	resp, err := api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, resp)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))
	request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	request.Metadata = make(map[string]string)
	mockFile.EXPECT().List(context.Background(), &file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{Files: nil, Marker: "hello", IsTruncated: true}, nil).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, err)
	assert.Equal(t, resp.Marker, "hello")
	assert.Equal(t, resp.IsTruncated, true)
	mockFile.EXPECT().List(context.Background(), &file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{}, errors.New("test fail")).Times(1)
	_, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.NotNil(t, err)
	info := &file.FilesInfo{FileName: "hello", Size: 10, LastModified: "2021.11.12"}
	files := make([]*file.FilesInfo, 0)
	files = append(files, info)
	mockFile.EXPECT().List(context.Background(), &file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{Files: files}, nil).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, err)
	assert.Equal(t, len(resp.Files), 1)
	assert.Equal(t, resp.Files[0].FileName, "hello")
}

func TestDelFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	request := &runtimev1pb.FileRequest{StoreName: "mock1"}
	request.Metadata = make(map[string]string)
	resp, err := api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.Nil(t, resp)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))
	request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	request.Metadata = make(map[string]string)
	mockFile.EXPECT().Del(context.Background(), &file.DelRequest{FileName: request.Name, Metadata: request.Metadata}).Return(nil).Times(1)
	_, err = api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.Nil(t, err)
	mockFile.EXPECT().Del(context.Background(), &file.DelRequest{FileName: request.Name, Metadata: request.Metadata}).Return(errors.New("test fail")).Times(1)
	_, err = api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.NotNil(t, err)
}

func TestGetFileMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	request := &runtimev1pb.GetFileMetaRequest{Request: nil}
	resp, err := api.GetFileMeta(context.Background(), request)
	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, st.Message(), "request can't be nil")
	request.Request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	meta := make(map[string]string)
	re := &file.FileMetaResp{
		Size:         10,
		LastModified: "123",
		Metadata: map[string][]string{
			"test": {},
		},
	}
	mockFile.EXPECT().Stat(context.Background(), &file.FileMetaRequest{FileName: request.Request.Name, Metadata: meta}).Return(re, nil).Times(1)
	resp, err = api.GetFileMeta(context.Background(), request)
	assert.Nil(t, err)
	assert.Equal(t, resp.LastModified, "123")
	assert.Equal(t, int(resp.Size), 10)
}
