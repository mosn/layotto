package main

import (
	"context"
	"flag"
	"fmt"
	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	key = "key666"
)

var storeName string

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
}

func main() {
	flag.Parse()
	if storeName == "" {
		panic("storeName is empty.")
	}

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
			return
		}
		fmt.Printf("Next id:%v \n", id)
	}
	fmt.Println("Demo success!")
}
