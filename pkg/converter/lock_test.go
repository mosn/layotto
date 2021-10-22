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
package converter

import (
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/lock"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"testing"
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
