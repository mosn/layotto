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
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/proxy-wasm-go-host/proxywasm/common"
)

func TestMapEncodeAndDecode(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	b := common.EncodeMap(m)
	n := common.DecodeMap(b)
	assert.Equal(t, m, n)
}
