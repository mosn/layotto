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
	// 1. invoke sayHello API
	client := runtimev1pb.NewRuntimeClient(conn)
	hello, err := client.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
		ServiceName: "helloworld",
		Name:        "eva",
	})
	if err != nil {
		panic(err)
	}
	// validate the response value
	if hello.Hello != "greeting, eva" {
		panic(fmt.Errorf("Assertion failed! Result is %v", hello.Hello))
	}
	fmt.Println(hello.Hello)

	// 2. invoke lifecycle API to modify the `sayHello` component configuration
	c := runtimev1pb.NewLifecycleClient(conn)
	metaData := make(map[string]string)
	// change the configuration from "greeting" to "goodbye"
	metaData["hello"] = "goodbye"
	req := &runtimev1pb.DynamicConfiguration{ComponentConfig: &runtimev1pb.ComponentConfig{
		Kind:     "hellos",
		Name:     "helloworld",
		Metadata: metaData,
	}}
	_, err = c.ApplyConfiguration(context.Background(), req)
	if err != nil {
		panic(err)
	}
	// 3. invoke sayHello API again
	client = runtimev1pb.NewRuntimeClient(conn)
	hello, err = client.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
		ServiceName: "helloworld",
		Name:        "eva",
	})
	if err != nil {
		panic(err)
	}
	// validate the response value. It should be different this time.
	if hello.Hello != "goodbye, eva" {
		panic(fmt.Errorf("Assertion failed! Result is %v", hello.Hello))
	}
	fmt.Println(hello.Hello)
}
