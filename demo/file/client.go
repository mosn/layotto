package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

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
	req := &runtimev1pb.GetFileRequest{StoreName: "localStore", Name: fileName}
	cli, err := c.GetFile(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		return
	}
	pic := make([]byte, 0, 0)
	for {
		resp, err := cli.Recv()
		if err != nil{
			break
		}
		pic = append(pic, resp.Data...)
	}
	fmt.Println("get file sucess, content is: ", string(pic))
}

func TestPut(fileName string,data []byte) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	mode := 0777
	meta["fileMode"] = strconv.Itoa(mode)
	meta["fileFlag"] = strconv.Itoa(os.O_RDWR | os.O_CREATE)
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.PutFileRequest{StoreName: "localStore", Name: fileName, Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}
	req.Data = data
	fmt.Println(string(req.Data))
	err = stream.Send(req)
	if err != nil{
		fmt.Println("send failed")
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
	}
	fmt.Println("finish file put")
}

func TestList(FileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.FileRequest{StoreName: "localStore", Name: FileName, Metadata: meta}
	listReq := &runtimev1pb.ListFileRequest{Request: req}
	resp, err := c.ListFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v", err)
		return
	}
	fmt.Printf("files under directory is: %+v \n", resp.FileName)
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
	req := &runtimev1pb.FileRequest{StoreName: "localStore", Name: fileName, Metadata: meta}
	listReq := &runtimev1pb.DelFileRequest{Request: req}
	_, err = c.DelFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v", err)
		return
	}
	fmt.Println("delete file success")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("no enought arguments")
		return
	}
	if os.Args[1] == "put" {
		TestPut(os.Args[2], []byte(os.Args[3]))
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
