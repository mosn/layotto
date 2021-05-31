package main

import (
	"context"
	"github.com/layotto/layotto/sdk/client"
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
	return cli.PublishEvent(context.Background(), in)
}
