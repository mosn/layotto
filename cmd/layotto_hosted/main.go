/*
 * Copyright 2026 Layotto Authors
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

// An independent Go host for a real Layotto API and a native TCP service.
// The TCP echo service demonstrates native hosting, not ZDAS/ZCache protocols.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/hello/helloworld"
	lgrpc "mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/pkg/grpc/default_api"
	"mosn.io/layotto/pkg/runtime/hosted"
	"mosn.io/layotto/pkg/runtime/hosted/grpcserver"
	"mosn.io/layotto/pkg/runtime/hosted/tcpserver"
)

// This adapter owns the component; the API borrows the same HelloService.
type helloModule struct{ service hello.HelloService }

func (m *helloModule) Init(ctx context.Context) error {
	return m.service.Init(&hello.HelloConfig{HelloString: "hello"})
}
func (m *helloModule) Start(context.Context) error { return nil }
func (m *helloModule) Status(context.Context) (hosted.Health, error) {
	return hosted.Health{Ready: true}, nil
}
func (m *helloModule) Drain(context.Context) error { return nil }
func (m *helloModule) Close(context.Context) error { return nil } // HelloWorld owns no external resources.

func assemble(mode, apiAddress, nativeAddress string) (*hosted.Runtime, *grpcserver.Module, *tcpserver.Module, error) {
	if mode != "api" && mode != "native" && mode != "both" {
		return nil, nil, nil, fmt.Errorf("unknown mode %q (use api, native or both)", mode)
	}
	var registrations []hosted.Registration
	var api *grpcserver.Module
	var native *tcpserver.Module
	if mode != "native" {
		var component *helloModule
		registrations = append(registrations, hosted.Registration{Name: "hello", New: func(context.Context) (hosted.Module, error) {
			component = &helloModule{service: helloworld.NewHelloWorld()}
			return component, nil
		}})
		api = grpcserver.New(grpcserver.Config{Address: apiAddress, Register: func(ctx context.Context, s *grpc.Server) error {
			ac := &lgrpc.ApplicationContext{AppId: "hosted-demo", Hellos: map[string]hello.HelloService{"hello": component.service}}
			service := default_api.NewGrpcAPI(ac)
			if err := service.Init(nil); err != nil {
				return err
			}
			return service.Register(s)
		}})
		registrations = append(registrations, hosted.Registration{Name: "api", DependsOn: []string{"hello"}, New: func(context.Context) (hosted.Module, error) { return api, nil }})
	}
	if mode != "api" {
		native = tcpserver.New(nativeAddress, func(ctx context.Context, c net.Conn) { _, _ = io.Copy(c, c) })
		registrations = append(registrations, hosted.Registration{Name: "native", New: func(context.Context) (hosted.Module, error) { return native, nil }})
	}
	r, err := hosted.New(registrations, hosted.Options{DrainTimeout: 3 * time.Second, CloseTimeout: 2 * time.Second})
	return r, api, native, err
}

func run() error {
	mode := flag.String("mode", "both", "api, native or both")
	apiAddress := flag.String("api-address", "127.0.0.1:34904", "Layotto API address")
	nativeAddress := flag.String("native-address", "127.0.0.1:16379", "native TCP echo address (not Redis)")
	flag.Parse()
	r, api, native, err := assemble(*mode, *apiAddress, *nativeAddress)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	startup, cancelStartup := context.WithTimeout(ctx, 10*time.Second)
	err = r.Start(startup)
	cancelStartup()
	if err != nil {
		return err
	}
	if api != nil {
		log.Printf("Layotto API: %s", api.Address())
	}
	if native != nil {
		log.Printf("native TCP echo: %s", native.Address())
	}
	<-ctx.Done()
	shutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	return r.Shutdown(shutdown)
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
