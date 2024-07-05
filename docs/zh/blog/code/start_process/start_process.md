# 源码解析  layotto启动流程

>作者简介：
>张立斌，https://github.com/ZLBer
>
>写作时间: 2022年5月4日


- [Overview](#Overview)
- [源码分析](#源码分析)
    * [1.cmd分析](#1.cmd分析)
    * [2.回调函数NewRuntimeGrpcServer分析](#2.回调函数NewRuntimeGrpcServer分析)
    * [3.runtime分析](#3.runtime分析)
- [总结](#总结)

## Overview
Layotto “寄生”在 MOSN 里，启动流程其实是先启动 MOSN, MOSN 在启动过程中回调 Layotto ，让 Layotto 启动。

## 源码分析
一切起源于我们的命令行: layotto start  -c  `configpath`

### 1.cmd分析

main 的 init 函数首先运行：

```
func init() {   
     //将layotto的初始化函数传给mosn，让mosn启动的时候进行回调
	mgrpc.RegisterServerHandler("runtime", NewRuntimeGrpcServer)
     ....
}
```

cmd 的 action 开始执行：

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

### 2.回调函数NewRuntimeGrpcServer分析

MOSN 启动的时候回调 NewRuntimeGrpcServer ，data 是未解析的配置文件，opts 是 grpc 的配置项，返回 Grpc server

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

### 3.runtime分析

看一下 runtime 的结构体，从整体上把握 runtime 的构成：

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

runtime 的 run 函数逻辑如下:

```
func (m *MosnRuntime) Run(opts ...Option) (mgrpc.RegisteredServer, error) {
	// 启动标志
	m.started = true
	// 新建runtime配置
	o := newRuntimeOptions()
	//这里运行我们之前传入的option函数，其实就是将各种组件Factory注册进来
	for _, opt := range opts {
		opt(o)
	}
	//初始化组件
	if err := m.initRuntime(o); err != nil {
		return nil, err
	}
	
	//初始化Grpc，api赋值
	var grpcOpts []grpc.Option
	if o.srvMaker != nil {
		grpcOpts = append(grpcOpts, grpc.WithNewServer(o.srvMaker))
	}
	var apis []grpc.GrpcAPI
	ac := &grpc.ApplicationContext{
		m.runtimeConfig.AppManagement.AppId,
		m.hellos,
		m.configStores,
		m.rpcs,
		m.pubSubs,
		m.states,
		m.files,
		m.locks,
		m.sequencers,
		m.sendToOutputBinding,
		m.secretStores,
		m.customComponent,
	}
     //调用组件的factory生成每个组件
	for _, apiFactory := range o.apiFactorys {
		api := apiFactory(ac)
		// init the GrpcAPI
		if err := api.Init(m.AppCallbackConn); err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	// 将api接口和配置传给grpc
	grpcOpts = append(grpcOpts,
		grpc.WithGrpcOptions(o.options...),
		grpc.WithGrpcAPIs(apis),
	)
	//启动grpc
	var err error = nil
	m.srv, err = grpc.NewGrpcServer(grpcOpts...)
	return m.srv, err
}

```

组件的初始化函数 initRuntime ：

```
func (m *MosnRuntime) initRuntime(r *runtimeOptions) error {
	st := time.Now()
	if len(m.initRuntimeStages) == 0 {
		m.initRuntimeStages = append(m.initRuntimeStages, DefaultInitRuntimeStage)
	}
	// 调用DefaultInitRuntimeStage
	for _, f := range m.initRuntimeStages {
		err := f(r, m)
		if err != nil {
			return err
		}
	}
    ...
	return nil
}
```

DefaultInitRuntimeStage 组件初始化逻辑，调用每个组件的 init 方法:

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

以 file 组件为例，看下初始化函数：

```
func (m *MosnRuntime) initFiles(files ...*file.FileFactory) error {

	//将配置的组件注册进去
	m.fileRegistry.Register(files...)
	for name, config := range m.runtimeConfig.Files {
	    //create调用NewXXX()函数新建一个组件实例
		c, err := m.fileRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create files component %s failed", name)
			return err
		}
		if err := c.Init(context.TODO(), &config); err != nil {
			m.errInt(err, "init files component %s failed", name)
			return err
		}
		//赋值给runtime
		m.files[name] = c
	}
	return nil
}
```

至此 MOSN、Grpc、Layotto 都已经启动完成，通过 Grpc 的接口就可以调用到组件的代码逻辑。

## 总结
总览整个启动流程，Layotto 结合 MOSN 来做启动，解析配置文件，生成配置文件中的组件类，将 Grpc 的 api 暴露出去。





