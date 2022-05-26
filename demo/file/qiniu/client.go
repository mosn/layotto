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
	"os"
	"strconv"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc"
)

const (
	storeName = "file_demo"
)

func TestGet(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.GetFileRequest{StoreName: storeName, Name: fileName}
	cli, err := c.GetFile(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		return
	}
	pic := make([]byte, 0)
	for {
		resp, err := cli.Recv()
		if err != nil {
			fmt.Printf("recv file failed")
			break
		}
		pic = append(pic, resp.Data...)
	}
	fmt.Println(string(pic))
}

func TestPut(fileName string, value string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	c := runtimev1pb.NewRuntimeClient(conn)
	data := []byte(value)
	meta["filesize"] = strconv.Itoa(len(data))
	req := &runtimev1pb.PutFileRequest{StoreName: storeName, Name: fileName, Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}
	req.Data = data
	stream.Send(req)
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
	}
}

func TestList(bucketName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	marker := ""
	for {
		req := &runtimev1pb.FileRequest{StoreName: storeName, Name: bucketName, Metadata: meta}
		listReq := &runtimev1pb.ListFileRequest{Request: req, PageSize: 1, Marker: marker}
		resp, err := c.ListFile(context.Background(), listReq)
		if err != nil {
			fmt.Printf("list file fail, err: %+v", err)
			return
		}
		marker = resp.Marker
		if !resp.IsTruncated {
			fmt.Printf("files under bucket is: %+v, %+v \n", resp.Files, marker)
			fmt.Printf("finish list \n")
			return
		}
		fmt.Printf("files under bucket is: %+v, %+v \n", resp.Files, marker)
	}

}

func TestDel(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.FileRequest{StoreName: storeName, Name: fileName, Metadata: meta}
	listReq := &runtimev1pb.DelFileRequest{Request: req}
	_, err = c.DelFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v \n", err)
		return
	}
	fmt.Printf("delete file success \n")
}

func TestStat(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.FileRequest{StoreName: storeName, Name: fileName, Metadata: meta}
	statReq := &runtimev1pb.GetFileMetaRequest{Request: req}
	data, err := c.GetFileMeta(context.Background(), statReq)
	//here use grpc error code check file exist or not.
	if m, ok := status.FromError(err); ok {
		if m.Code() == codes.NotFound {
			fmt.Println("file not exist")
			return
		}
		if m != nil {
			fmt.Printf("stat file fail,err:%+v \n", err)
			return
		}
	}
	fmt.Printf("get meta data of file: size:%+v, modifyTime:%+v \n", data.Size, data.LastModified)
	for k, v := range data.Response.Metadata {
		fmt.Printf("metadata:key:%+v,value:%+v \n", k, v)
	}

}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("you can use client like: client put/get/del/list fileName/directryName\n")
		fmt.Println("eg:")
		fmt.Println(" ./main put dir/a.txt aaa")
		fmt.Println(" ./main get dir/a.txt")
		fmt.Println(" ./main list dir/")
		fmt.Println(" ./main del dir/a.txt")
		return
	}
	if os.Args[1] == "put" {
		TestPut(os.Args[2], os.Args[3])
	}
	if os.Args[1] == "get" {
		TestGet(os.Args[2])
	}
	if os.Args[1] == "del" {
		TestDel(os.Args[2])
	}
	if os.Args[1] == "list" {
		TestList(os.Args[2])
	}
	if os.Args[1] == "stat" {
		TestStat(os.Args[2])
	}
}
