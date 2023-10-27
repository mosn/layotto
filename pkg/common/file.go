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
	"io/fs"
	"os"
	"path/filepath"
)

func GetFileSize(f string) int64 {
	fi, err := os.Stat(f)
	if err == nil {
		return fi.Size()
	}

	return -1
}

// IsSocket returns if the given file is a unix socket.
func IsSocket(f fs.FileInfo) bool {
	return f.Mode()&fs.ModeSocket != 0
}

// RemoveExt removes file extension
func RemoveExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
