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
	"fmt"
	"strings"
)

const (
	ETag = "ETag"
)

func GetBucketName(fileName string) (string, error) {
	index := strings.Index(fileName, "/")
	if index == -1 || index == 0 {
		return "", fmt.Errorf("invalid fileName format")
	}
	return fileName[:index], nil
}

func GetFilePrefixName(fileName string) string {
	index := strings.Index(fileName, "/")
	if index == -1 {
		return ""
	}
	return fileName[index+1:]
}

func GetFileName(fileName string) (string, error) {
	index := strings.Index(fileName, "/")
	if index == -1 {
		return "", fmt.Errorf("invalid fileName format")
	}
	name := fileName[index+1:]
	if name == "" {
		return "", fmt.Errorf("file name is empty")
	}
	return name, nil
}
