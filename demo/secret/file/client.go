package main

import (
	"context"
	"fmt"
	"mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	//1. get secret from local file
	resp, err := cli.GetSecret(ctx, &runtimev1pb.GetSecretRequest{
		StoreName: "local.file",
		Key:       "redisPassword",
	})

	//2. get bulk secret form local file
	fmt.Println(resp.Data)
	bulkSecrets, err := cli.GetBulkSecret(ctx, &runtimev1pb.GetBulkSecretRequest{
		StoreName: "local.file",
	})

	fmt.Println(bulkSecrets.Data)

}
