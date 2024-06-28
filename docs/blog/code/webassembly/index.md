# Layotto 源码解析 —— WebAssembly

>本文主要分析 Layotto 中 WASM 的相关实现和应用。
> 
>作者：[王志龙](https://github.com/rayowang) | 2022年5月18日

- [概述](#概述)
- [源码分析](#源码分析)
  - [框架INIT](#框架INIT)
  - [工作流程](#工作流程)
  - [FaaS模式](#FaaS模式)
- [总结](#总结)

## 概述
WebAssemly 简称 WASM，是一种运行在沙箱化的执行环境中的可移植、体积小、加载快的二进制格式，WASM最初设计是为了在网络浏览器中实现高性能应用，得益于它良好的隔离性和安全性、多语言支持、冷启动快等灵活性和敏捷性等特性，又被应用于嵌入其它应用程序中以获得较好的扩展能力，显然我们可以将它嵌入到 Layotto 中。Layotto 支持加载编译好的 WASM 文件，并通过 proxy_abi_version_0_2_0 的 API 与目标 WASM 进行交互;
另外 Layotto 也支持加载并运行以 WASM 为载体的 Function，并支持 Function 之间互相调用以及访问基础设施；同时 Layotto 社区也正在探索把 component 编译成 WASM 模块以此来增强模块间的隔离性。本文以 Layotto 官方 [quickstart](https://mosn.io/layotto/#/zh/start/wasm/start) 即访问redis相关示例为例来分析 Layotto 中 WebAssemly 相关的实现和应用。

## 源码分析
备注：本文基于 commit hash：f1cf350a52b5a1a0b3788a31681007a056e332ef

### 框架INIT
由于 Layotto 的底层是 Mosn，WASM 的扩展框架也是复用 Mosn 的 WASM 扩展框架，如图1 Layotto & Mosn WASM 框架 [1] 所示。

![mosn_wasm_ext_framework_module](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*jz4BSJmVQ3gAAAAAAAAAAAAAARQnAQ)
<center>图1 Layotto & Mosn WASM 框架 </center>

其中，Manager 负责对 WASM 插件进行管理和动态更新；VM 负责对 WASM 虚拟机、模块和实例进行管理；ABI 作为应用程序二进制接口，提供对外使用接口 [2]。

这里首先简单回顾下几个概念：\
[Proxy-Wasm](https://github.com/proxy-wasm) ：WebAssembly for Proxies (ABI specification) 是一个代理无关的 ABI 标准，它约定了代理和 WASM 模块如何以函数和回调的形式互动 [3]。\
[proxy-wasm-go-sdk](https://github.com/tetratelabs/proxy-wasm-go-sdk) ：定义了函数访问系统资源及基础设施服务的接口，基于 [proxy-wasm/spec](https://github.com/proxy-wasm/spec) 实现，在此基础上结合 Runtime API 增加了对基础设施访问的 ABI。\
[proxy-wasm-go-host](https://github.com/mosn/proxy-wasm-go-host) WebAssembly for Proxies (GoLang host implementation)：Proxy-Wasm 的 golang 实现，用以在 Layotto 中实现 Runtime ABI 的具体逻辑。\
VM：Virtual Machine 虚拟机，Runtime类型有：wasmtime、Wasmer、V8、 Lucet、WAMR、wasm3，本文例子中使用 wasmer

1、首先看 [quickstart例子](https://mosn.io/layotto/#/zh/start/wasm/start) 中 stream filter 的配置，如下可以看到配置中有两个 WASM 插件，使用 wasmer VM 分别启动一个实例，详见如下配置：

```json
 "stream_filters": [
            {
              "type": "Layotto",
              "config": {
                "function1": {
                  "name": "function1", // 插件名
                  "instance_num": 1, // 沙箱实例个数
                  "vm_config": {
                    "engine": "wasmer", // 虚拟机 Runtime 类型
                    "path": "demo/faas/code/golang/client/function_1.wasm" // wasm 文件路径
                  }
                },
                "function2": {
                  "name": "function2", // 插件名
                  "instance_num": 1, // 沙箱实例个数
                  "vm_config": {
                    "engine": "wasmer", // 虚拟机 Runtime 类型
                    "path": "demo/faas/code/golang/server/function_2.wasm" // wasm 文件路径
                  }
                }
              }
            }
          ]
```

上述配置中 function1 主要逻辑就是接收 HTTP 请求，然后通过 ABI 调用 function2，并返回 function2 结果，详见如下代码：

```go
func (ctx *httpHeaders) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	//1. get request body
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("GetHttpRequestBody failed: %v", err)
		return types.ActionPause
	}

	//2. parse request param
	bookName, err := getQueryParam(string(body), "name")
	if err != nil {
		proxywasm.LogErrorf("param not found: %v", err)
		return types.ActionPause
	}

	//3. request function2 through ABI
	inventories, err := proxywasm.InvokeService("id_2", "", bookName)
	if err != nil {
		proxywasm.LogErrorf("invoke service failed: %v", err)
		return types.ActionPause
	}

	//4. return result
	proxywasm.AppendHttpResponseBody([]byte("There are " + inventories + " inventories for " + bookName + "."))
	return types.ActionContinue
}
```

function2 主要逻辑就是接收 HTTP 请求，然后通过 ABI 调用 redis，并返回 redis 结果，详见如下代码：

```go
func (ctx *httpHeaders) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	//1. get request body
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("GetHttpRequestBody failed: %v", err)
		return types.ActionPause
	}
	bookName := string(body)

	//2. get request state from redis by specific key through ABI
	inventories, err := proxywasm.GetState("redis", bookName)
	if err != nil {
		proxywasm.LogErrorf("GetState failed: %v", err)
		return types.ActionPause
	}

	//3. return result
	proxywasm.AppendHttpResponseBody([]byte(inventories))
	return types.ActionContinue
}
```

2、对应图1 WASM 框架 中的 Manager 部分，在 Mosn filter Init 阶段进行初始化，详见如下代码：

```go
// Create a proxy factory for WasmFilter
func createProxyWasmFilterFactory(confs map[string]interface{}) (api.StreamFilterChainFactory, error) {
	factory := &FilterConfigFactory{
		config:        make([]*filterConfigItem, 0, len(confs)),
		RootContextID: 1,
		plugins:       make(map[string]*WasmPlugin),
		router:        &Router{routes: make(map[string]*Group)},
	}

	for configID, confIf := range confs {
		conf, ok := confIf.(map[string]interface{})
		if !ok {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory config not a map, configID: %s", configID)
			return nil, errors.New("config not a map")
		}
		// 解析 wasm filter 配置
		config, err := parseFilterConfigItem(conf)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to parse config, configID: %s, err: %v", configID, err)
			return nil, err
		}

		var pluginName string
		if config.FromWasmPlugin == "" {
			pluginName = utils.GenerateUUID()
            
			// 根据 stream filter 的配置初始化 WASM 插件配置，VmConfig 即 vm_config，InstanceNum 即 instance_num
			v2Config := v2.WasmPluginConfig{
				PluginName:  pluginName,
				VmConfig:    config.VmConfig,
				InstanceNum: config.InstanceNum,
			}
            
			// WasmManager 实例通过管理 PluginWrapper 对象对所有插件的配置进行统一管理，提供增删查改能力。下接3
			err = wasm.GetWasmManager().AddOrUpdateWasm(v2Config)
			if err != nil {
				config.PluginName = pluginName
				addWatchFile(config, factory)
				continue
			}

			addWatchFile(config, factory)
		} else {
			pluginName = config.FromWasmPlugin
		}
		config.PluginName = pluginName

		// PluginWrapper 在上面的 AddOrUpdateWasm 中对插件及配置进行封装完成初始化，这里根据插件名从 sync.Map 拿出，以管理并注册 PluginHandler
		pw := wasm.GetWasmManager().GetWasmPluginWrapperByName(pluginName)
		if pw == nil {
			return nil, errors.New("plugin not found")
		}

		config.VmConfig = pw.GetConfig().VmConfig
		factory.config = append(factory.config, config)

		wasmPlugin := &WasmPlugin{
			pluginName:    config.PluginName,
			plugin:        pw.GetPlugin(),
			rootContextID: config.RootContextID,
			config:        config,
		}
		factory.plugins[config.PluginName] = wasmPlugin
		// 注册 PluginHandler，以对插件的生命周期提供扩展回调能力，例如插件启动 OnPluginStart、更新 OnConfigUpdate。下接4
		pw.RegisterPluginHandler(factory)
	}

	return factory, nil
}
```

3、对应图1 WASM 框架中 VM 部分，NewWasmPlugin 用来创建初始化 WASM 插件，其中 VM、Module 和 Instance 分别对应 WASM 中的虚拟机、模块和实例，详见如下代码：

```go
func NewWasmPlugin(wasmConfig v2.WasmPluginConfig) (types.WasmPlugin, error) {
	// check instance num
	instanceNum := wasmConfig.InstanceNum
	if instanceNum <= 0 {
		instanceNum = runtime.NumCPU()
	}

	wasmConfig.InstanceNum = instanceNum

	// 根据配置获取 wasmer 编译和执行引擎
	vm := GetWasmEngine(wasmConfig.VmConfig.Engine)
	if vm == nil {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin fail to get wasm engine: %v", wasmConfig.VmConfig.Engine)
		return nil, ErrEngineNotFound
	}

	// load wasm bytes
	var wasmBytes []byte
	if wasmConfig.VmConfig.Path != "" {
		wasmBytes = loadWasmBytesFromPath(wasmConfig.VmConfig.Path)
	} else {
		wasmBytes = loadWasmBytesFromUrl(wasmConfig.VmConfig.Url)
	}

	if len(wasmBytes) == 0 {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin fail to load wasm bytes, config: %v", wasmConfig)
		return nil, ErrWasmBytesLoad
	}

	md5Bytes := md5.Sum(wasmBytes)
	newMd5 := hex.EncodeToString(md5Bytes[:])
	if wasmConfig.VmConfig.Md5 == "" {
		wasmConfig.VmConfig.Md5 = newMd5
	} else if newMd5 != wasmConfig.VmConfig.Md5 {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin the hash(MD5) of wasm bytes is incorrect, config: %v, real hash: %s",
			wasmConfig, newMd5)
		return nil, ErrWasmBytesIncorrect
	}

	// 创建 WASM 模块，WASM 模块是已被编译的无状态二进制代码
	module := vm.NewModule(wasmBytes)
	if module == nil {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin fail to create module, config: %v", wasmConfig)
		return nil, ErrModuleCreate
	}

	plugin := &wasmPluginImpl{
		config:    wasmConfig,
		vm:        vm,
		wasmBytes: wasmBytes,
		module:    module,
	}

	plugin.SetCpuLimit(wasmConfig.VmConfig.Cpu)
	plugin.SetMemLimit(wasmConfig.VmConfig.Mem)

	// 创建包含模块和运行时状态的实例，值得关注的是，这里最终会调用 proxywasm.RegisterImports 注册用户实现的 Imports 函数，比如示例中的 proxy_invoke_service 和 proxy_get_state
actual := plugin.EnsureInstanceNum(wasmConfig.InstanceNum)
	if actual == 0 {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin fail to ensure instance num, want: %v got 0", instanceNum)
		return nil, ErrInstanceCreate
	}

	return plugin, nil
}
```

4、 对应图1 WASM 框架 中的 ABI 部分，OnPluginStart 方法中会调用 proxy-wasm-go-host 的对应方法对 ABI 的 Exports 和 Imports 等进行相关设置。

```go
// Execute the plugin of FilterConfigFactory
func (f *FilterConfigFactory) OnPluginStart(plugin types.WasmPlugin) {
	plugin.Exec(func(instance types.WasmInstance) bool {
		wasmPlugin, ok := f.plugins[plugin.PluginName()]
		if !ok {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get wasm plugin, PluginName: %s",
				plugin.PluginName())
			return true
		}
        
		// 获取 proxy_abi_version_0_2_0 版本的与 WASM 交互的 API
		a := abi.GetABI(instance, AbiV2)
		a.SetABIImports(f)
		exports := a.GetABIExports().(Exports)
		f.LayottoHandler.Instance = instance

		instance.Lock(a)
		defer instance.Unlock()

		// 使用 exports 函数 proxy_get_id（对应到 WASM 插件中 GetID 函数）获取 WASM 的 ID
		id, err := exports.ProxyGetID()
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get wasm id, PluginName: %s, err: %v",
				plugin.PluginName(), err)
			return true
		}
		// 把ID 和 对应的插件注册到路由中，即可通过 http Header 中的键值对进行路由，比如 'id:id_1' 就会根据 id_1 路由到上面的 Function1 
		f.router.RegisterRoute(id, wasmPlugin)

		// 当第一个插件使用给定的根 ID 加载时通过 proxy_on_context_create 创建根上下文，并在虚拟机的整个生命周期中持续存在，直到 proxy_on_delete 删除 
		// 值得注意的是这里说的第一个插件指的是多个松散绑定的插件(通过 SDK 使用 Root ID 对 Root Context 访问）在同一已配置虚拟机内共享数据的使用场景 [4]
		err = exports.ProxyOnContextCreate(f.RootContextID, 0)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		vmConfigSize := 0
		if vmConfigBytes := wasmPlugin.GetVmConfig(); vmConfigBytes != nil {
			vmConfigSize = vmConfigBytes.Len()
		}

		// VM 伴随启动的插件启动时调用
		_, err = exports.ProxyOnVmStart(f.RootContextID, int32(vmConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		pluginConfigSize := 0
		if pluginConfigBytes := wasmPlugin.GetPluginConfig(); pluginConfigBytes != nil {
			pluginConfigSize = pluginConfigBytes.Len()
		}

		// 当插件加载或重新加载其配置时调用
		_, err = exports.ProxyOnConfigure(f.RootContextID, int32(pluginConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		return true
	})
}
```

### 工作流程
Layotto 中 WASM 的工作流程大致如下图2 Layotto & Mosn WASM 工作流程所示，其中配置更新在上述初始化环节基本已囊括，这里重点看一下请求处理流程。
![mosn_wasm_ext_framework_workflow](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*XTDeRq0alYsAAAAAAAAAAAAAARQnAQ)
<center>图2 Layotto & Mosn WASM 工作流程 </center>

1、由 Layotto 底层 Mosn 收到请求，经过 workpool 调度，在 proxy downstream 中按照配置依次执行 StreamFilterChain 到 Wasm StreamFilter 的 OnReceive 方法，具体逻辑详见如下代码：

```go
func (f *Filter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	// 获取 WASM 插件的 id
	id, ok := headers.Get("id")
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders no id in headers")
		return api.StreamFilterStop
	}
    
	// 从 router 中根据 id 获取对应的 WASM 插件
	wasmPlugin, err := f.router.GetRandomPluginByID(id)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders id, err: %v", err)
		return api.StreamFilterStop
	}
	f.pluginUsed = wasmPlugin

	plugin := wasmPlugin.plugin
	// 获取 WasmInstance 实例
	instance := plugin.GetInstance()
	f.instance = instance
	f.LayottoHandler.Instance = instance

	// ABI 包含 导出(Exports)和导入(Imports)两个部分，用户通过这它们与 WASM 扩展插件进行交互
	pluginABI := abi.GetABI(instance, AbiV2)
	if pluginABI == nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive fail to get instance abi")
		plugin.ReleaseInstance(instance)
		return api.StreamFilterStop
	}
	// 设置导入 Imports 部分，导入部分由用户提供，虚拟机的执行需要依赖宿主机 Layotto 提供的部分能力，例如获取请求信息，这些能力通过导入部分由用户提供，并由 WASM 扩展调用
	pluginABI.SetABIImports(f)

	// 导出 Exports 部分由 WASM 插件提供，用户可直接调用——唤醒 WASM 虚拟机，并在虚拟机中执行对应的 WASM 插件代码
	exports := pluginABI.GetABIExports().(Exports)
	f.exports = exports
	
	instance.Lock(pluginABI)
	defer instance.Unlock()
	
	// 根据 rootContextID 和 contextID 创建当前插件上下文
	err = exports.ProxyOnContextCreate(f.contextID, wasmPlugin.rootContextID)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to create context id: %v, rootContextID: %v, err: %v",
			f.contextID, wasmPlugin.rootContextID, err)
		return api.StreamFilterStop
	}

	endOfStream := 1
	if (buf != nil && buf.Len() > 0) || trailers != nil {
		endOfStream = 0
	}

	// 调用 proxy-wasm-go-host，编码请求头为规范指定的格式
	action, err := exports.ProxyOnRequestHeaders(f.contextID, int32(headerMapSize(headers)), int32(endOfStream))
	if err != nil || action != proxywasm.ActionContinue {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders err: %v", err)
		return api.StreamFilterStop
	}

	endOfStream = 1
	if trailers != nil {
		endOfStream = 0
	}

	if buf == nil {
		arg, _ := variable.GetString(ctx, types.VarHttpRequestArg)
		f.requestBuffer = buffer.NewIoBufferString(arg)
	} else {
		f.requestBuffer = buf
	}

	if f.requestBuffer != nil && f.requestBuffer.Len() > 0 {
		// 调用 proxy-wasm-go-host，编码请求体为规范指定的格式
		action, err = exports.ProxyOnRequestBody(f.contextID, int32(f.requestBuffer.Len()), int32(endOfStream))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestBody err: %v", err)
			return api.StreamFilterStop
		}
	}

	if trailers != nil {
        // 调用 proxy-wasm-go-host，编码请求尾为规范指定的格式
		action, err = exports.ProxyOnRequestTrailers(f.contextID, int32(headerMapSize(trailers)))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestTrailers err: %v", err)
			return api.StreamFilterStop
		}
	}

	return api.StreamFilterContinue
}
```

2、proxy-wasm-go-host 将 Mosn 请求三元组编码成规范指定的格式，并调用Proxy-Wasm ABI 规范中的 proxy_on_request_headers 等对应接口，调用 WASMER 虚拟机将请求信息传至 WASM 插件。

```go
func (a *ABIContext) CallWasmFunction(funcName string, args ...interface{}) (interface{}, Action, error) {
	ff, err := a.Instance.GetExportsFunc(funcName)
	if err != nil {
		return nil, ActionContinue, err
	}

	// 调用 wasmer 虚拟机（github.com/wasmerio/wasmer-go/wasmer.(*Function).Call at function.go）
	res, err := ff.Call(args...)
	if err != nil {
		a.Instance.HandleError(err)
		return nil, ActionContinue, err
	}

	// if we have sync call, e.g. HttpCall, then unlock the wasm instance and wait until it resp
	action := a.Imports.Wait()

	return res, action, nil
}
```

3、WASMER 虚拟机经过处理调用 WASM 插件的具体函数，比如例子中的 OnHttpRequestBody 函数
  // function, _:= instance.Exports.GetFunction("exported_function")
  // nativeFunction = function.Native()
  //_ = nativeFunction(1, 2, 3)
  // Native 会将 Function 转换为可以调用的原生 Go 函数

```go
func (self *Function) Native() NativeFunction {
	...
	self.lazyNative = func(receivedParameters ...interface{}) (interface{}, error) {
		numberOfReceivedParameters := len(receivedParameters)
		numberOfExpectedParameters := len(expectedParameters)
		...
		results := C.wasm_val_vec_t{}
		C.wasm_val_vec_new_uninitialized(&results, C.size_t(len(ty.Results())))
		defer C.wasm_val_vec_delete(&results)

		arguments := C.wasm_val_vec_t{}
		defer C.wasm_val_vec_delete(&arguments)

		if numberOfReceivedParameters > 0 {
			C.wasm_val_vec_new(&arguments, C.size_t(numberOfReceivedParameters), (*C.wasm_val_t)(unsafe.Pointer(&allArguments[0])))
		}

		// 调用 WASM 插件内函数
		trap := C.wasm_func_call(self.inner(), &arguments, &results)

		runtime.KeepAlive(arguments)
		runtime.KeepAlive(results)
		...
	}

	return self.lazyNative
}
```

4、proxy-wasm-go-sdk 将请求数据从规范格式转换为便于用户使用的格式，然后调用用户扩展代码。proxy-wasm-go-sdk 基于 proxy-wasm/spec 实现，定义了函数访问系统资源及基础设施服务的接口，并在此基础上结合 Runtime API 的思路，增加了对基础设施访问的ABI。

```go
// function1主要逻辑就是接收 HTTP 请求，然后通过 ABI 调用 function2，并返回 function2 结果，具体代码如下所示
func (ctx *httpHeaders) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	//1. get request body
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("GetHttpRequestBody failed: %v", err)
		return types.ActionPause
	}

	//2. parse request param
	bookName, err := getQueryParam(string(body), "name")
	if err != nil {
		proxywasm.LogErrorf("param not found: %v", err)
		return types.ActionPause
	}

	//3. request function2 through ABI
	inventories, err := proxywasm.InvokeService("id_2", "", bookName)
	if err != nil {
		proxywasm.LogErrorf("invoke service failed: %v", err)
		return types.ActionPause
	}

	//4. return result
	proxywasm.AppendHttpResponseBody([]byte("There are " + inventories + " inventories for " + bookName + "."))
	return types.ActionContinue
}
```

5、WASM 插件通过初始化时 RegisterFunc 注册的 ABI Imports 函数，比如例子中 Function1 RPC 调用 Function2 的 ProxyInvokeService，Function2 用以获取 Redis 中指定 Key 的 Valye 的 ProxyGetState，具体代码如下所示：

Function1 通过 ProxyInvokeService 调用 Function2，ProxyInvokeService 对应 Imports 函数 proxy_invoke_service

```go
func ProxyInvokeService(instance common.WasmInstance, idPtr int32, idSize int32, methodPtr int32, methodSize int32, paramPtr int32, paramSize int32, resultPtr int32, resultSize int32) int32 {
	id, err := instance.GetMemory(uint64(idPtr), uint64(idSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	method, err := instance.GetMemory(uint64(methodPtr), uint64(methodSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	param, err := instance.GetMemory(uint64(paramPtr), uint64(paramSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	ctx := getImportHandler(instance)
    
	// Laytto rpc calls
	ret, res := ctx.InvokeService(string(id), string(method), string(param))
	if res != WasmResultOk {
		return res.Int32()
	}

	return copyIntoInstance(instance, ret, resultPtr, resultSize).Int32()
}
```

Function2 通过 ProxyGetState 获取 Redis 中指定 Key 的 Valye， ProxyGetState 对应 Imports 函数 proxy_get_state

```go
func ProxyGetState(instance common.WasmInstance, storeNamePtr int32, storeNameSize int32, keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32 {
	storeName, err := instance.GetMemory(uint64(storeNamePtr), uint64(storeNameSize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	key, err := instance.GetMemory(uint64(keyPtr), uint64(keySize))
	if err != nil {
		return WasmResultInvalidMemoryAccess.Int32()
	}

	ctx := getImportHandler(instance)

	ret, res := ctx.GetState(string(storeName), string(key))
	if res != WasmResultOk {
		return res.Int32()
	}

	return copyIntoInstance(instance, ret, valuePtr, valueSize).Int32()
}
```

以上 Layotto rpc 流程简要说是通过两个虚拟连接借助 Dapr API 和 底层 Mosn 实现 [5],具体可参见前序文章[Layotto源码解析——处理RPC请求](https://mosn.io/layotto/#/zh/blog/code/layotto-rpc/index)，从 Redis 中获取数据可直接阅读 Dapr State 相关代码，在此不一一展开了。

### FaaS模式

回过头来再看 WASM 的特性：字节码有与机器码相匹敌的性能；沙箱中执行保证良好的隔离性和安全性；编译后跨平台、易分发和加载运行；具备轻量且多语言开发的灵活性，似乎天然的就适合做 FaaS。

所以 Layotto 也探索支持了 WASM FaaS 模式，即加载并运行以 WASM 为载体的 Function，并支持 Function 之间相互调用及访问基础设施。因加载 WASM 的核心逻辑并未变化，只是使用和部署方式上与上述方式有差别，故 Layotto 加载 WASM 部分逻辑不再赘述。

除 Wasm-Proxy 相关实现外，FaaS 模式核心逻辑是通过扩展 Containerd 实现多运行时插件 containerd-shim-layotto-v2 [6]，并借此"穿针引线"的巧妙的利用了 Docker 的镜像能力来管理 *.wasm 包和 Kubernetes 优秀的编排能力来调度函数，具体架构和工作流可见图3 Layotto FaaS Workflow。

![layotto_faas_workflow](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*XWmNT6-7FoEAAAAAAAAAAAAAARQnAQ)
<center>图3 Layotto FaaS Workflow </center>

这里简单看一下 containerd-shim-layotto-v2 的主函数，可以看到 shim.Run 设置的 WASM 的运行时为 io.containerd.layotto.v2，也就是 containerd 中 plugins.cri.containerd.runtimes 对应插件的 runtime_type。当创建 Pod 时，在 yaml 的 spec 中指定 runtimeClassName: layotto，经过调度，最终 kubelet 就会通过 cri-plugin 调用 containerd 中的 containerd-shim-layotto-v2 运行时来进行加载和运行等相关处理。

```go
func main() {
	startLayotto()
	// 解析输入参数，初始化运行时环境，调用 wasm.New 实例化 service 对象 
	shim.Run("io.containerd.layotto.v2", wasm.New)
}

func startLayotto() {
	conn, err := net.Dial("tcp", "localhost:2045")
	if err == nil {
		conn.Close()
		return
	}

	cmd := exec.Command("layotto", "start", "-c", "/home/docker/config.json")
	cmd.Start()
}
```

## 总结
Layotto WebAssemly 虽然涉及较多 WASM 相关的基础知识，但通过示例由浅入深，循序渐进也不难理解。最后整体看一下 WASM 技术，可以看到它已经被应用到Web前端、Serverless、游戏场景、边缘计算、服务网格等很多领域，甚至就连 Docker 之父 Solomon Hykes 在前不久都表示: "如果 WASM 这个技术在2008年就有的话，我就不搞Docker了"（后来又补充道：Docker 不会被替换，会与 WASM 并肩而行），不管怎么说，WASM 似乎在继 VM 和 Container 之后，正在成为更轻量及性能更好的云原生技术而被应用到更多的领域，与此同时，相信在 Mosn 社区的推动以及 Layotto 的继续探索中 WASM 也会有更多使用场景和用户，至此 Layotto WebAssemly 相关源码分析就完了，鉴于时间和篇幅，没有进行一些更全面和深入的剖析，如有纰漏之处，欢迎指正，联系方式：rayo.wangzl@gmail.com。

### 参考资料
- [1] [WebAssembly 在 MOSN 中的实践](https://mosn.io/blog/posts/mosn-wasm-framework/)
- [2] [feature: WASM plugin framework](https://github.com/mosn/mosn/pull/1589)
- [3] [WebAssembly for Proxies (ABI Spec)](https://github.com/proxy-wasm/spec)
- [4] [Proxy WebAssembly Architecture](https://techhenzy.com/proxy-webassembly-architecture/)
- [5] [Layotto源码解析——处理RPC请求](https://mosn.io/layotto/#/zh/blog/code/layotto-rpc/index)
- [6] [云原生运行时的下一个五年](https://www.sofastack.tech/blog/the-next-five-years-of-cloud-native-runtime/)
