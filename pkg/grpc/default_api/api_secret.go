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

	"github.com/dapr/components-contrib/secretstores"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mosn.io/layotto/pkg/messages"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/pkg/log"
)

func (a *api) GetSecret(ctx context.Context, in *runtimev1pb.GetSecretRequest) (*runtimev1pb.GetSecretResponse, error) {
	// check parameters
	if a.secretStores == nil || len(a.secretStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSecretStoreNotConfigured)
		log.DefaultLogger.Errorf("GetSecret fail,not configured err:%+v", err)
		return &runtimev1pb.GetSecretResponse{}, err
	}
	secretStoreName := in.StoreName
	if a.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		log.DefaultLogger.Errorf("GetSecret fail,not find err:%+v", err)
		return &runtimev1pb.GetSecretResponse{}, err
	}

	// TODO permission control
	if !a.isSecretAllowed(in.StoreName, in.Key) {
		err := status.Errorf(codes.PermissionDenied, messages.ErrPermissionDenied, in.Key, in.StoreName)
		return &runtimev1pb.GetSecretResponse{}, err
	}

	// delegate to components
	req := secretstores.GetSecretRequest{
		Name:     in.Key,
		Metadata: in.Metadata,
	}

	// parse result
	getResponse, err := a.secretStores[secretStoreName].GetSecret(req)
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrSecretGet, req.Name, secretStoreName, err.Error())
		log.DefaultLogger.Errorf("GetSecret fail,get secret err:%+v", err)
		return &runtimev1pb.GetSecretResponse{}, err
	}
	response := &runtimev1pb.GetSecretResponse{}
	if getResponse.Data != nil {
		response.Data = getResponse.Data
	}

	return &runtimev1pb.GetSecretResponse{Data: response.Data}, nil
}

func (a *api) GetBulkSecret(ctx context.Context, in *runtimev1pb.GetBulkSecretRequest) (*runtimev1pb.GetBulkSecretResponse, error) {
	// check parameters
	if a.secretStores == nil || len(a.secretStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSecretStoreNotConfigured)
		log.DefaultLogger.Errorf("GetBulkSecret fail,not configured err:%+v", err)
		return &runtimev1pb.GetBulkSecretResponse{}, err
	}
	secretStoreName := in.StoreName
	if a.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		log.DefaultLogger.Errorf("GetBulkSecret fail,not find err:%+v", err)
		return &runtimev1pb.GetBulkSecretResponse{}, err
	}

	// delegate to components
	req := secretstores.BulkGetSecretRequest{
		Metadata: in.Metadata,
	}
	getResponse, err := a.secretStores[secretStoreName].BulkGetSecret(req)

	// parse result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrBulkSecretGet, secretStoreName, err.Error())
		log.DefaultLogger.Errorf("GetBulkSecret fail,bulk secret err:%+v", err)
		return &runtimev1pb.GetBulkSecretResponse{}, err
	}

	// filter result
	filteredSecrets := map[string]map[string]string{}
	for key, v := range getResponse.Data {
		// TODO: permission control
		if a.isSecretAllowed(secretStoreName, key) {
			filteredSecrets[key] = v
		} else {
			log.DefaultLogger.Debugf(messages.ErrPermissionDenied, key, in.StoreName)
		}
	}
	response := &runtimev1pb.GetBulkSecretResponse{}
	if getResponse.Data != nil {
		response.Data = map[string]*runtimev1pb.SecretResponse{}
		for key, v := range filteredSecrets {
			response.Data[key] = &runtimev1pb.SecretResponse{Secrets: v}
		}
	}

	return response, nil
}

func (a *api) isSecretAllowed(storeName string, key string) bool {
	// TODO: add permission control
	return true
}
