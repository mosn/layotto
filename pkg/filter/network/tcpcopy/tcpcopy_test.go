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

package tcpcopy

import (
	"encoding/json"
	"testing"

	"mosn.io/api"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"

	"mosn.io/layotto/pkg/filter/network/tcpcopy/strategy"
)

func TestCreateTcpcopyFactory(t *testing.T) {
	data := "{\"port\":\"12220\"}"
	var rpcCfg map[string]interface{}
	json.Unmarshal([]byte(data), &rpcCfg)

	type args struct {
		cfg map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want api.NetworkFilterChainFactory
	}{
		{
			name: "TestCreateTcpcopyFactory",
			args: struct{ cfg map[string]interface{} }{cfg: nil},
		},
		{
			name: "TestCreateTcpcopyFactory",
			args: struct{ cfg map[string]interface{} }{cfg: rpcCfg},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateTcpcopyFactory(tt.args.cfg)
			if err != nil {
				t.Errorf("CreateTcpcopyFactory() error = %v", err)
				return
			}
		})
	}
}

func Test_tcpcopyFactory_OnData_switch_off(t *testing.T) {
	strategy.DumpSwitch = false

	type fields struct {
		tcpcopy *config
	}
	type args struct {
		data types.IoBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes api.FilterStatus
	}{
		{
			name:    "",
			fields:  struct{ tcpcopy *config }{tcpcopy: nil},
			args:    struct{ data types.IoBuffer }{data: nil},
			wantRes: api.Continue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tcpcopyFactory{
				cfg: tt.fields.tcpcopy,
			}
			if gotRes := f.OnData(tt.args.data); gotRes != tt.wantRes {
				t.Errorf("OnData() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_tcpcopyFactory_OnData_success(t *testing.T) {
	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1

	type fields struct {
		tcpcopy *config
	}
	tcpcopy_value := config{port: "12220"}

	type args struct {
		data types.IoBuffer
	}

	buffer := buffer.NewIoBufferString("test")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes api.FilterStatus
	}{
		{
			name:    "",
			fields:  struct{ tcpcopy *config }{tcpcopy: &tcpcopy_value},
			args:    struct{ data types.IoBuffer }{data: buffer},
			wantRes: api.Continue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tcpcopyFactory{
				cfg: tt.fields.tcpcopy,
			}
			if gotRes := f.OnData(tt.args.data); gotRes != tt.wantRes {
				t.Errorf("OnData() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_tcpcopyFactory_OnNewConnection(t *testing.T) {
	type fields struct {
		tcpcopy *config
	}
	tests := []struct {
		name   string
		fields fields
		want   api.FilterStatus
	}{
		{
			name:   "Test_tcpcopyFactory_OnNewConnection",
			fields: struct{ tcpcopy *config }{tcpcopy: nil},
			want:   api.Continue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &tcpcopyFactory{
				cfg: tt.fields.tcpcopy,
			}
			if got := f.OnNewConnection(); got != tt.want {
				t.Errorf("OnNewConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
