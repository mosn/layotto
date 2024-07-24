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

package dapr

import (
	"context"

	"github.com/dapr/components-contrib/secretstores"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
)

func (d *daprGrpcAPI) GetSecret(ctx context.Context, request *runtime.GetSecretRequest) (*runtime.GetSecretResponse, error) {
	// 1. check parameters
	if d.secretStores == nil || len(d.secretStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSecretStoreNotConfigured)
		log.DefaultLogger.Errorf("GetSecret fail,not configured err:%+v", err)
		return &runtime.GetSecretResponse{}, err
	}
	secretStoreName := request.StoreName

	if d.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		log.DefaultLogger.Errorf("GetSecret fail,not find err:%+v", err)
		return &runtime.GetSecretResponse{}, err
	}

	// 2. TODO permission control
	if !d.isSecretAllowed(request.StoreName, request.Key) {
		err := status.Errorf(codes.PermissionDenied, messages.ErrPermissionDenied, request.Key, request.StoreName)
		return &runtime.GetSecretResponse{}, err
	}

	// 3. delegate to components
	req := secretstores.GetSecretRequest{
		Name:     request.Key,
		Metadata: request.Metadata,
	}
	getResponse, err := d.secretStores[secretStoreName].GetSecret(req)
	// 4. parse result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrSecretGet, req.Name, secretStoreName, err.Error())
		log.DefaultLogger.Errorf("GetSecret fail,get secret err:%+v", err)
		return &runtime.GetSecretResponse{}, err
	}

	response := &runtime.GetSecretResponse{}
	if getResponse.Data != nil {
		response.Data = getResponse.Data
	}
	return response, nil
}

func (d *daprGrpcAPI) GetBulkSecret(ctx context.Context, in *runtime.GetBulkSecretRequest) (*runtime.GetBulkSecretResponse, error) {
	// 1. check parameters
	if d.secretStores == nil || len(d.secretStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSecretStoreNotConfigured)
		log.DefaultLogger.Errorf("GetBulkSecret fail,not configured err:%+v", err)
		return &runtime.GetBulkSecretResponse{}, err
	}
	secretStoreName := in.StoreName

	if d.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		log.DefaultLogger.Errorf("GetBulkSecret fail,not find err:%+v", err)
		return &runtime.GetBulkSecretResponse{}, err
	}
	// 2. delegate to components
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
