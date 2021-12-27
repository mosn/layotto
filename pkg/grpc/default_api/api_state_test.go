//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package default_api

import (
	"context"
	"fmt"
	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_state "mosn.io/layotto/pkg/mock/components/state"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"testing"
)

func TestSaveState(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).DoAndReturn(func(reqs []state.SetRequest) error {
			assert.Equal(t, 1, len(reqs))
			assert.Equal(t, "abc", reqs[0].Key)
			assert.Equal(t, []byte("mock data"), reqs[0].Value)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.Nil(t, err)
	})
	t.Run("with options last-write and eventual", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).DoAndReturn(func(reqs []state.SetRequest) error {
			assert.Equal(t, 1, len(reqs))
			assert.Equal(t, "abc", reqs[0].Key)
			assert.Equal(t, []byte("mock data"), reqs[0].Value)
			assert.Equal(t, "last-write", reqs[0].Options.Concurrency)
			assert.Equal(t, "eventual", reqs[0].Options.Consistency)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
					Options: &runtimev1pb.StateOptions{
						Concurrency: runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE,
						Consistency: runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL,
					},
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.Nil(t, err)
	})
	t.Run("with options first-write and strong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).DoAndReturn(func(reqs []state.SetRequest) error {
			assert.Equal(t, 1, len(reqs))
			assert.Equal(t, "abc", reqs[0].Key)
			assert.Equal(t, []byte("mock data"), reqs[0].Value)
			assert.Equal(t, "first-write", reqs[0].Options.Concurrency)
			assert.Equal(t, "strong", reqs[0].Options.Consistency)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
					Options: &runtimev1pb.StateOptions{
						Concurrency: runtimev1pb.StateOptions_CONCURRENCY_FIRST_WRITE,
						Consistency: runtimev1pb.StateOptions_CONSISTENCY_STRONG,
					},
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("save error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).Return(fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = failed saving state in state store mock: net error", err.Error())
	})
}

func TestGetResponse2GetStateResponse(t *testing.T) {
	resp := GetResponse2GetStateResponse(&state.GetResponse{
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: make(map[string]string),
	})
	assert.Equal(t, resp.Data, []byte("v"))
	assert.Equal(t, resp.Etag, "")
	assert.True(t, len(resp.Metadata) == 0)
}

func TestGetResponse2BulkStateItem(t *testing.T) {
	itm := GetResponse2BulkStateItem(&state.GetResponse{
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: make(map[string]string),
	}, "key")
	assert.Equal(t, itm.Key, "key")
	assert.Equal(t, itm.Data, []byte("v"))
	assert.Equal(t, itm.Etag, "")
	assert.Equal(t, itm.Error, "")
	assert.True(t, len(itm.Metadata) == 0)
}

func TestBulkGetResponse2BulkStateItem(t *testing.T) {
	t.Run("convert nil", func(t *testing.T) {
		itm := BulkGetResponse2BulkStateItem(nil)
		assert.NotNil(t, itm)
	})
	t.Run("normal", func(t *testing.T) {
		itm := BulkGetResponse2BulkStateItem(&state.BulkGetResponse{
			Key:      "key",
			Data:     []byte("v"),
			ETag:     nil,
			Metadata: nil,
			Error:    "",
		})
		assert.Equal(t, itm.Key, "key")
		assert.Equal(t, itm.Data, []byte("v"))
		assert.Equal(t, itm.Etag, "")
		assert.Equal(t, itm.Error, "")
		assert.True(t, len(itm.Metadata) == 0)
	})
}

func TestStateItem2SetRequest(t *testing.T) {
	req := StateItem2SetRequest(&runtimev1pb.StateItem{
		Key:      "",
		Value:    []byte("v"),
		Etag:     nil,
		Metadata: nil,
		Options: &runtimev1pb.StateOptions{
			Concurrency: runtimev1pb.StateOptions_CONCURRENCY_UNSPECIFIED,
			Consistency: runtimev1pb.StateOptions_CONSISTENCY_UNSPECIFIED,
		},
	}, "appid||key")
	assert.Equal(t, req.Key, "appid||key")
	assert.Equal(t, req.Value, []byte("v"))
	assert.Nil(t, req.ETag)
	assert.Equal(t, req.Options.Consistency, "")
	assert.Equal(t, req.Options.Concurrency, "")
}

func TestDeleteStateRequest2DeleteRequest(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		req := DeleteStateRequest2DeleteRequest(nil, "")
		assert.NotNil(t, req)
	})
	t.Run("normal", func(t *testing.T) {
		req := DeleteStateRequest2DeleteRequest(&runtimev1pb.DeleteStateRequest{
			StoreName: "redis",
			Key:       "",
			Etag:      nil,
			Options: &runtimev1pb.StateOptions{
				Concurrency: runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE,
				Consistency: runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL,
			},
			Metadata: nil,
		}, "appid||key")
		assert.Equal(t, req.Key, "appid||key")
		assert.Nil(t, req.ETag)
		assert.Equal(t, req.Options.Consistency, "eventual")
		assert.Equal(t, req.Options.Concurrency, "last-write")
	})
}

func TestStateItem2DeleteRequest(t *testing.T) {
	req := StateItem2DeleteRequest(&runtimev1pb.StateItem{
		Key:      "",
		Value:    []byte("v"),
		Etag:     nil,
		Metadata: nil,
		Options: &runtimev1pb.StateOptions{
			Concurrency: runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE,
			Consistency: runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL,
		},
	}, "appid||key")
	assert.Equal(t, req.Key, "appid||key")
	assert.Nil(t, req.ETag)
	assert.Nil(t, req.ETag)
	assert.Equal(t, req.Options.Consistency, "eventual")
	assert.Equal(t, req.Options.Concurrency, "last-write")
}
