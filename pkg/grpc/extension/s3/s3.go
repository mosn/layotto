package s3

import (
	"context"
	rawGRPC "google.golang.org/grpc"
	"mosn.io/layotto/components/custom"
	"mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/spec/proto/extension/v1"
)

var (
	s3Instance *S3Server
)

type S3Server struct {
	appId  string
	config custom.Config
	ctx    context.Context
}

func NewS3Component() custom.Component {
	s3Instance = &S3Server{}
	return s3Instance
}

func NewS3Server(ac *grpc.ApplicationContext) grpc.GrpcAPI {
	s3Instance.appId = ac.AppId
	return s3Instance
}

func (s *S3Server) Initialize(ctx context.Context, config custom.Config) error {
	s.config = config
	return nil
}

func (s *S3Server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *S3Server) Register(rawGrpcServer *rawGRPC.Server) error {
	s3.RegisterS3Server(rawGrpcServer, s)
	return nil
}

func (s *S3Server) GetObject(req *s3.GetObjectInput, stream s3.S3_GetObjectServer) error {
	return nil
}

func (s *S3Server) PutObject(s3.S3_PutObjectServer) error {
	return nil
}
