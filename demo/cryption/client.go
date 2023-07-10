/*
* Copyright 2021 Layotto Authors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"context"
	"fmt"

	"mosn.io/layotto/spec/proto/extension/v1/cryption"

	"google.golang.org/grpc"
)

const (
	storeName = "cryption_demo"
)

func TestEncrypt() []byte {
	// Connect to the object store.
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to establish connection: %+v", err)
		return nil
	}

	// Create a client for the object store.
	c := cryption.NewCryptionServiceClient(conn)

	// Create a request to get an object from the specified bucket.
	req := &cryption.EncryptRequest{ComponentName: storeName, PlainText: []byte("Hello, world")}

	// Retrieve the object using the client and the request.
	resp, err := c.Encrypt(context.Background(), req)
	if err != nil {
		fmt.Printf("failed to retrieve object: %+v", err)
		return nil
	}
	return resp.CipherText
}

func TestDecrypt(data []byte) []byte {
	// Connect to the object store.
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to establish connection: %+v", err)
		return nil
	}

	// Create a client for the object store.
	c := cryption.NewCryptionServiceClient(conn)

	// Create a request to get an object from the specified bucket.
	req := &cryption.DecryptRequest{ComponentName: storeName, CipherText: data}

	// Retrieve the object using the client and the request.
	resp, err := c.Decrypt(context.Background(), req)
	if err != nil {
		fmt.Printf("failed to Decrypt: %+v", err)
		return nil
	}
	fmt.Printf("Decrypt response: %+v", resp)
	return resp.PlainText
}

func main() {
	encyptContent := TestEncrypt()
	fmt.Printf("加密后的数据为: %+v \n", string(encyptContent))
	decryptContent := TestDecrypt(encyptContent)
	fmt.Printf("解密后的数据为: %+v \n", string(decryptContent))
}
