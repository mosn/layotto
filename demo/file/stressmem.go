package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"google.golang.org/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func GetFile(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.GetFileRequest{StoreName: "aliOSS", Name: "img.png"}
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
	ioutil.WriteFile("img2.png", pic, os.ModePerm)
	fmt.Printf("goroutine[%+v] finish get \n", id)
}

func PutFile(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = "Standard"
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.PutFileRequest{StoreName: "aliOSS", Name: "img.png", Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}
	fileHandle, err := os.Open("img.png")
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
		return
	}
	fmt.Printf("goroutine[%+v] finish put \n", id)
}

func main() {
	var wg sync.WaitGroup
	//Test when multi routine put big file, layotto memory cost
	for i := 0; i < 100; i++ {
		wg.Add(1)
		PutFile(&wg, i)
	}
	//Test when multi routine get file, layotto memory cost
	for {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			GetFile(&wg, i)
		}
		wg.Wait()
		fmt.Println("finish test")
		return
	}
}
