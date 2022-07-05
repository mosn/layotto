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

package secret

import (
	"os"
	"strings"

	"github.com/dapr/components-contrib/secretstores"
)

type wrapper struct {
	prefix string
	comp   secretstores.SecretStore
}

func (w *wrapper) Init(metadata secretstores.Metadata) error {
	if metadata.Properties != nil {
		metadata.Properties["secretsFile"] = w.prefix + metadata.Properties["secretsFile"]
	}
	return w.comp.Init(metadata)
}

func (w *wrapper) GetSecret(req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	return w.comp.GetSecret(req)
}

func (w *wrapper) BulkGetSecret(req secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	return w.comp.BulkGetSecret(req)
}

func Wrap(component secretstores.SecretStore) secretstores.SecretStore {
	return &wrapper{
		comp:   component,
		prefix: getPrefixConfigFilePath(),
	}
}

func getPrefixConfigFilePath() string {
	prefix := ""
	for i, str := range os.Args {
		if str == "-c" || str == "--config" {
			strs := strings.Split(os.Args[i+1], string(os.PathSeparator))
			for _, s := range strs[:len(strs)-1] {
				prefix = prefix + s + string(os.PathSeparator)
			}
			break
		}
	}
	return prefix
}
