package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// 1. invoke lifecycle API to modify the component configuration
	c := runtimev1pb.NewLifecycleClient(conn)
	metaData := make(map[string]string)
	metaData["hello"] = "goodbye"
	req := &runtimev1pb.DynamicConfiguration{ComponentConfig: &runtimev1pb.ComponentConfig{
		Kind:     "hellos",
		Name:     "quick_start_demo",
		Metadata: metaData,
	}}
	_, err = c.ApplyConfiguration(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 2. invoke sayHello API again
	client := runtimev1pb.NewRuntimeClient(conn)
	hello, err := client.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
		ServiceName: "quick_start_demo",
		Name:        "eva",
	})
	if err != nil {
		panic(err)
	}
	if hello.Hello != "goodbye, eva" {
		panic(errors.New(fmt.Sprintf("Assertion failed! Result is %v", hello.Hello)))
	}
	fmt.Println(hello.Hello)
}
