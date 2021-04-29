package persistence

import (
	"github.com/layotto/layotto/pkg/common"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/model"
	"github.com/layotto/layotto/pkg/filter/network/tcpcopy/strategy"
	"os"
	"sync/atomic"

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

var (
	TcpcopyPersistence      rlog.ErrorLogger
	MemPersistence          rlog.ErrorLogger
	StaticConfPersistence   rlog.ErrorLogger
	PortraitDataPersistence rlog.ErrorLogger

	md5ValueOfMemDump    string
	md5ValueOfStaticConf string

	memConfDumpFilePath    string
	staticConfDumpFilePath string
)

func init() {
	initLogger()
}

func initLogger() {
	tcpcopyDumpFilePath := common.GetLogPath(tcpcopyDumpFile)
	memConfDumpFilePath = common.GetLogPath(memConfDumpFile)
	staticConfDumpFilePath = common.GetLogPath(staticConfDumpFile)
	portraitDataDumpFilePath := common.GetLogPath(portraitDataDumpFile)

	tcpcopyLogger, err1 := log.GetOrCreateDefaultErrorLogger(tcpcopyDumpFilePath, log.INFO)
	if err1 != nil {
		log.StartLogger.Errorf("%s init tcpcopy logger error, err=&s", model.LogDumpKey, err1.Error())
	} else {
		TcpcopyPersistence = tcpcopyLogger
	}

	memDumpLogger, err2 := log.GetOrCreateDefaultErrorLogger(memConfDumpFilePath, log.INFO)
	if err2 != nil {
		log.StartLogger.Errorf("%s init mem dump logger error, err=&s", model.LogDumpKey, err2.Error())
	} else {
		MemPersistence = memDumpLogger
	}

	staticConfLogger, err3 := log.GetOrCreateDefaultErrorLogger(staticConfDumpFilePath, log.INFO)
	if err3 != nil {
		log.StartLogger.Errorf("%s init static config logger error, err=&s", model.LogDumpKey, err3.Error())
	} else {
		StaticConfPersistence = staticConfLogger
	}

	portraitDataLogger, err4 := log.GetOrCreateDefaultErrorLogger(portraitDataDumpFilePath, log.INFO)
	if err4 != nil {
		log.StartLogger.Errorf("%s init portrait data logger error, err=&s", model.LogDumpKey, err4.Error())
	} else {
		PortraitDataPersistence = portraitDataLogger
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
		if TcpcopyPersistence.GetLogLevel() >= log.INFO {
			TcpcopyPersistence.Infof("[%s][%s]% x", config.Unique_sample_window, config.Port, config.Binary_flow_data)
		}
	}
	if config.Portrait_data != "" && config.BusinessType != "" {
		// 2. Persistent user-defined data
		if PortraitDataPersistence.GetLogLevel() >= log.INFO {
			PortraitDataPersistence.Infof("[%s][%s][%s]%s", config.Unique_sample_window, config.BusinessType, config.Port, config.Portrait_data)
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
		if tmpMd5ValueOfMemDump != md5ValueOfMemDump || (tmpMd5ValueOfMemDump == md5ValueOfMemDump && common.GetFileSize(memConfDumpFilePath) <= 0) {
			md5ValueOfMemDump = tmpMd5ValueOfMemDump
			if MemPersistence.GetLogLevel() >= log.INFO {
				MemPersistence.Infof("[%s]%s", config.Unique_sample_window, buf)
			}
		} else {
			if MemPersistence.GetLogLevel() >= log.INFO {
				MemPersistence.Infof("[%s]%+v", config.Unique_sample_window, incrementLog)
			}
		}
	}
}
