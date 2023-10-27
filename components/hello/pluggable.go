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

package hello

import (
	"context"
	"fmt"
	"time"

	"mosn.io/layotto/components/pluggable"
	helloproto "mosn.io/layotto/spec/proto/pluggable/v1/hello"
)

func init() {
	// spec.proto.pluggable.v1.Hello
	pluggable.AddServiceDiscoveryCallback(helloproto.Hello_ServiceDesc.ServiceName, func(compType string, dialer pluggable.GRPCConnectionDialer) pluggable.Component {
		return NewHelloFactory(compType, func() HelloService {
			return NewGRPCHello(dialer)
		})
	})
}

type grpcHello struct {
	dialer pluggable.GRPCConnectionDialer
	client helloproto.HelloClient
}

func NewGRPCHello(dialer pluggable.GRPCConnectionDialer) HelloService {
	return &grpcHello{dialer: dialer}
}

// todo 优雅关闭时关闭 conn

func (g *grpcHello) Init(config *HelloConfig) error {
	// 1.dial grpc server
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	conn, err := g.dialer(ctx)
	if err != nil {
		return fmt.Errorf("dial hello pluggable component: %w", err)
	}

	// 2.init pluggable component
	g.client = helloproto.NewHelloClient(conn)
	if _, err := g.client.Init(ctx, &helloproto.HelloConfig{
		Config:      pluggable.ToProtoConfig(config.Config),
		Type:        config.Type,
		HelloString: config.HelloString,
		Metadata:    config.Metadata,
	}); err != nil {
		return fmt.Errorf("init hello pluggable component: %w", err)
	}

	return nil
}

func (g *grpcHello) Hello(ctx context.Context, request *HelloRequest) (*HelloResponse, error) {
	resp, err := g.client.SayHello(ctx, &helloproto.HelloRequest{
		Name: request.Name,
	})
	if err != nil {
		return nil, err
	}

	res := &HelloResponse{
		HelloString: resp.GetHelloString(),
	}
	return res, nil
}
