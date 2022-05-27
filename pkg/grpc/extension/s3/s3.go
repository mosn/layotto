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

package s3

import (
	"context"
	"encoding/json"
	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	l8s3 "mosn.io/layotto/components/file"
	"mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/spec/proto/extension/v1"
	"mosn.io/pkg/log"
	"sync"
)

var (
	s3Instance *S3Server
)

var (
	bytesPool = sync.Pool{
		New: func() interface{} {
			// set size to 100kb
			return new([]byte)
		},
	}
)

const (
	Provider        = "provider"
	Region          = "region"
	EndPoint        = "endpoint"
	AccessKeyID     = "accessKeyID"
	AccessKeySecret = "accessKeySecret"
)

type S3Server struct {
	appId       string
	ossInstance map[string]l8s3.Oss
}

func NewS3Server(ac *grpc.ApplicationContext) grpc.GrpcAPI {
	s3Instance = &S3Server{}
	s3Instance.appId = ac.AppId
	s3Instance.ossInstance = ac.Oss
	return s3Instance
}

func (s *S3Server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *S3Server) Register(rawGrpcServer *rawGRPC.Server) error {
	s3.RegisterS3Server(rawGrpcServer, s)
	return nil
}

func (s *S3Server) InitClient(ctx context.Context, req *s3.InitRequest) (*emptypb.Empty, error) {
	//if s.config.Metadata[Provider] == "" {
	//	return nil, errors.New("please specific the oss provider in configure file")
	//}

	return &emptypb.Empty{}, nil
}

func transferData(source interface{}, target interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, target)
	return err
}

func (s *S3Server) GetObject(req *s3.GetObjectInput, stream s3.S3_GetObjectServer) error {
	if s.ossInstance[req.StoreName] == nil {
		return status.Errorf(codes.InvalidArgument, "not supported store type: %+v", req.StoreName)
	}
	st := &l8s3.GetObjectInput{}
	err := transferData(req, st)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "transfer request data fail for GetObject,err: %+v", err)
	}
	data, err := s.ossInstance[req.StoreName].GetObject(stream.Context(), st)
	if err != nil {
		return status.Errorf(codes.Internal, "get file fail,err: %+v", err)
	}

	buffsPtr := bytesPool.Get().(*[]byte)
	buf := *buffsPtr
	if len(buf) == 0 {
		buf = make([]byte, 102400)
	}
	defer func() {
		data.Close()
		*buffsPtr = buf
		bytesPool.Put(buffsPtr)
	}()

	for {
		length, err := data.Read(buf)
		if err != nil && err != io.EOF {
			log.DefaultLogger.Warnf("get file fail, err: %+v", err)
			return status.Errorf(codes.Internal, "get file fail,err: %+v", err)
		}
		if err == nil || (err == io.EOF && length != 0) {
			resp := &s3.GetObjectOutput{Body: buf[:length]}
			if err = stream.Send(resp); err != nil {
				return status.Errorf(codes.Internal, "send file data fail,err: %+v", err)
			}
		}
		if err == io.EOF {
			return nil
		}
	}
	return nil
}

type putObjectStreamReader struct {
	data   []byte
	server s3.S3_PutObjectServer
}

func newPutObjectStreamReader(data []byte, server s3.S3_PutObjectServer) *putObjectStreamReader {
	return &putObjectStreamReader{data: data, server: server}
}

func (r *putObjectStreamReader) Read(p []byte) (int, error) {
	var count int
	total := len(p)
	for {
		if len(r.data) > 0 {
			n := copy(p[count:], r.data)
			r.data = r.data[n:]
			count += n
			if count == total {
				return count, nil
			}
		}
		req, err := r.server.Recv()
		if err != nil {
			if err != io.EOF {
				log.DefaultLogger.Errorf("recv data from grpc stream fail, err:%+v", err)
			}
			return count, err
		}
		r.data = req.Body
	}
}

func (s *S3Server) PutObject(stream s3.S3_PutObjectServer) error {
	req, err := stream.Recv()
	if err != nil {
		//if client send eof error directly, return nil
		if err == io.EOF {
			return nil
		}
		return status.Errorf(codes.Internal, "receive file data fail: err: %+v", err)
	}

	if s.ossInstance[req.StoreName] == nil {
		return status.Errorf(codes.InvalidArgument, "not support store type: %+v", req.StoreName)
	}
	fileReader := newPutObjectStreamReader(req.Body, stream)

	st := &l8s3.PutObjectInput{}
	err = transferData(req, st)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "transfer request data fail for PutObject,err: %+v", err)
	}
	st.DataStream = fileReader
	if resp, err := s.ossInstance[req.StoreName].PutObject(stream.Context(), st); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.PutObjectOutput{}
		err := transferData(resp, output)
		if err != nil {
			return status.Errorf(codes.Internal, "transfer response data fail for PutObject,err: %+v", err)
		}
		return stream.SendAndClose(output)
	}
}

func (s *S3Server) DeleteObject(context.Context, *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {

	return nil, nil
}
func (s *S3Server) PutObjectTagging(context.Context, *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	return nil, nil
}
func (s *S3Server) DeleteObjectTagging(context.Context, *s3.DeleteObjectTaggingInput) (*s3.DeleteObjectTaggingOutput, error) {
	return nil, nil
}
func (s *S3Server) GetObjectTagging(context.Context, *s3.GetObjectTaggingInput) (*s3.GetObjectTaggingOutput, error) {
	return nil, nil
}
func (s *S3Server) CopyObject(context.Context, *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return nil, nil
}
func (s *S3Server) DeleteObjects(context.Context, *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	return nil, nil
}
func (s *S3Server) ListObjects(context.Context, *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return nil, nil
}
func (s *S3Server) GetObjectAcl(context.Context, *s3.GetObjectAclInput) (*s3.GetObjectAclOutput, error) {
	return nil, nil
}
func (s *S3Server) PutObjectAcl(context.Context, *s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error) {
	return nil, nil
}
func (s *S3Server) RestoreObject(context.Context, *s3.RestoreObjectInput) (*s3.RestoreObjectOutput, error) {
	return nil, nil
}
func (s *S3Server) CreateMultipartUpload(context.Context, *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	return nil, nil
}
func (s *S3Server) UploadPart(s3.S3_UploadPartServer) error {
	return nil
}
func (s *S3Server) UploadPartCopy(context.Context, *s3.UploadPartCopyInput) (*s3.UploadPartCopyOutput, error) {
	return nil, nil
}
func (s *S3Server) CompleteMultipartUpload(context.Context, *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
	return nil, nil
}
func (s *S3Server) AbortMultipartUpload(context.Context, *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error) {
	return nil, nil
}
func (s *S3Server) ListMultipartUploads(context.Context, *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
	return nil, nil
}
func (s *S3Server) ListObjectVersions(context.Context, *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error) {
	return nil, nil
}
