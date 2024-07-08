# Layotto Source Parsing — WebAssembly

> This paper mainly analyses the relevant implementation and application of Layotto Middle WASM.
>
> by：[Wang Zhilong](https://github.com/rayowang) | 18 May 2022

- [overview](#overview)
- [source analysis](#source analysis)
    - [Frame INIT](#Frame INIT)
    - [workflow](#workflow)
    - [FaaSmode](#FaaS mode)
- [summary](#summary)

## General description

WebAssemly Abbreviations WASM, a portable, small and loaded binary format operating in sandboxing implementation environment, was originally designed to achieve high-performance applications in web browsers, benefiting from its good segregation and security, multilingual support, cool-start fast flexibility and agility and application to embed other applications for better expansion, and obviously we can embed it into Layotto.Layotto supports loading compiled WASM files and interacting with the Target WASM API via proxy_abi_version_0_2_0;
other Layotto also supports loading and running WASM carrier functions and supports interfaces between Function and access to infrastructure; and Layotto communities are also exploring the compilation of components into WASM modules to increase segregation between modules.This article uses the Layotto official [quickstart](https://mosn.io/layotto/docs/start/wasm/start) example of accessing redis as an example to analyze WebAssemly in Layotto Related implementation and application.

## Source analysis

Note：is based on commit hash：f1cf350a52b5a1a0b3788a31681007a056e332ef

### Frame INIT

As the bottom layer of Layotto is Mosn, the WASM extension framework is also the WASM extension framework that reuses Mosn, as shown in figure 1 Layotto & Mosn WASM framework [1].

![mosn\_wasm\_ext\_framework\_module](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*jz4BSJmVQ3gAAAAAAAAAAAAAARQnAQ)

<center>Figure 1 Layotto & Mosn WASM framework </center>

Among them, Manager is responsible for managing and dynamically updating WASM plugins;VM for managing WASM virtual machines, modules and instances;ABI serves as the application binary interface to provide an external interface [2].

Here a brief review of the following concepts：\
[Proxy-Wasm](https://github.com/proxy-wasm) ：WebAssembly for Proxies (ABI specification) is an unrelated ABI standard that defines how proxy and WASM modules interact [3] in functions and callbacks.
[proxy-wasm-go-sdk](https://github.com/tetratelabs/proxy-wasm-go-sdk) ：defines the interface of function access to system resources and infrastructure services based on [proxy-wasm/spec](https://github.com/proxy-wasm/spec) which brings together the Runtime API to increase access to infrastructure.\
[proxy-wasm-go-host](https://github.com/mosn/proxy-wasm-go-host) WebAssembly for Proxies (GoLang host implementation)：Proxy-Wasm golang implementation to implement Runtime ABI logic in Layotto.\
VM: Virtual Machine Virtual machine. The Runtime types are wasmtime, wasmer, V8, Lucet, WAMR, and wasm3

1, see first the configuration of stream filter in [quickstart例子](https://mosn.io/layotto/#/start/waste/start) as follows, two WASM plugins can be seen, using waste VM to start a separate instance with configuration： below

```json
 "stream_filters": [
            LO
              "type": "Layotto",
              "config": API
                "Function1": LOs
                  "name": "function1", // Plugin name
                  "instance_num": 1, // Number of sandbox instances
                  "vm_config": LO
                    "engine": "waste", // Virtual Machine Type Runtime Type
                    "path": "demo/faas/code/golang/client/function_1. asm" /waste file path
                  }
                },
                "Function2": LO
                  "name": "function2", // Plugin name
                  "instance_num": 1, // Number of sandbox instances
                  "vm_config": LO
                    "engine": "waste", // Virtual Machine Type Runtime Type
                    "path": "demo/faas/code/golang/server/function_2. asm" /wasm file path
                  }
                }
              }
            }
]
```

The primary logic in the configuration above is to receive HTTP requests, then call function2 through ABI, and return function2 as detailed below in code：

```go
func (Ctx *pHeaders) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action Led
	/1. get request body
	body, err := proxywasm. etHttpRequestBody(0, bodySize)
	if err != nil L/
		proxywasm.LogErrorf("GetHttpRequestBody failed: %v", err)
		return types. ctionPause
	}

	/2. parse request param
	bookName, err := getQueryParam(string(body), "name")
	if err != nil Led
		proxywasm. ogErrorf("param not found: %v", err)
		returns types. ctionPause
	}

	/3. Request function2 through ABI
	inventories, err := proxywasm. nvokeService("id_2", "", bookName)
	if err != nil LO
		proxywasm.Logrorf("invoke service failed: %v", err)
		return types. ctionPause
	}

	/4. return result
	proxywasm. ppendHttpResponseBody([]byte ("There are " + inventories + " inventories for " + bookName + ".")
	return types.ActionContinue
}
```

Function2 Primary logic is to receive HTTP requests, then call redisis through ABI and return to redis, as shown below in code：

```go
func (Ctx *pHeaders) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action 6
	//1. get requested body
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil Led
		proxywasm. ogErrorf("GetHttpRequestBody failed: %v", err)
		returns types.ActionPause
	}
	bookName:= string(body)

	/ 2. get request state from redis by specific key through ABI
	inventories, err := proxywastem. etState("redis", bookName)
	if err != nil LO
		proxywasm.LogErrorf("GetState failed: %v", err)
		returns types. ctionPause
	}

	/ 3. return result
	proxywasm.AppendHttpResponseBody([]byte(inventories))
	return types.ActionContinue
}
```

2. The Manager component of the Frame 1 WASM is initialized at Mosn filter Init stage as shown below in code：

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
		// Parse the wasm filter configuration
		config, err := parseFilterConfigItem(conf)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to parse config, configID: %s, err: %v", configID, err)
			return nil, err
		}

		var pluginName string
		if config.FromWasmPlugin == "" {
			pluginName = utils.GenerateUUID()
            
			// The WASM plug-in configuration is initialized according to the stream filter configuration. VmConfig is vm_config, and InstanceNum is instance_num
			v2Config := v2.WasmPluginConfig{
				PluginName:  pluginName,
				VmConfig:    config.VmConfig,
				InstanceNum: config.InstanceNum,
			}
            
			// The WasmManager instance manages the configuration of all plug-ins in a unified manner by managing the PluginWrapper object, providing the ability to add, delete, check and modify. Continue 3
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

		// PluginWrapper wraps the plug-in and configuration in AddOrUpdateWasm above to complete the initialization, which is pulled from sync.Map according to the plug-in name to manage and register the PluginHandler
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
		// Register PluginHandler to provide extended callback capabilities for the plug-in's life cycle, such as the plug-in starting OnPluginStart and updating OnConfigUpdate. Continue 4
		pw.RegisterPluginHandler(factory)
	}

	return factory, nil
}
```

3 Corresponding to Figure 1 WASM frame, NewWasmPlugin, for creating initialization of the WASM plugin, where VM, Module and Instance refer to virtual machines, modules and instances in WASM, as shown below in code：

```go
func NewWasmPlugin(wasmConfig v2.WasmPluginConfig) (types.WasmPlugin, error) {
	// check instance num
	instanceNum := wasmConfig.InstanceNum
	if instanceNum <= 0 {
		instanceNum = runtime.NumCPU()
	}

	wasmConfig.InstanceNum = instanceNum

	// Get the wasmer compilation and execution engine according to the configuration
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

	// Create the WASM module, which is the stateless binary code that has been compiled
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

	// Contains module and runtime state to create instance, notable is that here will call proxywasm. RegisterImports registered users realize the Imports of function, Examples include proxy_invoke_service and proxy_get_state
actual := plugin.EnsureInstanceNum(wasmConfig.InstanceNum)
	if actual == 0 {
		log.DefaultLogger.Errorf("[wasm][plugin] NewWasmPlugin fail to ensure instance num, want: %v got 0", instanceNum)
		return nil, ErrInstanceCreate
	}

	return plugin, nil
}
```

Corresponding to ABI components in Figure 1 WASM frames, the OnPluginStart method calls proxy-wasm-go-host corresponding to ABI Exports and Imports etc.

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

		// Use the exports function proxy_get_id (which corresponds to the GetID function in the WASM plug-in) to get the ID of WASM
		id, err := exports.ProxyGetID()
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get wasm id, PluginName: %s, err: %v",
				plugin.PluginName(), err)
			return true
		}
		// If you register the ID and the corresponding plug-in in the route, the route can be performed using the key-value pair in the http Header. For example, 'id:id_1' is routed to Function1 based on id_1 
		f.router.RegisterRoute(id, wasmPlugin)

		// The root context is created by proxy_on_context_create when the first plug-in is loaded with the given root ID and persists for the entire life of the virtual machine until proxy_on_delete is deleted
               // It is worth noting that the first plug-in here refers to a use case where multiple loosely bound plug-ins (accessed via the SDK using the Root ID to the Root Context) share data within the same configured virtual machine [4]
		err = exports.ProxyOnContextCreate(f.RootContextID, 0)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		vmConfigSize := 0
		if vmConfigBytes := wasmPlugin.GetVmConfig(); vmConfigBytes != nil {
			vmConfigSize = vmConfigBytes.Len()
		}

		// VM is called when the plug-in is started with the startup
		_, err = exports.ProxyOnVmStart(f.RootContextID, int32(vmConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		pluginConfigSize := 0
		if pluginConfigBytes := wasmPlugin.GetPluginConfig(); pluginConfigBytes != nil {
			pluginConfigSize = pluginConfigBytes.Len()
		}

		// Called when the plug-in loads or reloads its configuration
		_, err = exports.ProxyOnConfigure(f.RootContextID, int32(pluginConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		return true
	})
}
```

### Workflow

The workflow for Layotto Middle WASM is broadly as shown in figure 2 Layotto & Mosn WASM workflow, where the configuration is largely covered by the initial elements above, with a focus on the request processing.
![mosn\_wasm\_ext\_framework\_workflow](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*XTDeRq0alYsAAAAAAAAAAAAAARQnAQ)

<center>Figure 2 Layotto & Mosn WAS Workflow </center>

By Layotto underneath Mosn, as a workpool schedule, implement the OnReceive method of StreamFilterChain to Wasm StreamFilter in proxy downstream, as configured and detailed in code： below

```go
func (f *Filter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	// Gets the id of the WASM plug-in
	id, ok := headers.Get("id")
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders no id in headers")
		return api.StreamFilterStop
	}
    
	// Obtain the WASM plug-in from the router based on its id
	wasmPlugin, err := f.router.GetRandomPluginByID(id)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders id, err: %v", err)
		return api.StreamFilterStop
	}
	f.pluginUsed = wasmPlugin

	plugin := wasmPlugin.plugin
	// Obtain an instance of WasmInstance
	instance := plugin.GetInstance()
	f.instance = instance
	f.LayottoHandler.Instance = instance

	// The ABI consists of Exports and Imports, through which users interact with the WASM extension
	pluginABI := abi.GetABI(instance, AbiV2)
	if pluginABI == nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive fail to get instance abi")
		plugin.ReleaseInstance(instance)
		return api.StreamFilterStop
	}
	// Set the Imports section. The import section is provided by the user. The execution of the virtual machine depends on some of the capabilities provided by the host Layotto, such as obtaining request information, which are provided by the user through the import section and invoked by the WASM extension
	pluginABI.SetABIImports(f)

	// The Exports section is provided by the WASM plug-in and can be called directly by the user to wake up the WASM virtual machine and execute the corresponding WASM plug-in code in the virtual machine
	exports := pluginABI.GetABIExports().(Exports)
	f.exports = exports
	
	instance.Lock(pluginABI)
	defer instance.Unlock()
	
	// Create the current plug-in context according to rootContextID and contextID
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

	// Call proxy-wasm-go-host, encoding the request header in the format specified by the specification
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
		// Call proxy-wasm-go-host, encoding the request body in the format specified by the specification
		action, err = exports.ProxyOnRequestBody(f.contextID, int32(f.requestBuffer.Len()), int32(endOfStream))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestBody err: %v", err)
			return api.StreamFilterStop
		}
	}

	if trailers != nil {
        // Call proxy-wasm-go-host, encoding the request tail in the format specified by the specification
		action, err = exports.ProxyOnRequestTrailers(f.contextID, int32(headerMapSize(trailers)))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestTrailers err: %v", err)
			return api.StreamFilterStop
		}
	}

	return api.StreamFilterContinue
}
```

2, proxy-wasm-go-host encode Mosn requests for triplets into the specified format and call Proxy-Wasm ABI equivalent interface in Proxy_on_request_headers and call the WASMER virtual machine to pass the request information to the WASM plugin.

```go
func (a *ABIContext) CallWasmFunction (functionName string, args ..interface{}) (interface{}, Action, error) um
	ff, err := a.Instance. eExportsFunc(functionName)
	if err != nil {
		return nil, ActionContinue, err
	}

	// Call waste virtual machine (Github.com/wasmerio/wasmer-go/wasmer.(*Function).Call at function.go)
	res, err := ff. all(args....)
	if err != nil L/
		a.Instance.HandleError(err)
		return nil, ActionContinue, err
	}

	// if we have sync call, e. HttpCall, then unlocked the waste instance and wait until it resp
	action := a.Imports.Wait()

	return res, action, nil
}
```

3. The WASMER virtual machine is processed to call specific functions of the WASM plug-in, such as the OnHttpRequestBody function in the example
   // function, _:= instance.Exports.GetFunction("exported_function")
   // nativeFunction = function.Native()
   //_ = nativeFunction(1, 2, 3)
   // Native converts Function to a native Go function that can be called

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

		// Call functions inside the WASM plug-in
		trap := C.wasm_func_call(self.inner(), &arguments, &results)

		runtime.KeepAlive(arguments)
		runtime.KeepAlive(results)
		...
	}

	return self.lazyNative
}
```

4, proxy-wasm-go-sdk converts the requested data from the normative format to a user-friendly format and then calls the user extension code.Proxy-wasm-go-sdk, based on proxy-waste/spec implementation, defines the interface between function access to system resources and infrastructure services, and builds on this integration of the Runtime API, adding ABI to infrastructure access.

```go
// function1The main logic is to receive the HTTP request, call function2 using the ABI, and return the function2 result. The code is as follows
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

5, WASM plugin is registered at RegisterFunc initialization. For example, Function1 RPC calls Proxy InvokeService,Function2 to get ProxyGetState specified in Redis as shown below in：

Function1 Call Function2, Proxy InvokeService for Imports function proxy_invoke_service through the Proxy InvokeService

```go
func ProxyInvokeService(instance common). asmInstance, idPtr int32, idSize int32, methodPtr int32, methodPtr int32, paramPtr int32, resultPtr int32, resultSize int32) int32 56
	id, err := instance. etMemory(uint64(idPtr), uint64(idSize))
	if err != nil LO
		returnWasmResultInvalidMemoryAcces.Int32()
	}

	method, err := instance. etMemory(uint64 (methodPtr), uint64 (methodSize))
	if err != nil LO
		returnWasmResultInvalidMemoryAccess. nt32()
	}

	param, err := instance.GetMemory(uint64 (paramPtr), uint64 (paramSize))
	if err != nil Fe
		returnn WasmResultInvalidMemoryAccess. nt32()
	}

	ctx:= getImportHandler(instance)
    
	// Laytto rpc calls
	ret, res := ctx. nvokeService(string(id), string(param))
	if res != WasmResultOk 6
		return res.Int32()


	return copyIntoInstance(instance, ret, resultPtr, resultSize).Int32()
}
```

Function2 Get Redis via ProxyGetState to specify key Valye, ProxyGetState for Imports function proxy_get_state

```go
func ProxyGetState(instance common.WasmInstance, storeNamePtr int32, storeNameSize int32, keyPtr int32, valuePtr int32, valueSize int32) int32 Fe
	storeName, err := instance. etMemory(uint64 (storeNamePtr), uint64 (storeNameSize))
	if err != nil LO
		returnWasmResultInvalidMemoryAccess.Int32()
	}

	key, err := instance. etMemory(uint64(keyPtr), uint64(keySize))
	if err != nil LO
		returnWasmResultInvalidMemoryAccess.Int32()
	}

	ctx := getImportHandler(instance)

	ret, res := ctx. etState(string(storeName), string(key))
	if res != WasmResultOk 6
		return res.Int32()
	}

	return copyIntoInstance(instance, ret, valuePtr, valueSize). Int32()
}
```

More than the Layotto rpc process is briefly described as the implementation of [5]by two virtual connections using the Dapr API and underneath Mosn, see previous order articles [Layotto source parsing — processing RPC requests] (https://mosn.io/layotto/#/blog/code/layotto-rpc/index), where data from Redis can be obtained directly from Dapr State code and is not developed here.

### FaaS Mode

Look back back to the WASM features：bytes code that match the machine code; guarantee good segregation and security in the sandbox; compile cross-platforms, easily distributed, and load running; have lightweight and multilingual flexibilities and seem naturally suitable for FaaS.

So Layotto also explores support for WASM FaaS mode by loading and running WASM carrier functions and supporting interfaces and access to infrastructure between Function.Since the core logic of loading the WASM has not changed, except that there is a difference between usage and deployment methods and those described above, the Layotto load part of the ASM logic is not redundant.

In addition to the Wasm-Proxy implementation, the core logic of the FaaS mode is to manage the \*.wasm package and Kubernetes excellent structuring capabilities by expanding Containerd to multiple-run plugins containerd-shim-layotto-v2 [6]and using this "piercing wire" ingenuity to use Docker mirror capability. Specific structures and workflows can be found in Figure 3 Layotto FaaS Workflow.

![layotto_faas_workflow](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*XWmNT6-7FoEAAAAAAAAAAAAAARQnAQ)

<center>Figure 3 Layotto FaaS Workflow </center>

Here a simple look at the master function of containerd-shim-layotto-v2. It can be seen that shim.Run runs the WASM as io.containerd.layotto.v2, and runtime_type of the containerd plugins.crimerd.runtimes corresponding to the plugin.When creating a Pod, you specify runtimeClassName: layotto in yaml speed, and eventually kubelet will load and run them when cric-plugin calls containerd-shim-layotto-v2 is running.

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

## Summary

Layotto WebAssemly involves more basic WASM knowledge, but it is understandable that the examples are shallow deeper and gradual.At the end of the spectrum, the ASM technology can be seen to have been applied to many fields such as Web-Front, Serverlessness, Game Scene, Edge Computing, Service Grids, or even to the Docker parent Solomon Hykes recently said: "If the WASM technology is available in 2008, I will not be able to do the Docker" (later added that：Docker will not be replaced and will walk side by side with WASM) The ASM seems to be becoming lighter and better performing cloud-origin technology and being applied to more areas after the VM and Container, while believing that there will be more use scenes and users in Mosn community push and in Layotto continue exploration, here Layotto WebAssemly relevant source code analysis has been completed. Given time and length, some more comprehensive and in-depth profiles have not been carried out, and if there are flaws, welcome fingers, contact：rayo. angzl@gmail.com.

### References

- [1] [WebAssembly practice in MOSN](https://mosn.io/blog/posts/mosn-wasm-framework/)
- [2] [feature: WASM plugin framework](https://github.com/mosn/mosn/pull/1589)
- [3] [WebAssembly for Proxies (ABI Spec)](https://github.com/proxy-wasm/spec)
- [4] [Proxy WebAssembly Architecture](https://techhenzy.com/proxy-webassembly-archive/)
- [5] [Layotto source parse — processing RPC requests](https://mosn.io/layotto/#/blog/code/layotto-rpc/index)
- [6] [Cloud native runtime for the next five years](https://www.soft.tech/blog/the-next-fuve-years-of-cloud-native-runtime/)
