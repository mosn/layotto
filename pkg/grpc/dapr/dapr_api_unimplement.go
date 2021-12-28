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
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
	"mosn.io/pkg/log"
)

func (d *daprGrpcAPI) GetState(ctx context.Context, request *runtime.GetStateRequest) (*runtime.GetStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetBulkState(ctx context.Context, request *runtime.GetBulkStateRequest) (*runtime.GetBulkStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) SaveState(ctx context.Context, request *runtime.SaveStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) QueryStateAlpha1(ctx context.Context, request *runtime.QueryStateRequest) (*runtime.QueryStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) DeleteState(ctx context.Context, request *runtime.DeleteStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) DeleteBulkState(ctx context.Context, request *runtime.DeleteBulkStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) ExecuteStateTransaction(ctx context.Context, request *runtime.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) PublishEvent(ctx context.Context, request *runtime.PublishEventRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetSecret(ctx context.Context, request *runtime.GetSecretRequest) (*runtime.GetSecretResponse, error) {
	if d.secretStores == nil || len(d.secretStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSecretStoreNotConfigured)
		return &runtime.GetSecretResponse{}, err
	}

	secretStoreName := request.StoreName

	if d.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		return &runtime.GetSecretResponse{}, err
	}

	if !d.isSecretAllowed(request.StoreName, request.Key) {
		err := status.Errorf(codes.PermissionDenied, messages.ErrPermissionDenied, request.Key, request.StoreName)
		return &runtime.GetSecretResponse{}, err
	}

	req := secretstores.GetSecretRequest{
		Name:     request.Key,
		Metadata: request.Metadata,
	}

	getResponse, err := d.secretStores[secretStoreName].GetSecret(req)
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrSecretGet, req.Name, secretStoreName, err.Error())
		return &runtime.GetSecretResponse{}, err
	}

	response := &runtime.GetSecretResponse{}
	if getResponse.Data != nil {
		response.Data = getResponse.Data
	}
	return response, nil
}

func (d *daprGrpcAPI) GetBulkSecret(ctx context.Context, request *runtime.GetBulkSecretRequest) (*runtime.GetBulkSecretResponse, error) {
	if d.secretStores == nil || len(d.secretStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSecretStoreNotConfigured)
		return &runtime.GetBulkSecretResponse{}, err
	}
	secretStoreName := request.StoreName
	if d.secretStores[secretStoreName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSecretStoreNotFound, secretStoreName)
		return &runtime.GetBulkSecretResponse{}, err
	}

	req := secretstores.BulkGetSecretRequest{
		Metadata: request.Metadata,
	}

	getResponse, err := d.secretStores[secretStoreName].BulkGetSecret(req)
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrBulkSecretGet, secretStoreName, err.Error())
		return &runtime.GetBulkSecretResponse{}, err
	}

	filteredSecrets := map[string]map[string]string{}
	for key, v := range getResponse.Data {
		if d.isSecretAllowed(secretStoreName, key) {
			filteredSecrets[key] = v
		} else {
			log.DefaultLogger.Debugf(messages.ErrPermissionDenied, key, request.StoreName)
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

func (d *daprGrpcAPI) RegisterActorTimer(ctx context.Context, request *runtime.RegisterActorTimerRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) UnregisterActorTimer(ctx context.Context, request *runtime.UnregisterActorTimerRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) RegisterActorReminder(ctx context.Context, request *runtime.RegisterActorReminderRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) UnregisterActorReminder(ctx context.Context, request *runtime.UnregisterActorReminderRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetActorState(ctx context.Context, request *runtime.GetActorStateRequest) (*runtime.GetActorStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) ExecuteActorStateTransaction(ctx context.Context, request *runtime.ExecuteActorStateTransactionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) InvokeActor(ctx context.Context, request *runtime.InvokeActorRequest) (*runtime.InvokeActorResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetConfigurationAlpha1(ctx context.Context, request *runtime.GetConfigurationRequest) (*runtime.GetConfigurationResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) SubscribeConfigurationAlpha1(request *runtime.SubscribeConfigurationRequest, server runtime.Dapr_SubscribeConfigurationAlpha1Server) error {
	panic("implement me")
}

func (d *daprGrpcAPI) GetMetadata(ctx context.Context, empty *emptypb.Empty) (*runtime.GetMetadataResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) SetMetadata(ctx context.Context, request *runtime.SetMetadataRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) Shutdown(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	panic("implement me")
}
