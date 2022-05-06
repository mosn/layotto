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

	"mosn.io/layotto/components/lock"
	mock_lock "mosn.io/layotto/pkg/mock/components/lock"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestTryLockRequest2ComponentRequest(t *testing.T) {
	req := TryLockRequest2ComponentRequest(&runtimev1pb.TryLockRequest{
		StoreName:  "redis",
		ResourceId: "resourceId",
		LockOwner:  "owner1",
		Expire:     1000,
	})
	assert.True(t, req.ResourceId == "resourceId")
	assert.True(t, req.LockOwner == "owner1")
	assert.True(t, req.Expire == 1000)
	req = TryLockRequest2ComponentRequest(nil)
	assert.NotNil(t, req)
}

func TestTryLockResponse2GrpcResponse(t *testing.T) {
	resp := TryLockResponse2GrpcResponse(&lock.TryLockResponse{
		Success: true,
	})
	assert.True(t, resp.Success)
	resp2 := TryLockResponse2GrpcResponse(nil)
	assert.NotNil(t, resp2)
}

func TestUnlockGrpc2ComponentRequest(t *testing.T) {
	req := UnlockGrpc2ComponentRequest(&runtimev1pb.UnlockRequest{
		StoreName:  "redis",
		ResourceId: "resourceId",
		LockOwner:  "owner1",
	})
	assert.True(t, req.ResourceId == "resourceId")
	assert.True(t, req.LockOwner == "owner1")
	req = UnlockGrpc2ComponentRequest(nil)
	assert.NotNil(t, req)
}

func TestUnlockComp2GrpcResponse(t *testing.T) {
	resp := UnlockComp2GrpcResponse(&lock.UnlockResponse{Status: lock.SUCCESS})
	assert.True(t, resp.Status == runtimev1pb.UnlockResponse_SUCCESS)
	resp2 := UnlockComp2GrpcResponse(nil)
	assert.NotNil(t, resp2)
}

func TestTryLock(t *testing.T) {
	t.Run("lock store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName: "abc",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = lock store is not configured", err.Error())
	})

	t.Run("resourceid empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName: "abc",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = ResourceId is empty in lock store abc", err.Error())
	})

	t.Run("lock owner empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = LockOwner is empty in lock store abc", err.Error())
	})

	t.Run("lock expire is not positive", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
			LockOwner:  "owner",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Expire is not positive in lock store abc", err.Error())
	})

	t.Run("lock store not found", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
			LockOwner:  "owner",
			Expire:     1,
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = lock store abc not found", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		mockLockStore.EXPECT().TryLock(gomock.Any()).DoAndReturn(func(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
			assert.Equal(t, "lock|||resource", req.ResourceId)
			assert.Equal(t, "owner", req.LockOwner)
			assert.Equal(t, int32(1), req.Expire)
			return &lock.TryLockResponse{
				Success: true,
			}, nil
		})
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "mock",
			ResourceId: "resource",
			LockOwner:  "owner",
			Expire:     1,
		}
		resp, err := api.TryLock(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, true, resp.Success)
	})

}

func TestUnlock(t *testing.T) {
	t.Run("lock store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName: "abc",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = lock store is not configured", err.Error())
	})

	t.Run("resourceid empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName: "abc",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = ResourceId is empty in lock store abc", err.Error())
	})

	t.Run("lock owner empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = LockOwner is empty in lock store abc", err.Error())
	})

	t.Run("lock store not found", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
			LockOwner:  "owner",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = lock store abc not found", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		mockLockStore.EXPECT().Unlock(gomock.Any()).DoAndReturn(func(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
			assert.Equal(t, "lock|||resource", req.ResourceId)
			assert.Equal(t, "owner", req.LockOwner)
			return &lock.UnlockResponse{
				Status: lock.SUCCESS,
			}, nil
		})
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName:  "mock",
			ResourceId: "resource",
			LockOwner:  "owner",
		}
		resp, err := api.Unlock(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, runtimev1pb.UnlockResponse_SUCCESS, resp.Status)
	})
}
