package grpc

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const (
	// needed to load balance requests for target services with multiple endpoints, ie. multiple instances
	grpcServiceConfig = `{"loadBalancingPolicy":"round_robin"}`
	dialTimeout       = time.Second * 30
)

// Manager is a wrapper around gRPC connection pooling
type Manager struct {
	AppClientConn  *grpc.ClientConn
	lock           *sync.RWMutex
	connectionPool map[string]*grpc.ClientConn
}

// NewManager returns a new grpc manager
func NewManager() *Manager {
	return &Manager{
		lock:           &sync.RWMutex{},
		connectionPool: map[string]*grpc.ClientConn{},
	}
}

func (m *Manager) InitAppClient(port int) error {
	conn, err := m.GetGRPCConnection(fmt.Sprintf("127.0.0.1:%v", port), "", "", true, false, false)
	if err != nil {
		return errors.Errorf("error establishing connection to app grpc on port %v: %s", port, err)
	}

	m.AppClientConn = conn
	return nil
}

// GetGRPCConnection returns a new grpc connection for a given address and inits one if doesn't exist
func (m *Manager) GetGRPCConnection(address, id string, namespace string, skipTLS, recreateIfExists, sslEnabled bool) (*grpc.ClientConn, error) {
	// 1. read pool to check if exists.
	m.lock.RLock()
	if val, ok := m.connectionPool[address]; ok && !recreateIfExists {
		m.lock.RUnlock()
		return val, nil
	}
	m.lock.RUnlock()

	// 2. create connection
	m.lock.Lock()
	defer m.lock.Unlock()
	// read the value once again, as a concurrent writer could create it
	if val, ok := m.connectionPool[address]; ok && !recreateIfExists {
		return val, nil
	}

	opts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(grpcServiceConfig),
	}
	// TODO support TLS
	opts = append(opts, grpc.WithInsecure())
	// dial
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	if c, ok := m.connectionPool[address]; ok {
		c.Close()
	}

	m.connectionPool[address] = conn

	return conn, nil
}
