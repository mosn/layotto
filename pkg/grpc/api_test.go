package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/spec/proto/runtime/v1"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/layotto/layotto/pkg/mock"
	"github.com/layotto/layotto/pkg/services/hello"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	maxGRPCServerUptime = 100 * time.Millisecond
	testGRPCServerPort  = 19887
)

type MockGrpcServer struct {
	err error
	req *runtime.SubscribeConfigurationRequest
	grpc.ServerStream
}

func (m *MockGrpcServer) Send(res *runtime.SubscribeConfigurationResponse) error {
	return nil
}

func (m *MockGrpcServer) Recv() (*runtime.SubscribeConfigurationRequest, error) {
	return m.req, m.err
}

type mockGRPCAPI struct {
	API
}

func (m *mockGRPCAPI) SayHello(ctx context.Context, in *runtime.SayHelloRequest) (*runtime.SayHelloResponse, error) {
	return &runtime.SayHelloResponse{}, nil
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
		resp, err := api.SayHello(context.Background(), &runtime.SayHelloRequest{
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
		_, err := api.SayHello(context.Background(), &runtime.SayHelloRequest{
			ServiceName: "no register",
		})
		if err != ErrNoInstance {
			t.Fatalf("expected got a no instance error, but got %v", err)
		}
	})

	t.Run("empty say hello", func(t *testing.T) {
		api := &api{hellos: map[string]hello.HelloService{}}
		_, err := api.SayHello(context.Background(), &runtime.SayHelloRequest{
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
		runtime.RegisterRuntimeServer(server, testAPIServer)
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
	api := NewAPI(nil, map[string]configstores.Store{"mock": mockConfigStore})
	mockConfigStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]*configstores.ConfigurationItem{
		&configstores.ConfigurationItem{Key: "sofa", Content: "sofa1"},
	}, nil).Times(1)
	res, err := api.GetConfiguration(context.Background(), &runtime.GetConfigurationRequest{StoreName: "mock", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Nil(t, err)
	assert.Equal(t, res.Items[0].Key, "sofa")
	assert.Equal(t, res.Items[0].Content, "sofa1")
	_, err = api.GetConfiguration(context.Background(), &runtime.GetConfigurationRequest{StoreName: "etcd", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")

}

func TestSaveConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI(nil, map[string]configstores.Store{"mock": mockConfigStore})
	_, err := api.SaveConfiguration(context.Background(), &runtime.SaveConfigurationRequest{StoreName: "etcd"})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
}

func TestDeleteConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI(nil, map[string]configstores.Store{"mock": mockConfigStore})
	_, err := api.DeleteConfiguration(context.Background(), &runtime.DeleteConfigurationRequest{StoreName: "etcd"})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
}

func TestSubscribeConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	//test not support store type
	grpcServer := &MockGrpcServer{req: &runtime.SubscribeConfigurationRequest{}, err: nil}
	api := NewAPI(nil, map[string]configstores.Store{"mock": mockConfigStore})
	err := api.SubscribeConfiguration(grpcServer)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "configure store [] don't support now")

	//test
	mockConfigStore.EXPECT().StopSubscribe().Return().Times(1)
	grpcServer2 := &MockGrpcServer{req: &runtime.SubscribeConfigurationRequest{}, err: errors.New("exit")}
	err = api.SubscribeConfiguration(grpcServer2)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit")

}
