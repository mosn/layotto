package common

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"runtime/debug"
	"time"
)

func GetSystemUsageRate() (cpuRate float64, memRate float64, err error) {
	defer func() {
		if e := recover(); e != nil {
			cpuRate, memRate, err = 0, 0, fmt.Errorf("failed to get system usage, err msg: %v, stack: %s", e, debug.Stack())
		}
	}()

	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	if vm == nil {
		return 0, 0, fmt.Errorf("virtual memory info return nil")
	}

	cp, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, 0, err
	}
	if len(cp) != 1 {
		return 0, 0, fmt.Errorf("cpu used info return invalid")
	}

	return cp[0], vm.UsedPercent, nil
}
