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

package runtime

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/pkg/runtime/lifecycle"

	"mosn.io/layotto/components/oss"

	"mosn.io/layotto/pkg/runtime/ref"

	refconfig "mosn.io/layotto/components/ref"

	"github.com/dapr/components-contrib/secretstores"

	"mosn.io/layotto/components/custom"
	msecretstores "mosn.io/layotto/pkg/runtime/secretstores"

	"github.com/dapr/components-contrib/bindings"

	mbindings "mosn.io/layotto/pkg/runtime/bindings"

	"mosn.io/layotto/components/file"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	rawGRPC "google.golang.org/grpc"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/info"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/grpc"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"
	runtime_pubsub "mosn.io/layotto/pkg/runtime/pubsub"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	runtime_state "mosn.io/layotto/pkg/runtime/state"
)

type MosnRuntime struct {
	// configs
	runtimeConfig *MosnRuntimeConfig
	info          *info.RuntimeInfo
	srv           mgrpc.RegisteredServer
	// component registry
	helloRegistry           hello.Registry
	configStoreRegistry     configstores.Registry
	rpcRegistry             rpc.Registry
	pubSubRegistry          runtime_pubsub.Registry
	stateRegistry           runtime_state.Registry
	lockRegistry            runtime_lock.Registry
	sequencerRegistry       runtime_sequencer.Registry
	fileRegistry            file.Registry
	ossRegistry             oss.Registry
	bindingsRegistry        mbindings.Registry
	secretStoresRegistry    msecretstores.Registry
	customComponentRegistry custom.Registry
	Injector                *ref.DefaultInjector
	// component pool
	hellos map[string]hello.HelloService
	// config management system component
	configStores map[string]configstores.Store
	rpcs         map[string]rpc.Invoker
	pubSubs      map[string]pubsub.PubSub
	// state implementations store here are already initialized
	states            map[string]state.Store
	files             map[string]file.File
	oss               map[string]oss.Oss
	locks             map[string]lock.LockStore
	sequencers        map[string]sequencer.Store
	outputBindings    map[string]bindings.OutputBinding
	secretStores      map[string]secretstores.SecretStore
	customComponent   map[string]map[string]custom.Component
	dynamicComponents map[lifecycle.ComponentKey]common.DynamicComponent
	extensionComponents
	// app callback
	AppCallbackConn *rawGRPC.ClientConn
	// extends
	errInt            ErrInterceptor
	started           bool
	initRuntimeStages []initRuntimeStage
}

func (m *MosnRuntime) RuntimeConfig() *MosnRuntimeConfig {
	return m.runtimeConfig
}

type initRuntimeStage func(o *runtimeOptions, m *MosnRuntime) error

func NewMosnRuntime(runtimeConfig *MosnRuntimeConfig) *MosnRuntime {
	info := info.NewRuntimeInfo()
	return &MosnRuntime{
		runtimeConfig:           runtimeConfig,
		info:                    info,
		helloRegistry:           hello.NewRegistry(info),
		configStoreRegistry:     configstores.NewRegistry(info),
		rpcRegistry:             rpc.NewRegistry(info),
		pubSubRegistry:          runtime_pubsub.NewRegistry(info),
		stateRegistry:           runtime_state.NewRegistry(info),
		bindingsRegistry:        mbindings.NewRegistry(info),
		fileRegistry:            file.NewRegistry(info),
		ossRegistry:             oss.NewRegistry(info),
		lockRegistry:            runtime_lock.NewRegistry(info),
		sequencerRegistry:       runtime_sequencer.NewRegistry(info),
		secretStoresRegistry:    msecretstores.NewRegistry(info),
		customComponentRegistry: custom.NewRegistry(info),
		hellos:                  make(map[string]hello.HelloService),
		configStores:            make(map[string]configstores.Store),
		rpcs:                    make(map[string]rpc.Invoker),
		pubSubs:                 make(map[string]pubsub.PubSub),
		states:                  make(map[string]state.Store),
		files:                   make(map[string]file.File),
		oss:                     make(map[string]oss.Oss),
		locks:                   make(map[string]lock.LockStore),
		sequencers:              make(map[string]sequencer.Store),
		outputBindings:          make(map[string]bindings.OutputBinding),
		secretStores:            make(map[string]secretstores.SecretStore),
		customComponent:         make(map[string]map[string]custom.Component),
		dynamicComponents:       make(map[lifecycle.ComponentKey]common.DynamicComponent),
		extensionComponents:     *newExtensionComponents(),
		started:                 false,
	}
}

func (m *MosnRuntime) GetInfo() *info.RuntimeInfo {
	return m.info
}

func (m *MosnRuntime) sendToOutputBinding(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	if req.Operation == "" {
		return nil, errors.New("operation field is missing from request")
	}

	if binding, ok := m.outputBindings[name]; ok {
		ops := binding.Operations()
		for _, o := range ops {
			if o == req.Operation {
				return binding.Invoke(req)
			}
		}
		supported := make([]string, 0, len(ops))
		for _, o := range ops {
			supported = append(supported, string(o))
		}
		return nil, fmt.Errorf("binding %s does not support operation %s. supported operations:%s", name, req.Operation, strings.Join(supported, " "))
	}
	return nil, fmt.Errorf("couldn't find output binding %s", name)
}

func (m *MosnRuntime) Run(opts ...Option) (mgrpc.RegisteredServer, error) {
	// 0. mark already started
	m.started = true
	// 1. init runtime stage
	// prepare runtimeOptions
	o := newRuntimeOptions()
	for _, opt := range opts {
		opt(o)
	}
	// set ErrInterceptor
	if o.errInt != nil {
		m.errInt = o.errInt
	} else {
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
	}
	// init runtime with runtimeOptions
	if err := m.initRuntime(o); err != nil {
		return nil, err
	}
	// prepare grpcOpts
	var grpcOpts []grpc.Option
	if o.srvMaker != nil {
		grpcOpts = append(grpcOpts, grpc.WithNewServer(o.srvMaker))
	}
	// 2. init GrpcAPI stage
	var apis []grpc.GrpcAPI
	ac := newApplicationContext(m)

	for _, apiFactory := range o.apiFactorys {
		api := apiFactory(ac)
		// init the GrpcAPI
		if err := api.Init(m.AppCallbackConn); err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	// put them into grpc options
	grpcOpts = append(grpcOpts,
		grpc.WithGrpcOptions(o.options...),
		grpc.WithGrpcAPIs(apis),
	)
	// 3. create grpc server
	var err error
	m.srv, err = grpc.NewGrpcServer(grpcOpts...)
	return m.srv, err
}

func (m *MosnRuntime) Stop() {
	if m.srv != nil {
		m.srv.Stop()
	}
}

func (m *MosnRuntime) storeDynamicComponent(kind string, name string, store interface{}) {
	comp, ok := store.(common.DynamicComponent)
	if !ok {
		return
	}
	// put it in the components map
	m.dynamicComponents[lifecycle.ComponentKey{
		Kind: kind,
		Name: name,
	}] = lifecycle.ConcurrentDynamicComponent(comp)
}

func DefaultInitRuntimeStage(o *runtimeOptions, m *MosnRuntime) error {
	if m.runtimeConfig == nil {
		return errors.New("[runtime] init error:no runtimeConfig")
	}
	// init callback connection
	if err := m.initAppCallbackConnection(); err != nil {
		return err
	}

	// init all kinds of components with config
	//init secret & config first
	if err := m.initSecretStores(o.services.secretStores...); err != nil {
		return err
	}
	if err := m.initConfigStores(o.services.configStores...); err != nil {
		return err
	}
	m.Injector = ref.NewDefaultInjector(m.secretStores, m.configStores)
	if err := m.initCustomComponents(o.services.custom); err != nil {
		return err
	}
	if err := m.initHellos(o.services.hellos...); err != nil {
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
	if err := m.initOss(o.services.oss...); err != nil {
		return err
	}
	if err := m.initLocks(o.services.locks...); err != nil {
		return err
	}
	if err := m.initSequencers(o.services.sequencers...); err != nil {
		return err
	}
	if err := m.initExtensionComponent(o.services); err != nil {
		return err
	}
	return m.initInputBinding(o.services.inputBinding...)
}

func (m *MosnRuntime) initHellos(hellos ...*hello.HelloFactory) error {
	log.DefaultLogger.Infof("[runtime] init hello service")
	// register all hello services implementation
	m.helloRegistry.Register(hellos...)
	for name, config := range m.runtimeConfig.HelloServiceManagement {
		h, err := m.helloRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create hello's component %s failed", name)
			return err
		}
		//inject component
		if err := m.initComponentInject(h, config.ComponentRef); err != nil {
			return err
		}
		if err := h.Init(&config); err != nil {
			m.errInt(err, "init hello's component %s failed", name)
			return err
		}
		// register this component
		m.hellos[name] = h
		m.storeDynamicComponent(lifecycle.KindHello, name, h)
	}
	return nil
}

func (m *MosnRuntime) initConfigStores(configStores ...*configstores.StoreFactory) error {
	log.DefaultLogger.Infof("[runtime] init config service")
	// register all config store services implementation
	m.configStoreRegistry.Register(configStores...)
	for name, config := range m.runtimeConfig.ConfigStoreManagement {
		c, err := m.configStoreRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create configstore's component %s failed", name)
			return err
		}
		config.AppId = m.runtimeConfig.AppManagement.AppId
		config.StoreName = name
		if err := c.Init(&config); err != nil {
			m.errInt(err, "init configstore's component %s failed", name)
			return err
		}
		// register this component
		m.configStores[name] = c
		m.storeDynamicComponent(lifecycle.KindConfig, name, c)
	}
	return nil
}

func (m *MosnRuntime) initRpcs(rpcs ...*rpc.Factory) error {
	log.DefaultLogger.Infof("[runtime] init rpc service")
	// register all rpc components
	m.rpcRegistry.Register(rpcs...)
	for name, config := range m.runtimeConfig.RpcManagement {
		c, err := m.rpcRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create rpc's component %s failed", name)
			return err
		}
		if err := c.Init(config); err != nil {
			m.errInt(err, "init rpc's component %s failed", name)
			return err
		}
		// register this component
		m.rpcs[name] = c
		m.storeDynamicComponent(lifecycle.KindRPC, name, c)
	}
	return nil
}

func (m *MosnRuntime) initPubSubs(factorys ...*runtime_pubsub.Factory) error {
	// 1. init components
	log.DefaultLogger.Infof("[runtime] start initializing pubsub components")
	// register all components implementation
	m.pubSubRegistry.Register(factorys...)
	for name, config := range m.runtimeConfig.PubSubManagement {
		// create component
		comp, err := m.pubSubRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create pubsub component %s failed", name)
			return err
		}
		// check consumerID
		consumerID := strings.TrimSpace(config.Metadata["consumerID"])
		if consumerID == "" {
			if config.Metadata == nil {
				config.Metadata = make(map[string]string)
			}
			config.Metadata["consumerID"] = m.runtimeConfig.AppManagement.AppId
		}
		//inject secret to component
		if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
			return err
		}
		//inject component
		if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
			return err
		}
		// init this component with the config
		if err := comp.Init(pubsub.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init pubsub component %s failed", name)
			return err
		}
		// register this component
		m.pubSubs[name] = comp
		m.storeDynamicComponent(lifecycle.KindPubsub, name, comp)
	}
	return nil
}

func (m *MosnRuntime) initStates(factorys ...*runtime_state.Factory) error {
	log.DefaultLogger.Infof("[runtime] start initializing state components")
	// 1. register all the implementation
	m.stateRegistry.Register(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.StateManagement {
		// 2.1. create and store the component
		comp, err := m.stateRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create state component %s failed", name)
			return err
		}
		//inject secret to component
		if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
			return err
		}
		//inject component
		if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
			return err
		}
		if err := comp.Init(state.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init state component %s failed", name)
			return err
		}
		m.states[name] = comp
		m.storeDynamicComponent(lifecycle.KindState, name, comp)

		// 2.2. save prefix strategy
		err = runtime_state.SaveStateConfiguration(name, config.Metadata)
		if err != nil {
			log.DefaultLogger.Errorf("error save state keyprefix: %s", err.Error())
			return err
		}
	}
	return nil
}

func (m *MosnRuntime) initOss(factorys ...*oss.Factory) error {
	log.DefaultLogger.Infof("[runtime] init oss service")

	// 1. register all oss store services implementation
	m.ossRegistry.Register(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.Oss {
		// 2.1. create the component
		c, err := m.ossRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create oss component %s failed", name)
			return err
		}
		//inject component
		if err := m.initComponentInject(c, config.ComponentRef); err != nil {
			return err
		}
		// 2.2. init
		if err := c.Init(context.TODO(), &config); err != nil {
			m.errInt(err, "init oss component %s failed", name)
			return err
		}
		// register this component
		m.oss[name] = c
		m.storeDynamicComponent(lifecycle.KindOss, name, c)
	}
	return nil
}

func (m *MosnRuntime) initFiles(files ...*file.FileFactory) error {
	log.DefaultLogger.Infof("[runtime] init file service")

	// register all files store services implementation
	m.fileRegistry.Register(files...)

	for name, config := range m.runtimeConfig.Files {
		c, err := m.fileRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create files component %s failed", name)
			return err
		}
		//inject component
		if err := m.initComponentInject(c, config.ComponentRef); err != nil {
			return err
		}
		if err := c.Init(context.TODO(), &config); err != nil {
			m.errInt(err, "init files component %s failed", name)
			return err
		}
		m.files[name] = c
		m.storeDynamicComponent(lifecycle.KindFile, name, c)
	}
	return nil
}

func (m *MosnRuntime) initLocks(factorys ...*runtime_lock.Factory) error {
	log.DefaultLogger.Infof("[runtime] start initializing lock components")
	// 1. register all the implementation
	m.lockRegistry.Register(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.LockManagement {
		// 2.1. create the component
		comp, err := m.lockRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create lock component %s failed", name)
			return err
		}
		//inject secret to component
		if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
			return err
		}
		//inject component
		if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
			return err
		}
		// 2.2. init
		if err := comp.Init(lock.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init lock component %s failed", name)
			return err
		}
		// 2.3. save runtime related configs
		err = runtime_lock.SaveLockConfiguration(name, config.Metadata)
		if err != nil {
			m.errInt(err, "save lock configuration %s failed", name)
			return err
		}
		m.locks[name] = comp
		m.storeDynamicComponent(lifecycle.KindLock, name, comp)
	}
	return nil
}

func (m *MosnRuntime) initSequencers(factorys ...*runtime_sequencer.Factory) error {
	log.DefaultLogger.Infof("[runtime] start initializing sequencer components")
	// 1. register all the implementation
	m.sequencerRegistry.Register(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.SequencerManagement {
		// 2.1. create the component
		comp, err := m.sequencerRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create sequencer component %s failed", name)
			return err
		}
		//inject secret to component
		if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
			return err
		}
		//inject component
		if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
			return err
		}
		// 2.2. init
		if err = comp.Init(sequencer.Configuration{
			Properties: config.Metadata,
			BiggerThan: config.BiggerThan,
		}); err != nil {
			m.errInt(err, "init sequencer component %s failed", name)
			return err
		}
		// 2.3. save runtime related configs
		err = runtime_sequencer.SaveSeqConfiguration(name, config.Metadata)
		if err != nil {
			m.errInt(err, "save sequencer configuration %s failed", name)
			return err
		}
		// register this component
		m.sequencers[name] = comp
		m.storeDynamicComponent(lifecycle.KindSequencer, name, comp)
	}
	return nil
}

func (m *MosnRuntime) initAppCallbackConnection() error {
	// init the client connection for calling app
	if m.runtimeConfig == nil || m.runtimeConfig.AppManagement.GrpcCallbackPort == 0 {
		return nil
	}
	// get callback address
	host := m.runtimeConfig.AppManagement.GrpcCallbackHost
	if host == "" {
		host = "127.0.0.1"
	}
	port := m.runtimeConfig.AppManagement.GrpcCallbackPort
	address := fmt.Sprintf("%v:%v", host, port)
	opts := []rawGRPC.DialOption{
		rawGRPC.WithInsecure(),
	}
	// dial
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := rawGRPC.DialContext(ctx, address, opts...)
	if err != nil {
		log.DefaultLogger.Warnf("[runtime]failed to init callback client to address %v : %s", address, err)
		return err
	}
	m.AppCallbackConn = conn
	return nil
}

func (m *MosnRuntime) initOutputBinding(factorys ...*mbindings.OutputBindingFactory) error {
	log.DefaultLogger.Infof("[runtime] start initializing OutputBinding components")
	// 1. register all factory methods.
	m.bindingsRegistry.RegisterOutputBinding(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.Bindings {
		// 2.1. create the component
		comp, err := m.bindingsRegistry.CreateOutputBinding(config.Type)
		if err != nil {
			m.errInt(err, "create outbinding component %s failed", name)
			return err
		}
		//inject secret to component
		if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
			return err
		}
		//inject component
		if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
			return err
		}
		// 2.2. init
		if err := comp.Init(bindings.Metadata{Name: name, Properties: config.Metadata}); err != nil {
			m.errInt(err, "init outbinding component %s failed", name)
			return err
		}
		// 2.3. put it into the runtime component pool
		m.outputBindings[name] = comp
		m.storeDynamicComponent(lifecycle.KindBinding, name, comp)
	}
	return nil
}

// TODO: implement initInputBinding
func (m *MosnRuntime) initInputBinding(factorys ...*mbindings.InputBindingFactory) error {
	return nil
}

func (m *MosnRuntime) initSecretStores(factorys ...*msecretstores.SecretStoresFactory) error {
	log.DefaultLogger.Infof("[runtime] start initializing SecretStores components")
	// 1. register all factory methods.
	m.secretStoresRegistry.Register(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.SecretStoresManagement {
		// 2.1. create the component
		comp, err := m.secretStoresRegistry.Create(config.Type)
		if err != nil {
			m.errInt(err, "create secretStore component %s failed", name)
			return err
		}
		//inject component
		if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
			return err
		}
		// 2.2. init
		if err := comp.Init(secretstores.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init secretStore component %s failed", name)
			return err
		}

		// 2.3. save runtime related configs
		m.secretStores[name] = comp
		m.storeDynamicComponent(lifecycle.KindSecret, name, comp)
	}
	return nil
}

func (m *MosnRuntime) AppendInitRuntimeStage(f initRuntimeStage) {
	if f == nil || m.started {
		log.DefaultLogger.Errorf("[runtime] invalid initRuntimeStage or already started")
		return
	}
	m.initRuntimeStages = append(m.initRuntimeStages, f)
}

func (m *MosnRuntime) initRuntime(r *runtimeOptions) error {
	st := time.Now()
	// check default handler
	if len(m.initRuntimeStages) == 0 {
		m.initRuntimeStages = append(m.initRuntimeStages, DefaultInitRuntimeStage)
	}
	// do initialization
	for _, f := range m.initRuntimeStages {
		err := f(r, m)
		if err != nil {
			return err
		}
	}

	log.DefaultLogger.Infof("[runtime] initRuntime stages cost: %v", time.Since(st))
	return nil
}

func (m *MosnRuntime) SetCustomComponent(kind string, name string, component custom.Component) {
	if _, ok := m.customComponent[kind]; !ok {
		m.customComponent[kind] = make(map[string]custom.Component)
	}
	m.customComponent[kind][name] = component
}

func (m *MosnRuntime) initCustomComponents(kind2factorys map[string][]*custom.ComponentFactory) error {
	log.DefaultLogger.Infof("[runtime] start initializing custom components")
	// loop all configured custom components.
	for kind, name2Config := range m.runtimeConfig.CustomComponent {
		factorys, ok := kind2factorys[kind]
		if !ok || len(factorys) == 0 {
			log.DefaultLogger.Errorf("[runtime] Your required component kind %s is not supported.", kind)
			continue
		}
		// register all the factorys
		m.customComponentRegistry.Register(kind, factorys...)
		// loop initializing component instances
		for name, config := range name2Config {
			// create the component
			comp, err := m.customComponentRegistry.Create(kind, config.Type)
			if err != nil {
				m.errInt(err, "create custom component %s failed", name)
				return err
			}
			//inject secret to component
			if config.Metadata, err = m.Injector.InjectSecretRef(config.SecretRef, config.Metadata); err != nil {
				return err
			}
			//inject component
			if err := m.initComponentInject(comp, config.ComponentRef); err != nil {
				return err
			}
			// init
			if err := comp.Initialize(context.TODO(), config); err != nil {
				m.errInt(err, "init custom component %s failed", name)
				return err
			}
			// initialization finish
			m.SetCustomComponent(kind, name, comp)
			m.storeDynamicComponent(fmt.Sprintf("%s.%s", lifecycle.KindCustom, kind), name, comp)
		}

	}

	return nil
}

func (m *MosnRuntime) initComponentInject(comp interface{}, config *refconfig.ComponentRefConfig) error {
	if setComp, ok := comp.(common.SetComponent); ok {
		configRef, err := m.Injector.GetConfigStore(config)
		if err != nil {
			return err
		}
		if configRef != nil {
			if err = setComp.SetConfigStore(configRef); err != nil {
				return err
			}
		}

		secretRef, err := m.Injector.GetSecretStore(config)
		if err != nil {
			return err
		}
		if secretRef != nil {
			if err = setComp.SetSecretStore(secretRef); err != nil {
				return err
			}
		}
	}
	return nil
}
