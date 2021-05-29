package main

import (
	"context"
	"fmt"
	"github.com/layotto/layotto/spec/proto/runtime/v1"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := runtime.NewRuntimeClient(conn)

	for i := 0; i < 10; i++ {
		r, err := c.SayHello(context.Background(), &runtime.SayHelloRequest{
			ServiceName: "helloworld",
		})
		if err != nil {
			fmt.Println("get an error: ", err)
		} else {
			fmt.Println("get a message: ", r.GetHello())
		}
	}

}
