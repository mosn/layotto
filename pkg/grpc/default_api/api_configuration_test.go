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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"errors"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/pkg/mock"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestGetConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
	mockConfigStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]*configstores.ConfigurationItem{
		{Key: "sofa", Content: "sofa1"},
	}, nil).Times(1)
	res, err := api.GetConfiguration(context.Background(), &runtimev1pb.GetConfigurationRequest{StoreName: "mock", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Nil(t, err)
	assert.Equal(t, res.Items[0].Key, "sofa")
	assert.Equal(t, res.Items[0].Content, "sofa1")
	_, err = api.GetConfiguration(context.Background(), &runtimev1pb.GetConfigurationRequest{StoreName: "etcd", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")

}

func TestSaveConfiguration(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		mockConfigStore.EXPECT().Set(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req *configstores.SetRequest) error {
			assert.Equal(t, "appid", req.AppId)
			assert.Equal(t, "mock", req.StoreName)
			assert.Equal(t, 1, len(req.Items))
			return nil
		})
		req := &runtimev1pb.SaveConfigurationRequest{
			StoreName: "mock",
			AppId:     "appid",
			Items: []*runtimev1pb.ConfigurationItem{
				{
					Key:      "key",
					Content:  "value",
					Group:    "  ",
					Label:    "  ",
					Tags:     nil,
					Metadata: nil,
				},
			},
			Metadata: nil,
		}
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.SaveConfiguration(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("unsupport configstore", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.SaveConfiguration(context.Background(), &runtimev1pb.SaveConfigurationRequest{StoreName: "etcd"})
		assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
	})

}

func TestDeleteConfiguration(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		mockConfigStore.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req *configstores.DeleteRequest) error {
			assert.Equal(t, "appid", req.AppId)
			assert.Equal(t, 1, len(req.Keys))
			assert.Equal(t, "key", req.Keys[0])
			return nil
		})
		req := &runtimev1pb.DeleteConfigurationRequest{
			StoreName: "mock",
			AppId:     "appid",
			Keys:      []string{"key"},
			Metadata:  nil,
		}
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.DeleteConfiguration(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("unsupport configstore", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.DeleteConfiguration(context.Background(), &runtimev1pb.DeleteConfigurationRequest{StoreName: "etcd"})
		assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
	})

}

func TestSubscribeConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)

	//test not support store type
	grpcServer := &MockGrpcServer{req: &runtimev1pb.SubscribeConfigurationRequest{}, err: nil}
	err := api.SubscribeConfiguration(grpcServer)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "configure store [] don't support now")

	//test
	grpcServer2 := &MockGrpcServer{req: &runtimev1pb.SubscribeConfigurationRequest{}, err: errors.New("exit")}
	err = api.SubscribeConfiguration(grpcServer2)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit")
}

type MockGrpcServer struct {
	err error
	req *runtimev1pb.SubscribeConfigurationRequest
	grpc.ServerStream
}

func (m *MockGrpcServer) Send(res *runtimev1pb.SubscribeConfigurationResponse) error {
	return nil
}

func (m *MockGrpcServer) Recv() (*runtimev1pb.SubscribeConfigurationRequest, error) {
	return m.req, m.err
}
