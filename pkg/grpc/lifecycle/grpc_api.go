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

package lifecycle

import (
	"context"

	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/pkg/runtime/lifecycle"

	"mosn.io/layotto/pkg/grpc"
	grpc_api "mosn.io/layotto/pkg/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func NewLifecycleAPI(ac *grpc_api.ApplicationContext) grpc.GrpcAPI {
	return &server{
		components: ac.DynamicComponents,
	}
}

// server implements runtimev1pb.LifecycleServer
type server struct {
	components map[lifecycle.ComponentKey]common.DynamicComponent
}

func (s *server) ApplyConfiguration(ctx context.Context, in *runtimev1pb.DynamicConfiguration) (*runtimev1pb.ApplyConfigurationResponse, error) {
	// 1. validate parameters
	if in.ComponentConfig == nil || in.ComponentConfig.Kind == "" {
		return &runtimev1pb.ApplyConfigurationResponse{}, invalidArgumentError(grpc_api.ErrNoField, "kind")
	}
	kind := in.ComponentConfig.Kind
	name := in.ComponentConfig.Name
	if name == "" {
		return &runtimev1pb.ApplyConfigurationResponse{}, invalidArgumentError(grpc_api.ErrNoField, "name")
	}
	if len(in.ComponentConfig.Metadata) == 0 {
		return &runtimev1pb.ApplyConfigurationResponse{}, invalidArgumentError(grpc_api.ErrNoField, "metadata")
	}
	// 2. find the component
	key := lifecycle.ComponentKey{
		Kind: kind,
		Name: name,
	}
	holder, ok := s.components[key]
	if !ok {
		return &runtimev1pb.ApplyConfigurationResponse{}, invalidArgumentError(grpc_api.ErrComponentNotFound, kind, name)
	}

	// 3. delegate to the components
	err := holder.ApplyConfig(ctx, in.GetComponentConfig().Metadata)
	return &runtimev1pb.ApplyConfigurationResponse{}, err
}

func invalidArgumentError(format string, a ...interface{}) error {
	err := status.Errorf(codes.InvalidArgument, format, a...)
	log.DefaultLogger.Errorf("ApplyConfiguration fail: %+v", err)
	return err
}

func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	runtimev1pb.RegisterLifecycleServer(rawGrpcServer, s)
	return nil
}
