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
