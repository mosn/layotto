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
	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/common"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
)

// Whether it has been fused
func IsAvaliable() (ava bool) {

	cpuRate, memRate, err := common.GetSystemUsageRate()
	if err != nil {
		log.DefaultLogger.Errorf(model.AlertDumpKey + " failed to get system usage rate info.")
		return false
	}

	if cpuRate < DumpCpuMaxRate && memRate < DumpMemMaxRate {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s cpuRate:%s is less than max rate %s, memRate:%s is less than max rate %s", model.LogDumpKey, cpuRate, memRate, DumpCpuMaxRate, DumpMemMaxRate)
		}
		return true
	}

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("%s cpuRate:%s, memRate:%s, one or both of them are larger than max rate. Max cpu rate %s. Max mem rate %s", model.LogDumpKey, cpuRate, memRate, DumpCpuMaxRate, DumpMemMaxRate)
	}
	return false
}
