package main

import (
	"context"
	"fmt"
	client "mosn.io/layotto/sdk/go-sdk/client"
)

const topicName = "topic1"

func main() {
	// 1. construct client
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	// 2. publish a new event
	testPublish(cli)
	cli.Close()
}

func testPublish(cli client.Client) error {
	in := &client.PublishEventRequest{
		PubsubName:      "redis",
		Topic:           topicName,
		Data:            []byte("value1"),
		DataContentType: "",
		Metadata:        nil,
	}
	err := cli.PublishEvent(context.Background(), in)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published a new event.Topic: %s ,Data: %s \n", in.Topic, in.Data)
	return err
}
