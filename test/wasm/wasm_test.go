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

package wasm

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"mosn.io/layotto/components/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	cli := utils.NewRedisClient(utils.RedisMetadata{
		Host: "localhost:6379",
	})
	err := cli.Set(context.Background(), "book1", "100", 0).Err()
	if err != nil {
		t.Fatal("set inventories error")
	}

	ids := []string{"id_1"}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:2045?name=book1", nil)

	for _, id := range ids {
		req.Header.Set("id", id)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed, err: %s", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Read body failed, err: %s", err)
		}
		assert.Equal(t, "There are 100 inventories for book1.", string(body))
	}
}
