package strategy

import (
	"mosn.io/layotto/pkg/common"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
	"mosn.io/pkg/log"
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
