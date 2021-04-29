package client

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"gitlab.alipay-inc.com/ant-mesh/runtime/proto/runtime/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"os"
	"testing"
)

const (
	testBufSize = 1024 * 1024
)

var (
	testClient Client
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	c, f := getTestClient(ctx)
	testClient = c
	r := m.Run()
	f()
	os.Exit(r)
}

func TestNewClient(t *testing.T) {
	t.Run("no arg for with port", func(t *testing.T) {
		_, err := NewClientWithPort("")
		assert.Error(t, err)
	})

	t.Run("no arg for with address", func(t *testing.T) {
		_, err := NewClientWithAddress("")
		assert.Error(t, err)
	})

	t.Run("new client closed with empty token", func(t *testing.T) {
		c, err := NewClient()
		assert.NoError(t, err)
		defer c.Close()
	})

	t.Run("new client with trace ID", func(t *testing.T) {
		c, err := NewClient()
		assert.NoError(t, err)
		defer c.Close()
	})
}

func getTestClient(ctx context.Context) (client Client, closer func()) {
	s := grpc.NewServer()
	runtime.RegisterMosnRuntimeServer(s, &testRuntimeServer{kv: make(map[string]string), subscribed: make(map[string]bool)})

	l := bufconn.Listen(testBufSize)
	go func() {
		if err := s.Serve(l); err != nil && err.Error() != "closed" {
			logger.Fatalf("test server exited with error: %v", err)
		}
	}()

	d := grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return l.Dial()
	})

	c, err := grpc.DialContext(ctx, "", d, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("failed to dial test context: %v", err)
	}

	closer = func() {
		l.Close()
		s.Stop()
	}

	client = NewClientWithConnection(c)
	return
}

type testRuntimeServer struct {
	runtime.UnimplementedMosnRuntimeServer
	kv         map[string]string
	subscribed map[string]bool
}

func (t *testRuntimeServer) GetConfiguration(ctx context.Context, req *runtime.GetConfigurationRequest) (*runtime.GetConfigurationResponse, error) {
	resp := &runtime.GetConfigurationResponse{}
	for _, v := range req.Keys {
		if _, ok := t.kv[v]; ok {
			item := &runtime.ConfigurationItem{Key: v, Content: t.kv[v]}
			resp.Items = append(resp.Items, item)
		}
	}
	return resp, nil
}
func (t *testRuntimeServer) SaveConfiguration(ctx context.Context, req *runtime.SaveConfigurationRequest) (*empty.Empty, error) {
	for _, v := range req.Items {
		t.kv[v.Key] = v.Content
	}
	return &empty.Empty{}, nil
}
func (t *testRuntimeServer) DeleteConfiguration(ctx context.Context, req *runtime.DeleteConfigurationRequest) (*empty.Empty, error) {
	for _, v := range req.Keys {
		delete(t.kv, v)
	}
	return &empty.Empty{}, nil
}
func (t *testRuntimeServer) SubscribeConfiguration(srv runtime.MosnRuntime_SubscribeConfigurationServer) error {
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	for _, key := range req.Keys {
		t.subscribed[key] = true
	}
	resp := &runtime.SubscribeConfigurationResponse{}
	for key, _ := range t.subscribed {
		item := &runtime.ConfigurationItem{Key: key, Content: "Test"}
		resp.Items = append(resp.Items, item)
	}
	err = srv.Send(resp)
	return err
}

func (*testRuntimeServer) SayHello(ctx context.Context, req *runtime.SayHelloRequest) (*runtime.SayHelloResponse, error) {
	resp := &runtime.SayHelloResponse{Hello: "world"}
	return resp, nil
}
