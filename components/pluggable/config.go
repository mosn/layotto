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
	"mosn.io/layotto/components/ref"
	pb "mosn.io/layotto/spec/proto/pluggable/v1/common"
)

func ToProtoConfig(config ref.Config) *pb.Config {
	return &pb.Config{
		SecretRef:    ToProtoSecretRef(config.SecretRef),
		ComponentRef: ToProtoComponentRef(config.ComponentRef),
	}
}

func ToProtoSecretRef(secrets []*ref.SecretRefConfig) []*pb.SecretRefConfig {
	if secrets == nil {
		return nil
	}
	res := make([]*pb.SecretRefConfig, len(secrets))
	for i, s := range secrets {
		res[i] = &pb.SecretRefConfig{
			StoreName: s.StoreName,
			Key:       s.Key,
			SubKey:    s.SubKey,
			InjectAs:  s.InjectAs,
		}
	}
	return res
}

func ToProtoComponentRef(config *ref.ComponentRefConfig) *pb.ComponentRefConfig {
	if config == nil {
		return nil
	}
	return &pb.ComponentRefConfig{
		SecretStore: config.SecretStore,
		ConfigStore: config.ConfigStore,
	}
}
