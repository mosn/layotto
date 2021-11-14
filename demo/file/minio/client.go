package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc"
)

func TestGet(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.GetFileRequest{StoreName: "minioOSS", Name: fileName}
	cli, err := c.GetFile(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		return
	}
	pic := make([]byte, 0, 0)
	for {
		resp, err := cli.Recv()
		if err != nil {
			fmt.Errorf("recv file failed")
			break
		}
		pic = append(pic, resp.Data...)
	}
	ioutil.WriteFile(fileName, pic, os.ModePerm)
}

func TestPut(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	// meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.PutFileRequest{StoreName: "minioOSS", Name: fileName, Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}
	fileHandle, err := os.Open(fileName)
	if err != nil {
		fmt.Println("open file fail")
		return
	}
	defer fileHandle.Close()
	//Upload in multiples, the minimum size is 100kb
	buffer := make([]byte, 102400)

	for {
		n, err := fileHandle.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("read file failed, err:%+v", err)
			break
		}
		if n == 0 {
			//stream.CloseSend()
			break
		}
		req.Data = buffer[:n]
		err = stream.Send(req)
		if err != nil {
			fmt.Printf("send request failed: err: %+v", err)
			break
		}
	}
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
	req := &runtimev1pb.FileRequest{StoreName: "minioOSS", Name: bucketName, Metadata: meta}
	listReq := &runtimev1pb.ListFileRequest{Request: req}
	resp, err := c.ListFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v", err)
		return
	}
	for _, file := range resp.Files {
		fmt.Printf("files under bucket is: %+v", file.FileName)
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
	req := &runtimev1pb.FileRequest{StoreName: "minioOSS", Name: fileName, Metadata: meta}
	listReq := &runtimev1pb.DelFileRequest{Request: req}
	_, err = c.DelFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v", err)
		return
	}
	fmt.Printf("delete file success")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("you can use client like: client put/get/del/list fileName/directryName")
		return
	}
	if os.Args[1] == "put" {
		TestPut(os.Args[2])
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
}
