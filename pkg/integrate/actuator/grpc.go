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

package actuator

import (
	"net"

	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
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
