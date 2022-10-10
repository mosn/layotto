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

package ref

// Config is ref json config
type Config struct {
	SecretRef    []*SecretRefConfig  `json:"secret_ref"`
	ComponentRef *ComponentRefConfig `json:"component_ref"`
}

type SecretRefConfig struct {
	//secret component name, such as : local.file
	StoreName string `json:"store_name"`
	// key in the secret component
	Key string `json:"key"`
	//sub key in the secret component
	SubKey string `json:"sub_key"`
	// key need to inject into metadata
	InjectAs string `json:"inject_as"`
}

type ComponentRefConfig struct {
	SecretStore string `json:"secret_store"`
	ConfigStore string `json:"config_store"`
}
