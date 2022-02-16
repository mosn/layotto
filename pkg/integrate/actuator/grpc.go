package actuator

import (
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	"net"
)

type GrpcServerWithActuator struct {
	srv mgrpc.RegisteredServer
}

func (s *GrpcServerWithActuator) Serve(ln net.Listener) (err error) {
	err = s.srv.Serve(ln)
	if err != nil {
		GetRuntimeReadinessIndicator().SetUnhealthy(err.Error())
		GetRuntimeLivenessIndicator().SetUnhealthy(err.Error())
	}
	return
}
func (s *GrpcServerWithActuator) Stop() {
	GetRuntimeReadinessIndicator().SetUnhealthy("shutdown")
	GetRuntimeLivenessIndicator().SetUnhealthy("shutdown")
	s.srv.Stop()
}

func (s *GrpcServerWithActuator) GracefulStop() {
	GetRuntimeReadinessIndicator().SetUnhealthy("shutdown")
	GetRuntimeLivenessIndicator().SetUnhealthy("shutdown")
	s.srv.GracefulStop()
}

func NewGrpcServerWithActuator(srv mgrpc.RegisteredServer) (mgrpc.RegisteredServer, error) {
	return &GrpcServerWithActuator{
		srv: srv,
	}, nil
}
