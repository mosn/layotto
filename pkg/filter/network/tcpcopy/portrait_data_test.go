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

package tcpcopy

import (
	"testing"

	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/filter/network/tcpcopy/strategy"
	_type "mosn.io/layotto/pkg/filter/network/tcpcopy/type"
)

func Test_isHandle_switch_off(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = false

	want := false
	if got := isHandle(_type.RPC); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func Test_isHandle_mutex(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	getAndSwapDumpBusinessCache(_type.RPC, 1)

	want := false
	if got := isHandle(_type.RPC); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func Test_isHandle_normal(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpCpuMaxRate = 200
	strategy.DumpMemMaxRate = 200

	want := true
	if got := isHandle(_type.CONFIGURATION); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func Test_getAndSwapDumpBusinessCache(t *testing.T) {
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	getAndSwapDumpBusinessCache(_type.RPC, 0)
	want := 0
	if got := getAndSwapDumpBusinessCache(_type.RPC, 1); got != want {
		t.Errorf("isHandle() = %v, want %v", got, want)
	}
}

func TestUploadPortraitData_switch_off(t *testing.T) {

	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = false

	want := false
	if got := UploadPortraitData(_type.RPC, nil, nil); got != want {
		t.Errorf("UploadPortraitData() = %v, want %v", got, want)
	}
}

func TestUploadPortraitData_normal(t *testing.T) {

	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpCpuMaxRate = 200
	strategy.DumpMemMaxRate = 200
	data := "{\"key\":\"value\"}"

	want := true
	if got := UploadPortraitData(_type.STATE, data, nil); got != want {
		t.Errorf("UploadPortraitData() = %v, want %v", got, want)
	}
}
