// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluggable

import (
	"reflect"
	"testing"

	"mosn.io/layotto/components/ref"
	pb "mosn.io/layotto/spec/proto/pluggable/v1/common"
)

func TestToProtoConfig(t *testing.T) {
	type args struct {
		config ref.Config
	}
	tests := []struct {
		name string
		args args
		want *pb.Config
	}{
		{
			name: "test normal params",
			args: args{
				config: ref.Config{
					SecretRef: []*ref.SecretRefConfig{
						{
							StoreName: "store",
							Key:       "key",
							SubKey:    "subkey",
							InjectAs:  "injectas",
						},
					},
					ComponentRef: &ref.ComponentRefConfig{
						SecretStore: "secret",
						ConfigStore: "config",
					},
				},
			},
			want: &pb.Config{
				SecretRef: []*pb.SecretRefConfig{
					{
						StoreName: "store",
						Key:       "key",
						SubKey:    "subkey",
						InjectAs:  "injectas",
					},
				},
				ComponentRef: &pb.ComponentRefConfig{
					SecretStore: "secret",
					ConfigStore: "config",
				},
			},
		},
		{
			name: "with nil secret config and nil component config",
			args: args{
				config: ref.Config{},
			},
			want: &pb.Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToProtoConfig(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToProtoConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToProtoSecretRef(t *testing.T) {
	type args struct {
		secrets []*ref.SecretRefConfig
	}
	tests := []struct {
		name string
		args args
		want []*pb.SecretRefConfig
	}{
		{
			name: "with normal params input",
			args: args{
				secrets: []*ref.SecretRefConfig{
					{
						StoreName: "store",
						Key:       "key",
						SubKey:    "subkey",
						InjectAs:  "injectas",
					},
				},
			},
			want: []*pb.SecretRefConfig{
				{
					StoreName: "store",
					Key:       "key",
					SubKey:    "subkey",
					InjectAs:  "injectas",
				},
			},
		},
		{
			name: "with nil input params",
			args: args{secrets: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToProtoSecretRef(tt.args.secrets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToProtoSecretRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToProtoComponentRef(t *testing.T) {
	type args struct {
		config *ref.ComponentRefConfig
	}
	tests := []struct {
		name string
		args args
		want *pb.ComponentRefConfig
	}{
		{
			name: "test normal config input",
			args: args{
				config: &ref.ComponentRefConfig{
					SecretStore: "secret",
					ConfigStore: "config",
				},
			},
			want: &pb.ComponentRefConfig{
				SecretStore: "secret",
				ConfigStore: "config",
			},
		},
		{
			name: "test with nil params",
			args: args{config: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToProtoComponentRef(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToProtoComponentRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
