package grpc

import (
	"testing"
)

func TestNewGrpcServer(t *testing.T) {
	apiInterface := &api{}
	NewGrpcServer(WithAPI(apiInterface), WithNewServer(NewDefaultServer), WithGrpcOptions())
}
