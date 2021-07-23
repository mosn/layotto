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
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/contenttype"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	jsoniter "github.com/json-iterator/go"
	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/actuators"
	"mosn.io/layotto/components/pkg/info"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/actuator/health"
	"mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/pkg/integrate/actuator"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"
	runtime_pubsub "mosn.io/layotto/pkg/runtime/pubsub"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	runtime_state "mosn.io/layotto/pkg/runtime/state"
	"mosn.io/layotto/pkg/wasm"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	"mosn.io/pkg/log"
	"strings"
)

type MosnRuntime struct {
	// configs
	runtimeConfig *MosnRuntimeConfig
	info          *info.RuntimeInfo
	srv           mgrpc.RegisteredServer
	// component registry
	helloRegistry       hello.Registry
	configStoreRegistry configstores.Registry
	rpcRegistry         rpc.Registry
	pubSubRegistry      runtime_pubsub.Registry
	stateRegistry       runtime_state.Registry
	lockRegistry        runtime_lock.Registry
	sequencerRegistry   runtime_sequencer.Registry
	// component pool
	hellos            map[string]hello.HelloService
	configStores      map[string]configstores.Store
	rpcs              map[string]rpc.Invoker
	pubSubs           map[string]pubsub.PubSub
	topicPerComponent map[string]TopicSubscriptions
	states            map[string]state.Store
	locks             map[string]lock.LockStore
	sequencers        map[string]sequencer.Store
	// app callback
	AppCallbackConn *rawGRPC.ClientConn
	// extends
	errInt ErrInterceptor
	json   jsoniter.API
}

type Details struct {
	metadata map[string]string
}

type TopicSubscriptions struct {
	topic2Details map[string]Details
}

func NewMosnRuntime(runtimeConfig *MosnRuntimeConfig) *MosnRuntime {
	info := info.NewRuntimeInfo()
	return &MosnRuntime{
		runtimeConfig:       runtimeConfig,
		info:                info,
		helloRegistry:       hello.NewRegistry(info),
		configStoreRegistry: configstores.NewRegistry(info),
		rpcRegistry:         rpc.NewRegistry(info),
		pubSubRegistry:      runtime_pubsub.NewRegistry(info),
		stateRegistry:       runtime_state.NewRegistry(info),
		lockRegistry:        runtime_lock.NewRegistry(info),
		sequencerRegistry:   runtime_sequencer.NewRegistry(info),
		hellos:              make(map[string]hello.HelloService),
		configStores:        make(map[string]configstores.Store),
		rpcs:                make(map[string]rpc.Invoker),
		pubSubs:             make(map[string]pubsub.PubSub),
		states:              make(map[string]state.Store),
		locks:               make(map[string]lock.LockStore),
		sequencers:          make(map[string]sequencer.Store),
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
	wasm.Layotto = grpc.NewAPI(
		m.runtimeConfig.AppManagement.AppId,
		m.hellos,
		m.configStores,
		m.rpcs,
		m.pubSubs,
		m.states,
		m.locks,
		m.sequencers,
	)
	grpcOpts = append(grpcOpts,
		grpc.WithGrpcOptions(o.options...),
		grpc.WithAPI(wasm.Layotto),
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
	if m.runtimeConfig == nil {
		return errors.New("[runtime] init error:no runtimeConfig")
	}
	// init callback connection
	if err := m.initAppCallbackConnection(); err != nil {
		return err
	}
	// init all kinds of components with config
	if err := m.initHellos(o.services.hellos...); err != nil {
		return err
	}
	if err := m.initConfigStores(o.services.configStores...); err != nil {
		return err
	}
	if err := m.initRpcs(o.services.rpcs...); err != nil {
		return err
	}
	if err := m.initPubSubs(o.services.pubSubs...); err != nil {
		return err
	}
	if err := m.initStates(o.services.states...); err != nil {
		return err
	}
	if err := m.initLocks(o.services.locks...); err != nil {
		return err
	}
	if err := m.initSequencers(o.services.sequencers...); err != nil {
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
		v := actuators.GetIndicatorWithName(name)
		//Now don't force user implement actuator of components
		if v != nil {
			health.AddLivenessIndicator(name, v.LivenessIndicator)
			health.AddReadinessIndicator(name, v.ReadinessIndicator)
		}
	}
	return nil
}

func (m *MosnRuntime) initRpcs(rpcs ...*rpc.Factory) error {
	log.DefaultLogger.Infof("[runtime] init rpc service")
	// register all config store services implementation
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
		m.rpcs[name] = c
	}
	return nil
}

func (m *MosnRuntime) initPubSubs(factorys ...*runtime_pubsub.Factory) error {
	// 1. init components
	log.DefaultLogger.Infof("[runtime] start initializing pubsub components")
	// register all config store services implementation
	m.pubSubRegistry.Register(factorys...)
	for name, config := range m.runtimeConfig.PubSubManagement {
		comp, err := m.pubSubRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create pubsub component %s failed", name)
			return err
		}
		consumerID := strings.TrimSpace(config.Metadata["consumerID"])
		if consumerID == "" {
			config.Metadata["consumerID"] = m.runtimeConfig.AppManagement.AppId
		}

		if err := comp.Init(pubsub.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init pubsub component %s failed", name)
			return err
		}
		m.pubSubs[name] = comp
	}
	// 2. start subscribing
	return m.startSubscribing()
}

func (m *MosnRuntime) initStates(factorys ...*runtime_state.Factory) error {
	log.DefaultLogger.Infof("[runtime] start initializing state components")
	// 1. register all the implementation
	m.stateRegistry.Register(factorys...)
	// 2. loop initializing
	for name, config := range m.runtimeConfig.StateManagement {
		// 2.1. create and store the component
		comp, err := m.stateRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create state component %s failed", name)
			return err
		}
		if err := comp.Init(state.Metadata{Properties: config.Metadata}); err != nil {
			m.errInt(err, "init state component %s failed", name)
			return err
		}
		m.states[name] = comp
		// 2.2. save prefix strategy
		err = runtime_state.SaveStateConfiguration(name, config.Metadata)
		if err != nil {
			log.DefaultLogger.Errorf("error save state keyprefix: %s", err.Error())
			return err
		}
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
		comp, err := m.lockRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create lock component %s failed", name)
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
		comp, err := m.sequencerRegistry.Create(name)
		if err != nil {
			m.errInt(err, "create sequencer component %s failed", name)
			return err
		}
		// 2.2. init
		//// parse common config `biggerThan`
		//biggerThan := sequencer.DefaultBiggerThan
		//if config.Metadata[sequencer.BiggerThanKey] != "" {
		//	parseInt, err := strconv.ParseInt(config.Metadata[sequencer.BiggerThanKey], 10, 64)
		//	if err != nil {
		//		return err
		//	}
		//	biggerThan = parseInt
		//}
		// do init
		if err := comp.Init(sequencer.Configuration{
			Properties: config.Metadata,
			BiggerThan: config.BiggerThan,
		}); err != nil {
			m.errInt(err, "init sequencer component %s failed", name)
			return err
		}
		m.sequencers[name] = comp
	}
	return nil
}

func (m *MosnRuntime) startSubscribing() error {
	// 1. check if there is no need to do it
	if len(m.pubSubs) == 0 {
		return nil
	}
	topicRoutes, err := m.getInterestedTopics()
	if err != nil {
		return err
	}
	if len(topicRoutes) == 0 {
		//	no need
		return nil
	}
	// 2. loop subscribe
	for name, pubsub := range m.pubSubs {
		if err := m.beginPubSub(name, pubsub, topicRoutes); err != nil {
			return err
		}
	}
	return nil
}

func (m *MosnRuntime) beginPubSub(pubsubName string, ps pubsub.PubSub, topicRoutes map[string]TopicSubscriptions) error {
	// 1. call app to find topic topic2Details.
	v, ok := topicRoutes[pubsubName]
	if !ok {
		return nil
	}
	// 2. loop subscribing every <topic, route>
	for topic, route := range v.topic2Details {
		// TODO limit topic scope
		log.DefaultLogger.Debugf("[runtime][beginPubSub]subscribing to topic=%s on pubsub=%s", topic, pubsubName)
		// ask component to subscribe
		if err := ps.Subscribe(pubsub.SubscribeRequest{
			Topic:    topic,
			Metadata: route.metadata,
		}, func(ctx context.Context, msg *pubsub.NewMessage) error {
			if msg.Metadata == nil {
				msg.Metadata = make(map[string]string, 1)
			}
			msg.Metadata[Metadata_key_pubsubName] = pubsubName
			return m.publishMessageGRPC(ctx, msg)
		}); err != nil {
			log.DefaultLogger.Warnf("[runtime][beginPubSub]failed to subscribe to topic %s: %s", topic, err)
			return err
		}
	}

	return nil
}

func (m *MosnRuntime) getInterestedTopics() (map[string]TopicSubscriptions, error) {
	// 1. check
	if m.topicPerComponent != nil {
		return m.topicPerComponent, nil
	}
	if m.AppCallbackConn == nil {
		return make(map[string]TopicSubscriptions), nil
	}
	comp2Topic := make(map[string]TopicSubscriptions)
	var subscriptions []*runtimev1pb.TopicSubscription

	// 2. handle app subscriptions
	client := runtimev1pb.NewAppCallbackClient(m.AppCallbackConn)
	subscriptions = runtime_pubsub.ListTopicSubscriptions(client, log.DefaultLogger)
	// TODO handle declarative subscriptions

	// 3. prepare result
	for _, s := range subscriptions {
		if s == nil {
			continue
		}
		if _, ok := comp2Topic[s.PubsubName]; !ok {
			comp2Topic[s.PubsubName] = TopicSubscriptions{topic2Details: make(map[string]Details)}
		}
		comp2Topic[s.PubsubName].topic2Details[s.Topic] = Details{metadata: s.Metadata}
	}

	// 4. log
	if len(comp2Topic) > 0 {
		for pubsubName, v := range comp2Topic {
			topics := []string{}
			for topic := range v.topic2Details {
				topics = append(topics, topic)
			}
			log.DefaultLogger.Infof("[runtime][getInterestedTopics]app is subscribed to the following topics: %v through pubsub=%s", topics, pubsubName)
		}
	}
	// 5. cache the result
	m.topicPerComponent = comp2Topic
	return comp2Topic, nil
}

func (m *MosnRuntime) publishMessageGRPC(ctx context.Context, msg *pubsub.NewMessage) error {
	// 1. Unmarshal to cloudEvent model
	var cloudEvent map[string]interface{}
	err := m.json.Unmarshal(msg.Data, &cloudEvent)
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
			envelope.Data, _ = m.json.Marshal(data)
		}
	}
	// TODO tracing

	// 4. Call appcallback
	clientV1 := runtimev1pb.NewAppCallbackClient(m.AppCallbackConn)
	res, err := clientV1.OnTopicEvent(ctx, envelope)

	// 5. Check result
	return retryStrategy(err, res, cloudEvent)
}

// retryStrategy returns error when the message should be redelivered
func retryStrategy(err error, res *runtimev1pb.TopicEventResponse, cloudEvent map[string]interface{}) error {
	if err != nil {
		errStatus, hasErrStatus := status.FromError(err)
		if hasErrStatus && (errStatus.Code() == codes.Unimplemented) {
			// DROP
			log.DefaultLogger.Warnf("[runtime]non-retriable error returned from app while processing pub/sub event %v: %s", cloudEvent[pubsub.IDField].(string), err)
			return nil
		}

		err = errors.New(fmt.Sprintf("error returned from app while processing pub/sub event %v: %s", cloudEvent[pubsub.IDField].(string), err))
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
		return errors.New(fmt.Sprintf("RETRY status returned from app while processing pub/sub event %v", cloudEvent[pubsub.IDField].(string)))
	case runtimev1pb.TopicEventResponse_DROP:
		log.DefaultLogger.Warnf("[runtime]DROP status returned from app while processing pub/sub event %v", cloudEvent[pubsub.IDField].(string))
		return nil
	}
	// Consider unknown status field as error and retry
	return errors.New(fmt.Sprintf("unknown status returned from app while processing pub/sub event %v: %v", cloudEvent[pubsub.IDField].(string), res.GetStatus()))
}

func (m *MosnRuntime) initAppCallbackConnection() error {
	// init the client connection for calling app
	if m.runtimeConfig == nil || m.runtimeConfig.AppManagement.GrpcCallbackPort == 0 {
		return nil
	}
	port := m.runtimeConfig.AppManagement.GrpcCallbackPort
	opts := []rawGRPC.DialOption{
		rawGRPC.WithInsecure(),
	}
	// dial
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := rawGRPC.DialContext(ctx, fmt.Sprintf("127.0.0.1:%v", port), opts...)
	if err != nil {
		log.DefaultLogger.Warnf("[runtime]failed to init callback client at port %v : %s", port, err)
		return err
	}
	m.AppCallbackConn = conn
	return nil
}
