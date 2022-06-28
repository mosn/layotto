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

package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerToString(t *testing.T) {
	assert.Equal(t, PointerToString(nil), "")
	s := ""
	assert.Equal(t, PointerToString(&s), "")
}

func TestGetPrefixConfigFilePathNormalCase(t *testing.T) {
	os.Args = append(os.Args, "-c", "../../configs/config_standalone.json")
	assert.Equal(t, "../../configs/", GetPrefixConfigFilePath())
}

func TestGetPrefixConfigFilePathSimpleCase(t *testing.T) {
	os.Args = append(os.Args, "-c", "./config_standalone.json")
	assert.Equal(t, "./", GetPrefixConfigFilePath())
}
