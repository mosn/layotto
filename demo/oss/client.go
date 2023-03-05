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
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"mosn.io/layotto/spec/proto/extension/v1/s3"

	"google.golang.org/grpc"
)

const (
	storeName = "oss_demo"
)

// TestGetObjectInput retrieves an object from an S3-compatible object storage service.
// Parameters:
// - bucket: the name of the bucket to which the object belongs.
// - fileName: the name of the object to retrieve
func TestGetObjectInput(bucket, fileName string) {
	// Connect to the object store.
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to establish connection: %+v", err)
		return
	}

	// Create a client for the object store.
	c := s3.NewObjectStorageServiceClient(conn)

	// Create a request to get an object from the specified bucket.
	req := &s3.GetObjectInput{StoreName: storeName, Bucket: bucket, Key: fileName}

	// Retrieve the object using the client and the request.
	cli, err := c.GetObject(context.Background(), req)
	if err != nil {
		fmt.Printf("failed to retrieve object: %+v", err)
		return
	}

	// Read the object data into a byte array.
	pic := make([]byte, 0)
	for {
		resp, err := cli.Recv()
		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
			}
			break
		}
		pic = append(pic, resp.Body...)
	}

	// Convert and print the byte array as a string.
	fmt.Println(string(pic))
}

// TestPutObject puts an object into an S3 bucket with the given filename and value.
func TestPutObject(bucket, fileName string, value string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())

	// Check if there's an error in connecting to the server
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	// Create an S3 object storage service client
	c := s3.NewObjectStorageServiceClient(conn)

	// Create a PutObjectInput object with the given parameters
	req := &s3.PutObjectInput{StoreName: storeName, Bucket: bucket, Key: fileName}

	// Send a stream of PutObject requests to the server
	stream, err := c.PutObject(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}

	// Set the body of the request to the value provided
	req.Body = []byte(value)

	// Send the request to the server
	stream.Send(req)

	// Close the stream and wait for the response
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
	}
}

// TestListObjects connects to a gRPC service and lists objects in a bucket
// under a specified store, by iterating through the objects in the bucket
// using markers.
func TestListObjects(bucket string) {
	// Connect to the gRPC service.
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	// Create a new ObjectStorageServiceClient.
	c := s3.NewObjectStorageServiceClient(conn)

	// Initialize the marker to start at the beginning of the list.
	marker := ""

	// Iterate through the objects in the bucket using markers.
	for {
		// Create a new ListObjectsInput request with the store name, bucket name, max keys, and marker.
		req := &s3.ListObjectsInput{StoreName: storeName, Bucket: bucket, MaxKeys: 2, Marker: marker}

		// Send the ListObjects request to the service and get the response.
		resp, err := c.ListObjects(context.Background(), req)
		if err != nil {
			fmt.Printf("list file fail, err: %+v", err)
			return
		}

		// Set the marker to the value of NextMarker in the response.
		marker = resp.NextMarker

		// If the response is not truncated, print the objects in the bucket and return.
		if !resp.IsTruncated {
			fmt.Printf("files under bucket is: %+v, %+v \n", resp.Contents, marker)
			fmt.Printf("finish list \n")
			return
		}

		// If the response is truncated, print the objects in the bucket and continue iterating.
		fmt.Printf("files under bucket is: %+v, %+v \n", resp.Contents, marker)
	}
}

func TestDeleteObject(bucket, fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.DeleteObjectInput{StoreName: storeName, Bucket: bucket, Key: fileName}
	resp, err := c.DeleteObject(context.Background(), req)
	if err != nil {
		fmt.Printf("DeleteObject fail, err: %+v \n", err)
		return
	}
	fmt.Printf("delete file success with resp: %+v \n", resp)
}

func TestDeleteObjects(bucket, fileName1, fileName2 string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req2 := &s3.DeleteObjectsInput{StoreName: storeName, Bucket: bucket, Delete: &s3.Delete{}}
	object1 := &s3.ObjectIdentifier{Key: fileName1}
	object2 := &s3.ObjectIdentifier{Key: fileName2}
	req2.Delete.Objects = append(req2.Delete.Objects, object1)
	req2.Delete.Objects = append(req2.Delete.Objects, object2)
	resp2, err := c.DeleteObjects(context.Background(), req2)
	if err != nil {
		fmt.Printf("DeleteObjects fail, err: %+v \n", err)
		return
	}
	fmt.Printf("DeleteObjects success with resp: %+v \n", resp2)
}

func TestTagging(bucket, name string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.PutObjectTaggingInput{StoreName: storeName, Bucket: bucket, Key: name, Tags: map[string]string{"Abc": "123", "Def": "456"}}
	_, err = c.PutObjectTagging(context.Background(), req)
	if err != nil {
		fmt.Printf("PutObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("PutObjectTagging success, try get tagging\n")

	req2 := &s3.GetObjectTaggingInput{StoreName: storeName, Bucket: bucket, Key: name}
	getResp, err := c.GetObjectTagging(context.Background(), req2)
	if err != nil {
		fmt.Printf("GetObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("GetObjectTagging: %+v \n", getResp.Tags)

	req3 := &s3.DeleteObjectTaggingInput{StoreName: storeName, Bucket: bucket, Key: name}
	delResp, err := c.DeleteObjectTagging(context.Background(), req3)
	if err != nil {
		fmt.Printf("DeleteObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("DeleteObjectTagging success: %+v \n", delResp.VersionId)

	req4 := &s3.GetObjectTaggingInput{StoreName: storeName, Bucket: bucket, Key: name}
	getResp4, err := c.GetObjectTagging(context.Background(), req4)
	if err != nil {
		fmt.Printf("GetObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("GetObjectTagging after delete tag: %+v \n", getResp4.Tags)
}

func TestAcl(bucket, name string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.GetObjectCannedAclInput{StoreName: storeName, Bucket: bucket, Key: name}
	resp, err := c.GetObjectCannedAcl(context.Background(), req)
	if err != nil {
		fmt.Printf("GetObjectAcl fail, err: %+v \n", err)
	} else {
		fmt.Printf("get acl success, resp: %+v\n", resp)
	}

	putRequest := &s3.PutObjectCannedAclInput{StoreName: storeName, Bucket: bucket, Key: name, Acl: "public-read-write"}
	resp2, err := c.PutObjectCannedAcl(context.Background(), putRequest)
	if err != nil {
		fmt.Printf("TestAcl fail, err: %+v \n", err)
		return
	}
	fmt.Printf("put acl public-read-write success with resp: %+v \n", resp2)

}

func TestCopyObject(bucket, name string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.CopyObjectInput{StoreName: storeName, Bucket: bucket, Key: name + ".copy", CopySource: &s3.CopySource{CopySourceBucket: bucket, CopySourceKey: name}}
	resp, err := c.CopyObject(context.Background(), req)
	if err != nil {
		fmt.Printf("CopyObject fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CopyObject success, resp: %+v\n", resp)

}

func TestPart(bucket, name string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.CreateMultipartUploadInput{StoreName: storeName, Bucket: bucket, Key: "multicopy.jpg"}
	resp, err := c.CreateMultipartUpload(context.Background(), req)
	if err != nil {
		fmt.Printf("CreateMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CreateMultipartUpload success, resp: %+v\n", resp)

	req1 := &s3.ListMultipartUploadsInput{StoreName: storeName, Bucket: bucket, MaxUploads: 1000, KeyMarker: "multicopy.jpg", UploadIdMarker: resp.UploadId}
	resp1, err := c.ListMultipartUploads(context.Background(), req1)
	if err != nil {
		fmt.Printf("ListMultipartUploads fail, err: %+v \n", err)
		return
	}
	fmt.Printf("ListMultipartUploads success, resp: %+v \n", resp1)

	req2 := &s3.UploadPartCopyInput{StoreName: storeName, Bucket: bucket, PartNumber: 1, UploadId: resp.UploadId, Key: "multicopy.jpg", StartPosition: 0, PartSize: 1000, CopySource: &s3.CopySource{CopySourceBucket: bucket, CopySourceKey: name}}
	resp2, err := c.UploadPartCopy(context.Background(), req2)
	if err != nil {
		fmt.Printf("UploadPartCopy fail, err: %+v \n", err)
		return
	}
	fmt.Printf("UploadPartCopy success, resp: %+v \n", resp2)

	req3 := &s3.CompleteMultipartUploadInput{StoreName: storeName, Bucket: bucket, Key: "multicopy.jpg", UploadId: resp.UploadId, MultipartUpload: &s3.CompletedMultipartUpload{Parts: []*s3.CompletedPart{{Etag: resp2.CopyPartResult.Etag, PartNumber: req2.PartNumber}}}}
	resp3, err := c.CompleteMultipartUpload(context.Background(), req3)
	if err != nil {
		fmt.Printf("CompleteMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CompleteMultipartUpload success, resp: %+v \n", resp3)

	//req4 := &s3.AbortMultipartUploadInput{StoreName: storeName, Bucket: bucket, Key: "海贼王.jpeg", UploadId: "EEE5317D0EB841AC9B80D0B6A2F811AA"}
	//resp4, err := c.AbortMultipartUpload(context.Background(), req4)
	//if err != nil {
	//	fmt.Printf("AbortMultipartUpload fail, err: %+v \n", err)
	//	return
	//}
	//fmt.Printf("AbortMultipartUpload success, resp: %+v \n", resp4)

	req5 := &s3.CreateMultipartUploadInput{StoreName: storeName, Bucket: bucket, Key: "海贼王.jpg"}
	resp5, err := c.CreateMultipartUpload(context.Background(), req5)
	if err != nil {
		fmt.Printf("CreateMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CreateMultipartUpload success, resp: %+v\n", resp5)

	req6 := &s3.UploadPartInput{
		StoreName:  storeName,
		Bucket:     bucket,
		Key:        "海贼王.jpg",
		UploadId:   resp5.UploadId,
		PartNumber: 0,
	}
	f, err := os.Open("海贼王.jpg")
	if err != nil {
		fmt.Printf("open file fail, err: %+v\n", err)
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var parts []*s3.CompletedPart
	for {
		dataByte := make([]byte, 120*1024)
		var n int
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}
		req6.Body = dataByte[:n]
		req6.ContentLength = int64(n)
		req6.PartNumber = req6.PartNumber + 1
		stream, err := c.UploadPart(context.TODO())
		if err != nil {
			fmt.Printf("UploadPart fail, err: %+v \n", err)
			return
		}
		err = stream.Send(req6)
		if err != nil {
			fmt.Printf("UploadPart send fail, err: %+v \n", err)
			return
		}
		resp6, err := stream.CloseAndRecv()
		if err != nil {
			fmt.Printf("UploadPart CloseAndRecv fail, err: %+v \n", err)
			return
		}
		part := &s3.CompletedPart{Etag: resp6.Etag, PartNumber: req6.PartNumber}
		parts = append(parts, part)
	}
	fmt.Printf("UploadPart success, parts: %+v \n", parts)
	req8 := &s3.ListPartsInput{StoreName: storeName, Bucket: bucket, Key: "海贼王.jpg", UploadId: resp5.UploadId}
	resp8, err := c.ListParts(context.Background(), req8)
	if err != nil {
		fmt.Printf("ListPartsInput fail, err: %+v \n", err)
	} else {
		fmt.Printf("ListPartsInput success, resp: %+v \n", resp8)
	}
	req7 := &s3.CompleteMultipartUploadInput{StoreName: storeName, Bucket: bucket, Key: "海贼王.jpg", UploadId: resp5.UploadId, MultipartUpload: &s3.CompletedMultipartUpload{Parts: parts}}
	resp7, err := c.CompleteMultipartUpload(context.Background(), req7)
	if err != nil {
		fmt.Printf("CompleteMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CompleteMultipartUpload success, resp: %+v \n", resp7)
}

func TestListVersion(bucket string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.ListObjectVersionsInput{StoreName: storeName, Bucket: bucket}
	resp, err := c.ListObjectVersions(context.Background(), req)
	if err != nil {
		fmt.Printf("ListObjectVersions fail, err: %+v \n", err)
		return
	}
	fmt.Printf("ListObjectVersions success, resp: %+v\n", resp)

}

func TestRestore(bucket, name string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.RestoreObjectInput{StoreName: storeName, Bucket: bucket, Key: name, RestoreRequest: &s3.RestoreRequest{Days: 1}}
	resp, err := c.RestoreObject(context.Background(), req)
	if err != nil {
		fmt.Printf("RestoreObject fail, err: %+v \n", err)
		return
	}
	fmt.Printf("RestoreObject success, resp: %+v\n", resp)

}

func TestObjectExist(bucket, name string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.IsObjectExistInput{StoreName: storeName, Bucket: bucket, Key: name}
	resp, err := c.IsObjectExist(context.Background(), req)
	if err != nil {
		fmt.Printf("TestObjectExist fail, err: %+v \n", err)
		return
	}
	fmt.Printf("TestObjectExist success, resp: %+v\n", resp.FileExist)

}

func TestAbortMultipartUpload(bucket string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)

	req := &s3.ListMultipartUploadsInput{StoreName: storeName, Bucket: bucket, MaxUploads: 1000}
	resp, err := c.ListMultipartUploads(context.Background(), req)
	if err != nil {
		fmt.Printf("ListMultipartUploads fail, err: %+v \n", err)
		return
	}
	fmt.Printf("ListMultipartUploads success, resp: %+v \n", resp)

	for _, v := range resp.Uploads {
		req4 := &s3.AbortMultipartUploadInput{StoreName: storeName, Bucket: bucket, Key: v.Key, UploadId: v.UploadId}
		resp4, err := c.AbortMultipartUpload(context.Background(), req4)
		if err != nil {
			fmt.Printf("AbortMultipartUpload fail, err: %+v \n", err)
			return
		}
		fmt.Printf("AbortMultipartUpload success, resp: %+v \n", resp4)
	}

	fmt.Printf("AbortMultipartUpload test success")
}

func TestSign(bucket, name, method string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.SignURLInput{StoreName: storeName, Bucket: bucket, Key: name, ExpiredInSec: int64(10), Method: method}
	resp, err := c.SignURL(context.Background(), req)
	if err != nil {
		fmt.Printf("SignURLInput fail, err: %+v \n", err)
		return
	}
	fmt.Printf("SignURLInput success, resp: %+v\n", resp.SignedUrl)

}

func TestAppend(bucket, fileName, data, position string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	ps, _ := strconv.Atoi(position)
	req := &s3.AppendObjectInput{StoreName: storeName, Bucket: bucket, Key: fileName, Body: []byte(data), Position: int64(ps)}
	stream, err := c.AppendObject(context.Background())
	if err != nil {
		fmt.Printf("AppendObject fail,err:%+v", err)
		return
	}
	err = stream.Send(req)
	if err != nil {
		fmt.Printf("AppendObject fail,err:%+v", err)
		return
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("AppendObject fail,err:%+v", err)
		return
	}
	fmt.Printf("AppendObject success,resp: %+v \n", resp)
}

func TestHeadObject(bucket, fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.HeadObjectInput{StoreName: storeName, Bucket: bucket, Key: fileName}
	resp, err := c.HeadObject(context.Background(), req)
	if err != nil {
		fmt.Printf("HeadObjectInput fail,err:%+v", err)
		return
	}

	fmt.Printf("HeadObjectInput success,resp: %+v \n", resp)
}

func main() {

	if os.Args[1] == "put" {
		TestPutObject(os.Args[2], os.Args[3], os.Args[4])
	}
	if os.Args[1] == "get" {
		TestGetObjectInput(os.Args[2], os.Args[3])
	}
	if os.Args[1] == "del" {
		TestDeleteObject(os.Args[2], os.Args[3])
	}
	if os.Args[1] == "dels" {
		TestDeleteObjects(os.Args[2], os.Args[3], os.Args[4])
	}

	if os.Args[1] == "list" {
		TestListObjects(os.Args[2])
	}

	if os.Args[1] == "tag" {
		TestTagging(os.Args[2], os.Args[3])
	}

	if os.Args[1] == "acl" {
		TestAcl(os.Args[2], os.Args[3])
	}

	if os.Args[1] == "copy" {
		TestCopyObject(os.Args[2], os.Args[3])
	}

	if os.Args[1] == "part" {
		TestPart(os.Args[2], os.Args[3])
	}

	if os.Args[1] == "version" {
		TestListVersion(os.Args[2])
	}

	if os.Args[1] == "restore" {
		TestRestore(os.Args[2], os.Args[3])
	}
	if os.Args[1] == "exist" {
		TestObjectExist(os.Args[2], os.Args[3])
	}
	if os.Args[1] == "abort" {
		TestAbortMultipartUpload(os.Args[2])
	}

	if os.Args[1] == "sign" {
		TestSign(os.Args[2], os.Args[3], os.Args[4])
	}

	if os.Args[1] == "append" {
		TestAppend(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	}

	if os.Args[1] == "head" {
		TestHeadObject(os.Args[2], os.Args[3])
	}
}
