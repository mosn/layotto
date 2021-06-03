package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	data := flag.String("d", "hello runtime", "-d")
	flag.Parse()

	conn, err := grpc.Dial("localhost:34904", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	cli := runtimev1pb.NewRuntimeClient(conn)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	resp, err := cli.InvokeService(
		ctx,
		&runtimev1pb.InvokeServiceRequest{
			Id: "HelloService:1.0",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:      "hello",
				ContentType: "",
				Data:        &anypb.Any{Value: []byte(*data)}}},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp.Data.GetValue()))
}
