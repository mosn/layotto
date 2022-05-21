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

package runtime

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

var componentName = "state_demo"

func TestHelloApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:11104")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	helloReq := &client.SayHelloRequest{
		ServiceName: "quick_start_demo",
	}
	helloResp, err := cli.SayHello(ctx, helloReq)
	assert.Nil(t, err)
	assert.Equal(t, "welcome layotto", helloResp.Hello)
}

func TestStateApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:11104")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	stateKey := "MyKey"
	stateValue := []byte("Hello Layotto!")
	err = cli.SaveState(ctx, componentName, stateKey, stateValue)
	assert.Nil(t, err)

	stateResp, err := cli.GetState(ctx, "state_demo", stateKey)
	assert.Nil(t, err)
	assert.Equal(t, stateValue, stateResp.Value)
}

func TestLockApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:11104")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	owner1 := uuid.New().String()
	owner2 := uuid.New().String()
	resourceID := "MyLock"

	// 1. client1 tryLock
	resp, err := cli.TryLock(ctx, &runtimev1pb.TryLockRequest{
		StoreName:  storeName,
		ResourceId: resourceID,
		LockOwner:  owner1,
		Expire:     100000,
	})
	assert.Nil(t, err)
	assert.True(t, resp.Success)

	var wg sync.WaitGroup
	wg.Add(1)
	//	2. client2 tryLock
	go func() {
		resp, err := cli.TryLock(ctx, &runtimev1pb.TryLockRequest{
			StoreName:  storeName,
			ResourceId: resourceID,
			LockOwner:  owner2,
			Expire:     1000,
		})
		assert.Nil(t, err)
		assert.False(t, resp.Success)
		wg.Done()
	}()
	wg.Wait()
	// 3. client1 unlock
	unlockResp, err := cli.Unlock(ctx, &runtimev1pb.UnlockRequest{
		StoreName:  storeName,
		ResourceId: resourceID,
		LockOwner:  owner1,
	})
	assert.Nil(t, err)
	assert.Equal(t, runtimev1pb.UnlockResponse_SUCCESS, unlockResp.Status)

	// 4. client2 get lock
	wg.Add(1)
	go func() {
		resp, err := cli.TryLock(ctx, &runtimev1pb.TryLockRequest{
			StoreName:  storeName,
			ResourceId: resourceID,
			LockOwner:  owner2,
			Expire:     10,
		})
		assert.Nil(t, err)
		assert.True(t, true, resp.Success)
		// 5. client2 unlock
		unlockResp, err := cli.Unlock(ctx, &runtimev1pb.UnlockRequest{
			StoreName:  storeName,
			ResourceId: resourceID,
			LockOwner:  owner2,
		})
		assert.Nil(t, err)
		assert.Equal(t, runtimev1pb.UnlockResponse_SUCCESS, unlockResp.Status)
		wg.Done()
	}()
	wg.Wait()
}

func TestSequencerApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:11104")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()
	sequencerKey := "MyKey"

	for i := 1; i < 10; i++ {
		resp, err := cli.GetNextId(ctx, &runtimev1pb.GetNextIdRequest{
			StoreName: "sequencer_demo",
			Key:       sequencerKey,
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(i), resp.NextId)
	}
}
