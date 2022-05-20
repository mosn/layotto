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

package actuator

import (
	"sync/atomic"
	"time"

	"mosn.io/layotto/pkg/actuator/info"
)

type AppContributor info.Contributor

func GetAppContributor() AppContributor {
	return info.ContributorAdapter(func() (interface{}, error) {
		return GetAppInfoSingleton(), nil
	})
}

var appInfoSingleton atomic.Value

type AppInfo struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string `json:"name"`
	// Version of the program
	Version string `json:"version"`
	// Compilation date
	Compiled time.Time `json:"compiled,omitempty"`
}

func NewAppInfo() *AppInfo {
	return &AppInfo{}
}

func GetAppInfoSingleton() AppInfo {
	info, ok := appInfoSingleton.Load().(AppInfo)
	if ok {
		return info
	}
	return AppInfo{}
}

func SetAppInfoSingleton(a *AppInfo) {
	if a == nil {
		appInfoSingleton.Store(AppInfo{})
		return
	}
	appInfoSingleton.Store(*a)
}
