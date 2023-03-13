package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	// Dial to the gRPC server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}

	// Create a new client
	c := runtimev1pb.NewRuntimeClient(conn)

	// Make a request to get a file
	req := &runtimev1pb.GetFileRequest{StoreName: storeName, Name: fileName}
	cli, err := c.GetFile(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		panic(err)
	}

	// Receive data from the server and store in pic
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

	// Print the result of the GetFile request
	fmt.Println("GetFile successfully. Result:")
	fmt.Println(string(pic))
}

func TestPut(fileName string, value string) {
	// Dial to the gRPC server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}

	// Create metadata and a new client
	meta := make(map[string]string)
	meta["storageType"] = storageType
	c := runtimev1pb.NewRuntimeClient(conn)

	// Make a request to put a file
	req := &runtimev1pb.PutFileRequest{StoreName: storeName, Name: fileName, Metadata: meta}
	stream, err := c.PutFile(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		panic(err)
	}
	req.Data = []byte(value)
	stream.Send(req)
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
		panic(err)
	}
}

func TestList(bucketName string) {
	// Dial to the gRPC server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}

	// Create metadata and a new client
	meta := make(map[string]string)
	meta["storageType"] = storageType
	c := runtimev1pb.NewRuntimeClient(conn)

	// Make a request to list files
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

// TestDel deletes a file with the given fileName from the server
func TestDel(fileName string) {
	// Dial a connection to the server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	// Define metadata to be sent with the request
	meta := make(map[string]string)
	meta["storageType"] = storageType

	// Create a new runtime client using the connection
	c := runtimev1pb.NewRuntimeClient(conn)

	// Create a file request with the given storeName, fileName and metadata
	req := &runtimev1pb.FileRequest{StoreName: storeName, Name: fileName, Metadata: meta}

	// Create a delete file request using the file request created above
	listReq := &runtimev1pb.DelFileRequest{Request: req}

	// Send the delete file request to the server and check for errors
	_, err = c.DelFile(context.Background(), listReq)
	if err != nil {
		fmt.Printf("list file fail, err: %+v \n", err)
		panic(err)
	}

	// If successful, print a message indicating success
	fmt.Printf("delete file success \n")
}

// TestStat retrieves metadata for the file with the given fileName from the server
func TestStat(fileName string) {
	// Dial a connection to the server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		panic(err)
	}

	// Define metadata to be sent with the request
	meta := make(map[string]string)
	meta["storageType"] = storageType

	// Create a new runtime client using the connection
	c := runtimev1pb.NewRuntimeClient(conn)

	// Create a file request with the given storeName, fileName and metadata
	req := &runtimev1pb.FileRequest{StoreName: storeName, Name: fileName, Metadata: meta}

	// Create a get file metadata request using the file request created above
	statReq := &runtimev1pb.GetFileMetaRequest{Request: req}

	// Send the get file metadata request to the server and check for errors
	data, err := c.GetFileMeta(context.Background(), statReq)

	// Check if the error returned is a "not found" error and print a message if so
	if m, ok := status.FromError(err); ok {
		if m.Code() == codes.NotFound {
			fmt.Println("file not exist")
			return
		}

		// If it's not a "not found" error and not nil, print an error message
		if m != nil {
			fmt.Printf("stat file fail,err:%+v \n", err)
			return
		}
	}

	// If successful, print metadata for the file
	fmt.Printf("get meta data of file: size:%+v, modifyTime:%+v \n", data.Size, data.LastModified)
	for k, v := range data.Response.Metadata {
		fmt.Printf("metadata:key:%+v,value:%+v \n", k, v)
	}

}

// CreateBucket function creates a new bucket in the specified S3-compatible object storage service.
// Parameters:
//   - bn (string): the name of the new bucket to be created
func CreateBucket(bn string) {
	// Set the connection parameters for the S3-compatible object storage service.
	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "layotto"
	secretAccessKey := "layotto_secret"
	useSSL := false

	// Initialize minio client object with the connection parameters.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket with the specified name and location.
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
