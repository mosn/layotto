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
	"os"
	"path/filepath"
	"time"

	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"

	"mosn.io/layotto/components/pluggable"
	"mosn.io/layotto/pkg/common"
)

const (
	// the default folder to store pluggable component socket files.
	defaultSocketFolder = "/tmp/runtime/component-sockets"

	// SocketFolderEnvVar uses to set the path of folder to store pluggable component socket files replacing defaultSocketFolder
	SocketFolderEnvVar = "LAYOTTO_COMPONENTS_SOCKETS_FOLDER"
)

// GetSocketFolderPath gets the path of folder storing pluggable component socket files.
func GetSocketFolderPath() string {
	if v, ok := os.LookupEnv(SocketFolderEnvVar); ok {
		return v
	}
	return defaultSocketFolder
}

// Discover discovers pluggable component.
// At present, layotto only support register component from unix domain socket connection,
// and not compatible with windows.
func Discover() ([]pluggable.Component, error) {
	// 1. discover pluggable component
	serviceList, err := discover()
	if err != nil {
		return nil, err
	}

	// 2. callback to register factory into MosnRuntime
	res := callback(serviceList)
	return res, nil
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
	dialer pluggable.GRPCConnectionDialer
}

type grpcConnectionCloser interface {
	grpc.ClientConnInterface
	Close() error
}

// discover use grpc reflect to get services' information.
func discover() ([]grpcService, error) {
	// set grpc connection timeout to prevent block
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	serviceList, err := serviceDiscovery(func(socket string) (client reflectServiceClient, closer func(), err error) {
		conn, err := pluggable.SocketDial(
			ctx,
			socket,
		)
		if err != nil {
			return nil, nil, err
		}

		client = grpcreflect.NewClientV1Alpha(ctx, reflectpb.NewServerReflectionClient(conn))
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
				dialer:        pluggable.SocketDialer(socket, grpc.FailOnNonTempDialError(true)),
				componentName: common.RemoveExt(f.Name()),
			})
		}
	}
	return res, nil
}

// callback use callback function to register pluggable component factories into MosnRuntime
func callback(services []grpcService) []pluggable.Component {
	res := make([]pluggable.Component, 0, len(services))
	mapper := pluggable.GetServiceDiscoveryMapper()
	for _, service := range services {
		fn, ok := mapper[service.protoRef]
		if !ok {
			continue
		}

		res = append(res, fn(service.componentName, service.dialer))
	}
	return res
}
