package main

import (
	"context"
	"fmt"
	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	key       = "key666"
	storeName = "redis"
)

func main() {

	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	fmt.Printf("Try to get next id.Key:%s \n", key)
	for i := 0; i < 10; i++ {
		id, err := cli.GetNextId(ctx, &runtimev1pb.GetNextIdRequest{
			StoreName: storeName,
			Key:       key,
			Options:   nil,
			Metadata:  nil,
		})
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Printf("Next id:%v \n", id)
	}
	fmt.Println("Demo success!")
}
