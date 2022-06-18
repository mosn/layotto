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

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"

	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	appid      = "testApplication_yang"
	group      = "application"
	writeTimes = 4
)

var (
	storeName string
	mode      string
)

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
	flag.StringVar(&mode, "mode", "raw", "set `mode`")
}

func main() {
	// parse command arguments
	flag.Parse()
	if storeName == "" {
		panic("storeName is empty.")
	}

	// create a layotto client
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()

	// Belows are CRUD examples
	// 1. set
	testSet(ctx, cli)

	// 2. get after set
	// Since configuration data might be cached and eventual-consistent,we need to sleep a while before querying new data
	time.Sleep(time.Second * 2)
	testGet(ctx, cli)

	// 3. delete
	testDelete(ctx, cli)

	// 4. get after delete
	//sleep a while before querying deleted data
	time.Sleep(time.Second * 2)
	testGet(ctx, cli)

	// 5. show how to use subscribe API
	// with sdk
	if mode == "sdk" {
		testSubscribeWithSDK(ctx, cli)
	} else {
		// besides sdk,u can also call layotto with grpc
		testSubscribeWithGrpc(ctx)
	}

}

func testSubscribeWithSDK(ctx context.Context, cli client.Client) {
	var wg sync.WaitGroup
	wg.Add(1)
	// 1. subscribe
	ch := cli.SubscribeConfiguration(ctx, &client.ConfigurationRequestItem{
		StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"haha", "heihei"},
	})
	// 2. client loop receiving changes in another gorountine
	go func() {
		for resp := range ch {
			if resp.Err != nil {
				panic(resp.Err)
			}
			marshal, err := json.Marshal(resp.Item)
			if err != nil {
				fmt.Println("marshal resp.Item failed.")
				return
			}
			fmt.Printf("receive subscribe resp: %+v\n", string(marshal))
		}
		fmt.Println("subscribe channel has been closed.")
	}()

	// 3. loop setting configuration in another gorountine
	go func() {
		idx := 0
		for {
			time.Sleep(1 * time.Second)
			idx++
			newItmes := make([]*client.ConfigurationItem, 0, 10)
			newItmes = append(newItmes, &client.ConfigurationItem{Group: "application", Label: "prod", Key: "heihei", Content: "heihei" + strconv.Itoa(idx)})
			fmt.Println("write start")
			cli.SaveConfiguration(ctx, &client.SaveConfigurationRequest{StoreName: storeName, AppId: appid, Items: newItmes})
		}
	}()
	wg.Wait()
}

func testSubscribeWithGrpc(ctx context.Context) {
	// 1. connect with grpc
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := runtimev1pb.NewRuntimeClient(conn)

	// get client for subscribe
	var wg sync.WaitGroup
	wg.Add(2)
	cli, err := c.SubscribeConfiguration(ctx)
	if err != nil {
		panic(err)
	}
	// client receive changes
	go func() {
		for {
			data, err := cli.Recv()
			if err != nil {
				fmt.Println("subscribe failed", err)
				continue
			}
			fmt.Println("receive subscribe resp", data)
			// if it's the last item, break and end this demo
			if data.Items[0].Content == "heihei4" {
				break
			}
		}
		wg.Done()
	}()
	// client send subscribe request
	go func() {
		cli.Send(&runtimev1pb.SubscribeConfigurationRequest{StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"haha", "key1"}})
		cli.Send(&runtimev1pb.SubscribeConfigurationRequest{StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"heihei"}})
	}()

	// loop write in another gorountine
	go func() {
		for idx := 1; idx <= writeTimes; idx++ {
			time.Sleep(1 * time.Second)
			newItmes := make([]*runtimev1pb.ConfigurationItem, 0, 10)
			newItmes = append(newItmes, &runtimev1pb.ConfigurationItem{Group: "application", Label: "prod", Key: "heihei", Content: "heihei" + strconv.Itoa(idx)})
			fmt.Println("write start")
			c.SaveConfiguration(ctx, &runtimev1pb.SaveConfigurationRequest{StoreName: storeName, AppId: appid, Items: newItmes})
		}
		wg.Done()
	}()
	wg.Wait()
}

func testSet(ctx context.Context, cli client.Client) {
	item1 := &client.ConfigurationItem{Group: group, Label: "prod", Key: "key1", Content: "value1", Tags: map[string]string{
		"release": "1.0.0",
		"feature": "print",
	}}
	item2 := &client.ConfigurationItem{Group: group, Label: "prod", Key: "haha", Content: "heihei", Tags: map[string]string{
		"release": "1.0.0",
		"feature": "haha",
	}}
	saveRequest := &client.SaveConfigurationRequest{StoreName: storeName, AppId: appid}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)
	if err := cli.SaveConfiguration(ctx, saveRequest); err != nil {
		panic(err)
	}
	fmt.Println("save key success")
}

func testGet(ctx context.Context, cli client.Client) {
	getRequest := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "prod", Keys: []string{"key1", "haha"}}
	items, err := cli.GetConfiguration(ctx, getRequest)
	//validate
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("get configuration after save, %+v \n", item)
	}
}

func testDelete(ctx context.Context, cli client.Client) {
	request := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "prod", Keys: []string{"key1", "haha"}}
	if err := cli.DeleteConfiguration(ctx, request); err != nil {
		panic(err)
	}
	fmt.Printf("delete keys success\n")
}
