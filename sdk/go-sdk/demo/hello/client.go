package main

import (
	"context"
	"fmt"
	"log"

	"mosn.io/layotto/sdk/go-sdk/client"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	res, err := cli.SayHello(context.Background(), &client.SayHelloRequest{
		ServiceName: "helloworld",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Hello)
}
