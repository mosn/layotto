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
	"sync"

	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/pkg/grpc"
	grpc_api "mosn.io/layotto/pkg/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func NewLifecycleAPI(ac *grpc_api.ApplicationContext) grpc.GrpcAPI {
	components := make(map[componentKey]*dynamicComponentHolder)
	for k, store := range ac.Hellos {
		checkDynamicComponent(kindHello, k, store, components)
	}
	for k, store := range ac.ConfigStores {
		checkDynamicComponent(kindConfig, k, store, components)
	}
	for k, store := range ac.Rpcs {
		checkDynamicComponent(kindRPC, k, store, components)
	}
	for k, store := range ac.PubSubs {
		checkDynamicComponent(kindPubsub, k, store, components)
	}
	for k, store := range ac.StateStores {
		checkDynamicComponent(kindState, k, store, components)
	}
	for k, store := range ac.Files {
		checkDynamicComponent(kindFile, k, store, components)
	}
	for k, store := range ac.Oss {
		checkDynamicComponent(kindOss, k, store, components)
	}
	for k, store := range ac.LockStores {
		checkDynamicComponent(kindLock, k, store, components)
	}
	for k, store := range ac.Sequencers {
		checkDynamicComponent(kindSequencer, k, store, components)
	}
	for k, store := range ac.SecretStores {
		checkDynamicComponent(kindSecret, k, store, components)
	}
	for k, store := range ac.CustomComponent {
		checkDynamicComponent(kindCustom, k, store, components)
	}
	return &server{
		components: components,
		ac:         ac,
	}
}

func checkDynamicComponent(kind string, name string, store interface{}, components map[componentKey]*dynamicComponentHolder) {
	comp, ok := store.(common.DynamicComponent)
	if !ok {
		return
	}
	// put it in the components map
	components[componentKey{
		kind: kind,
		name: name,
	}] = &dynamicComponentHolder{
		DynamicComponent: comp,
		mu:               sync.Mutex{},
	}
}

// server implements runtimev1pb.LifecycleServer
type server struct {
	components map[componentKey]*dynamicComponentHolder
	ac         *grpc_api.ApplicationContext
}

type componentKey struct {
	kind string
	name string
}

type dynamicComponentHolder struct {
	common.DynamicComponent
	mu sync.Mutex
}

func (s *server) ApplyConfiguration(ctx context.Context, in *runtimev1pb.DynamicConfiguration) (*runtimev1pb.ApplyConfigurationResponse, error) {
	// 1. validate parameters
	if in.ComponentConfig == nil || in.ComponentConfig.Kind == "" {
		err := status.Errorf(codes.InvalidArgument, ErrNoField, "kind")
		log.DefaultLogger.Errorf("ApplyConfiguration fail: %+v", err)
		return &runtimev1pb.ApplyConfigurationResponse{}, err
	}
	kind := in.ComponentConfig.Kind
	name := in.ComponentConfig.Name
	if name == "" {
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
	key := componentKey{
		kind: kind,
		name: name,
	}
	holder, ok := s.components[key]
	if !ok {
		err := status.Errorf(codes.InvalidArgument, ErrComponentNotFound, kind, name)
		log.DefaultLogger.Errorf("ApplyConfiguration fail: %+v", err)
		return &runtimev1pb.ApplyConfigurationResponse{}, err
	}

	// 3. lock
	holder.mu.Lock()
	defer holder.mu.Unlock()

	// 4. delegate to components
	err := holder.ApplyConfig(ctx, in.GetComponentConfig().Metadata)
	return &runtimev1pb.ApplyConfigurationResponse{}, err
}

func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	runtimev1pb.RegisterLifecycleServer(rawGrpcServer, s)
	return nil
}
