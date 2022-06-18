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
		panic(err)
	}

	c := runtimev1pb.NewRuntimeClient(conn)
	metaData := make(map[string]string)
	metaData["token"] = "123"
	req := &runtimev1pb.InvokeBindingRequest{Name: "bindings_demo", Operation: "get", Metadata: metaData, Data: []byte("auth data")}
	resp, err := c.InvokeBinding(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		return
	}
	fmt.Println(resp)
}
