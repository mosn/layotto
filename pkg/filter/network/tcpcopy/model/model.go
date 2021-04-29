package model

import (
	_type "github.com/layotto/layotto/pkg/filter/network/tcpcopy/type"
)

const (
	AlertDumpKey = "DUMP"
	LogDumpKey   = "[DUMP]"
)

type DumpConfig struct {
	Switch     string  `json:"switch"`       // dump switch.'ON' or 'OFF'
	Interval   int     `json:"interval"`     // dump sampling interval, unit: second
	Duration   int     `json:"duration"`     // Single sampling duration,unit: second
	CpuMaxRate float64 `json:"cpu_max_rate"` // cpu max rate.When cpu rate bigger than this threshold,dump function will be fused
	MemMaxRate float64 `json:"mem_max_rate"` // mem max rate.When memory rate bigger than this threshold,dump function will be fused
}

type DumpUploadDynamicConfig struct {
	Unique_sample_window string             // Specific sampling window
	BusinessType         _type.BusinessType // business type
	Port                 string             // Port
	Binary_flow_data     []byte             // Binary data
	Portrait_data        string             // Portrait data reported by users
}

func NewDumpUploadDynamicConfig(unique_sample_window string, businessType _type.BusinessType, port string, binary_flow_data []byte, portrait_data string) *DumpUploadDynamicConfig {
	dynamicConfig := &DumpUploadDynamicConfig{
		Unique_sample_window: unique_sample_window,
		BusinessType:         businessType,
		Port:                 port,
		Binary_flow_data:     binary_flow_data,
		Portrait_data:        portrait_data,
	}

	return dynamicConfig
}
