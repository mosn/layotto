package main

import (
	"context"
	"fmt"

	"mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

var (
	storeName = "sequencer-demo"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		fmt.Println("connect to layotto failed")
		panic(err)
	}

	defer cli.Close()

	var id *runtimev1pb.GetNextIdResponse
	for i := 0; i < 10; i++ {
		id, err = cli.GetNextId(context.Background(), &runtimev1pb.GetNextIdRequest{
			StoreName: storeName,
			Key:       "key",
		})

		if err != nil {
			fmt.Println("get next id failed")
			panic(err)
		}
	}

	fmt.Println(id)
}
