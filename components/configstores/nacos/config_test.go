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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNacosMetadata(t *testing.T) {
	// without metadata
	_, err := ParseNacosMetadata(nil)
	assert.Error(t, err)

	properties := make(map[string]string)

	// success
	_, err = ParseNacosMetadata(properties)
	assert.Nil(t, err)

	// test set namespace
	namespaceId := "namespace"
	properties[namespaceIdKey] = namespaceId
	metadata, err := ParseNacosMetadata(properties)
	assert.Nil(t, err)
	assert.Equal(t, namespaceId, metadata.NameSpaceId)
}
