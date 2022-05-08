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

package persistence

import (
	"os"
	"path"
	"testing"

	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/filter/network/tcpcopy/strategy"
)

func InitForTest() {
	getLogPath = getPathForTest
}

func getPathForTest(fileName string) string {
	logFolder, _ := os.Getwd()
	logFolder = path.Join(logFolder, "logs/mosn")
	logPath := logFolder + string(os.PathSeparator) + fileName
	return logPath
}

func TestIsPersistence_switch_off(t *testing.T) {
	InitForTest()
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	tests := []struct {
		name string
		want bool
	}{
		{
			name: "test switch off",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPersistence(); got != tt.want {
				t.Errorf("IsPersistence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPersistence_not_in_sample_duration(t *testing.T) {
	InitForTest()
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 0

	got := IsPersistence()

	want := false
	if got != want {
		t.Errorf("IsPersistence() = %v, want %v", got, want)
	}
}

func TestIsPersistence_cpu_exceed_max_rate(t *testing.T) {
	InitForTest()
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpCpuMaxRate = 0

	got := IsPersistence()

	want := false
	if got != want {
		t.Errorf("IsPersistence() = %v, want %v", got, want)
	}
}

func TestIsPersistence_mem_exceed_max_rate(t *testing.T) {
	InitForTest()
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpMemMaxRate = 0

	got := IsPersistence()

	want := false
	if got != want {
		t.Errorf("IsPersistence() = %v, want %v", got, want)
	}
}

func TestIsPersistence_success(t *testing.T) {
	InitForTest()
	log.DefaultLogger.SetLogLevel(log.DEBUG)

	strategy.DumpSwitch = true
	strategy.DumpSampleFlag = 1
	strategy.DumpCpuMaxRate = 100
	strategy.DumpMemMaxRate = 100

	got := IsPersistence()

	want := true
	if got != want {
		t.Errorf("IsPersistence() = %v, want %v", got, want)
	}
}
