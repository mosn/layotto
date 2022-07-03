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
)

type LocalFileWrapper struct {
	prefix string
}

func NewLocalFileWrapper() Wrapper {
	return &LocalFileWrapper{}
}

func (w *LocalFileWrapper) Init() error {
	w.prefix = getPrefixConfigFilePath()
	return nil
}

func (w *LocalFileWrapper) Wrap(metadata map[string]string) error {
	metadata["secretsFile"] = w.prefix + metadata["secretsFile"]
	return nil
}

func getPrefixConfigFilePath() string {
	prefix := ""
	for i, str := range os.Args {
		if str == "-c" {
			strs := strings.Split(os.Args[i+1], "/")
			for _, s := range strs[:len(strs)-1] {
				prefix = prefix + s + "/"
			}
			break
		}
	}
	return prefix
}
