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

package nacos

const (
	defaultNamespaceId = "" // if this is not set, then nacos will use the default namespaceId.
	defaultGroup       = "default"
	defaultLabel       = "default"
	defaultLogDir      = "/tmp/layotto/nacos/log"
	defaultCacheDir    = "/tmp/layotto/nacos/cache"
	defaultLogLevel    = "debug"
	defaultTimeout     = 10 // second
)

type NacosMetadata struct {
	AppName string
}

func ParseNacosMetadata(properties map[string]string) (*NacosMetadata, error) {
	config := &NacosMetadata{}
	config.AppName = properties["app_name"]
	if config.AppName == "" {
		return nil, errConfigMissingField("app_name")
	}

	return config, nil
}
