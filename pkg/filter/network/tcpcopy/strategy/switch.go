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

package strategy

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/utils"

	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
)

const (
	minInterval = 30
	maxInterval = 60 * 60

	defaultCpuMaxRate = 80
	defaultMemMaxRate = 70

	defaultDuration = 1

	kindOn       = "ON"        // ON
	kindOff      = "OFF"       // OFF
	kindForceOff = "FORCE_OFF" // Forced shutdown
)

var (
	appDumpConfig = &model.DumpConfig{
		Switch:     kindOff,
		Interval:   minInterval,
		Duration:   defaultDuration,
		CpuMaxRate: defaultCpuMaxRate,
		MemMaxRate: defaultMemMaxRate,
	}

	globalDumpConfig = &model.DumpConfig{
		Switch:     kindOff,
		Interval:   minInterval,
		Duration:   defaultDuration,
		CpuMaxRate: defaultCpuMaxRate,
		MemMaxRate: defaultMemMaxRate,
	}

	// switch status
	DumpSwitch = true

	// Sampling Flag, 0 means no sampling, 1 means sampling
	DumpSampleFlag int32

	// cpu fuse threshold
	DumpCpuMaxRate float64 = defaultCpuMaxRate

	// mem fuse threshold
	DumpMemMaxRate float64 = defaultMemMaxRate

	// Dump Interval
	DumpInterval = minInterval

	// Single sampling duration
	DumpDuration = defaultDuration

	// Dump uuid
	DumpSampleUuid = "inituuid"

	// Sampling status of different Business
	DumpBusinessCache = new(sync.Map)

	initOnce = new(sync.Once)
)

// For hot reloading app-level dumpConfig
func UpdateAppDumpConfig(value string) bool {
	if value == "" {
		return false
	}

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[dumpConfig] update app dump config, value=%s", value)
	}
	// unmarshal
	var temp model.DumpConfig
	if err := json.Unmarshal([]byte(value), &temp); err != nil {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update app dump config failed, value=%s is illegal.", value)
		return false
	}
	// validate
	if temp.Switch != kindOn && temp.Switch != kindOff {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update app dump config failed, the switch is illegal, value=%s", value)
		return false
	}
	if temp.Interval < minInterval || temp.Interval > maxInterval {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update app dump config failed, the interval should be between %v and %v, value=%s", minInterval, maxInterval, value)
		return false
	}
	if temp.Duration <= 0 || temp.Duration >= temp.Interval {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update app dump config failed, the duration should be between %v and %v, value=%s", 0, temp.Interval, value)
		return false
	}
	if temp.CpuMaxRate <= 0 || temp.CpuMaxRate >= 100 {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update app dump config failed, the cpu_max_rate should be between %v and %v, value=%s", 0, 100, value)
		return false
	}
	if temp.MemMaxRate <= 0 || temp.MemMaxRate >= 100 {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update app dump config failed, the mem_max_rate should be between %v and %v, value=%s", 0, 100, value)
		return false
	}
	// publish config
	appDumpConfig = &temp
	updateDumpConfig()

	return true
}

// For hot reloading global dumpConfig
func UpdateGlobalDumpConfig(value string) bool {
	if "" == value {
		return false
	}
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[dumpConfig] update global dump config, value=%s", value)
	}
	// unmarshal
	var temp model.DumpConfig
	if err := json.Unmarshal([]byte(value), &temp); err != nil {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update global dump config failed, value=%s is illegal.", value)
		return false
	}
	// validate
	if temp.Switch != kindOn && temp.Switch != kindOff && temp.Switch != kindForceOff {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update global dump config failed, the switch is illegal, value=%s", value)
		return false
	}

	if temp.Interval < minInterval || temp.Interval > maxInterval {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update global dump config failed, the interval should be between %v and %v, value=%s", minInterval, maxInterval, value)
		return false
	}
	if temp.Duration <= 0 || temp.Duration >= temp.Interval {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update global dump config failed, the duration should be between %v and %v, value=%s", 0, temp.Interval, value)
		return false
	}
	if temp.CpuMaxRate <= 0 || temp.CpuMaxRate >= 100 {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update global dump config failed, the cpu_max_rate should be between %v and %v, value=%s", 0, 100, value)
		return false
	}
	if temp.MemMaxRate <= 0 || temp.MemMaxRate >= 100 {
		log.DefaultLogger.Alertf("dump", "[dumpConfig] update global dump config failed, the mem_max_rate should be between %v and %v, value=%s", 0, 100, value)
		return false
	}
	// publish config
	globalDumpConfig = &temp
	updateDumpConfig()

	return true
}

func updateDumpConfig() {
	DumpSwitch = isDumpSwitchOpen()
	DumpCpuMaxRate = getDumpCpuMaxRate()
	DumpMemMaxRate = getDumpMemMaxRate()
	DumpInterval = getDumpInterval()
	DumpDuration = getDumpDuration()

	if DumpSwitch {
		initOnce.Do(func() {
			if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
				log.DefaultLogger.Debugf("%s start updateSampleFlag.", model.LogDumpKey)
			}
			utils.GoWithRecover(updateSampleFlag, nil)
		})
	}
}

func isDumpSwitchOpen() bool {
	global := globalDumpConfig.Switch
	if global == kindForceOff {
		return false
	}

	app := appDumpConfig.Switch
	if app == kindOff {
		return global == kindOn
	}
	return app == kindOn
}

func getDumpInterval() int {
	global := globalDumpConfig.Switch
	if global == kindForceOff {
		return maxInterval
	}

	app := appDumpConfig.Switch
	if app == kindOn {
		return appDumpConfig.Interval
	}
	return globalDumpConfig.Interval
}

func getDumpDuration() int {
	global := globalDumpConfig.Switch
	if global == kindForceOff {
		return defaultDuration
	}

	app := appDumpConfig.Switch
	if app == kindOn {
		return appDumpConfig.Duration
	}
	return globalDumpConfig.Duration
}

func getDumpCpuMaxRate() float64 {
	global := globalDumpConfig.Switch
	if global == kindForceOff {
		return defaultCpuMaxRate
	}

	app := appDumpConfig.Switch
	if app == kindOn {
		return appDumpConfig.CpuMaxRate
	}
	return globalDumpConfig.CpuMaxRate
}

func getDumpMemMaxRate() float64 {
	global := globalDumpConfig.Switch
	if global == kindForceOff {
		return defaultMemMaxRate
	}

	app := appDumpConfig.Switch
	if app == kindOn {
		return appDumpConfig.MemMaxRate
	}
	return globalDumpConfig.MemMaxRate
}

func updateSampleFlag() {
	for {
		// Default sampling interval is 30s
		st := time.Duration(DumpInterval) * time.Second
		time.Sleep(st)

		// Update the sampling flag
		if IsAvaliable() {
			if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
				log.DefaultLogger.Debugf("%s open sample window", model.LogDumpKey)
			}
			DumpSampleUuid = utils.GenerateUUID()
			atomic.StoreInt32(&DumpSampleFlag, 1)
			// Continuous sampling
			dst := time.Duration(DumpDuration) * time.Second
			time.Sleep(dst)
		}
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s close sample window", model.LogDumpKey)
		}
		atomic.StoreInt32(&DumpSampleFlag, 0)
		// Send a sampling token (reset the counter) to each business. 0 means the number of samplings in this sampling period is 0, i.e., available for sampling
		DumpBusinessCache.Range(func(key, value interface{}) bool {
			DumpBusinessCache.Store(key, 0)
			return true
		})
	}
}
