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

	s3 "mosn.io/layotto/spec/proto/extension/v1"

	"google.golang.org/grpc"
)

const (
	storeName = "oss_demo"
)

func TestGetObjectInput(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}

	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.GetObjectInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: fileName}
	cli, err := c.GetObject(context.Background(), req)
	if err != nil {
		fmt.Printf("get file error: %+v", err)
		return
	}
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
	fmt.Println(string(pic))
}

func TestPutObject(fileName string, value string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.PutObjectInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: fileName}
	stream, err := c.PutObject(context.TODO())
	if err != nil {
		fmt.Printf("put file failed:%+v", err)
		return
	}
	req.Body = []byte(value)
	stream.Send(req)
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("cannot receive response: %+v", err)
	}
}

func TestListObjects() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	marker := ""
	for {
		req := &s3.ListObjectsInput{StoreName: storeName, Bucket: "antsys-wenxuwan", MaxKeys: 2, Marker: marker}
		resp, err := c.ListObjects(context.Background(), req)
		if err != nil {
			fmt.Printf("list file fail, err: %+v", err)
			return
		}
		marker = resp.NextMarker
		if !resp.IsTruncated {
			fmt.Printf("files under bucket is: %+v, %+v \n", resp.Contents, marker)
			fmt.Printf("finish list \n")
			return
		}
		fmt.Printf("files under bucket is: %+v, %+v \n", resp.Contents, marker)
	}

}

func TestDeleteObject(fileName string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.DeleteObjectInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: fileName}
	resp, err := c.DeleteObject(context.Background(), req)
	if err != nil {
		fmt.Printf("DeleteObject fail, err: %+v \n", err)
		return
	}
	fmt.Printf("delete file success with resp: %+v \n", resp)
}

func TestDeleteObjects(fileName1, fileName2 string) {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req2 := &s3.DeleteObjectsInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Delete: &s3.Delete{}}
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

func TestTagging() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.PutObjectTaggingInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg", Tags: map[string]string{"Abc": "123", "Def": "456"}}
	_, err = c.PutObjectTagging(context.Background(), req)
	if err != nil {
		fmt.Printf("PutObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("PutObjectTagging success, try get tagging\n")

	req2 := &s3.GetObjectTaggingInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg"}
	getResp, err := c.GetObjectTagging(context.Background(), req2)
	if err != nil {
		fmt.Printf("GetObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("GetObjectTagging: %+v \n", getResp.Tags)

	req3 := &s3.DeleteObjectTaggingInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg"}
	delResp, err := c.DeleteObjectTagging(context.Background(), req3)
	if err != nil {
		fmt.Printf("DeleteObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("DeleteObjectTagging success: %+v \n", delResp.VersionId)

	req4 := &s3.GetObjectTaggingInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg"}
	getResp4, err := c.GetObjectTagging(context.Background(), req4)
	if err != nil {
		fmt.Printf("GetObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("GetObjectTagging after delete tag: %+v \n", getResp4.Tags)
}

func TestAcl() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.GetObjectAclInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg"}
	resp, err := c.GetObjectAcl(context.Background(), req)
	if err != nil {
		fmt.Printf("GetObjectAcl fail, err: %+v \n", err)
		return
	}
	fmt.Printf("get acl success, resp: %+v\n", resp)

	putRequest := &s3.PutObjectAclInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg", Acl: "public-read-write"}
	resp2, err := c.PutObjectAcl(context.Background(), putRequest)
	if err != nil {
		fmt.Printf("PutObjectTagging fail, err: %+v \n", err)
		return
	}
	fmt.Printf("put acl public-read-write success with resp: %+v \n", resp2)

}

func TestCopyObject() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.CopyObjectInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "王文学.jpg.copy", CopySource: &s3.CopySource{CopySourceKey: "王文学.jpg"}}
	resp, err := c.CopyObject(context.Background(), req)
	if err != nil {
		fmt.Printf("CopyObject fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CopyObject success, resp: %+v\n", resp)

}

func TestPart() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.CreateMultipartUploadInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "multicopy.jpg"}
	resp, err := c.CreateMultipartUpload(context.Background(), req)
	if err != nil {
		fmt.Printf("CreateMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CreateMultipartUpload success, resp: %+v\n", resp)

	req1 := &s3.ListMultipartUploadsInput{StoreName: storeName, Bucket: "antsys-wenxuwan", MaxUploads: 1000, KeyMarker: "multicopy.jpg", UploadIdMarker: resp.UploadId}
	resp1, err := c.ListMultipartUploads(context.Background(), req1)
	if err != nil {
		fmt.Printf("ListMultipartUploads fail, err: %+v \n", err)
		return
	}
	fmt.Printf("ListMultipartUploads success, resp: %+v \n", resp1)

	req2 := &s3.UploadPartCopyInput{StoreName: storeName, Bucket: "antsys-wenxuwan", PartNumber: 1, UploadId: resp.UploadId, Key: "multicopy.jpg", StartPosition: 0, PartSize: 1000, CopySource: &s3.CopySource{CopySourceBucket: "antsys-wenxuwan", CopySourceKey: "王文学.jpg"}}
	resp2, err := c.UploadPartCopy(context.Background(), req2)
	if err != nil {
		fmt.Printf("UploadPartCopy fail, err: %+v \n", err)
		return
	}
	fmt.Printf("UploadPartCopy success, resp: %+v \n", resp2)

	req3 := &s3.CompleteMultipartUploadInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "multicopy.jpg", UploadId: resp.UploadId, MultipartUpload: &s3.CompletedMultipartUpload{Parts: []*s3.CompletedPart{{Etag: resp2.CopyPartResult.Etag, PartNumber: req2.PartNumber}}}}
	resp3, err := c.CompleteMultipartUpload(context.Background(), req3)
	if err != nil {
		fmt.Printf("CompleteMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CompleteMultipartUpload success, resp: %+v \n", resp3)

	//req4 := &s3.AbortMultipartUploadInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "海贼王.jpeg", UploadId: "EEE5317D0EB841AC9B80D0B6A2F811AA"}
	//resp4, err := c.AbortMultipartUpload(context.Background(), req4)
	//if err != nil {
	//	fmt.Printf("AbortMultipartUpload fail, err: %+v \n", err)
	//	return
	//}
	//fmt.Printf("AbortMultipartUpload success, resp: %+v \n", resp4)

	req5 := &s3.CreateMultipartUploadInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "client"}
	resp5, err := c.CreateMultipartUpload(context.Background(), req5)
	if err != nil {
		fmt.Printf("CreateMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CreateMultipartUpload success, resp: %+v\n", resp5)

	req6 := &s3.UploadPartInput{
		StoreName:  storeName,
		Bucket:     "antsys-wenxuwan",
		Key:        "client",
		UploadId:   resp5.UploadId,
		PartNumber: 0,
	}
	f, err := os.Open("haizei.jpg")
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
		req6.Body = dataByte
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
	req7 := &s3.CompleteMultipartUploadInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "client", UploadId: resp5.UploadId, MultipartUpload: &s3.CompletedMultipartUpload{Parts: parts}}
	resp7, err := c.CompleteMultipartUpload(context.Background(), req7)
	if err != nil {
		fmt.Printf("CompleteMultipartUpload fail, err: %+v \n", err)
		return
	}
	fmt.Printf("CompleteMultipartUpload success, resp: %+v \n", resp7)
}

func TestListVersion() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.ListObjectVersionsInput{StoreName: storeName, Bucket: "antsys-wenxuwan"}
	resp, err := c.ListObjectVersions(context.Background(), req)
	if err != nil {
		fmt.Printf("ListObjectVersions fail, err: %+v \n", err)
		return
	}
	fmt.Printf("ListObjectVersions success, resp: %+v\n", resp)

}

func TestRestore() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.RestoreObjectInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "client"}
	resp, err := c.RestoreObject(context.Background(), req)
	if err != nil {
		fmt.Printf("RestoreObject fail, err: %+v \n", err)
		return
	}
	fmt.Printf("RestoreObject success, resp: %+v\n", resp)

}

func TestObjectExist() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.IsObjectExistInput{StoreName: storeName, Bucket: "antsys-wenxuwan", Key: "client"}
	resp, err := c.IsObjectExist(context.Background(), req)
	if err != nil {
		fmt.Printf("TestObjectExist fail, err: %+v \n", err)
		return
	}
	fmt.Printf("TestObjectExist success, resp: %+v\n", resp.FileExist)

}

func main() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err:%+v", err)
		return
	}
	c := s3.NewObjectStorageServiceClient(conn)
	req := &s3.InitInput{StoreName: storeName}
	_, err = c.InitClient(context.Background(), req)
	if err != nil {
		fmt.Printf("Init client fail,err:%+v", err)
		return
	}

	if os.Args[1] == "put" {
		TestPutObject(os.Args[2], os.Args[3])
	}
	if os.Args[1] == "get" {
		TestGetObjectInput(os.Args[2])
	}
	if os.Args[1] == "del" {
		TestDeleteObject(os.Args[2])
	}
	if os.Args[1] == "dels" {
		TestDeleteObjects(os.Args[2], os.Args[3])
	}

	if os.Args[1] == "list" {
		TestListObjects()
	}

	if os.Args[1] == "tag" {
		TestTagging()
	}

	if os.Args[1] == "acl" {
		TestAcl()
	}

	if os.Args[1] == "copy" {
		TestCopyObject()
	}

	if os.Args[1] == "part" {
		TestPart()
	}

	if os.Args[1] == "version" {
		TestListVersion()
	}

	if os.Args[1] == "restore" {
		TestRestore()
	}
	if os.Args[1] == "exist" {
		TestObjectExist()
	}
}
