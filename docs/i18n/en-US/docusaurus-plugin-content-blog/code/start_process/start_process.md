# Source parsing layotto startup process

> Author Intro to：
> Libin, https://github.com/ZLBer
>
> Writing: 4 May 2022

- [Overview](#Overview)
- [source analysis](#source analysis)
  - [1.cmd analysis](#1.cmd analysis)
  - [2.Callback functionNewRuntimeGrpcServer分析](#2.callback function NewRuntimeGrpcServer analysis)
  - [3.runtimeanalyze](#3.runtime analyse)
- [summary](#summary)

## Overview

Layotto "Parasite" in MOSN. The start process is in effect starting MOSN, MOSN back Layotto during startup to get Layotto start.

## Source analysis

Everything originating from our command line: layotto start -c `configpath`

### 1.cmd analysis

Main init function starts with：

```
func init() {   
     //将layotto的初始化函数传给mosn，让mosn启动的时候进行回调
	mgrpc.RegisterServerHandler("runtime", NewRuntimeGrpcServer)
     ....
}
```

cmd action starts to execute：

```
	Action: func(c *cli.Context) error {
		app := mosn.NewMosn()
		//stagemanager用于管理mosn启动的每个阶段，可以添加相应的阶段函数，比如下面的ParamsParsedStage、InitStage、PreStartStage、AfterStartStage
		//这里是将configpath传给mosn，下面都是mosn相关的逻辑
		stm := stagemanager.InitStageManager(c, c.String("config"), app) 
		stm.AppendParamsParsedStage(ExtensionsRegister)
		stm.AppendParamsParsedStage(func(c *cli.Context) {
			err := featuregate.Set(c.String("feature-gates"))
			if err != nil {
				os.Exit(1)
			}
		})·
		stm.AppendInitStage(mosn.DefaultInitStage)
		stm.AppendPreStartStage(mosn.DefaultPreStartStage)
		stm.AppendStartStage(mosn.DefaultStartStage)
		//这里添加layotto的健康检查机制
		stm.AppendAfterStartStage(SetActuatorAfterStart)
		stm.Run()
		// wait mosn finished
		stm.WaitFinish()
		return nil
	},
```

### NewRuntimeGrpcServer Analysis

Returns NewRuntimeGrpcServer when MOSN is launched, data is an unparsed configuration, opts is a grpc configuration, returning Gpc server

```
func NewRuntimeGrpcServer(data json.RawMessage, opts ...grpc.ServerOption) (mgrpc.RegisteredServer, error) {
	// 将原始的配置文件解析成结构体形式。
	cfg, err := runtime.ParseRuntimeConfig(data)
    // 新建layotto runtime， runtime包含各种组件的注册器和各种组件的实例。
	rt := runtime.NewMosnRuntime(cfg)
	// 3.runtime开始启动
	server, err := rt.Run(
	       ...
        // 4. 添加所有组件的初始化函数
	 	// 我们只看下File组件的，将NewXXX()添加到组件Factory里
		runtime.WithFileFactory(
			file.NewFileFactory("aliyun.oss", alicloud.NewAliCloudOSS),
			file.NewFileFactory("minio", minio.NewMinioOss),
			file.NewFileFactory("aws.s3", aws.NewAwsOss),
			file.NewFileFactory("tencent.oss", tencentcloud.NewTencentCloudOSS),
			file.NewFileFactory("local", local.NewLocalStore),
			file.NewFileFactory("qiniu.oss", qiniu.NewQiniuOSS),
		),
	     ...
   return server, err		 
	
	)
	
	//
}

```

### runtime analysis

Look at the structure of runtime, the composition of the `runtime' at the aggregate level of the`：'

```
type MosnRuntime struct {
	// 包括组件的config
	runtimeConfig *MosnRuntimeConfig
	info          *info.RuntimeInfo
	srv           mgrpc.RegisteredServer
	// 组件注册器，用来注册和新建组件，里面有组件的NewXXX()函数
	helloRegistry           hello.Registry
	configStoreRegistry     configstores.Registry
	rpcRegistry             rpc.Registry
	pubSubRegistry          runtime_pubsub.Registry
	stateRegistry           runtime_state.Registry
	lockRegistry            runtime_lock.Registry
	sequencerRegistry       runtime_sequencer.Registry
	fileRegistry            file.Registry
	bindingsRegistry        mbindings.Registry
	secretStoresRegistry    msecretstores.Registry
	customComponentRegistry custom.Registry
	hellos map[string]hello.HelloService
	// 各种组件
	configStores map[string]configstores.Store
	rpcs         map[string]rpc.Invoker
	pubSubs      map[string]pubsub.PubSub
	states          map[string]state.Store
	files           map[string]file.File
	locks           map[string]lock.LockStore
	sequencers      map[string]sequencer.Store
	outputBindings  map[string]bindings.OutputBinding
	secretStores    map[string]secretstores.SecretStore
	customComponent map[string]map[string]custom.Component
	AppCallbackConn *rawGRPC.ClientConn
	errInt            ErrInterceptor
	started           bool
	//初始化函数
	initRuntimeStages []initRuntimeStage
}
```

runtime is the run function logic as follows:

```
func (m *MosnRuntime) Run(opts..Option) (mgrpc.RegisteredServer, error) um
	// launch flag
	m. targeted = true
	// newly created runtime configuration
	o := newRuntimeOptions()
	// run our previously imported option,. Really register various components Factory with
	for _, opt := range opts {
		opt(o)
	}
	//initialization component
	if err := m. nitRuntime(o); err != nil {
		return nil, err
	}
	
	//initialize Grpc,api assignment
	var grpcOpts[]grpc. Absorption
	if o.srvMaker != nil LO
		grpcOpts = append(grpcOpts, grpc.GithNewServer(o.srvMaker))
	}
	var apis []grpc.GrpcAPI
	ac := &grpc. pimplicationContextFe
		m.runtimeConfig.AppManagement.AppId,
		m.hellos,
		m.configStories,
		m.rpcs,
		m.pubSubs,
		m. tates,
		m.files,
		m.locks,
		m.sequencers,
		m.sendToOutputBinding,
		m.secretStories,
		m. ustomCompany,
	}
     // Factor generation of each component
	for _, apiFactory := range o. piFactorys LOR
		api := apiFactory(ac)
		// init the GrpcAPI
		if err := api.Init(m. ppCallbackCon); err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	// pass the api interface and configuration to grpc
	grpcOpts = append(grpcOpts,
		grpc.GrpOptions(o.options... ,
		grpc.MithGrpcAPIs(apis),
	)
	//start grpc
	var err error = nil
	m. rv, err = grpc.NewGrpServer (grpcOpts...)
	return m.srv, err
}

```

Component initialization function initRuntime ：

```
func (m *MosnRuntime) initRuntime (r *runtimeOptions) errant error LO
	st := time.Now()
	if len(m.initRuntimeStages) === 0 56
		m.initRuntimeStages = append(m. nitRuntimeStages, DefaultInitRuntimeStage
	}
	// Call DefaultInitRuntimeStage
	for _, f := range m. nitRuntime Stages FEM
		err := f(r, m)
		if err != nil {
			return err
		}
	}
    . .
	return nil
}
```

DefaultInitRuntimeStage component initialization logic, call init method for each component:

```
func DefaultInitRuntimeStage(o *runtimeOptions, m *MosnRuntime) error {
	 ...
	 //初始化config/state/file/lock/sequencer/secret等各种组件
	if err := m.initCustomComponents(o.services.custom); err != nil {
		return err
	}
	if err := m.initHellos(o.services.hellos...); err != nil {
		return err
	}
	if err := m.initConfigStores(o.services.configStores...); err != nil {
		return err
	}
	if err := m.initStates(o.services.states...); err != nil {
		return err
	}
	if err := m.initRpcs(o.services.rpcs...); err != nil {
		return err
	}
	if err := m.initOutputBinding(o.services.outputBinding...); err != nil {
		return err
	}
	if err := m.initPubSubs(o.services.pubSubs...); err != nil {
		return err
	}
	if err := m.initFiles(o.services.files...); err != nil {
		return err
	}
	if err := m.initLocks(o.services.locks...); err != nil {
		return err
	}
	if err := m.initSequencers(o.services.sequencers...); err != nil {
		return err
	}
	if err := m.initInputBinding(o.services.inputBinding...); err != nil {
		return err
	}
	if err := m.initSecretStores(o.services.secretStores...); err != nil {
		return err
	}
	return nil
}
```

Example file component, see initialization function：

```
func (m *MosnRuntime) initFiles(files ...file.FileFactory) ERRORY ERROR LO

	//register configured components on
	m.fileRegistry.Register(files...)
	for name, config := range m. untimesConfig.Files Fact
	    //create/create a new component instance
		c, err := m.fileRegistry.Create(name)
		if err !=nil L/
			m. rrInt(err, "creation files component %s failed", name)
			return err
		}
		if err := c. nit(context.TODO(), &config); err != nil LO
			m. rrInt(err, "init files component %s failed", name)
			return err
		}
		//assignment to runtime
		m. files[name] = c
	}
	return nil
}
```

Here MOSN, Grpc and Layotto are all started, and the code logic of the component can be called through the Gypc interface.

## Summary

Overall view of the entire startup process, Layotto integration with MOSN to start, parse configuration files, generate component classes in the configuration file and expose the api of Grpc.
