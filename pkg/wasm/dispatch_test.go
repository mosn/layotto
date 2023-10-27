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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	wasmPluginName1 = "test1"
	wasmPluginName2 = "test2"

	idValid   = "wasm_test1"
	idInvalid = "wasm_test2"
)

func mockRouters() map[string]*Group {
	wasmPlugin := &WasmPlugin{
		pluginName: wasmPluginName1,
	}
	group := &Group{
		count:   1,
		plugins: []*WasmPlugin{wasmPlugin},
	}

	return map[string]*Group{
		idValid: group,
	}
}

func TestRouter_GetRandomPluginByID(t *testing.T) {
	type fields struct {
		routes map[string]*Group
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *WasmPlugin
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			fields: fields{
				routes: mockRouters(),
			},
			args: args{
				id: idValid,
			},
			want: &WasmPlugin{
				pluginName: wasmPluginName1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "not found",
			fields: fields{
				routes: mockRouters(),
			},
			args: args{
				id: idInvalid,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "empty",
			fields: fields{
				routes: make(map[string]*Group),
			},
			args: args{
				id: idInvalid,
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := &Router{
				routes: tt.fields.routes,
			}
			got, err := route.GetRandomPluginByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRandomPluginByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRandomPluginByID(%v)", tt.args.id)
		})
	}
}

func TestRouter_RegisterRoute(t *testing.T) {
	type fields struct {
		routes map[string]*Group
	}
	type args struct {
		id     string
		plugin *WasmPlugin
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		groupCnt int
	}{
		{
			name: "add",
			fields: fields{
				routes: mockRouters(),
			},
			args: args{
				id: idValid,
				plugin: &WasmPlugin{
					pluginName: wasmPluginName2,
				},
			},
			groupCnt: 2,
		},
		{
			name: "replace",
			fields: fields{
				routes: mockRouters(),
			},
			args: args{
				id: idValid,
				plugin: &WasmPlugin{
					pluginName: wasmPluginName1,
				},
			},
			groupCnt: 1,
		},
		{
			name: "empty",
			fields: fields{
				routes: make(map[string]*Group),
			},
			args: args{
				id: idValid,
				plugin: &WasmPlugin{
					pluginName: wasmPluginName1,
				},
			},
			groupCnt: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := &Router{
				routes: tt.fields.routes,
			}
			route.RegisterRoute(tt.args.id, tt.args.plugin)
			_, err := route.GetRandomPluginByID(tt.args.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.groupCnt, route.routes[tt.args.id].count)
		})
	}
}

func TestRouter_RemoveRoute(t *testing.T) {
	type fields struct {
		routes map[string]*Group
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "normal",
			fields: fields{
				routes: mockRouters(),
			},
			args: args{
				id: idValid,
			},
		},
		{
			name: "not exist",
			fields: fields{
				routes: mockRouters(),
			},
			args: args{
				id: idInvalid,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := &Router{
				routes: tt.fields.routes,
			}
			route.RemoveRoute(tt.args.id)
		})
	}
}
