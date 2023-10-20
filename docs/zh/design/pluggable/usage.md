# pluggable component 使用文档

## 编写组件

下面以 go 实现 hello 组件为例

在 `layotto/spec/proto/pluggable` 中找到对应组件的 proto 文件，生成对应实现语言的 grpc 文件。
go 语言的 pb 文件已经生成并放在了 `spec/proto/pluggable/v1` 下，用户在使用时直接引用即可。

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
	AuthToken      = "123456" // token 校验
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

1. 实现对应组件 proto 文件的 grpc 服务。
2. 启动 socket 服务，sock 文件需放置在 `/tmp/runtime/component-sockets` 下，也可以使用 `LAYOTTO_COMPONENTS_SOCKETS_FOLDER` 环境变量进行设置。
3. 注册 grpc 服务，除了注册 hello 服务外，还需要注册 reflection 服务。该服务用于 layotto 服务发现时，获取该 socket 服务具体实现了哪些 proto 文件定义的服务。
4. 启动服务

## 注册组件

填写配置文件，在对应组件下添加相关配置项，以上述的 hello 组件为例。

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

组件的 type 为 `hello-grpc-demo`，由 socket 文件的前缀名决定。

配置项与注册普通 hello 组件一致。提供 metadata 项，便于用户设置自定义配置需求。

