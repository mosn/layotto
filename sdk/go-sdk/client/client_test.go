package client

import (
	"context"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	empty "google.golang.org/protobuf/types/known/emptypb"
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
	runtimev1pb.RegisterRuntimeServer(s, &testRuntimeServer{kv: make(map[string]string), subscribed: make(map[string]bool)})

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
	runtimev1pb.UnimplementedRuntimeServer
	kv         map[string]string
	subscribed map[string]bool
}

func (t *testRuntimeServer) GetConfiguration(ctx context.Context, req *runtimev1pb.GetConfigurationRequest) (*runtimev1pb.GetConfigurationResponse, error) {
	resp := &runtimev1pb.GetConfigurationResponse{}
	for _, v := range req.Keys {
		if _, ok := t.kv[v]; ok {
			item := &runtimev1pb.ConfigurationItem{Key: v, Content: t.kv[v]}
			resp.Items = append(resp.Items, item)
		}
	}
	return resp, nil
}
func (t *testRuntimeServer) SaveConfiguration(ctx context.Context, req *runtimev1pb.SaveConfigurationRequest) (*empty.Empty, error) {
	for _, v := range req.Items {
		t.kv[v.Key] = v.Content
	}
	return &empty.Empty{}, nil
}
func (t *testRuntimeServer) DeleteConfiguration(ctx context.Context, req *runtimev1pb.DeleteConfigurationRequest) (*empty.Empty, error) {
	for _, v := range req.Keys {
		delete(t.kv, v)
	}
	return &empty.Empty{}, nil
}
func (t *testRuntimeServer) SubscribeConfiguration(srv runtimev1pb.Runtime_SubscribeConfigurationServer) error {
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	for _, key := range req.Keys {
		t.subscribed[key] = true
	}
	resp := &runtimev1pb.SubscribeConfigurationResponse{}
	for key, _ := range t.subscribed {
		item := &runtimev1pb.ConfigurationItem{Key: key, Content: "Test"}
		resp.Items = append(resp.Items, item)
	}
	err = srv.Send(resp)
	return err
}

func (*testRuntimeServer) SayHello(ctx context.Context, req *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error) {
	resp := &runtimev1pb.SayHelloResponse{Hello: "world"}
	return resp, nil
}
