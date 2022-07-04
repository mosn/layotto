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
	"flag"
	"fmt"

	client "mosn.io/layotto/sdk/go-sdk/client"
)

const (
	key1 = "key1"
	key2 = "key2"
	key3 = "key3"
	key4 = "key4"
	key5 = "key5"
)

var storeName string

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
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
	value := []byte("hello world")
	fmt.Printf("Start testing %v\n", storeName)

	// Belows are CRUD examples.
	// save state
	testSave(ctx, cli, storeName, key1, value)

	// get state
	testGet(ctx, cli, storeName, key1)

	// SaveBulkState with options and metadata
	testSaveBulkState(ctx, cli, storeName, key1, value, key2)

	keyTostate := testGetBulkState(ctx, cli, storeName, key1, key2)

	// delete state
	testDelete(ctx, cli, storeName, key1, keyTostate[key1].Etag)
	testDelete(ctx, cli, storeName, key2, keyTostate[key2].Etag)
}

func testGetBulkState(ctx context.Context, cli client.Client, store string, key1 string, key2 string) map[string]*client.BulkStateItem {
	state, err := cli.GetBulkState(ctx, store, []string{key1, key2, key3, key4, key5}, nil, 3)
	if err != nil {
		panic(err)
	}
	m := make(map[string]*client.BulkStateItem)
	for _, item := range state {
		fmt.Printf("GetBulkState succeeded.key:%v ,value:%v ,etag:%v ,metadata:%v \n", item.Key, string(item.Value), item.Etag, item.Metadata)
		m[item.Key] = item
	}
	return m
}

func testDelete(ctx context.Context, cli client.Client, store string, key string, etag string) {
	if err := cli.DeleteStateWithETag(ctx, store, key, &client.ETag{Value: etag}, nil, nil); err != nil {
		panic(err)
	}
	fmt.Printf("DeleteState succeeded.key:%v\n", key)
}

func testSaveBulkState(ctx context.Context, cli client.Client, store string, key string, value []byte, key2 string) {
	item := &client.SetStateItem{
		// etag is used to implement Optimistic Concurrency Control (OCC)
		//	see https://docs.dapr.io/developing-applications/building-blocks/state-management/state-management-overview/#concurrency
		Etag: &client.ETag{
			Value: "2",
		},
		Key: key,
		Metadata: map[string]string{
			"some-key-for-component": "some-value",
		},
		Value: value,
		Options: &client.StateOptions{
			Concurrency: client.StateConcurrencyLastWrite,
			Consistency: client.StateConsistencyStrong,
		},
	}
	item2 := *item
	item2.Key = key2

	if err := cli.SaveBulkState(ctx, store, item, &item2); err != nil {
		panic(err)
	}
	fmt.Printf("SaveBulkState succeeded.[key:%s etag:%s]: %s\n", item.Key, item.Etag.Value, string(item.Value))
	fmt.Printf("SaveBulkState succeeded.[key:%s etag:%s]: %s\n", item2.Key, item2.Etag.Value, string(item2.Value))
}

func testGet(ctx context.Context, cli client.Client, store string, key string) {
	item, err := cli.GetState(ctx, store, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetState succeeded.[key:%s etag:%s]: %s\n", item.Key, item.Etag, string(item.Value))
}

func testSave(ctx context.Context, cli client.Client, store string, key string, value []byte) {
	if err := cli.SaveState(ctx, store, key, value); err != nil {
		panic(err)
	}
	fmt.Printf("SaveState succeeded.key:%v , value: %v \n", key, string(value))
}
