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

package client

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
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

type runtimeAPI interface {
	SayHello(ctx context.Context, in *SayHelloRequest) (*SayHelloResp, error)

	GetConfiguration(ctx context.Context, in *ConfigurationRequestItem) ([]*ConfigurationItem, error)

	// InvokeMethod invokes service without raw data
	InvokeMethod(ctx context.Context, appID, methodName, verb string) (out []byte, err error)

	// InvokeMethodWithContent invokes service with content
	InvokeMethodWithContent(ctx context.Context, appID, methodName, verb string, content *DataContent) (out []byte, err error)

	// InvokeMethodWithCustomContent invokes app with custom content (struct + content type).
	InvokeMethodWithCustomContent(ctx context.Context, appID, methodName, verb string, contentType string, content interface{}) (out []byte, err error)

	// PublishEvent publishes data onto topic in specific pubsub component.
	PublishEvent(ctx context.Context, pubsubName, topicName string, data []byte) error

	// PublishEventfromCustomContent serializes an struct and publishes its contents as data (JSON) onto topic in specific pubsub component.
	PublishEventfromCustomContent(ctx context.Context, pubsubName, topicName string, data interface{}) error

	// SaveConfiguration saves configuration into configuration store.
	SaveConfiguration(ctx context.Context, in *SaveConfigurationRequest) error

	// DeleteConfiguration deletes configuration from configuration store.
	DeleteConfiguration(ctx context.Context, in *ConfigurationRequestItem) error

	// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
	SubscribeConfiguration(ctx context.Context, in *ConfigurationRequestItem) WatchChan

	// SaveState saves the raw data into store using default state options.
	SaveState(ctx context.Context, storeName, key string, data []byte, so ...StateOption) error

	// SaveBulkState saves multiple state item to store with specified options.
	SaveBulkState(ctx context.Context, storeName string, items ...*SetStateItem) error

	// GetState retrieves state from specific store using default consistency option.
	GetState(ctx context.Context, storeName, key string) (item *StateItem, err error)

	// GetStateWithConsistency retrieves state from specific store using provided state consistency.
	GetStateWithConsistency(ctx context.Context, storeName, key string, meta map[string]string, sc StateConsistency) (item *StateItem, err error)

	// GetBulkState retrieves state for multiple keys from specific store.
	GetBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string, parallelism int32) ([]*BulkStateItem, error)

	// DeleteState deletes content from store using default state options.
	DeleteState(ctx context.Context, storeName, key string) error

	// DeleteStateWithETag deletes content from store using provided state options and etag.
	DeleteStateWithETag(ctx context.Context, storeName, key string, etag *ETag, meta map[string]string, opts *StateOptions) error

	// ExecuteStateTransaction provides way to execute multiple operations on a specified store.
	ExecuteStateTransaction(ctx context.Context, storeName string, meta map[string]string, ops []*StateOperation) error

	// DeleteBulkState deletes content for multiple keys from store.
	DeleteBulkState(ctx context.Context, storeName string, keys []string) error

	// DeleteBulkState deletes content for multiple keys from store.
	DeleteBulkStateItems(ctx context.Context, storeName string, items []*DeleteStateItem) error

	// Distributed Lock API
	TryLock(context.Context, *runtimev1pb.TryLockRequest) (*runtimev1pb.TryLockResponse, error)
	Unlock(context.Context, *runtimev1pb.UnlockRequest) (*runtimev1pb.UnlockResponse, error)

	// Sequencer API
	// Get next unique id with some auto-increment guarantee
	GetNextId(ctx context.Context, in *runtimev1pb.GetNextIdRequest) (*runtimev1pb.GetNextIdResponse, error)

	// Secret API
	GetSecret(ctx context.Context, in *runtimev1pb.GetSecretRequest, opts ...grpc.CallOption) (*runtimev1pb.GetSecretResponse, error)
	GetBulkSecret(ctx context.Context, in *runtimev1pb.GetBulkSecretRequest, opts ...grpc.CallOption) (*runtimev1pb.GetBulkSecretResponse, error)

	// Close cleans up all resources created by the client.
	Close()
}

// NewClient instantiates runtime client using runtime_GRPC_PORT environment variable as port.
// Note, this default factory function creates runtime client only once. All subsequent invocations
// will return the already created instance. To create multiple instances of the runtime client,
// use one of the parameterized factory functions:
//
// NewClientWithPort(port string) (client Client, err error)
// NewClientWithAddress(address string) (client Client, err error)
// NewClientWithConnection(conn *grpc.ClientConn) Client
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

// Close cleans up all resources created by the client.
func (c *GRPCClient) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}
