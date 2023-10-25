package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "mosn.io/layotto/spec/proto/pluggable/v1/hello"
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

	return &empty.Empty{}, nil
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
