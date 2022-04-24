package main

import (
	"context"
	"fmt"
	"mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func main() {
	/*
	 1. make sure exec below command before start server
	  kubectl create secret generic db-user-pass \
	   --from-literal=username=devuser \
	   --from-literal=password='S!S\*d$zDsb='
	*/
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	//2. get the secret from kubernetes
	//set the namespace of secret
	metadata := map[string]string{
		"namespace": "default",
	}
	resp, err := cli.GetSecret(ctx, &runtimev1pb.GetSecretRequest{
		StoreName: "kubernetes",
		Key:       "db-user-pass",
		Metadata:  metadata,
	})
	fmt.Println(resp)

	//3. get the bulk secret form env, kubernetes this  will output all your kubernetes secrets
	bulkSecrets, err := cli.GetBulkSecret(ctx, &runtimev1pb.GetBulkSecretRequest{
		StoreName: "kubernetes",
		Metadata:  metadata,
	})

	fmt.Println(bulkSecrets)
}
