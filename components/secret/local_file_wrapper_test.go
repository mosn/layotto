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
	"testing"

	"github.com/stretchr/testify/assert"
)

var metadata = make(map[string]string)

func TestLocalFileWrapper(t *testing.T) {
	initConfigFilePath()
	initMetaData()
	wrapper := NewLocalFileWrapper()
	wrapper.Init()
	wrapper.Wrap(metadata)
	assert.Equal(t, "../../configs/secret_local_file.json", metadata["secretsFile"])
}

func initConfigFilePath() {
	os.Args = append(os.Args, "-c", "../../configs/config_standalone.json")
}

func initMetaData() {
	metadata["secretsFile"] = "secret_local_file.json"
}
