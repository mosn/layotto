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

package pluggable

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"testing"

	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/protobuf/types/known/structpb"
)

type fakeSvc struct {
	onHandlerCalled func(context.Context)
}

func (f *fakeSvc) handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	f.onHandlerCalled(ctx)
	return structpb.NewNullValue(), nil
}

func TestSocketDialer(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}
	const (
		fakeSvcName    = "layotto.service.fake"
		fakeMethodName = "MyMethod"
	)

	handlerCalled := 0
	fakeSvc := &fakeSvc{
		onHandlerCalled: func(ctx context.Context) {
			handlerCalled++
		},
	}

	// start grpc socket server
	const fakeSocketFolder = "/tmp"
	uniqueID := guuid.New().String()
	socket := fmt.Sprintf("%s/%s.sock", fakeSocketFolder, uniqueID)
	lis, err := net.Listen("unix", socket)
	require.NoError(t, err)
	defer lis.Close()

	s := grpc.NewServer()
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: fakeSvcName,
		HandlerType: (*interface{})(nil),
		Methods: []grpc.MethodDesc{{
			MethodName: fakeMethodName,
			Handler:    fakeSvc.handler,
		}},
	}
	s.RegisterService(fakeDesc, fakeSvc)

	go func() {
		if serveErr := s.Serve(lis); serveErr != nil {
			t.Log(serveErr)
		}
	}()

	dialer := SocketDialer(socket)
	conn, err := dialer(context.TODO())
	assert.NoError(t, err)
	defer conn.Close()

	acceptedStatus := []connectivity.State{
		connectivity.Ready,
		connectivity.Idle,
	}

	assert.Contains(t, acceptedStatus, conn.GetState())
	require.NoError(t, conn.Invoke(context.Background(), fmt.Sprintf("/%s/%s", fakeSvcName, fakeMethodName), structpb.NewNullValue(), structpb.NewNullValue()))
	assert.Equal(t, 1, handlerCalled)
}
