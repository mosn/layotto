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

// map keys
const (
	namespaceIdKey = "namespace_id"
	appNameKey     = "app_name"
	userNameKey    = "username"
	passwordKey    = "password"
	endPointKey    = "end_point"
	regionIdKey    = "region_id"
	accessKey      = "access_key"
	secretKey      = "secret_key"
)

type Metadata struct {
	AppName     string
	NameSpaceId string
	Username    string
	Password    string
	// ACM & KMS
	Endpoint  string
	RegionId  string
	AccessKey string
	SecretKey string
	OpenKMS   bool
}

func ParseNacosMetadata(properties map[string]string) (*Metadata, error) {
	if properties == nil {
		return nil, errConfigMissingField("metadata")
	}

	config := &Metadata{}
	config.AppName = properties[appNameKey]
	if config.AppName == "" {
		return nil, errConfigMissingField(appNameKey)
	}

	// the namespace of config, not required
	config.NameSpaceId = properties[namespaceIdKey]
	if config.NameSpaceId == "" {
		config.NameSpaceId = defaultNamespaceId
	}

	if v, ok := properties[userNameKey]; ok && v != "" {
		config.Username = v
	}

	if v, ok := properties[passwordKey]; ok && v != "" {
		config.Password = v
	}

	// ACM & KMS
	if v, ok := properties[endPointKey]; ok && v != "" {
		config.Endpoint = v
		config.OpenKMS = true
	}

	if v, ok := properties[regionIdKey]; ok && v != "" {
		config.RegionId = v
		config.OpenKMS = true
	}

	if v, ok := properties[accessKey]; ok && v != "" {
		config.AccessKey = v
		config.OpenKMS = true
	}

	if v, ok := properties[secretKey]; ok && v != "" {
		config.SecretKey = v
		config.OpenKMS = true
	}

	return config, nil
}
