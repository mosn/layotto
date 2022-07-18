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

package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pathResolver_isEmpty(t *testing.T) {
	p := NewPathResolver("/")
	assert.True(t, !p.HasNext())
	p = NewPathResolver("")
	assert.True(t, !p.HasNext())
	p = NewPathResolver("/a")
	assert.True(t, p.HasNext())
}

func Test_pathResolver_next(t *testing.T) {
	p := NewPathResolver("/a/b/c")
	assert.True(t, p.Next() == "a")
	assert.True(t, p.Next() == "b")
	assert.True(t, p.Next() == "c")
	assert.True(t, !p.HasNext())
	assert.True(t, p.Next() == "")
}

func Test_pathResolver_unresolvedPath(t *testing.T) {
	p := NewPathResolver("/a/b/c")
	assert.True(t, p.Next() == "a")
	assert.True(t, p.UnresolvedPath() == "/b/c")
	assert.True(t, p.Next() == "b")
	assert.True(t, p.UnresolvedPath() == "/c")
	assert.True(t, p.Next() == "c")
	assert.True(t, p.UnresolvedPath() == "")
	assert.True(t, !p.HasNext())
	assert.True(t, p.Next() == "")
	assert.True(t, p.UnresolvedPath() == "")
}
