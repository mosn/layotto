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

package strategy

import (
	"testing"

	"mosn.io/pkg/log"
)

func TestUpdateAppDumpConfig_invalid_value(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty value",
			args: struct{ value string }{value: ""},
			want: false,
		},
		{
			name: "invalid value",
			args: struct{ value string }{value: "not a json string"},
			want: false,
		},
		{
			name: "invalid switch value",
			args: struct{ value string }{value: "{\"switch\":\"test\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "less than min interval value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":0,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "larger than max interval value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":3601,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "duration value equals 0",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":0,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "duration value larger than interval value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":40,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "cpu max rate equals 0",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":0,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "cpu max rate larger than 100",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":101,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "mem max rate equals 0",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":0}"},
			want: false,
		},
		{
			name: "mem max rate larger than 100",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":101}"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateAppDumpConfig(tt.args.value)
			if result != tt.want {
				t.Errorf("UpdateAppDumpConfig() result = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestUpdateGlobalDumpConfig_invalid_value(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty value",
			args: struct{ value string }{value: ""},
			want: false,
		},
		{
			name: "invalid value",
			args: struct{ value string }{value: "not a json string"},
			want: false,
		},
		{
			name: "invalid switch value",
			args: struct{ value string }{value: "{\"switch\":\"test\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "less than min interval value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":0,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "larger than max interval value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":3601,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "duration value equals 0",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":0,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "duration value larger than interval value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":40,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "cpu max rate equals 0",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":0,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "cpu max rate larger than 100",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":101,\"mem_max_rate\":80}"},
			want: false,
		},
		{
			name: "mem max rate equals 0",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":0}"},
			want: false,
		},
		{
			name: "mem max rate larger than 100",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":101}"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateGlobalDumpConfig(tt.args.value)
			if result != tt.want {
				t.Errorf("UpdateAppDumpConfig() result = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestUpdateAppDumpConfig_success(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: true,
		},
		{
			name: "normal value",
			args: struct{ value string }{value: "{\"switch\":\"OFF\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateAppDumpConfig(tt.args.value)
			if result != tt.want {
				t.Errorf("UpdateAppDumpConfig() result = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestUpdateGlobalDumpConfig_success(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal value",
			args: struct{ value string }{value: "{\"switch\":\"ON\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: true,
		},
		{
			name: "normal value",
			args: struct{ value string }{value: "{\"switch\":\"OFF\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: true,
		},
		{
			name: "normal value",
			args: struct{ value string }{value: "{\"switch\":\"FORCE_OFF\",\"interval\":30,\"duration\":10,\"cpu_max_rate\":80,\"mem_max_rate\":80}"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateGlobalDumpConfig(tt.args.value)
			if result != tt.want {
				t.Errorf("UpdateAppDumpConfig() result = %v, want %v", result, tt.want)
			}
		})
	}
}

func Test_isDumpSwitchOpen(t *testing.T) {
	globalDumpConfig.Switch = "OFF"
	appDumpConfig.Switch = "invalid"

	tests := []struct {
		name string
		want bool
	}{
		{
			name: "test isDumpSwitchOpen",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDumpSwitchOpen(); got != tt.want {
				t.Errorf("isDumpSwitchOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}
