# 源码解析 4层流量治理，tcp流量dump

>作者简介：
>龚中强，是开源社区的爱好者，致力于拥抱开源。
> 
>写作时间: 2022年4月26日


## Overview
此文档的目的在于分析 tcp 流量 dump 的实现

## 前提：
文档内容所涉及代码版本如下

[https://github.com/mosn/layotto](https://github.com/mosn/layotto)

Layotto   0e97e970dc504e0298017bd956d2841c44c0810b（main分支）

## 源码分析

### 代码均在： [tcpcopy代码](https://github.com/mosn/layotto/tree/main/pkg/filter/network/tcpcopy)

### model.go分析
此类是 tcpcopy 的配置对象的核心类

```go
type DumpConfig struct {
	Switch     string  `json:"switch"`       // dump 开关.配置值：'ON' 或 'OFF'
	Interval   int     `json:"interval"`     // dump 采样间隔， 单位: 秒
	Duration   int     `json:"duration"`     // 单个采样周期， 单位: 秒
	CpuMaxRate float64 `json:"cpu_max_rate"` // cpu 最大使用率。当超过此阈值,dump 功能将停止
	MemMaxRate float64 `json:"mem_max_rate"` // mem 最大使用率。当超过此阈值,dump 功能将停止
}

type DumpUploadDynamicConfig struct {
	Unique_sample_window string             // 指定采样窗口
	BusinessType         _type.BusinessType // 业务类型
	Port                 string             // 端口
	Binary_flow_data     []byte             // 二进制数据
	Portrait_data        string             // 用户上传的数据
}
```

### persistence.go分析
此类是 tcpcopy 的 dump 持久化核心处理类

```go
// 该方法在 tcpcopy.go 中 OnData 中调用
func IsPersistence() bool {
	// 判断 dump 开关是否开启
	if !strategy.DumpSwitch {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the dump switch is %t", model.LogDumpKey, strategy.DumpSwitch)
		}
		return false
	}

	// 判断是否在采样窗口中
	if atomic.LoadInt32(&strategy.DumpSampleFlag) == 0 {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the dump sample flag is %d", model.LogDumpKey, strategy.DumpSampleFlag)
		}
		return false
	}

	// 判断是否 dump 功能停止（获取系统负载判断处理器和内存是否超过 tcpcopy 的阈值，如果超过则停止）
	if !strategy.IsAvaliable() {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("%s the system usages are beyond max rate.", model.LogDumpKey)
		}
		return false
	}

	return true
}

// 根据配置信息持久化数据
func persistence(config *model.DumpUploadDynamicConfig) {
	// 1.持久化二进制数据
	if config.Binary_flow_data != nil && config.Port != "" {
		if GetTcpcopyLogger().GetLogLevel() >= log.INFO {
			GetTcpcopyLogger().Infof("[%s][%s]% x", config.Unique_sample_window, config.Port, config.Binary_flow_data)
		}
	}
	if config.Portrait_data != "" && config.BusinessType != "" {
		// 2. 持久化用户定义的数据
		if GetPortraitDataLogger().GetLogLevel() >= log.INFO {
			GetPortraitDataLogger().Infof("[%s][%s][%s]%s", config.Unique_sample_window, config.BusinessType, config.Port, config.Portrait_data)
		}

		// 3. 增量持久化内存中的配置信息的变动内容
		buf, err := configmanager.DumpJSON()
		if err != nil {
			if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
				log.DefaultLogger.Debugf("[dump] Failed to load mosn config mem.")
			}
			return
		}
		// 3.1. 如果数据变化则 dump 
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

### tcpcopy.go分析
此类为是 tcpcopy 的核心类。

```go
// 向 Mosn 注册 NetWork 
func init() {
	api.RegisterNetwork("tcpcopy", CreateTcpcopyFactory)
}

// 返回 tcpcopy 工厂
func CreateTcpcopyFactory(cfg map[string]interface{}) (api.NetworkFilterChainFactory, error) {
	tcpConfig := &config{}
	// dump 策略转静态配置
	if stg, ok := cfg["strategy"]; ok {
	...
	}
	// TODO extract some other fields
	return &tcpcopyFactory{
		cfg: tcpConfig,
	}, nil
}

// 供 pkg/configmanager/parser.go 调用添加或者更新Network filter factory
func (f *tcpcopyFactory) Init(param interface{}) error {
	// 设置监听的地址和端口配置
	...
	return nil
}

// 实现的是 ReadFilter 的 OnData 接口，每次从连接拿到数据都进方法进行处理
func (f *tcpcopyFactory) OnData(data types.IoBuffer) (res api.FilterStatus) {
	// 判断当前请求数据是否需要采样 dump 
	if !persistence.IsPersistence() {
		return api.Continue
	}

	// 异步的采样 dump
	config := model.NewDumpUploadDynamicConfig(strategy.DumpSampleUuid, "", f.cfg.port, data.Bytes(), "")
	persistence.GetDumpWorkPoolInstance().Schedule(config)
	return api.Continue
}
```


最后我们再来回顾一下整体流程走向:

1. 从 tcpcopy.go 的初始化函数init() 开始,程序向 CreateGRPCServerFilterFactory 传入 CreateTcpcopyFactory.

2. Mosn 创建出一个filter chain(代码位置[factory.go](https://github.com/mosn/mosn/tree/master/pkg/filter/network/proxy/factory.go)) ,通过循环调用CreateFilterChain将所有的filter加入到链路结构包括本文的 tcpcopy.

3. 当流量通过 mosn 将会进入到 tcpcopy.go 的 OnData 方法进行 tcpdump 的逻辑处理.
