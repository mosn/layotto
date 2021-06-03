package client

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	runtimev1pb "github.com/layotto/layotto/proto/runtime/v1"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	runtimePortDefault    = "34904"
	runtimePortEnvVarName = "runtime_GRPC_PORT"
)

var (
	logger               = log.New(os.Stdout, "", 0)
	_             Client = (*GRPCClient)(nil)
	defaultClient Client
	doOnce        sync.Once
)

// Client is the interface for runtime client implementation.

type Client interface {
	SayHello(ctx context.Context, in *SayHelloRequest) (*SayHelloResp, error)
	GetConfiguration(ctx context.Context, in *ConfigurationRequestItem) ([]*ConfigurationItem, error)
	// SaveConfiguration saves configuration into configuration store.
	SaveConfiguration(ctx context.Context, in *SaveConfigurationRequest) error
	// DeleteConfiguration deletes configuration from configuration store.
	DeleteConfiguration(ctx context.Context, in *ConfigurationRequestItem) error
	// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
	SubscribeConfiguration(ctx context.Context, in *ConfigurationRequestItem) WatchChan
	// Close cleans up all resources created by the client.
	Close()
}

// NewClient instantiates runtime client using runtime_GRPC_PORT environment variable as port.
// Note, this default factory function creates runtime client only once. All subsequent invocations
// will return the already created instance. To create multiple instances of the runtime client,
// use one of the parameterized factory functions:
//   NewClientWithPort(port string) (client Client, err error)
//   NewClientWithAddress(address string) (client Client, err error)
//   NewClientWithConnection(conn *grpc.ClientConn) Client
func NewClient() (client Client, err error) {
	port := os.Getenv(runtimePortEnvVarName)
	if port == "" {
		port = runtimePortDefault
	}
	var onceErr error
	doOnce.Do(func() {
		c, err := NewClientWithPort(port)
		onceErr = errors.Wrap(err, "error creating default client")
		defaultClient = c
	})

	return defaultClient, onceErr
}

// NewClientWithPort instantiates runtime using specific port.
func NewClientWithPort(port string) (client Client, err error) {
	if port == "" {
		return nil, errors.New("nil port")
	}
	return NewClientWithAddress(net.JoinHostPort("127.0.0.1", port))
}

// NewClientWithAddress instantiates runtime using specific address (including port).
func NewClientWithAddress(address string) (client Client, err error) {
	if address == "" {
		return nil, errors.New("nil address")
	}
	logger.Printf("runtime client initializing for: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "error creating connection to '%s': %v", address, err)
	}

	return NewClientWithConnection(conn), nil
}

// NewClientWithConnection instantiates runtime client using specific connection.
func NewClientWithConnection(conn *grpc.ClientConn) Client {
	return &GRPCClient{
		connection:  conn,
		protoClient: runtimev1pb.NewMosnRuntimeClient(conn),
	}
}

// GRPCClient is the gRPC implementation of runtime client.
type GRPCClient struct {
	connection  *grpc.ClientConn
	protoClient runtimev1pb.MosnRuntimeClient
	mux         sync.Mutex
}

// Close cleans up all resources created by the client.
func (c *GRPCClient) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}
