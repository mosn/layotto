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
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	appid = "testApplication"
	group = "application"
)

var PubsubStoreName string

func TestHelloApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
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
	assert.Equal(t, "greeting", helloResp.Hello)
}

func TestConfigurationApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	item1 := &client.ConfigurationItem{Group: group, Label: "test", Key: "key1", Content: "value1"}
	item2 := &client.ConfigurationItem{Group: group, Label: "test", Key: "key2", Content: "value2"}

	componentArray := [...]string{"etcd_config_demo"}
	for _, storeName := range componentArray {
		saveRequest := &client.SaveConfigurationRequest{StoreName: storeName, AppId: appid}
		saveRequest.Items = append(saveRequest.Items, item1)
		saveRequest.Items = append(saveRequest.Items, item2)
		err = cli.SaveConfiguration(ctx, saveRequest)
		assert.Nil(t, err)

		// Since configuration data might be cached and eventual-consistent,we need to sleep a while before querying new data
		time.Sleep(time.Second * 2)
		getRequest := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "test", Keys: []string{"key1", "key2"}}
		resp, err := cli.GetConfiguration(ctx, getRequest)
		assert.Nil(t, err)
		assert.Equal(t, item1, resp[0])
		assert.Equal(t, item2, resp[1])
	}
}

func TestStateApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	stateKey := "MyStateKey"
	stateValue := []byte("Hello Layotto!")

	componentArray := [...]string{"redis_state_demo", "zookeeper_state_demo"}
	for _, storeName := range componentArray {
		err = cli.SaveState(ctx, storeName, stateKey, stateValue)
		assert.Nil(t, err)

		stateResp, err := cli.GetState(ctx, storeName, stateKey)
		assert.Nil(t, err)
		assert.Equal(t, stateValue, stateResp.Value)
	}
}

func TestLockApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	owner1 := uuid.New().String()
	owner2 := uuid.New().String()
	resourceID := "MyLock"

	componentArray := [...]string{"redis_lock_demo", "etcd_lock_demo", "zookeeper_lock_demo"}
	for _, storeName := range componentArray {
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
}

func TestSequencerApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	sequencerKey := "MySequencerKey"

	componentArray := [...]string{"redis_sequencer_demo", "etcd_sequencer_demo", "zookeeper_sequencer_demo"}
	for _, storeName := range componentArray {
		for i := 1; i < 10; i++ {
			resp, err := cli.GetNextId(ctx, &runtimev1pb.GetNextIdRequest{
				StoreName: storeName,
				Key:       sequencerKey,
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(i), resp.NextId)
		}
	}
}

func TestPubsubApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	componentArray := [...]string{"redis_pub_subs_demo"}
	for _, storeName := range componentArray {
		PubsubStoreName = storeName
		//start a subscriber to subscribe to events
		go func() {
			port := 8888
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			assert.Nil(t, err)
			grpcServer := grpc.NewServer()
			runtimev1pb.RegisterAppCallbackServer(grpcServer, &AppCallbackServerImpl{})
			err = grpcServer.Serve(lis)
			assert.Nil(t, err)
		}()

		//publisher publishes an event
		err := cli.PublishEventfromCustomContent(ctx, storeName, "topic1", "value1")
		assert.Nil(t, err)

		//sleep for a while before subscriber get messages
		time.Sleep(time.Second)
	}
}

func TestSecretApi(t *testing.T) {
	cli, err := client.NewClientWithAddress("127.0.0.1:34904")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()

	componentArray := [...]string{"local_file_secret_demo"}
	for _, storeName := range componentArray {
		//get the secret
		secret, err := cli.GetSecret(ctx, &runtimev1pb.GetSecretRequest{
			StoreName: storeName,
			Key:       "testPassword",
		})
		assert.Nil(t, err)
		assert.Equal(t, "pw2", secret.Data["testPassword"])

		//get the bulk secret
		bulkSecret, err := cli.GetBulkSecret(ctx, &runtimev1pb.GetBulkSecretRequest{
			StoreName: storeName,
		})
		assert.Nil(t, err)
		assert.Equal(t, "admin", bulkSecret.Data["db-user-pass:username"].Secrets["db-user-pass:username"])
		assert.Equal(t, "pw1", bulkSecret.Data["db-user-pass:password"].Secrets["db-user-pass:password"])
		assert.Equal(t, "pw2", bulkSecret.Data["testPassword"].Secrets["testPassword"])
	}
}

type AppCallbackServerImpl struct {
}

func (a *AppCallbackServerImpl) ListTopicSubscriptions(ctx context.Context, empty *empty.Empty) (*runtimev1pb.ListTopicSubscriptionsResponse, error) {
	result := &runtimev1pb.ListTopicSubscriptionsResponse{}
	ts := &runtimev1pb.TopicSubscription{
		PubsubName: PubsubStoreName,
		Topic:      "topic1",
		Metadata:   nil,
	}
	result.Subscriptions = append(result.Subscriptions, ts)
	return result, nil
}

func (a *AppCallbackServerImpl) OnTopicEvent(ctx context.Context, request *runtimev1pb.TopicEventRequest) (*runtimev1pb.TopicEventResponse, error) {
	var t *testing.T
	assert.Equal(t, "topic1", request.Topic)
	assert.Equal(t, "value1", request.Data)
	return &runtimev1pb.TopicEventResponse{}, nil
}
