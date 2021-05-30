package runtime

import (
	"github.com/dapr/components-contrib/pubsub"
	"github.com/layotto/layotto/pkg/grpc"
	"github.com/layotto/layotto/pkg/info"
	"github.com/layotto/layotto/pkg/integrate/actuator"
	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/pkg/services/hello"
	pubsub_service "github.com/layotto/layotto/pkg/services/pubsub"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	"mosn.io/pkg/log"
)

type MosnRuntime struct {
	// configs
	runtimeConfig *MosnRuntimeConfig
	info          *info.RuntimeInfo
	srv           mgrpc.RegisteredServer
	// services
	helloRegistry       hello.Registry
	configStoreRegistry configstores.Registry
	pubSubRegistry      pubsub_service.Registry
	hellos              map[string]hello.HelloService
	configStores        map[string]configstores.Store
	pubSubs             map[string]pubsub.PubSub
	// extends
	errInt ErrInterceptor
}

func NewMosnRuntime(runtimeConfig *MosnRuntimeConfig) *MosnRuntime {
	info := info.NewRuntimeInfo()
	return &MosnRuntime{
		runtimeConfig:       runtimeConfig,
		info:                info,
		helloRegistry:       hello.NewRegistry(info),
		configStoreRegistry: configstores.NewRegistry(info),
		hellos:              make(map[string]hello.HelloService),
		configStores:        make(map[string]configstores.Store),
		pubSubs:             make(map[string]pubsub.PubSub),
	}
}

func (m *MosnRuntime) GetInfo() *info.RuntimeInfo {
	return m.info
}

func (m *MosnRuntime) Run(opts ...Option) (mgrpc.RegisteredServer, error) {
	var o runtimeOptions
	for _, opt := range opts {
		opt(&o)
	}
	if o.errInt != nil {
		m.errInt = o.errInt
	} else {
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
	}

	if err := m.initRuntime(&o); err != nil {
		return nil, err
	}
	var grpcOpts []grpc.Option
	if o.srvMaker != nil {
		grpcOpts = append(grpcOpts, grpc.WithNewServer(o.srvMaker))
	}
	// TODO: support NewAPI extends
	grpcOpts = append(grpcOpts,
		grpc.WithGrpcOptions(o.options...),
		grpc.WithAPI(grpc.NewAPI(
			m.hellos,
			m.configStores,
			m.pubSubs,
		)),
	)
	m.srv = grpc.NewGrpcServer(grpcOpts...)
	return m.srv, nil
}

func (m *MosnRuntime) Stop() {
	if m.srv != nil {
		m.srv.Stop()
	}
	actuator.GetRuntimeReadinessIndicator().SetUnhealthy("shutdown")
	actuator.GetRuntimeLivenessIndicator().SetUnhealthy("shutdown")
}

func (m *MosnRuntime) initRuntime(o *runtimeOptions) error {
	// init hello implementation by config
	if err := m.initHellos(o.services.hellos...); err != nil {
		return err
	}
	if err := m.initConfigStores(o.services.configStores...); err != nil {
		return err
	}
	if err := m.initPubSubs(o.services.pubSubs...); err != nil {
		return err
	}
	return nil
}

func (m *MosnRuntime) initHellos(hellos ...*hello.HelloFactory) error {
	log.DefaultLogger.Infof("[runtime] init hello service")
	// register all hello services implementation
	m.helloRegistry.Register(hellos...)
	for name, config := range m.runtimeConfig.HelloServiceManagement {
		h, err := m.helloRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create hello's component %s failed", name)
			return err
		}
		if err := h.Init(&config); err != nil {
			m.errInt(err, "init hello's component %s failed", name)
			return err
		}
		m.hellos[name] = h
	}
	return nil
}

func (m *MosnRuntime) initConfigStores(configStores ...*configstores.StoreFactory) error {
	log.DefaultLogger.Infof("[runtime] init config service")
	// register all config store services implementation
	m.configStoreRegistry.Register(configStores...)
	for name, config := range m.runtimeConfig.ConfigStoreManagement {
		c, err := m.configStoreRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create configstore's component %s failed", name)
			return err
		}
		if err := c.Init(&config); err != nil {
			m.errInt(err, "init configstore's component %s failed", name)
			return err
		}
		m.configStores[name] = c
	}
	return nil
}

func (m *MosnRuntime) initPubSubs(factorys ...*pubsub_service.Factory) error {
	//	1. init components
	log.DefaultLogger.Infof("[runtime] init config service")
	// register all config store services implementation
	m.pubSubRegistry.Register(factorys...)
	for name, config := range m.runtimeConfig.PubSubManagement {
		comp, err := m.pubSubRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create configstore's component %s failed", name)
			return err
		}
		if err := comp.Init(pubsub.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init configstore's component %s failed", name)
			return err
		}
		m.pubSubs[name] = comp
	}
	//	2. start subscribing

}
