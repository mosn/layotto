package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.InvokeBindingRequest{Name: "http", Operation: "get"}
	resp, err := c.InvokeBinding(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		return
	}
	fmt.Println(resp)
}
