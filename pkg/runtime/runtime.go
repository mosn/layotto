package runtime

import (
	"context"
	"encoding/base64"
	"github.com/dapr/components-contrib/contenttype"
	"github.com/dapr/components-contrib/pubsub"
	jsoniter "github.com/json-iterator/go"
	"github.com/layotto/layotto/pkg/grpc"
	"github.com/layotto/layotto/pkg/info"
	"github.com/layotto/layotto/pkg/integrate/actuator"
	runtime_pubsub "github.com/layotto/layotto/pkg/runtime/pubsub"
	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/pkg/services/hello"
	pubsub_service "github.com/layotto/layotto/pkg/services/pubsub"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	topicRoutes         map[string]TopicRoute
	grpc                *grpc.Manager
	// extends
	errInt ErrInterceptor
	json   jsoniter.API
}

type Route struct {
	path     string
	metadata map[string]string
}

type TopicRoute struct {
	routes map[string]Route
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
		grpc:                grpc.NewManager(),
		json:                jsoniter.ConfigFastest,
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
	// 1. init components
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
	// 2. init the client for calling app
	if m.runtimeConfig != nil && m.runtimeConfig.AppManagement.GrpcCallbackPort > 0 {
		port := m.runtimeConfig.AppManagement.GrpcCallbackPort
		err := m.grpc.InitAppClient(port)
		if err != nil {
			log.DefaultLogger.Warnf("[runtime]failed to init callback client at port %v : %s", port, err)
		}
	}
	// 3. start subscribing
	return m.startSubscribing()
}

func (m *MosnRuntime) startSubscribing() error {
	for name, pubsub := range m.pubSubs {
		if err := m.beginPubSub(name, pubsub); err != nil {
			return err
		}
	}
	return nil
}

func (m *MosnRuntime) beginPubSub(name string, ps pubsub.PubSub) error {
	publishFunc := m.publishMessageGRPC
	// 1. call app to find topic routes.
	topicRoutes, err := m.getTopicRoutes()
	if err != nil {
		return err
	}
	v, ok := topicRoutes[name]
	if !ok {
		return nil
	}
	// 2. loop subscribing every <topic, route>
	for topic, route := range v.routes {
		// TODO limit topic scope
		log.DefaultLogger.Debugf("[runtime][beginPubSub]subscribing to topic=%s on pubsub=%s", topic, name)

		// ask component to subscribe
		if err := ps.Subscribe(pubsub.SubscribeRequest{
			Topic:    topic,
			Metadata: route.metadata,
		}, func(msg *pubsub.NewMessage) error {
			if msg.Metadata == nil {
				msg.Metadata = make(map[string]string, 1)
			}

			msg.Metadata[Metadata_key_pubsubName] = name
			return publishFunc(msg)
		}); err != nil {
			log.DefaultLogger.Warnf("[runtime][beginPubSub]failed to subscribe to topic %s: %s", topic, err)
			return err
		}
	}

	return nil
}

func (m *MosnRuntime) getTopicRoutes() (map[string]TopicRoute, error) {
	if m.topicRoutes != nil {
		return m.topicRoutes, nil
	}

	topicRoutes := make(map[string]TopicRoute)
	var subscriptions []runtime_pubsub.Subscription

	// handle app subscriptions
	client := runtimev1pb.NewAppCallbackClient(m.grpc.AppClientConn)
	subscriptions = runtime_pubsub.GetSubscriptionsGRPC(client, log.DefaultLogger)
	// TODO handle declarative subscriptions

	// prepare result
	for _, s := range subscriptions {
		if _, ok := topicRoutes[s.PubsubName]; !ok {
			topicRoutes[s.PubsubName] = TopicRoute{routes: make(map[string]Route)}
		}

		topicRoutes[s.PubsubName].routes[s.Topic] = Route{path: s.Route, metadata: s.Metadata}
	}

	// log
	if len(topicRoutes) > 0 {
		for pubsubName, v := range topicRoutes {
			topics := []string{}
			for topic := range v.routes {
				topics = append(topics, topic)
			}
			log.DefaultLogger.Infof("[runtime][getTopicRoutes]app is subscribed to the following topics: %v through pubsub=%s", topics, pubsubName)
		}
	}
	m.topicRoutes = topicRoutes
	return topicRoutes, nil
}

func (a *MosnRuntime) publishMessageGRPC(msg *pubsub.NewMessage) error {
	// 1. Unmarshal to cloudEvent model
	var cloudEvent map[string]interface{}
	err := a.json.Unmarshal(msg.Data, &cloudEvent)
	if err != nil {
		log.DefaultLogger.Debugf("[runtime]error deserializing cloud events proto: %s", err)
		return err
	}

	// 2. Drop msg if the current cloud event has expired
	if pubsub.HasExpired(cloudEvent) {
		log.DefaultLogger.Warnf("[runtime]dropping expired pub/sub event %v as of %v", cloudEvent[pubsub.IDField].(string), cloudEvent[pubsub.ExpirationField].(string))
		return nil
	}

	// 3. Convert to proto domain struct
	envelope := &runtimev1pb.TopicEventRequest{
		Id:              cloudEvent[pubsub.IDField].(string),
		Source:          cloudEvent[pubsub.SourceField].(string),
		DataContentType: cloudEvent[pubsub.DataContentTypeField].(string),
		Type:            cloudEvent[pubsub.TypeField].(string),
		SpecVersion:     cloudEvent[pubsub.SpecVersionField].(string),
		Topic:           msg.Topic,
		PubsubName:      msg.Metadata[Metadata_key_pubsubName],
	}

	// set data field
	if data, ok := cloudEvent[pubsub.DataBase64Field]; ok && data != nil {
		decoded, decodeErr := base64.StdEncoding.DecodeString(data.(string))
		if decodeErr != nil {
			log.DefaultLogger.Debugf("unable to base64 decode cloudEvent field data_base64: %s", decodeErr)
			return err
		}

		envelope.Data = decoded
	} else if data, ok := cloudEvent[pubsub.DataField]; ok && data != nil {
		envelope.Data = nil

		if contenttype.IsStringContentType(envelope.DataContentType) {
			envelope.Data = []byte(data.(string))
		} else if contenttype.IsJSONContentType(envelope.DataContentType) {
			envelope.Data, _ = a.json.Marshal(data)
		}
	}
	// TODO tracing

	// 4. Call appcallback
	ctx := context.Background()
	clientV1 := runtimev1pb.NewAppCallbackClient(a.grpc.AppClientConn)
	res, err := clientV1.OnTopicEvent(ctx, envelope)

	// 5. Check result
	if err != nil {
		errStatus, hasErrStatus := status.FromError(err)
		if hasErrStatus && (errStatus.Code() == codes.Unimplemented) {
			// DROP
			log.DefaultLogger.Warnf("[runtime]non-retriable error returned from app while processing pub/sub event %v: %s", cloudEvent[pubsub.IDField].(string), err)
			return nil
		}

		err = errors.Errorf("error returned from app while processing pub/sub event %v: %s", cloudEvent[pubsub.IDField].(string), err)
		log.DefaultLogger.Debugf("%s", err)
		// on error from application, return error for redelivery of event
		return err
	}

	switch res.GetStatus() {
	case runtimev1pb.TopicEventResponse_SUCCESS:
		// on uninitialized status, this is the case it defaults to as an uninitialized status defaults to 0 which is
		// success from protobuf definition
		return nil
	case runtimev1pb.TopicEventResponse_RETRY:
		return errors.Errorf("RETRY status returned from app while processing pub/sub event %v", cloudEvent[pubsub.IDField].(string))
	case runtimev1pb.TopicEventResponse_DROP:
		log.DefaultLogger.Warnf("[runtime]DROP status returned from app while processing pub/sub event %v", cloudEvent[pubsub.IDField].(string))
		return nil
	}
	// Consider unknown status field as error and retry
	return errors.Errorf("unknown status returned from app while processing pub/sub event %v: %v", cloudEvent[pubsub.IDField].(string), res.GetStatus())
}
