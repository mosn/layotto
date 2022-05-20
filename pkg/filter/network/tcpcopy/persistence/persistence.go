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
	"sync"
	"sync/atomic"

	"mosn.io/layotto/pkg/common"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/strategy"

	"mosn.io/mosn/pkg/configmanager"
	"mosn.io/mosn/pkg/log"
	rlog "mosn.io/pkg/log"
)

const (
	dumpBasePath         = "dump"
	tcpcopyDumpFile      = dumpBasePath + string(os.PathSeparator) + "dump_tcp_copy.log"
	memConfDumpFile      = dumpBasePath + string(os.PathSeparator) + "dump_mem_dump.log"
	staticConfDumpFile   = dumpBasePath + string(os.PathSeparator) + "dump_static_conf.log"
	portraitDataDumpFile = dumpBasePath + string(os.PathSeparator) + "dump_portrait_data.log"

	incrementLog = "no_change"
)

type GetLogPath func(fileName string) string

var (
	getLogPath GetLogPath = common.GetLogPath
	//Logger
	tcpcopyPersistence      rlog.ErrorLogger
	memPersistence          rlog.ErrorLogger
	staticConfPersistence   rlog.ErrorLogger
	portraitDataPersistence rlog.ErrorLogger
	//md5 for diff
	md5ValueOfMemDump string
	// md5ValueOfStaticConf string

	memConfDumpFilePath string

	initLoggerOnce sync.Once
)

func getMemConfDumpFilePath() string {
	InitLogger()
	return memConfDumpFilePath
}

func GetTcpcopyLogger() rlog.ErrorLogger {
	InitLogger()
	return tcpcopyPersistence
}

func GetMemLogger() rlog.ErrorLogger {
	InitLogger()
	return memPersistence
}

func GetStaticConfLogger() rlog.ErrorLogger {
	InitLogger()
	return staticConfPersistence
}

func GetPortraitDataLogger() rlog.ErrorLogger {
	InitLogger()
	return portraitDataPersistence
}

func InitLogger() {
	initLoggerOnce.Do(doInitLogger)
}
func doInitLogger() {
	// local variable
	tcpcopyDumpFilePath := getLogPath(tcpcopyDumpFile)
	portraitDataDumpFilePath := getLogPath(portraitDataDumpFile)
	staticConfDumpFilePath := getLogPath(staticConfDumpFile)
	// write global variable
	memConfDumpFilePath = getLogPath(memConfDumpFile)

	// init logger using these path variables.
	tcpcopyLogger, err1 := log.GetOrCreateDefaultErrorLogger(tcpcopyDumpFilePath, log.INFO)
	if err1 != nil {
		log.StartLogger.Errorf("%s init tcpcopy logger error, err=&s", model.LogDumpKey, err1.Error())
	} else {
		tcpcopyPersistence = tcpcopyLogger
	}

	memDumpLogger, err2 := log.GetOrCreateDefaultErrorLogger(memConfDumpFilePath, log.INFO)
	if err2 != nil {
		log.StartLogger.Errorf("%s init mem dump logger error, err=&s", model.LogDumpKey, err2.Error())
	} else {
		memPersistence = memDumpLogger
	}

	staticConfLogger, err3 := log.GetOrCreateDefaultErrorLogger(staticConfDumpFilePath, log.INFO)
	if err3 != nil {
		log.StartLogger.Errorf("%s init static config logger error, err=&s", model.LogDumpKey, err3.Error())
	} else {
		staticConfPersistence = staticConfLogger
	}

	portraitDataLogger, err4 := log.GetOrCreateDefaultErrorLogger(portraitDataDumpFilePath, log.INFO)
	if err4 != nil {
		log.StartLogger.Errorf("%s init portrait data logger error, err=&s", model.LogDumpKey, err4.Error())
	} else {
		portraitDataPersistence = portraitDataLogger
	}
}

func IsPersistence() bool {
	// Determine the switch state
	if !strategy.DumpSwitch {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the dump switch is %t", model.LogDumpKey, strategy.DumpSwitch)
		}
		return false
	}

	// Determine whether it is within the sampling period
	if atomic.LoadInt32(&strategy.DumpSampleFlag) == 0 {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the dump sample flag is %d", model.LogDumpKey, strategy.DumpSampleFlag)
		}
		return false
	}

	// Determine whether it is fused
	if !strategy.IsAvaliable() {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the system usages are beyond max rate.", model.LogDumpKey)
		}
		return false
	}

	return true
}

func persistence(config *model.DumpUploadDynamicConfig) {
	// 1.Persist binary data
	if config.Binary_flow_data != nil && config.Port != "" {
		if GetTcpcopyLogger().GetLogLevel() >= log.INFO {
			GetTcpcopyLogger().Infof("[%s][%s]% x", config.Unique_sample_window, config.Port, config.Binary_flow_data)
		}
	}
	if config.Portrait_data != "" && config.BusinessType != "" {
		// 2. Persistent user-defined data
		if GetPortraitDataLogger().GetLogLevel() >= log.INFO {
			GetPortraitDataLogger().Infof("[%s][%s][%s]%s", config.Unique_sample_window, config.BusinessType, config.Port, config.Portrait_data)
		}

		// 3. Persistent memory configuration data, only make incremental changes
		buf, err := configmanager.DumpJSON()
		if err != nil {
			if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
				log.DefaultLogger.Debugf("[dump] Failed to load mosn config mem.")
			}
			return
		}
		// 3.1. dump if the data has been changed
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
