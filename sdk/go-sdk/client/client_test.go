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
	"fmt"

	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/anypb"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "mosn.io/layotto/spec/proto/runtime/v1"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
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
	RegisterRecoverLogger(debugIgnoreLogger)
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
	runtimev1pb.RegisterRuntimeServer(s, &testRuntimeServer{
		kv:         make(map[string]string),
		subscribed: make(map[string]bool),
		state:      make(map[string][]byte),
		lock:       make(map[string]string),
	})

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
	state      map[string][]byte
	lock       map[string]string
}

func (t *testRuntimeServer) InvokeService(ctx context.Context, req *runtimev1pb.InvokeServiceRequest) (*runtimev1pb.InvokeResponse, error) {
	if req.Message == nil {
		return &runtimev1pb.InvokeResponse{
			ContentType: "text/plain",
			Data: &anypb.Any{
				Value: []byte("pong"),
			},
		}, nil
	}
	return &runtimev1pb.InvokeResponse{
		ContentType: req.Message.ContentType,
		Data:        req.Message.Data,
	}, nil
}

func (t *testRuntimeServer) GetState(ctx context.Context, req *runtimev1pb.GetStateRequest) (*runtimev1pb.GetStateResponse, error) {
	return &pb.GetStateResponse{
		Data: t.state[req.Key],
		Etag: "1",
	}, nil
}

func (t *testRuntimeServer) GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error) {
	items := make([]*runtimev1pb.BulkStateItem, 0)
	for _, k := range in.GetKeys() {
		if v, found := t.state[k]; found {
			item := &pb.BulkStateItem{
				Key:  k,
				Etag: "1",
				Data: v,
			}
			items = append(items, item)
		}
	}
	return &pb.GetBulkStateResponse{
		Items: items,
	}, nil
}

func (t *testRuntimeServer) SaveState(ctx context.Context, req *runtimev1pb.SaveStateRequest) (*empty.Empty, error) {
	if req == nil {
		return &empty.Empty{}, nil
	}
	for _, item := range req.States {
		if item == nil {
			continue
		}
		t.state[item.Key] = item.Value
	}
	return &empty.Empty{}, nil
}

func (t *testRuntimeServer) DeleteState(ctx context.Context, req *runtimev1pb.DeleteStateRequest) (*empty.Empty, error) {
	delete(t.state, req.Key)
	return &empty.Empty{}, nil
}

func (t *testRuntimeServer) DeleteBulkState(ctx context.Context, req *runtimev1pb.DeleteBulkStateRequest) (*empty.Empty, error) {
	for _, item := range req.States {
		delete(t.state, item.Key)
	}
	return &empty.Empty{}, nil
}

func (t *testRuntimeServer) ExecuteStateTransaction(ctx context.Context, in *runtimev1pb.ExecuteStateTransactionRequest) (*empty.Empty, error) {
	for _, op := range in.GetOperations() {
		item := op.GetRequest()
		switch opType := op.GetOperationType(); opType {
		case "upsert":
			t.state[item.Key] = item.Value
		case "delete":
			delete(t.state, item.Key)
		default:
			return &empty.Empty{}, fmt.Errorf("invalid operation type: %s", opType)
		}
	}
	return &empty.Empty{}, nil
}

func (t *testRuntimeServer) PublishEvent(ctx context.Context, req *runtimev1pb.PublishEventRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
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
	for key := range t.subscribed {
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

func (*testRuntimeServer) GetSecret(ctx context.Context, in *runtimev1pb.GetSecretRequest) (*runtimev1pb.GetSecretResponse, error) {

	secrets := make(map[string]string)
	secrets[secretKey] = secretValue
	resp := &runtimev1pb.GetSecretResponse{
		Data: secrets,
	}
	return resp, nil
}
func (*testRuntimeServer) GetBulkSecret(ctx context.Context, in *runtimev1pb.GetBulkSecretRequest) (*runtimev1pb.GetBulkSecretResponse, error) {

	secrets := make(map[string]string)
	secrets[secretKey] = secretValue
	data := make(map[string]*runtimev1pb.SecretResponse)
	data[secretKey] = &runtimev1pb.SecretResponse{
		Secrets: secrets,
	}
	resp := &runtimev1pb.GetBulkSecretResponse{
		Data: data,
	}
	return resp, nil
}

func (t *testRuntimeServer) TryLock(ctx context.Context, in *runtimev1pb.TryLockRequest) (*runtimev1pb.TryLockResponse, error) {
	if len(t.lock[in.ResourceId]) == 0 {
		t.lock[in.ResourceId] = in.LockOwner
		return &runtimev1pb.TryLockResponse{
			Success: true,
		}, nil
	}
	// lock exist
	return &runtimev1pb.TryLockResponse{
		Success: false,
	}, nil
}

func (t *testRuntimeServer) Unlock(ctx context.Context, in *runtimev1pb.UnlockRequest) (*runtimev1pb.UnlockResponse, error) {
	// LOCK_UNEXIST
	if len(t.lock[in.ResourceId]) == 0 {
		return &runtimev1pb.UnlockResponse{
			Status: pb.UnlockResponse_LOCK_UNEXIST,
		}, nil
	}
	// SUCCESS
	if t.lock[in.ResourceId] == in.LockOwner {
		delete(t.lock, in.ResourceId)
		return &runtimev1pb.UnlockResponse{
			Status: pb.UnlockResponse_SUCCESS,
		}, nil
	}
	// LOCK_BELONG_TO_OTHERS
	return &runtimev1pb.UnlockResponse{
		Status: pb.UnlockResponse_LOCK_BELONG_TO_OTHERS,
	}, nil
}
