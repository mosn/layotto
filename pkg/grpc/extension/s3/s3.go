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
	"io"
	"sync"

	"mosn.io/layotto/spec/proto/extension/v1/s3"

	l8s3 "mosn.io/layotto/components/oss"

	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/grpc"
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
	s3.RegisterObjectStorageServiceServer(rawGrpcServer, s)
	return nil
}

func transferData(source interface{}, target interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, target)
	return err
}

func (s *S3Server) GetObject(req *s3.GetObjectInput, stream s3.ObjectStorageService_GetObjectServer) error {
	// 1. validate
	if s.ossInstance[req.StoreName] == nil {
		return status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	// 2. convert request
	st := &l8s3.GetObjectInput{}
	err := transferData(req, st)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "transfer request data fail for GetObject,err: %+v", err)
	}
	// 3. find the component
	result, err := s.ossInstance[req.StoreName].GetObject(stream.Context(), st)
	if err != nil {
		return status.Errorf(codes.Internal, "get file fail,err: %+v", err)
	}

	buffsPtr := bytesPool.Get().(*[]byte)
	buf := *buffsPtr
	if len(buf) == 0 {
		buf = make([]byte, 102400)
	}
	defer func() {
		result.DataStream.Close()
		*buffsPtr = buf
		bytesPool.Put(buffsPtr)
	}()

	for {
		length, err := result.DataStream.Read(buf)
		if err != nil && err != io.EOF {
			log.DefaultLogger.Warnf("get file fail, err: %+v", err)
			return status.Errorf(codes.Internal, "get file fail,err: %+v", err)
		}
		if err == nil || (err == io.EOF && length != 0) {
			resp := &s3.GetObjectOutput{}
			err := transferData(result, resp)
			if err != nil {
				return status.Errorf(codes.InvalidArgument, "transfer request data fail for GetObject,err: %+v", err)
			}
			resp.Body = buf[:length]
			if err = stream.Send(resp); err != nil {
				return status.Errorf(codes.Internal, "send file data fail,err: %+v", err)
			}
		}
		if err == io.EOF {
			return nil
		}
	}
}

type putObjectStreamReader struct {
	data   []byte
	server s3.ObjectStorageService_PutObjectServer
}

func newPutObjectStreamReader(data []byte, server s3.ObjectStorageService_PutObjectServer) *putObjectStreamReader {
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

func (s *S3Server) PutObject(stream s3.ObjectStorageService_PutObjectServer) error {
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
	var resp *l8s3.PutObjectOutput
	if resp, err = s.ossInstance[req.StoreName].PutObject(stream.Context(), st); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.PutObjectOutput{}
	err = transferData(resp, output)
	if err != nil {
		return status.Errorf(codes.Internal, "transfer response data fail for PutObject,err: %+v", err)
	}
	return stream.SendAndClose(output)
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
	var resp *l8s3.DeleteObjectOutput
	if resp, err = s.ossInstance[req.StoreName].DeleteObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.DeleteObjectOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for DeleteObject,err: %+v", err)
	}
	return output, nil
}
func (s *S3Server) PutObjectTagging(ctx context.Context, req *s3.PutObjectTaggingInput) (*s3.PutObjectTaggingOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}

	st := &l8s3.PutObjectTaggingInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for PutObjectTagging,err: %+v", err)
	}
	var resp *l8s3.PutObjectTaggingOutput
	if resp, err = s.ossInstance[req.StoreName].PutObjectTagging(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.PutObjectTaggingOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for PutObjectTagging,err: %+v", err)
	}
	return output, nil
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
	var resp *l8s3.DeleteObjectTaggingOutput
	if resp, err = s.ossInstance[req.StoreName].DeleteObjectTagging(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.DeleteObjectTaggingOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for DeleteObjectTagging,err: %+v", err)
	}
	return output, nil
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
	var resp *l8s3.GetObjectTaggingOutput
	if resp, err = s.ossInstance[req.StoreName].GetObjectTagging(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.GetObjectTaggingOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for GetObjectTagging,err: %+v", err)
	}
	return output, nil
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
	var resp *l8s3.CopyObjectOutput
	if resp, err = s.ossInstance[req.StoreName].CopyObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.CopyObjectOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for CopyObject,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.DeleteObjectsOutput
	if resp, err = s.ossInstance[req.StoreName].DeleteObjects(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.DeleteObjectsOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for DeleteObjects,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.ListObjectsOutput
	if resp, err = s.ossInstance[req.StoreName].ListObjects(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.ListObjectsOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for ListObjects,err: %+v", err)
	}
	return output, nil

}
func (s *S3Server) GetObjectCannedAcl(ctx context.Context, req *s3.GetObjectCannedAclInput) (*s3.GetObjectCannedAclOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.GetObjectCannedAclInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for GetObjectAcl,err: %+v", err)
	}
	var resp *l8s3.GetObjectCannedAclOutput
	if resp, err = s.ossInstance[req.StoreName].GetObjectCannedAcl(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.GetObjectCannedAclOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for GetObjectAcl,err: %+v", err)
	}
	return output, nil

}
func (s *S3Server) PutObjectCannedAcl(ctx context.Context, req *s3.PutObjectCannedAclInput) (*s3.PutObjectCannedAclOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.PutObjectCannedAclInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for PutObjectAcl,err: %+v", err)
	}
	var resp *l8s3.PutObjectCannedAclOutput
	if resp, err = s.ossInstance[req.StoreName].PutObjectCannedAcl(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.PutObjectCannedAclOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for PutObjectAcl,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.RestoreObjectOutput
	if resp, err = s.ossInstance[req.StoreName].RestoreObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.RestoreObjectOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for RestoreObject,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.CreateMultipartUploadOutput
	if resp, err = s.ossInstance[req.StoreName].CreateMultipartUpload(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.CreateMultipartUploadOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for CreateMultipartUpload,err: %+v", err)
	}
	return output, nil

}

type uploadPartStreamReader struct {
	data   []byte
	server s3.ObjectStorageService_UploadPartServer
}

func newUploadPartStreamReader(data []byte, server s3.ObjectStorageService_UploadPartServer) *uploadPartStreamReader {
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

func (s *S3Server) UploadPart(stream s3.ObjectStorageService_UploadPartServer) error {
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
	var resp *l8s3.UploadPartOutput
	if resp, err = s.ossInstance[req.StoreName].UploadPart(stream.Context(), st); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.UploadPartOutput{}
	err = transferData(resp, output)
	if err != nil {
		return status.Errorf(codes.Internal, "transfer response data fail for UploadPart,err: %+v", err)
	}
	return stream.SendAndClose(output)

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
	var resp *l8s3.UploadPartCopyOutput
	if resp, err = s.ossInstance[req.StoreName].UploadPartCopy(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.UploadPartCopyOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for UploadPartCopy,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.CompleteMultipartUploadOutput
	if resp, err = s.ossInstance[req.StoreName].CompleteMultipartUpload(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.CompleteMultipartUploadOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for CompleteMultipartUpload,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.AbortMultipartUploadOutput
	if resp, err = s.ossInstance[req.StoreName].AbortMultipartUpload(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.AbortMultipartUploadOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for AbortMultipartUpload,err: %+v", err)
	}
	return output, nil

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
	var resp *l8s3.ListMultipartUploadsOutput
	if resp, err = s.ossInstance[req.StoreName].ListMultipartUploads(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.ListMultipartUploadsOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for AbortMultipartUpload,err: %+v", err)
	}
	return output, nil
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
	var resp *l8s3.ListObjectVersionsOutput
	if resp, err = s.ossInstance[req.StoreName].ListObjectVersions(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.ListObjectVersionsOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for ListObjectVersions,err: %+v", err)
	}
	return output, nil

}

func (s *S3Server) HeadObject(ctx context.Context, req *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.HeadObjectInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for ListObjectVersions,err: %+v", err)
	}
	var resp *l8s3.HeadObjectOutput
	if resp, err = s.ossInstance[req.StoreName].HeadObject(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.HeadObjectOutput{}
	err = transferData(resp, output)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer response data fail for ListObjectVersions,err: %+v", err)
	}
	return output, nil

}

func (s *S3Server) IsObjectExist(ctx context.Context, req *s3.IsObjectExistInput) (*s3.IsObjectExistOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.IsObjectExistInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for IsObjectExist,err: %+v", err)
	}
	var resp *l8s3.IsObjectExistOutput
	if resp, err = s.ossInstance[req.StoreName].IsObjectExist(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.IsObjectExistOutput{}
	output.FileExist = resp.FileExist
	return output, nil
}

func (s *S3Server) SignURL(ctx context.Context, req *s3.SignURLInput) (*s3.SignURLOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.SignURLInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for SignURL,err: %+v", err)
	}
	var resp *l8s3.SignURLOutput
	if resp, err = s.ossInstance[req.StoreName].SignURL(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.SignURLOutput{}
	output.SignedUrl = resp.SignedUrl
	return output, nil
}

func (s *S3Server) UpdateDownloadBandwidthRateLimit(ctx context.Context, req *s3.UpdateBandwidthRateLimitInput) (*emptypb.Empty, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.UpdateBandwidthRateLimitInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for UpdateDownloadBandwidthRateLimit,err: %+v", err)
	}
	if err := s.ossInstance[req.StoreName].UpdateDownloadBandwidthRateLimit(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *S3Server) UpdateUploadBandwidthRateLimit(ctx context.Context, req *s3.UpdateBandwidthRateLimitInput) (*emptypb.Empty, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.UpdateBandwidthRateLimitInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for UpdateUploadBandwidthRateLimit,err: %+v", err)
	}
	if err := s.ossInstance[req.StoreName].UpdateUploadBandwidthRateLimit(ctx, st); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

type appendObjectStreamReader struct {
	data   []byte
	server s3.ObjectStorageService_AppendObjectServer
}

func newAppendObjectStreamReader(data []byte, server s3.ObjectStorageService_AppendObjectServer) *appendObjectStreamReader {
	return &appendObjectStreamReader{data: data, server: server}
}

func (r *appendObjectStreamReader) Read(p []byte) (int, error) {
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

func (s *S3Server) AppendObject(stream s3.ObjectStorageService_AppendObjectServer) error {
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
	fileReader := newAppendObjectStreamReader(req.Body, stream)

	st := &l8s3.AppendObjectInput{}
	err = transferData(req, st)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "transfer request data fail for AppendObject,err: %+v", err)
	}
	st.DataStream = fileReader
	var resp *l8s3.AppendObjectOutput
	if resp, err = s.ossInstance[req.StoreName].AppendObject(stream.Context(), st); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	output := &s3.AppendObjectOutput{}
	output.AppendPosition = resp.AppendPosition
	return stream.SendAndClose(output)

}

func (s *S3Server) ListParts(ctx context.Context, req *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
	if s.ossInstance[req.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, NotSupportStoreName, req.StoreName)
	}
	st := &l8s3.ListPartsInput{}
	err := transferData(req, st)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer request data fail for ListParts,err: %+v", err)
	}
	resp, err := s.ossInstance[req.StoreName].ListParts(ctx, st)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	out := &s3.ListPartsOutput{}
	err = transferData(resp, out)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "transfer response data fail for ListParts,err: %+v", err)
	}
	return out, nil
}
