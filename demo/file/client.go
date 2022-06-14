package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc"
)

const (
	storeName   = "file_demo"
	storageType = "Standard"
)

func TestGet(fileName string) {
	// dial
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}

	// new client
	c := runtimev1pb.NewRuntimeClient(conn)

	// getFile
	req := &runtimev1pb.GetFileRequest{StoreName: storeName, Name: fileName}
	cli, err := c.GetFile(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		panic(err)
	}
	pic := make([]byte, 0)
	for {
		resp, err := cli.Recv()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("recv file failed")
				panic(err)
			}
			break
		}
		pic = append(pic, resp.Data...)
	}
	fmt.Println("GetFile successfully. Result:")
	fmt.Println(string(pic))
}

func TestPut(fileName string, value string) {
	conn, err := grpc.Dial("127.0.0.1:11004", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}
	meta := make(map[string]string)
	meta["storageType"] = storageType
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.PutFileRequest{StoreName: storeName, Name: fileName, Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		panic(err)
	}
	req.Data = []byte(value)
	meta["length"] = strconv.Itoa(len(value))
	stream.Send(req)
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
		panic(err)
	}
}

func TestList(bucketName string) {
	conn, err := grpc.Dial("127.0.0.1:11004", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}
	meta := make(map[string]string)
	meta["storageType"] = storageType
	c := runtimev1pb.NewRuntimeClient(conn)
	marker := ""
	for {
		req := &runtimev1pb.FileRequest{StoreName: storeName, Name: bucketName, Metadata: meta}
		listReq := &runtimev1pb.ListFileRequest{Request: req, PageSize: 2, Marker: marker}
		resp, err := c.ListFile(context.Background(), listReq)
		if err != nil {
			fmt.Printf("list file fail, err: %+v", err)
			panic(err)
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
	conn, err := grpc.Dial("127.0.0.1:11004", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	meta := make(map[string]string)
	meta["storageType"] = storageType
	c := runtimev1pb.NewRuntimeClient(conn)
	req := &runtimev1pb.FileRequest{StoreName: storeName, Name: fileName, Metadata: meta}
	listReq := &runtimev1pb.DelFileRequest{Request: req}
	_, err = c.DelFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v \n", err)
		panic(err)
	}
	fmt.Printf("delete file success \n")
}

func TestStat(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:11004", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}
	meta := make(map[string]string)
	meta["storageType"] = storageType
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

func CreateBucket(bn string) {

	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "layotto"
	secretAccessKey := "layotto_secret"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := bn
	location := "us-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("you can use client like: client put/get/del/list fileName/directryName")
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
	if os.Args[1] == "bucket" {
		CreateBucket(os.Args[2])
	}
}
