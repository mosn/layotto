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

package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mosn.io/layotto/components/file"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/pkg/mock"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	maxGRPCServerUptime = 100 * time.Millisecond
	testGRPCServerPort  = 19887
)

type MockGrpcServer struct {
	err error
	req *runtimev1pb.SubscribeConfigurationRequest
	grpc.ServerStream
}

func (m *MockGrpcServer) Send(res *runtimev1pb.SubscribeConfigurationResponse) error {
	return nil
}

func (m *MockGrpcServer) Recv() (*runtimev1pb.SubscribeConfigurationRequest, error) {
	return m.req, m.err
}

type mockGRPCAPI struct {
	API
}

func (m *mockGRPCAPI) SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error) {
	return &runtimev1pb.SayHelloResponse{}, nil
}

func TestStartServerAPI(t *testing.T) {
	port := testGRPCServerPort
	server := startTestRuntimeAPIServer(port, &mockGRPCAPI{})
	defer server.Stop()
}

func TestSayHello(t *testing.T) {
	t.Run("request ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockHello := mock.NewMockHelloService(ctrl)
		api := &api{hellos: map[string]hello.HelloService{
			"mock": mockHello,
		}}
		mockHello.EXPECT().Hello(gomock.Any()).Return(&hello.HelloReponse{
			HelloString: "mock hello",
		}, nil).Times(1)
		resp, err := api.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "mock",
		})
		if err != nil {
			t.Fatalf("say hello request failed: %v", err)
		}
		if resp.Hello != "mock hello" {
			t.Fatalf("say hello response is not expected: %v", resp)
		}
	})

	t.Run("no hello stored", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockHello := mock.NewMockHelloService(ctrl)
		api := &api{hellos: map[string]hello.HelloService{
			"mock": mockHello,
		}}
		_, err := api.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "no register",
		})
		if err != ErrNoInstance {
			t.Fatalf("expected got a no instance error, but got %v", err)
		}
	})

	t.Run("empty say hello", func(t *testing.T) {
		api := &api{hellos: map[string]hello.HelloService{}}
		_, err := api.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "mock",
		})
		if err != ErrNoInstance {
			t.Fatalf("expected got a no instance error, but got %v", err)
		}
	})
}

func startTestRuntimeAPIServer(port int, testAPIServer API) *grpc.Server {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	opts := []grpc.ServerOption{grpc.WriteBufferSize(1)}

	server := grpc.NewServer(opts...)
	go func() {
		runtimev1pb.RegisterRuntimeServer(server, testAPIServer)
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	time.Sleep(maxGRPCServerUptime)

	return server
}

func TestGetConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil)
	mockConfigStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]*configstores.ConfigurationItem{
		&configstores.ConfigurationItem{Key: "sofa", Content: "sofa1"},
	}, nil).Times(1)
	res, err := api.GetConfiguration(context.Background(), &runtimev1pb.GetConfigurationRequest{StoreName: "mock", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Nil(t, err)
	assert.Equal(t, res.Items[0].Key, "sofa")
	assert.Equal(t, res.Items[0].Content, "sofa1")
	_, err = api.GetConfiguration(context.Background(), &runtimev1pb.GetConfigurationRequest{StoreName: "etcd", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")

}

func TestSaveConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil)
	_, err := api.SaveConfiguration(context.Background(), &runtimev1pb.SaveConfigurationRequest{StoreName: "etcd"})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
}

func TestDeleteConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil)
	_, err := api.DeleteConfiguration(context.Background(), &runtimev1pb.DeleteConfigurationRequest{StoreName: "etcd"})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
}

func TestSubscribeConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil)

	//test not support store type
	grpcServer := &MockGrpcServer{req: &runtimev1pb.SubscribeConfigurationRequest{}, err: nil}
	err := api.SubscribeConfiguration(grpcServer)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "configure store [] don't support now")

	//test
	grpcServer2 := &MockGrpcServer{req: &runtimev1pb.SubscribeConfigurationRequest{}, err: errors.New("exit")}
	err = api.SubscribeConfiguration(grpcServer2)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit")

}

func SendData(w net.Conn) {
	w.Write([]byte("testFile"))
	w.Close()
}
func TestGetFile(t *testing.T) {
	r, w := net.Pipe()
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	mockStream := mock.NewMockRuntime_GetFileServer(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil)
	err := api.GetFile(&runtimev1pb.GetFileRequest{StoreName: "mock1"}, mockStream)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not supported store type: mock1"))
	mockFile.EXPECT().Get(&file.GetFileStu{ObjectName: "", Metadata: nil}).Return(r, nil).Times(1)
	mockStream.EXPECT().Send(&runtimev1pb.GetFileResponse{Data: []byte("testFile")}).Times(1)
	go SendData(w)
	api.GetFile(&runtimev1pb.GetFileRequest{StoreName: "mock"}, mockStream)
}

func putFile(t *testing.T, api API, wg *sync.WaitGroup, mockStream runtimev1pb.Runtime_PutFileServer) {
	err := api.PutFile(mockStream)
	assert.Nil(t, err)
	wg.Done()
}
func TestPutFile(t *testing.T) {
	var wg sync.WaitGroup
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	mockStream := mock.NewMockRuntime_PutFileServer(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil)

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	err := api.PutFile(mockStream)
	assert.Nil(t, err)

	mockStream.EXPECT().Recv().Return(&runtimev1pb.PutFileRequest{StoreName: "mock1"}, nil).Times(1)
	err = api.PutFile(mockStream)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))

	mockStream.EXPECT().Recv().Return(&runtimev1pb.PutFileRequest{StoreName: "mock", Name: "fileName", Data: []byte("fileContent")}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	mockStream.EXPECT().SendAndClose(&emptypb.Empty{}).Times(1)
	mockFile.EXPECT().CompletePut(int64(3)).Return(nil).Times(1)
	mockFile.EXPECT().Put(&file.PutFileStu{FileName: "fileName", Data: []byte("fileContent"), Metadata: nil, StreamId: 3, ChunkNumber: 1}).Return(nil).Times(1)
	wg.Add(1)
	go putFile(t, api, &wg, mockStream)
	wg.Wait()
}

func TestListFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil)
	request := &runtimev1pb.FileRequest{StoreName: "mock1"}
	resp, err := api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, resp)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))
	request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	mockFile.EXPECT().List(&file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{FilesName: []string{"file1", "file2"}}, nil).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, err)
	assert.Equal(t, resp.FileName[0], "file1")
	assert.Equal(t, resp.FileName[1], "file2")
	mockFile.EXPECT().List(&file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{FilesName: []string{"file1", "file2"}}, errors.New("test fail")).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.NotNil(t, err)
}

func TestDelFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil)
	request := &runtimev1pb.FileRequest{StoreName: "mock1"}
	resp, err := api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.Nil(t, resp)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))
	request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	mockFile.EXPECT().Del(&file.DelRequest{FileName: request.Name, Metadata: request.Metadata}).Return(nil).Times(1)
	_, err = api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.Nil(t, err)
	mockFile.EXPECT().Del(&file.DelRequest{FileName: request.Name, Metadata: request.Metadata}).Return(errors.New("test fail")).Times(1)
	_, err = api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.NotNil(t, err)
}
