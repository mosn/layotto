// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcserver

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The outer interceptors cover handlers and middleware cleanup. Admission and
// completion use one mutex, so no late handler can race an empty-group wait.
type calls struct {
	mu      sync.Mutex
	closing bool
	active  int
	done    chan struct{}
}

func (c *calls) enter() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return false
	}
	c.active++
	return true
}
func (c *calls) leave() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.active--
	if c.closing && c.active == 0 {
		close(c.done)
	}
}
func (c *calls) stop() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closing {
		c.closing = true
		c.done = make(chan struct{})
		if c.active == 0 {
			close(c.done)
		}
	}
	return c.done
}
func (c *calls) unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !c.enter() {
		return nil, status.Error(codes.Unavailable, "server is stopping")
	}
	defer c.leave()
	return handler(ctx, req)
}
func (c *calls) stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if !c.enter() {
		return status.Error(codes.Unavailable, "server is stopping")
	}
	defer c.leave()
	return handler(srv, stream)
}

func newServer(options []grpc.ServerOption, c *calls) (server *grpc.Server, err error) {
	// gRPC prepends single interceptors to chained interceptors. Reserve those
	// slots to ensure user middleware defers finish before our active count drops.
	// Translate only the known conflicting-option panics into configuration errors.
	defer func() {
		if v := recover(); v != nil {
			if message, ok := v.(string); ok && (message == "The unary server interceptor was already set and may not be reset." ||
				message == "The stream server interceptor was already set and may not be reset.") {
				err = fmt.Errorf("hosted gRPC: use ChainUnaryInterceptor/ChainStreamInterceptor: %v", v)
				return
			}
			panic(v)
		}
	}()
	opts := []grpc.ServerOption{grpc.UnaryInterceptor(c.unary), grpc.StreamInterceptor(c.stream)}
	opts = append(opts, options...)
	return grpc.NewServer(opts...), nil
}
