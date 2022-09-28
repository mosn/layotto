/*
 * Copyright 2022 Layotto Authors
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

package redis

import (
	"context"
	"testing"

	"github.com/helbing/testtools/flow"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/sdk/go-sdk/client"
	"mosn.io/layotto/tests/tools/envrunner"
)

func TestComponentsStateRedis(t *testing.T) {
	cli, err := client.NewClient()
	assert.Nil(t, err)
	defer cli.Close()

	var (
		storeName = "test_components_state_redis"
		ctx       = context.Background()
		value     = []byte("hello world")
	)

	flow.New(envrunner.NewLayotto(t, "./docker-compose.yaml")).
		Case("Test state", func(t *testing.T) {
			key := "test_state_key"

			// save
			err = cli.SaveState(ctx, storeName, key, value)
			assert.Nil(t, err)

			// get
			item, err := cli.GetState(ctx, storeName, key)
			assert.Nil(t, err)
			assert.Equal(t, key, item.Key)
			assert.Equal(t, value, item.Value)

			// delete
			err = cli.DeleteStateWithETag(ctx, storeName, key, &client.ETag{
				Value: item.Etag,
			}, nil, nil)
			assert.Nil(t, err)

			// check deleted
			item, err = cli.GetState(ctx, storeName, key)
			assert.Nil(t, err)
			assert.Equal(t, key, item.Key)
			assert.Equal(t, []byte(nil), item.Value)
		}).
		Case("Test bulkstate", func(t *testing.T) {
			key := "test_bulkstate_key"

			item := &client.SetStateItem{
				Key:   key,
				Value: value,
			}

			// save
			err = cli.SaveBulkState(ctx, storeName, item)
			assert.Nil(t, err)

			// get
			items, err := cli.GetBulkState(ctx, storeName, []string{key}, nil, 3)
			assert.Nil(t, err)
			assert.Len(t, items, 1)
			assert.Equal(t, key, items[0].Key)
			assert.Equal(t, value, items[0].Value)
		}).
		Run(t)
}
