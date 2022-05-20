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
	"fmt"
	"testing"

	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_state "mosn.io/layotto/pkg/mock/components/state"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestSaveState(t *testing.T) {
	t.Run("error when request is nil", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.SaveState(context.Background(), nil)
		assert.NotNil(t, err)
	})
	t.Run("error when no state store registered", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{}, nil, nil, nil, nil, nil)
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
	})
	t.Run("error when store name wrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock1",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.NotNil(t, err)
	})
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
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
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
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
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
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
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
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
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

	t.Run("save error: ETagInvalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).Return(state.NewETagError(state.ETagInvalid, nil))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
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
		assert.Equal(t, "rpc error: code = InvalidArgument desc = failed saving state in state store mock: invalid etag value", err.Error())
	})

	t.Run("save error: ETagMismatch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).Return(state.NewETagError(state.ETagMismatch, nil))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
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
		assert.Equal(t, "rpc error: code = Aborted desc = failed saving state in state store mock: possible etag mismatch. error from state store", err.Error())
	})
}

func TestGetBulkState(t *testing.T) {
	t.Run("state store not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "abc",
		}
		_, err := api.GetBulkState(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = state store abc is not found", err.Error())
	})

	t.Run("get state error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkGet(gomock.Any()).Return(false, nil, fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "mock",
			Keys:      []string{"mykey"},
		}
		_, err := api.GetBulkState(context.Background(), req)
		assert.Equal(t, "net error", err.Error())
	})

	t.Run("support bulk get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		compResp := []state.BulkGetResponse{
			{
				Data:     []byte("mock data"),
				Metadata: nil,
			},
		}
		mockStore.EXPECT().BulkGet(gomock.Any()).Return(true, compResp, nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "mock",
			Keys:      []string{"mykey"},
		}
		rsp, err := api.GetBulkState(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, []byte("mock data"), rsp.GetItems()[0].GetData())
	})

	t.Run("don't support bulk get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		resp1 := &state.GetResponse{
			Data:     []byte("mock data"),
			Metadata: nil,
		}

		resp2 := &state.GetResponse{
			Data:     []byte("mock data2"),
			Metadata: nil,
		}
		mockStore.EXPECT().BulkGet(gomock.Any()).Return(false, nil, nil)
		mockStore.EXPECT().Get(gomock.Any()).Return(resp1, nil)
		mockStore.EXPECT().Get(gomock.Any()).Return(resp2, nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "mock",
			Keys:      []string{"mykey", "mykey2"},
		}
		rsp, err := api.GetBulkState(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(rsp.GetItems()))
	})

}

func TestGetState(t *testing.T) {
	t.Run("state store not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "abc",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = state store abc is not found", err.Error())
	})

	t.Run("state store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "abc",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = state store is not configured", err.Error())
	})

	t.Run("get modified state key error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "mock",
			Key:       "mykey||abc",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "input key/keyPrefix 'mykey||abc' can't contain '||'", err.Error())
	})

	t.Run("get state error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "mock",
			Key:       "mykey",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "rpc error: code = Internal desc = fail to get mykey from state store mock: net error", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		compResp := &state.GetResponse{
			Data:     []byte("mock data"),
			Metadata: nil,
		}
		mockStore.EXPECT().Get(gomock.Any()).Return(compResp, nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "mock",
			Key:       "mykey",
		}
		rsp, err := api.GetState(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, []byte("mock data"), rsp.GetData())
	})

}

func TestDeleteState(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().Delete(gomock.Any()).DoAndReturn(func(req *state.DeleteRequest) error {
			assert.Equal(t, "abc", req.Key)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteStateRequest{
			StoreName: "mock",
			Key:       "abc",
		}
		_, err := api.DeleteState(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteStateRequest{
			StoreName: "mock",
			Key:       "abc",
		}
		_, err := api.DeleteState(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = failed deleting state with key abc: net error", err.Error())
	})
}

func TestDeleteBulkState(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkDelete(gomock.Any()).DoAndReturn(func(reqs []state.DeleteRequest) error {
			assert.Equal(t, "abc", reqs[0].Key)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteBulkStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key: "abc",
				},
			},
		}
		_, err := api.DeleteBulkState(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkDelete(gomock.Any()).Return(fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteBulkStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key: "abc",
				},
			},
		}
		_, err := api.DeleteBulkState(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "net error", err.Error())
	})
}

type MockTxStore struct {
	state.Store
	state.TransactionalStore
}

func (m *MockTxStore) Init(metadata state.Metadata) error {
	return m.Store.Init(metadata)
}

func TestExecuteStateTransaction(t *testing.T) {
	t.Run("state store not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "abc",
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = state store abc is not found", err.Error())
	})

	t.Run("state store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "abc",
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = state store is not configured", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return([]state.Feature{state.FeatureTransactional})

		mockTxStore := mock_state.NewMockTransactionalStore(gomock.NewController(t))
		mockTxStore.EXPECT().Multi(gomock.Any()).DoAndReturn(func(req *state.TransactionalStateRequest) error {
			assert.Equal(t, 2, len(req.Operations))
			assert.Equal(t, "mosn", req.Metadata["runtime"])
			assert.Equal(t, state.Upsert, req.Operations[0].Operation)
			assert.Equal(t, state.Delete, req.Operations[1].Operation)
			return nil
		})

		store := &MockTxStore{
			mockStore,
			mockTxStore,
		}

		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": store}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "mock",
			Operations: []*runtimev1pb.TransactionalStateOperation{
				{
					OperationType: string(state.Upsert),
					Request: &runtimev1pb.StateItem{
						Key:   "upsert",
						Value: []byte("mock data"),
					},
				},
				{
					OperationType: string(state.Delete),
					Request: &runtimev1pb.StateItem{
						Key: "delete_abc",
					},
				},
				{
					OperationType: string(state.Delete),
				},
			},
			Metadata: map[string]string{
				"runtime": "mosn",
			},
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return([]state.Feature{state.FeatureTransactional})

		mockTxStore := mock_state.NewMockTransactionalStore(gomock.NewController(t))
		mockTxStore.EXPECT().Multi(gomock.Any()).Return(fmt.Errorf("net error"))

		store := &MockTxStore{
			mockStore,
			mockTxStore,
		}
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": store}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "mock",
			Operations: []*runtimev1pb.TransactionalStateOperation{
				{
					OperationType: string(state.Upsert),
					Request: &runtimev1pb.StateItem{
						Key:   "upsert",
						Value: []byte("mock data"),
					},
				},
				{
					OperationType: string(state.Delete),
					Request: &runtimev1pb.StateItem{
						Key: "delete_abc",
					},
				},
				{
					OperationType: string(state.Delete),
				},
			},
			Metadata: map[string]string{
				"runtime": "mosn",
			},
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = error while executing state transaction: net error", err.Error())
	})
}
