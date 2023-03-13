package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

// This is the entry point of the program
func main() {
	// Establishes a connection with a gRPC server at the specified address
	// The communication between the client and server is not encrypted
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		// If the connection cannot be established, panic and stop the program
		panic(err)
	}

	// Creates a new instance of a client that can be used to call the gRPC methods provided by the server
	c := runtimev1pb.NewRuntimeClient(conn)

	// Creates a new map and adds a key-value pair to it
	// This is used to add metadata information to the gRPC request
	metaData := make(map[string]string)
	metaData["token"] = "123"

	// Creates a new InvokeBindingRequest message with the necessary data and metadata fields set
	// The request includes the name of the binding, the operation to be performed, the metadata map created earlier, and the authentication data
	req := &runtimev1pb.InvokeBindingRequest{Name: "bindings_demo", Operation: "get", Metadata: metaData, Data: []byte("auth data")}

	// Calls the InvokeBinding() method on the client with the InvokeBindingRequest message as a parameter
	// This sends the request to the server and returns a response
	resp, err := c.InvokeBinding(context.Background(), req)
	if err != nil {
		// If there is an error, print it to the console
		fmt.Printf("get file error: %+v", err)
		return
	}

	// Prints the response from the server to the console
	fmt.Println(resp)
}
