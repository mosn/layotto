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
	"os/user"
	"path"
	"runtime"
)

func GetLogPath(fileName string) string {
	var logFolder string
	if u, err := user.Current(); err != nil {
		logFolder = "/home/admin/logs/mosn"
	} else if runtime.GOOS == "darwin" {
		logFolder = path.Join(u.HomeDir, "logs/mosn")
	} else if runtime.GOOS == "windows" {
		logFolder = path.Join(u.HomeDir, "logs/mosn")
	} else {
		logFolder = "/home/admin/logs/mosn"
	}

	logPath := logFolder + string(os.PathSeparator) + fileName
	return logPath
}
