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
	"github.com/dapr/components-contrib/secretstores"
	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/layotto/pkg/grpc"
	grpc_api "mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/pkg/log"
)

func NewLifecycleAPI(ac *grpc_api.ApplicationContext) grpc.GrpcAPI {
	return &server{
		ac: ac,
	}
}

// server implements runtimev1pb.LifecycleServer
type server struct {
	ac *grpc_api.ApplicationContext
}

func (s *server) ApplyConfiguration(ctx context.Context, in *runtimev1pb.DynamicConfiguration) (*runtimev1pb.ApplyConfigurationResponse, error) {
	// 1. validate parameters
	if in.ComponentConfig == nil || in.ComponentConfig.Kind == "" {
		err := status.Errorf(codes.InvalidArgument, ErrNoField, "kind")
		log.DefaultLogger.Errorf("ApplyConfiguration fail: %+v", err)
		return &runtimev1pb.ApplyConfigurationResponse{}, err
	}
	if in.ComponentConfig.Name == "" {
		err := status.Errorf(codes.InvalidArgument, ErrNoField, "name")
		log.DefaultLogger.Errorf("ApplyConfiguration fail: %+v", err)
		return &runtimev1pb.ApplyConfigurationResponse{}, err
	}
	if len(in.ComponentConfig.Metadata) == 0 {
		err := status.Errorf(codes.InvalidArgument, ErrNoField, "metadata")
		log.DefaultLogger.Errorf("ApplyConfiguration fail: %+v", err)
		return &runtimev1pb.ApplyConfigurationResponse{}, err
	}
	// 2. find the component
	var component interface{}
	switch in.ComponentConfig.Kind {
	case ""
	}


	secretStoreName := in.StoreName

	if d.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		log.DefaultLogger.Errorf("GetBulkSecret fail,not find err:%+v", err)
		return &runtime.GetBulkSecretResponse{}, err
	}
	// 3. lock
	// 4. delegate to components
	req := secretstores.BulkGetSecretRequest{
		Metadata: in.Metadata,
	}
	getResponse, err := d.secretStores[secretStoreName].BulkGetSecret(req)
	// 3. parse result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrBulkSecretGet, secretStoreName, err.Error())
		log.DefaultLogger.Errorf("GetBulkSecret fail,bulk secret err:%+v", err)
		return &runtime.GetBulkSecretResponse{}, err
	}

	// 4. filter result
	filteredSecrets := map[string]map[string]string{}
	for key, v := range getResponse.Data {
		// TODO: permission control
		if d.isSecretAllowed(secretStoreName, key) {
			filteredSecrets[key] = v
		} else {
			log.DefaultLogger.Debugf(messages.ErrPermissionDenied, key, in.StoreName)
		}
	}
	response := &runtime.GetBulkSecretResponse{}
	if getResponse.Data != nil {
		response.Data = map[string]*runtime.SecretResponse{}
		for key, v := range filteredSecrets {
			response.Data[key] = &runtime.SecretResponse{Secrets: v}
		}
	}
	return response, nil
}

func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	runtimev1pb.RegisterLifecycleServer(rawGrpcServer, s)
	return nil
}
