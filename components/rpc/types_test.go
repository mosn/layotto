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

package rpc

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCHeader_Get(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		header := RPCHeader{
			"h1": []string{"v1", "v2"},
			"h2": []string{"v3"},
		}

		value := header.Get("h1")
		assert.Equal(t, "v1,v2", value)
	})

	t.Run("not exist", func(t *testing.T) {
		header := RPCHeader{
			"h1": []string{"v1", "v2"},
			"h2": []string{"v3"},
		}
		value := header.Get("h3")
		assert.Equal(t, "", value)
	})

}

func TestRPCHeader_Range(t *testing.T) {
	header := RPCHeader{
		"h1": []string{"v1", "v2"},
		"h2": []string{"v3"},
	}

	var finalResult []string
	f := func(key string, value string) bool {
		finalResult = append(finalResult, fmt.Sprintf("%s:%s", key, value))
		return true
	}

	header.Range(f)

	sort.Slice(finalResult, func(i, j int) bool {
		return finalResult[i] < finalResult[j]
	})
	assert.Equal(t, "h1:v1,v2|h2:v3", strings.Join(finalResult, "|"))
}
