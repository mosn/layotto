package main

import (
	"context"
	"fmt"
	"mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func main() {

	// 1. make sure exec `export redisPassword=envPassword`  before start server
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	//2. get the secret from env
	resp, err := cli.GetSecret(ctx, &runtimev1pb.GetSecretRequest{
		StoreName: "local.env",
		Key:       "redisPassword",
	})
	fmt.Println(resp)

	//3. get the bulk secret form env, maybe this  will output all your environment variable
	bulkSecrets, err := cli.GetBulkSecret(ctx, &runtimev1pb.GetBulkSecretRequest{
		StoreName: "local.env",
	})

	fmt.Println(bulkSecrets)

}
