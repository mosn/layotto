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
