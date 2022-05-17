package client

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
