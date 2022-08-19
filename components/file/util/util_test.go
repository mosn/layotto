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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBucketName(t *testing.T) {
	_, err := GetBucketName("/")
	assert.Equal(t, "invalid fileName format", err.Error())
	_, err = GetBucketName("")
	assert.Equal(t, "invalid fileName format", err.Error())
	_, err = GetBucketName("bucketName")
	assert.Equal(t, "invalid fileName format", err.Error())
	name, err := GetBucketName("bucketName/")
	assert.Nil(t, err)
	assert.Equal(t, name, "bucketName")
}

func TestGetFileName(t *testing.T) {
	name, err := GetFileName("/")
	assert.Equal(t, err.Error(), "file name is empty")
	assert.Equal(t, "", name)
	name, err = GetFileName("aaa")
	assert.Equal(t, err.Error(), "invalid fileName format")
	assert.Equal(t, "", name)
	name, err = GetFileName("/aaa")
	assert.Nil(t, err)
	assert.Equal(t, "aaa", name)
}
