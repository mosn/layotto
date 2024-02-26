// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"google.golang.org/grpc"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

var (
	storeName string
	FileName  = "file"
	FileValue = "value"
)

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
}

func main() {
	flag.Parse()
	if storeName == "" {
		panic("storeName is empty.")
	}

	// conn to layotto grpc server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	c := runtimev1pb.NewRuntimeClient(conn)

	// Set metadata, the specific configuration needs to see what configuration is required
	// to put the relevant source code in the local file.

	// Create metadata and a new client
	meta := make(map[string]string)
	meta["FileMode"] = "0777"
	meta["FileFlag"] = strconv.Itoa(os.O_CREATE | os.O_RDWR)

	// Make a request to put a file
	// storeName is the name of file instance.
	req := &runtimev1pb.PutFileRequest{StoreName: storeName, Name: FileName, Metadata: meta}
	stream, err := c.PutFile(context.Background())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		panic(err)
	}
	req.Data = []byte(FileValue)
	_ = stream.Send(req)
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
		panic(err)
	}

	// Get File
	file, err := c.GetFile(context.Background(), &runtimev1pb.GetFileRequest{StoreName: storeName, Name: FileName, Metadata: nil})
	if err != nil {
		fmt.Printf("get file failed: %+v", err)
		panic(err)
	}

	// Receive data from the server and store in pic
	pic := make([]byte, 0)
	for {
		resp, err := file.Recv()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("recv file failed")
				panic(err)
			}
			break
		}
		pic = append(pic, resp.Data...)
	}

	if string(pic) != FileValue {
		fmt.Printf("the file is not the same as the value we put before. %s", string(pic))
		return
	}

	fmt.Println("test file operate success")
}
