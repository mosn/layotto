# pluggable Component Usage Guide

## Complete component

let's take the example of implementing the hello component in Go

在 `layotto/spec/proto/pluggable` 中找到对应组件的 proto 文件，生成对应实现语言的 grpc 文件。
go 语言的 pb 文件已经生成并放在了 `spec/proto/pluggable/v1` 下，用户在使用时直接引用即可。

Find the proto file for the corresponding component in `layotto/spec/proto/pluggable` and generate the grpc files for the corresponding implementation language.

The pb files for the Go language have already been generated and are located in `spec/proto/pluggable/v1`. Users can directly reference them when using.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "mosn.io/layotto/spec/proto/pluggable/v1/hello"
	"net"
	"os"
)

const (
	AuthToken      = "123456"  
	TokenConfigKey = "token"
	SocketFilePath = "/tmp/runtime/component-sockets/hello-grpc-demo.sock"
)

type HelloService struct {
	pb.UnimplementedHelloServer
	hello string
	token string
}

func (h *HelloService) Init(ctx context.Context, config *pb.HelloConfig) (*empty.Empty, error) {
	h.hello = config.GetHelloString()
	h.token = config.Metadata[TokenConfigKey]
	if h.token != AuthToken {
		return nil, errors.New("auth failed")
	}

	return nil, nil
}

func (h *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	res := &pb.HelloResponse{
		HelloString: h.hello,
	}
	return res, nil
}

func main() {
	listen, err := net.Listen("unix", SocketFilePath)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(SocketFilePath)

	server := grpc.NewServer()
	srv := &HelloService{}
	pb.RegisterHelloServer(server, srv)
	reflection.Register(server)

	fmt.Println("start grpc server")
	if err := server.Serve(listen); err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Println(err)
	}
}
```

1. Implement the gRPC service for the corresponding component's proto file.
2. Start the socket service. The sock file should be placed under `/tmp/runtime/component-sockets`, or you can also set the `LAYOTTO_COMPONENTS_SOCKETS_FOLDER` environment variable for configuration.
3. Register the gRPC service. In addition to registering the "hello" service, it is also necessary to register the reflection service. This service is used for Layotto service discovery to determine which services defined in the proto files are implemented by this socket service. 。
5. Start service and wait for layotto registering it.

## Component Register

Fill in the configuration file and add the relevant configuration items under the corresponding component. Taking the "hello" component mentioned above as an example.

```json
"grpc_config": {
  "hellos": {
    "helloworld": {
      "type": "hello-grpc-demo",
      "hello": "hello",
      "metadata": {
        "token": "123456"
      }
    }
  }
}
```

The component's type is `hello-grpc-demo`, determined by the prefix name of the socket file.

The configuration items are the same as registering a regular "hello" component. Provide the `metadata` field to allow users to set custom configuration requirements.

