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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileSize(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "TestGetFileSize",
			args: struct{ f string }{f: "/home/admin/logs/mosn/default.log"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileSize(tt.args.f); got != tt.want {
				t.Errorf("GetFileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveExt(t *testing.T) {
	t.Run("remove ext should remove file extension when it has one", func(t *testing.T) {
		assert.Equal(t, RemoveExt("a.sock"), "a")
	})
	t.Run("remove ext should not change file name when it has no extension", func(t *testing.T) {
		assert.Equal(t, RemoveExt("a"), "a")
	})
}
