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
	rawGRPC "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	l8s3 "mosn.io/layotto/components/file"
	"mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/spec/proto/extension/v1"
)

var (
	s3Instance *S3Server
)

const (
	AliyunOSS = "aliyun"
	MinioOSS  = "minio"
	AwsOSS    = "aws"
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

func (s *S3Server) GetObject(req *s3.GetObjectInput, stream s3.S3_GetObjectServer) error {
	return nil
}

func (s *S3Server) PutObject(s3.S3_PutObjectServer) error {
	return nil
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
