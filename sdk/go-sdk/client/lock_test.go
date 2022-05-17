package client

import (
	"context"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"testing"
)

func TestTryLock(t *testing.T) {
	ctx := context.Background()
	t.Run("try lock", func(t *testing.T) {
		request := runtimev1pb.TryLockRequest{}
		testClient.TryLock(ctx, &request)
	})
}

func TestUnLock(t *testing.T) {
	ctx := context.Background()
	t.Run("try lock", func(t *testing.T) {
		request := runtimev1pb.UnlockRequest{}
		testClient.Unlock(ctx, &request)
	})
}
