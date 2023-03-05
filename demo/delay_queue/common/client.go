package main

import (
	"context"
	"fmt"
	"log"

	pb "mosn.io/layotto/spec/proto/extension/v1/delay_queue"

	"google.golang.org/grpc"
)

func main() {
	// set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create a new client for the DelayQueue service
	client := pb.NewDelayQueueClient(conn)

	// create a new context
	ctx := context.Background()

	// create a new delay message request
	request := &pb.DelayMessageRequest{
		ComponentName:   "my_delay_queue",
		Topic:           "my_topic",
		Data:            []byte("hello, world"),
		DataContentType: "text/plain",
		DelayInSeconds:  300,
		Metadata: map[string]string{
			"key": "value",
		},
	}

	// call the PublishDelayMessage method of the DelayQueue service
	response, err := client.PublishDelayMessage(ctx, request)
	if err != nil {
		log.Fatalf("could not publish delay message: %v", err)
	}

	// print the message ID from the response
	fmt.Printf("Message ID: %s\n", response.MessageId)
}
