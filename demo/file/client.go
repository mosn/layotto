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

func TestGet() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.GetFileRequest{StoreName: "aliOSS", Name: "海贼王.jpeg"}
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
	ioutil.WriteFile("海贼王.jpeg", pic, os.ModePerm)
}

func TestPut() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.PutFileRequest{StoreName: "aliOSS", Name: "海贼王.jpeg", Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}
	fileHandle, err := os.Open("海贼王.jpeg")
	defer fileHandle.Close()
	//分片上传，片最小为100kb
	buffer := make([]byte, 102400)

	for {
		n, err := fileHandle.Read(buffer)
		// 控制条件,根据实际调整
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

func main() {
	TestGet()
	//TestPut()
}
