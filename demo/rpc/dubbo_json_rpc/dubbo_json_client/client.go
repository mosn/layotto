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
	data := flag.String("d", `{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}`, "-d")
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
			Id: "com.ikurento.user.UserProvider",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:        "GetUser",
				ContentType:   "",
				Data:          &anypb.Any{Value: []byte(*data)},
				HttpExtension: &runtimev1pb.HTTPExtension{Verb: runtimev1pb.HTTPExtension_POST},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp.Data.GetValue()))
}
