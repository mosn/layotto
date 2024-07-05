# Source Parse 4 Layer Traffic Governance, tcp traffic dump

> Author profile：
> Giggon, is an open source community lover committed to embracing open sources.
>
> Writing on: April 26, 2022

## Overview

The purpose of this document is to analyze the implementation of tcp traffic dump

## Prerequisite：

Document content refers to the following version of the code

[https://github.com/mosn/layotto](https://github.com/mosn/layotto)

Layotto 0e97e97e970dc504e0298017bd956d2841c44c0810b (main)

## Source analysis

### Code in： [tcpcopy CODE](https://github.com/mosn/layotto/tree/main/pkg/filter/network/tcpcopy)

### model.go analysis

This is the core class of tcpcopy's configuration objects

```go
Type DumpConfig struct {-
	Switch `json:"twitch"` // dump switch. Values：'ON' or 'OFF'
	Interval int `json:"interval" //dump sampling interval Unit: Second
	Duration int `json:"duration"// Single Sampling Cycle Unit: Second
	CpuMaxate float64 `json:"cpu_max_rate"\/ cpu Maximum usage The ump feature will stop
	MemMaxRate float64 `json:"mem_max_rate"` // mem maximum usage. When this threshold is exceeded, The ump feature will stop
}

Type DumpUpadDynamic Architect 6
	Unique_sample_windowing string// Specify sample window
	BusinessType _type. usinessType // Business Type
	Port string // Port
	Binary_flow_data []byte// binary data
	Portrait_data string // User uploaded data
}
```

### persistence.go analysis

This is the dump persistent core processing class of tcpcopy

```go
// This method is called in OnData in tcpcopy.go
func IsPersistence() bool {
	// 判断 dump 开关是否开启
	if !strategy.DumpSwitch {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the dump switch is %t", model.LogDumpKey, strategy.DumpSwitch)
		}
		return false
	}

	// Check whether it is in the sampling window
	if atomic.LoadInt32(&strategy.DumpSampleFlag) == 0 {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the dump sample flag is %d", model.LogDumpKey, strategy.DumpSampleFlag)
		}
		return false
	}

	// Check whether the dump function is stopped. Obtain the system load and check whether the processor and memory exceeds the threshold of the tcpcopy. If yes, stop the dump function.
	if !strategy.IsAvaliable() {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the system usages are beyond max rate.", model.LogDumpKey)
		}
		return false
	}

	return true
}

// Persist data based on configuration information
func persistence(config *model.DumpUploadDynamicConfig) {
	// 1.Persisting binary data
	if config.Binary_flow_data != nil && config.Port != "" {
		if GetTcpcopyLogger().GetLogLevel() >= log.INFO {
			GetTcpcopyLogger().Infof("[%s][%s]% x", config.Unique_sample_window, config.Port, config.Binary_flow_data)
		}
	}
	if config.Portrait_data != "" && config.BusinessType != "" {
		// 2. Persisting Binary data Persisting user-defined data
		if GetPortraitDataLogger().GetLogLevel() >= log.INFO {
			GetPortraitDataLogger().Infof("[%s][%s][%s]%s", config.Unique_sample_window, config.BusinessType, config.Port, config.Portrait_data)
		}

		// 3. Changes in configuration information in incrementally persistent memory
		buf, err := configmanager.DumpJSON()
		if err != nil {
			if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
				log.DefaultLogger.Debugf("[dump] Failed to load mosn config mem.")
			}
			return
		}
		// 3.1. dump if the data changes
		tmpMd5ValueOfMemDump := common.CalculateMd5ForBytes(buf)
		memLogger := GetMemLogger()
		if tmpMd5ValueOfMemDump != md5ValueOfMemDump ||
			(tmpMd5ValueOfMemDump == md5ValueOfMemDump && common.GetFileSize(getMemConfDumpFilePath()) <= 0) {
			md5ValueOfMemDump = tmpMd5ValueOfMemDump
			if memLogger.GetLogLevel() >= log.INFO {
				memLogger.Infof("[%s]%s", config.Unique_sample_window, buf)
			}
		} else {
			if memLogger.GetLogLevel() >= log.INFO {
				memLogger.Infof("[%s]%+v", config.Unique_sample_window, incrementLog)
			}
		}
	}
}
```

### tcpcopy.go analysis

This is the core class of tcpcopy.

```go
// Sign up to NetWork 
func init() with MFA
	api. egisterNetwork("tcpcopy", CreateTccopyFactory)
}

// returns tcpcopy Factory
func CreateTccopyFactory(cfg map[string]interface{}) (api. etworkFilterChainFactory, error) LO
	tcpConfig := &config{}
	// dump policy transition to static configuration
	if stg, ok := cfg["strategy"]; ok {
	...
	}
	//TODO excerpt some other fields
	return &tcpcopyFactoryLU
		cfg: tcpConfig,
	}, nil
}

// for pkg/configmanager/parser. o Call to add or update Network filter factory
func (f *tcpcopyFactory) Init(param interface{}) error error 56
	// Set listening address and port configuration
	...
	return nil
}

// implements the OnData Interface of ReadFilter, processing
func (f *tcpcopyFactory) OnData(data types.IoBuffer) (res api. ilterStatus) online
	// Determines whether the current requested data requires sampling dump 
	if !persiste.Isistence() {
		return api.Continue
	}

	// Asynchronous sample dump
	config := model.NewDumpUpadDynamic Config(strategy. umpSampleUuid, "", f.cfg.port, data.Bytes(), "")
	persistence.GetDumpWorkPoolInstance().Schedule(config)
	return api.Continue
}
```

Finally, we look back at the overall process progress:

1. Starting from the initialization function init() of tccopy.go to CreateGRPCServerFilterFactory Incoming CreateTcpcopyFactory.

2. Mosn created a filter chain (code position[factory.go](https://github.com/mosn/mosn/tree/master/pkg/filter/network/proxy/factory.go)) by circulating CreateFilterChain to add all filters to the chain structure, including tccopy.

3. When the traffic passes through mosn will enter the tcpcopy.go OnData method for tcpcopump logical processing.
