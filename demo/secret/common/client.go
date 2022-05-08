package main

import (
	"context"
	"flag"
	"fmt"
	"mosn.io/layotto/sdk/go-sdk/client"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
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
	//2. get the secret
	resp, err := cli.GetSecret(ctx, &runtimev1pb.GetSecretRequest{
		StoreName: storeName,
		Key:       "redisPassword",
	})
	fmt.Println(resp)

	//3. get the bulk secret
	bulkSecrets, err := cli.GetBulkSecret(ctx, &runtimev1pb.GetBulkSecretRequest{
		StoreName: storeName,
	})

	fmt.Println(bulkSecrets)

}
