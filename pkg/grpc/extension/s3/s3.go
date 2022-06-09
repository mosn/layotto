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
	Region   = "region"
	EndPoint = "endpoint"
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

func (s *S3Server) InitClient(ctx context.Context, req *s3.InitInput) (*emptypb.Empty, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	err := s.ossInstance[req.StoreName].InitClient(ctx, &l8s3.InitRequest{Metadata: req.Metadata})
	if err != nil {
		log.DefaultLogger.Errorf("InitClient fail, err: %+v", err)
	}
	return &emptypb.Empty{}, err
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
		return status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
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
		return status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
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

func (s *S3Server) DeleteObject(ctx context.Context, req *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.DeleteObjectInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for DeleteObject,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].DeleteObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.DeleteObjectOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for DeleteObject,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) PutObjectTagging(ctx context.Context, req *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.PutBucketTaggingInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for PutObjectTagging,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].PutObjectTagging(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.PutBucketTaggingOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for PutObjectTagging,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) DeleteObjectTagging(ctx context.Context, req *s3.DeleteObjectTaggingInput) (*s3.DeleteObjectTaggingOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.DeleteObjectTaggingInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for DeleteObjectTagging,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].DeleteObjectTagging(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.DeleteObjectTaggingOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for DeleteObjectTagging,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) GetObjectTagging(ctx context.Context, req *s3.GetObjectTaggingInput) (*s3.GetObjectTaggingOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.GetObjectTaggingInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for GetObjectTagging,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].GetObjectTagging(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.GetObjectTaggingOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for GetObjectTagging,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) CopyObject(ctx context.Context, req *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.CopyObjectInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for CopyObject,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].CopyObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.CopyObjectOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for CopyObject,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) DeleteObjects(ctx context.Context, req *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.DeleteObjectsInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for DeleteObjects,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].DeleteObjects(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.DeleteObjectsOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for DeleteObjects,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) ListObjects(ctx context.Context, req *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.ListObjectsInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for ListObjects,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].ListObjects(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.ListObjectsOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for ListObjects,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) GetObjectAcl(ctx context.Context, req *s3.GetObjectAclInput) (*s3.GetObjectAclOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.GetObjectAclInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for GetObjectAcl,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].GetObjectAcl(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.GetObjectAclOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for GetObjectAcl,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) PutObjectAcl(ctx context.Context, req *s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.PutObjectAclInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for PutObjectAcl,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].PutObjectAcl(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.PutObjectAclOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for PutObjectAcl,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) RestoreObject(ctx context.Context, req *s3.RestoreObjectInput) (*s3.RestoreObjectOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.RestoreObjectInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for RestoreObject,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].RestoreObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.RestoreObjectOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for RestoreObject,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) CreateMultipartUpload(ctx context.Context, req *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.CreateMultipartUploadInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for CreateMultipartUpload,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].CreateMultipartUpload(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.CreateMultipartUploadOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for CreateMultipartUpload,err: %+v", err)
		}
		return output, nil
	}
}

type uploadPartStreamReader struct {
	data   []byte
	server s3.S3_UploadPartServer
}

func newUploadPartStreamReader(data []byte, server s3.S3_UploadPartServer) *uploadPartStreamReader {
	return &uploadPartStreamReader{data: data, server: server}
}

func (r *uploadPartStreamReader) Read(p []byte) (int, error) {
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

func (s *S3Server) UploadPart(stream s3.S3_UploadPartServer) error {
	req, err := stream.Recv()
	if err != nil {
		//if client send eof error directly, return nil
		if err == io.EOF {
			return nil
		}
		return status.Errorf(codes.Internal, "receive file data fail: err: %+v", err)
	}

	if s.ossInstance[req.StoreName] == nil {
		return status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	fileReader := newUploadPartStreamReader(req.Body, stream)

	st := &l8s3.UploadPartInput{}
	err = transferData(req, st)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "transfer request data fail for UploadPart,err: %+v", err)
	}
	st.DataStream = fileReader
	if resp, err := s.ossInstance[req.StoreName].UploadPart(stream.Context(), st); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.UploadPartOutput{}
		err := transferData(resp, output)
		if err != nil {
			return status.Errorf(codes.Internal, "transfer response data fail for UploadPart,err: %+v", err)
		}
		return stream.SendAndClose(output)
	}
}
func (s *S3Server) UploadPartCopy(ctx context.Context, req *s3.UploadPartCopyInput) (*s3.UploadPartCopyOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.UploadPartCopyInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for UploadPartCopy,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].UploadPartCopy(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.UploadPartCopyOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for UploadPartCopy,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) CompleteMultipartUpload(ctx context.Context, req *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.CompleteMultipartUploadInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for CompleteMultipartUpload,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].CompleteMultipartUpload(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.CompleteMultipartUploadOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for CompleteMultipartUpload,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) AbortMultipartUpload(ctx context.Context, req *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.AbortMultipartUploadInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for AbortMultipartUpload,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].AbortMultipartUpload(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.AbortMultipartUploadOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for AbortMultipartUpload,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) ListMultipartUploads(ctx context.Context, req *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.ListMultipartUploadsInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for AbortMultipartUpload,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].ListMultipartUploads(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.ListMultipartUploadsOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for AbortMultipartUpload,err: %+v", err)
		}
		return output, nil
	}
}
func (s *S3Server) ListObjectVersions(ctx context.Context, req *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.ListObjectVersionsInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for ListObjectVersions,err: %+v", err)
	}
	if resp, err := s.ossInstance[req.StoreName].ListObjectVersions(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		output := &s3.ListObjectVersionsOutput{}
		err := transferData(resp, output)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "transfer response data fail for ListObjectVersions,err: %+v", err)
		}
		return output, nil
	}
}
