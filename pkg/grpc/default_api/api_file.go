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

package default_api

import (
	"context"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mosn.io/layotto/components/file"

	"mosn.io/pkg/log"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (a *api) GetFile(req *runtimev1pb.GetFileRequest, stream runtimev1pb.Runtime_GetFileServer) error {
	if a.fileOps[req.StoreName] == nil {
		return status.Errorf(codes.InvalidArgument, "not supported store type: %+v", req.StoreName)
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}
	st := &file.GetFileStu{FileName: req.Name, Metadata: req.Metadata}
	data, err := a.fileOps[req.StoreName].Get(stream.Context(), st)
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
			resp := &runtimev1pb.GetFileResponse{Data: buf[:length]}
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
	server runtimev1pb.Runtime_PutFileServer
}

func newPutObjectStreamReader(data []byte, server runtimev1pb.Runtime_PutFileServer) *putObjectStreamReader {
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
		r.data = req.Data
	}
}

func (a *api) PutFile(stream runtimev1pb.Runtime_PutFileServer) error {
	req, err := stream.Recv()
	if err != nil {
		//if client send eof error directly, return nil
		if err == io.EOF {
			return nil
		}
		return status.Errorf(codes.Internal, "receive file data fail: err: %+v", err)
	}

	if a.fileOps[req.StoreName] == nil {
		return status.Errorf(codes.InvalidArgument, "not support store type: %+v", req.StoreName)
	}
	fileReader := newPutObjectStreamReader(req.Data, stream)
	if req.Metadata == nil {
		req.Metadata = make(map[string]string)
	}
	st := &file.PutFileStu{DataStream: fileReader, FileName: req.Name, Metadata: req.Metadata}
	if err = a.fileOps[req.StoreName].Put(stream.Context(), st); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	stream.SendAndClose(&empty.Empty{})
	return nil
}

// ListFile list all files
func (a *api) ListFile(ctx context.Context, in *runtimev1pb.ListFileRequest) (*runtimev1pb.ListFileResp, error) {
	if in.Request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request can't be nil")
	}
	if in.Request.Metadata == nil {
		in.Request.Metadata = make(map[string]string)
	}

	if a.fileOps[in.Request.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, "not support store type: %+v", in.Request.StoreName)
	}
	resp, err := a.fileOps[in.Request.StoreName].List(ctx, &file.ListRequest{DirectoryName: in.Request.Name, PageSize: in.PageSize, Marker: in.Marker, Metadata: in.Request.Metadata})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	files := make([]*runtimev1pb.FileInfo, 0)
	for _, v := range resp.Files {
		file := &runtimev1pb.FileInfo{}
		file.FileName = v.FileName
		file.LastModified = v.LastModified
		file.Size = v.Size
		file.Metadata = v.Meta
		files = append(files, file)
	}
	return &runtimev1pb.ListFileResp{Files: files, Marker: resp.Marker, IsTruncated: resp.IsTruncated}, nil
}

// DelFile delete specific file
func (a *api) DelFile(ctx context.Context, in *runtimev1pb.DelFileRequest) (*emptypb.Empty, error) {
	errCode := codes.Internal
	if in.Request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request can't be nil")
	}
	if in.Request.Metadata == nil {
		in.Request.Metadata = make(map[string]string)
	}
	if a.fileOps[in.Request.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, "not support store type: %+v", in.Request.StoreName)
	}
	err := a.fileOps[in.Request.StoreName].Del(ctx, &file.DelRequest{FileName: in.Request.Name, Metadata: in.Request.Metadata})
	if err != nil {
		if code, ok := FileErrMap2GrpcErr[err]; ok {
			errCode = code
		}
		return nil, status.Errorf(errCode, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// GetFileMeta get meta of file
func (a *api) GetFileMeta(ctx context.Context, in *runtimev1pb.GetFileMetaRequest) (*runtimev1pb.GetFileMetaResponse, error) {
	errCode := codes.Internal
	if in.Request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request can't be nil")
	}
	if in.Request.Metadata == nil {
		in.Request.Metadata = make(map[string]string)
	}
	if a.fileOps[in.Request.StoreName] == nil {
		return nil, status.Errorf(codes.InvalidArgument, "not support store type: %+v", in.Request.StoreName)
	}
	resp, err := a.fileOps[in.Request.StoreName].Stat(ctx, &file.FileMetaRequest{FileName: in.Request.Name, Metadata: in.Request.Metadata})
	if err != nil {
		if code, ok := FileErrMap2GrpcErr[err]; ok {
			errCode = code
		}
		return nil, status.Errorf(errCode, err.Error())
	}
	meta := &runtimev1pb.FileMeta{}
	meta.Metadata = make(map[string]*runtimev1pb.FileMetaValue)
	for k, v := range resp.Metadata {
		meta.Metadata[k] = &runtimev1pb.FileMetaValue{}
		meta.Metadata[k].Value = append(meta.Metadata[k].Value, v...)
	}
	return &runtimev1pb.GetFileMetaResponse{Size: resp.Size, LastModified: resp.LastModified, Response: meta}, nil
}
