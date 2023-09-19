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
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"mosn.io/layotto/components/hello"

	"mosn.io/layotto/pkg/common"
	grpcdial "mosn.io/layotto/pkg/grpc"

	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func init() {
	onServiceDiscovered = make(map[string]CallbackFunc)
}

const (
	// the default folder to store pluggable component socket files.
	defaultSocketFolder = "/tmp/runtime/component-sockets"

	// SocketFolderEnvVar uses to set the path of folder to store pluggable component socket files replacing defaultSocketFolder
	SocketFolderEnvVar = "LAYOTTO_COMPONENTS_SOCKETS_FOLDER"
)

var (
	// register the callback function of pluggable component, using grpc connection creates factory and registers into MosnRuntime
	onServiceDiscovered map[string]CallbackFunc
)

type CallbackFunc func(compType string, dialer grpcdial.GRPCConnectionDialer, m *DiscoverFactory)

type DiscoverFactory struct {
	Hellos []*hello.HelloFactory
}

// GetSocketFolderPath gets the path of folder storing pluggable component socket files.
func GetSocketFolderPath() string {
	if v, ok := os.LookupEnv(SocketFolderEnvVar); ok {
		return v
	}
	return defaultSocketFolder
}

// AddServiceDiscoveryCallback register callback function, not concurrent secure
func AddServiceDiscoveryCallback(serviceName string, callback CallbackFunc) {
	onServiceDiscovered[serviceName] = callback
}

func NewDefaultDiscoverFactory() *DiscoverFactory {
	factory := new(DiscoverFactory)
	factory.Hellos = make([]*hello.HelloFactory, 0)
	return factory
}

func Discover() (*DiscoverFactory, error) {
	// At present, layotto only support register component from unix domain socket connection,
	// and not compatible with windows.
	if runtime.GOOS == "windows" {
		return nil, errors.New("windows haven't support register pluggable component")
	}

	// 1. discover pluggable component
	factory := NewDefaultDiscoverFactory()
	serviceList, err := discover()
	if err != nil {
		return nil, err
	}

	// 2. callback to register factory into MosnRuntime
	callback(serviceList, factory)
	return factory, nil
}

// get service form socket files.
type reflectServiceClient interface {
	ListServices() ([]string, error)
	Reset()
}

type grpcService struct {
	// protoRef is the proto service name
	protoRef string
	// componentName is the component name that implements such service.
	componentName string
	// dialer is the used grpc connectiondialer.
	dialer grpcdial.GRPCConnectionDialer
}

type grpcConnectionCloser interface {
	grpc.ClientConnInterface
	Close() error
}

// discover use grpc reflect to get services' information.
func discover() ([]grpcService, error) {
	ctx := context.TODO()
	serviceList, err := serviceDiscovery(func(socket string) (client reflectServiceClient, closer func(), err error) {
		conn, err := grpcdial.SocketDial(
			ctx,
			socket,
			grpc.WithBlock(),
		)
		if err != nil {
			return nil, nil, err
		}

		client = grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(conn))
		return client, reflectServiceConnectionCloser(conn, client), nil
	})
	if err != nil {
		return nil, err
	}

	return serviceList, nil
}

// reflectServiceConnectionCloser is used for cleanup the stream created to be used for the reflection service.
func reflectServiceConnectionCloser(conn grpcConnectionCloser, client reflectServiceClient) func() {
	return func() {
		client.Reset()
		conn.Close()
	}
}

// discover service socket files and get service information from factory factor function.
func serviceDiscovery(reflectClientFactory func(socket string) (client reflectServiceClient, cleanup func(), err error)) ([]grpcService, error) {
	res := make([]grpcService, 0)

	// 1. get socket folder files
	path := GetSocketFolderPath()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return res, nil
	}
	if err != nil {
		return nil, err
	}

	// 2. read socket files
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range files {
		if dirEntry.IsDir() { // skip dirs
			continue
		}

		f, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		if isSocket := common.IsSocket(f); !isSocket { // check is socket files
			continue
		}

		// 3. using reflectClientFactory gets service information
		socket := filepath.Join(path, f.Name())
		client, cleanup, err := reflectClientFactory(socket)
		if err != nil {
			return nil, err
		}
		defer cleanup()

		serviceList, err := client.ListServices()
		if err != nil {
			return nil, err
		}

		for _, s := range serviceList {
			res = append(res, grpcService{
				protoRef:      s,
				dialer:        grpcdial.SocketDialer(socket, grpc.WithBlock(), grpc.FailOnNonTempDialError(true)),
				componentName: common.RemoveExt(f.Name()),
			})
		}
	}
	return res, nil
}

// callback use callback function to register pluggable component factories into MosnRuntime
func callback(services []grpcService, factory *DiscoverFactory) {
	for _, service := range services {
		callback, ok := onServiceDiscovered[service.protoRef]
		if !ok {
			continue
		}

		callback(service.componentName, service.dialer, factory)
	}
}
