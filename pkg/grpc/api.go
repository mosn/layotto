package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/state"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/layotto/layotto/pkg/common"
	"strings"
	"sync"

	contrib_contenttype "github.com/dapr/components-contrib/contenttype"
	"github.com/dapr/components-contrib/pubsub"
	contrib_pubsub "github.com/dapr/components-contrib/pubsub"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/layotto/layotto/components/configstores"
	"github.com/layotto/layotto/components/hello"
	"github.com/layotto/layotto/components/rpc"
	mosninvoker "github.com/layotto/layotto/components/rpc/invoker/mosn"
	"github.com/layotto/layotto/pkg/messages"
	runtime_state "github.com/layotto/layotto/pkg/runtime/state"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/rpc"
	mosninvoker "mosn.io/layotto/components/rpc/invoker/mosn"
	"mosn.io/layotto/pkg/messages"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/pkg/log"
)

var (
	ErrNoInstance = errors.New("no instance found")
)

type API interface {
	SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error)
	// GetConfiguration gets configuration from configuration store.
	GetConfiguration(context.Context, *runtimev1pb.GetConfigurationRequest) (*runtimev1pb.GetConfigurationResponse, error)
	// InvokeService do rpc calls.
	InvokeService(ctx context.Context, in *runtimev1pb.InvokeServiceRequest) (*runtimev1pb.InvokeResponse, error)
	// SaveConfiguration saves configuration into configuration store.
	SaveConfiguration(context.Context, *runtimev1pb.SaveConfigurationRequest) (*emptypb.Empty, error)
	// DeleteConfiguration deletes configuration from configuration store.
	DeleteConfiguration(context.Context, *runtimev1pb.DeleteConfigurationRequest) (*emptypb.Empty, error)
	// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
	SubscribeConfiguration(runtimev1pb.Runtime_SubscribeConfigurationServer) error
	// Publishes events to the specific topic.
	PublishEvent(context.Context, *runtimev1pb.PublishEventRequest) (*emptypb.Empty, error)
	// State
	GetState(ctx context.Context, in *runtimev1pb.GetStateRequest) (*runtimev1pb.GetStateResponse, error)
	GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error)
	SaveState(ctx context.Context, in *runtimev1pb.SaveStateRequest) (*emptypb.Empty, error)
	DeleteState(ctx context.Context, in *runtimev1pb.DeleteStateRequest) (*emptypb.Empty, error)
	DeleteBulkState(ctx context.Context, in *runtimev1pb.DeleteBulkStateRequest) (*emptypb.Empty, error)
	ExecuteStateTransaction(ctx context.Context, in *runtimev1pb.ExecuteStateTransactionRequest) (*emptypb.Empty, error)
}

// api is a default implementation for MosnRuntimeServer.
type api struct {
	appId                    string
	hellos                   map[string]hello.HelloService
	configStores             map[string]configstores.Store
	rpcs                     map[string]rpc.Invoker
	pubSubs                  map[string]pubsub.PubSub
	stateStores              map[string]state.Store
	transactionalStateStores map[string]state.TransactionalStore
}

func NewAPI(
	appId string,
	hellos map[string]hello.HelloService,
	configStores map[string]configstores.Store,
	rpcs map[string]rpc.Invoker,
	pubSubs map[string]pubsub.PubSub,
	stateStores map[string]state.Store,
) API {
	transactionalStateStores := map[string]state.TransactionalStore{}
	for key, store := range stateStores {
		if state.FeatureTransactional.IsPresent(store.Features()) {
			transactionalStateStores[key] = store.(state.TransactionalStore)
		}
	}
	return &api{
		appId:                    appId,
		hellos:                   hellos,
		configStores:             configStores,
		rpcs:                     rpcs,
		pubSubs:                  pubSubs,
		stateStores:              stateStores,
		transactionalStateStores: transactionalStateStores,
	}
}

func (a *api) SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error) {
	h, err := a.getHello(in.ServiceName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.say_hello] get hello error: %v", err)
		return nil, err
	}
	// create hello request based on pb.go struct
	req := &hello.HelloRequest{
		Name: in.Name,
	}
	resp, err := h.Hello(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.say_hello] request hello error: %v", err)
		return nil, err
	}
	// create response base on hello.Response
	return &runtimev1pb.SayHelloResponse{
		Hello: resp.HelloString,
	}, nil

}

func (a *api) getHello(name string) (hello.HelloService, error) {
	if len(a.hellos) == 0 {
		return nil, ErrNoInstance
	}
	h, ok := a.hellos[name]
	if !ok {
		return nil, ErrNoInstance
	}
	return h, nil
}

func (a *api) InvokeService(ctx context.Context, in *runtimev1pb.InvokeServiceRequest) (*runtimev1pb.InvokeResponse, error) {
	msg := in.GetMessage()
	req := &rpc.RPCRequest{
		Ctx:         ctx,
		Id:          in.GetId(),
		Method:      msg.GetMethod(),
		ContentType: msg.GetContentType(),
		Data:        msg.GetData().GetValue(),
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		req.Header = rpc.RPCHeader(md)
	} else {
		req.Header = rpc.RPCHeader(map[string][]string{})
	}
	if ext := msg.GetHttpExtension(); ext != nil {
		req.Header["verb"] = []string{ext.Verb.String()}
		req.Header["query_string"] = []string{ext.GetQuerystring()}
	}

	invoker, ok := a.rpcs[mosninvoker.Name]
	if !ok {
		return nil, errors.New("invoker not init")
	}

	resp, err := invoker.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Header != nil {
		header := metadata.Pairs()
		for k, values := range resp.Header {
			for _, v := range values {
				header.Append(k, v)
			}
		}
		grpc.SetHeader(ctx, header)
	}
	return &runtimev1pb.InvokeResponse{
		ContentType: resp.ContentType,
		Data:        &anypb.Any{Value: resp.Data},
	}, nil
}

// GetConfiguration gets configuration from configuration store.
func (a *api) GetConfiguration(ctx context.Context, req *runtimev1pb.GetConfigurationRequest) (*runtimev1pb.GetConfigurationResponse, error) {
	resp := &runtimev1pb.GetConfigurationResponse{}
	// check store type supported or not
	store, ok := a.configStores[req.StoreName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("configure store [%+v] don't support now", req.StoreName))
	}
	//here protect user use space for sting, eg: " ", "de fault"
	if strings.ReplaceAll(req.Group, " ", "") == "" {
		req.Group = store.GetDefaultGroup()
	}
	if strings.ReplaceAll(req.Label, " ", "") == "" {
		req.Label = store.GetDefaultLabel()
	}
	items, err := store.Get(ctx, &configstores.GetRequest{AppId: req.AppId, Group: req.Group, Label: req.Label, Keys: req.Keys, Metadata: req.Metadata})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get configuration failed with error: %+v", err))
	}
	for _, item := range items {
		resp.Items = append(resp.Items, &runtimev1pb.ConfigurationItem{Group: item.Group, Label: item.Label, Key: item.Key, Content: item.Content, Tags: item.Tags, Metadata: item.Metadata})
	}
	return resp, err
}

// SaveConfiguration saves configuration into configuration store.
func (a *api) SaveConfiguration(ctx context.Context, req *runtimev1pb.SaveConfigurationRequest) (*emptypb.Empty, error) {
	store, ok := a.configStores[req.StoreName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("configure store [%+v] don't support now", req.StoreName))
	}
	setReq := &configstores.SetRequest{}
	setReq.AppId = req.AppId
	setReq.StoreName = req.StoreName

	for index, item := range req.Items {
		if strings.ReplaceAll(item.Group, " ", "") == "" {
			req.Items[index].Group = store.GetDefaultGroup()
		}
		if strings.ReplaceAll(item.Label, " ", "") == "" {
			req.Items[index].Label = store.GetDefaultLabel()
		}
		setReq.Items = append(setReq.Items, &configstores.ConfigurationItem{Group: item.Group, Label: item.Label, Key: item.Key, Content: item.Content, Tags: item.Tags, Metadata: item.Metadata})
	}
	err := store.Set(ctx, setReq)
	return &emptypb.Empty{}, err
}

// DeleteConfiguration deletes configuration from configuration store.
func (a *api) DeleteConfiguration(ctx context.Context, req *runtimev1pb.DeleteConfigurationRequest) (*emptypb.Empty, error) {
	store, ok := a.configStores[req.StoreName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("configure store [%+v] don't support now", req.StoreName))
	}
	if strings.ReplaceAll(req.Group, " ", "") == "" {
		req.Group = store.GetDefaultGroup()
	}
	if strings.ReplaceAll(req.Label, " ", "") == "" {
		req.Label = store.GetDefaultLabel()
	}
	err := store.Delete(ctx, &configstores.DeleteRequest{AppId: req.AppId, Group: req.Group, Label: req.Label, Keys: req.Keys, Metadata: req.Metadata})
	return &emptypb.Empty{}, err
}

// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
func (a *api) SubscribeConfiguration(sub runtimev1pb.Runtime_SubscribeConfigurationServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var subErr error
	respCh := make(chan *configstores.SubscribeResp)
	recvExitCh := make(chan struct{})
	subscribedStore := make([]configstores.Store, 0, 1)
	go func() {
		defer wg.Done()
		for {
			req, err := sub.Recv()
			if err != nil {
				log.DefaultLogger.Errorf("occur error in subscribe, err: %+v", err)
				for _, store := range subscribedStore {
					store.StopSubscribe()
				}
				subErr = err
				if len(subscribedStore) == 0 {
					close(recvExitCh)
				}
				return
			}
			store, ok := a.configStores[req.StoreName]
			if !ok {
				log.DefaultLogger.Errorf("configure store [%+v] don't support now", req.StoreName)
				subErr = errors.New(fmt.Sprintf("configure store [%+v] don't support now", req.StoreName))
				close(recvExitCh)
				return
			}
			if strings.ReplaceAll(req.Group, " ", "") == "" {
				req.Group = store.GetDefaultGroup()
			}
			if strings.ReplaceAll(req.Label, " ", "") == "" {
				req.Label = store.GetDefaultLabel()
			}
			store.Subscribe(&configstores.SubscribeReq{AppId: req.AppId, Group: req.Group, Label: req.Label, Keys: req.Keys, Metadata: req.Metadata}, respCh)
			subscribedStore = append(subscribedStore, store)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case resp, ok := <-respCh:
				if !ok {
					return
				}
				items := make([]*runtimev1pb.ConfigurationItem, 0, 10)
				for _, item := range resp.Items {
					items = append(items, &runtimev1pb.ConfigurationItem{Group: item.Group, Label: item.Label, Key: item.Key, Content: item.Content, Tags: item.Tags, Metadata: item.Metadata})
				}
				sub.Send(&runtimev1pb.SubscribeConfigurationResponse{StoreName: resp.StoreName, AppId: resp.StoreName, Items: items})
			case <-recvExitCh:
				return
			}
		}
	}()
	wg.Wait()
	log.DefaultLogger.Warnf("subscribe gorountine exit")
	return subErr
}

func (a *api) PublishEvent(ctx context.Context, in *runtimev1pb.PublishEventRequest) (*emptypb.Empty, error) {
	result, err := a.doPublishEvent(ctx, in.PubsubName, in.Topic, in.Data, in.DataContentType, in.Metadata)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.PublishEvent] %v", err)
	}
	return result, err
}

// doPublishEvent is a protocal irrelevant function to do event publishing.
// It's easy to add APIs for other protocals.Just move this func to a separate layer if you need.
func (a *api) doPublishEvent(ctx context.Context, pubsubName string, topic string, data []byte, contentType string, metadata map[string]string) (*emptypb.Empty, error) {
	// 1. validate
	if pubsubName == "" {
		err := status.Error(codes.InvalidArgument, messages.ErrPubsubEmpty)
		return &emptypb.Empty{}, err
	}
	if topic == "" {
		err := status.Errorf(codes.InvalidArgument, messages.ErrTopicEmpty, pubsubName)
		return &emptypb.Empty{}, err
	}
	// 2. get component
	component, ok := a.pubSubs[pubsubName]
	if !ok {
		err := status.Errorf(codes.InvalidArgument, messages.ErrPubsubNotFound, pubsubName)
		return &emptypb.Empty{}, err
	}

	// 3. new cloudevent request
	if data == nil {
		data = []byte{}
	}
	var envelope map[string]interface{}
	var err error = nil
	if contrib_contenttype.IsCloudEventContentType(contentType) {
		envelope, err = contrib_pubsub.FromCloudEvent(data, topic, pubsubName, "")
		if err != nil {
			err = status.Errorf(codes.InvalidArgument, messages.ErrPubsubCloudEventCreation, err.Error())
			return &emptypb.Empty{}, err
		}
	} else {
		envelope = contrib_pubsub.NewCloudEventsEnvelope(uuid.New().String(), "", contrib_pubsub.DefaultCloudEventType, "", topic, pubsubName,
			contentType, data, "")
	}
	features := component.Features()
	pubsub.ApplyMetadata(envelope, features, metadata)

	b, err := jsoniter.ConfigFastest.Marshal(envelope)
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, messages.ErrPubsubCloudEventsSer, topic, pubsubName, err.Error())
		return &emptypb.Empty{}, err
	}
	// 4. publish
	req := pubsub.PublishRequest{
		PubsubName: pubsubName,
		Topic:      topic,
		Data:       b,
		Metadata:   metadata,
	}

	// TODO limit topic scope
	err = component.Publish(&req)
	if err != nil {
		nerr := status.Errorf(codes.Internal, messages.ErrPubsubPublishMessage, topic, pubsubName, err.Error())
		return &emptypb.Empty{}, nerr
	}
	return &emptypb.Empty{}, nil
}

// GetState obtains the state for a specific key.
func (a *api) GetState(ctx context.Context, in *runtimev1pb.GetStateRequest) (*runtimev1pb.GetStateResponse, error) {
	// 1. get store
	store, err := a.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] error: %v", err)
		return nil, err
	}
	// 2. generate the actual key
	key, err := runtime_state.GetModifiedStateKey(in.Key, in.StoreName, a.appId)
	if err != nil {
		return &runtimev1pb.GetStateResponse{}, err
	}
	req := state.GetRequest{
		Key:      key,
		Metadata: in.Metadata,
		Options: state.GetStateOption{
			Consistency: stateConsistencyToString(in.Consistency),
		},
	}
	// 3. query
	getResponse, err := store.Get(&req)
	// 4. check result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateGet, in.Key, in.StoreName, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] %v", err)
		return &runtimev1pb.GetStateResponse{}, err
	}

	response := &runtimev1pb.GetStateResponse{}
	if getResponse != nil {
		response.Etag = stringValueOrEmpty(getResponse.ETag)
		response.Data = getResponse.Data
		response.Metadata = getResponse.Metadata
	}
	return response, nil
}

func (a *api) getStateStore(name string) (state.Store, error) {
	if a.stateStores == nil || len(a.stateStores) == 0 {
		return nil, status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
	}

	if a.stateStores[name] == nil {
		return nil, status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, name)
	}
	return a.stateStores[name], nil
}

func stringValueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func (a *api) GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error) {
	// 1. get store
	store, err := a.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetBulkState] error: %v", err)
		return &runtimev1pb.GetBulkStateResponse{}, err
	}

	bulkResp := &runtimev1pb.GetBulkStateResponse{}
	if len(in.Keys) == 0 {
		return bulkResp, nil
	}

	// 2. try bulk get first
	// 2.1. convert reqs
	reqs := make([]state.GetRequest, len(in.Keys))
	for i, k := range in.Keys {
		key, err1 := runtime_state.GetModifiedStateKey(k, in.StoreName, a.appId)
		if err1 != nil {
			return &runtimev1pb.GetBulkStateResponse{}, err1
		}
		r := state.GetRequest{
			Key:      key,
			Metadata: in.Metadata,
		}
		reqs[i] = r
	}
	// 2.2. query
	bulkGet, responses, err := store.BulkGet(reqs)
	// 2.3. parse and return result if store supports bulk get
	if bulkGet {
		if err != nil {
			return bulkResp, err
		}
		for i := 0; i < len(responses); i++ {
			item := &runtimev1pb.BulkStateItem{
				Key:      runtime_state.GetOriginalStateKey(responses[i].Key),
				Data:     responses[i].Data,
				Etag:     stringValueOrEmpty(responses[i].ETag),
				Metadata: responses[i].Metadata,
				Error:    responses[i].Error,
			}
			bulkResp.Items = append(bulkResp.Items, item)
		}
		return bulkResp, nil
	}

	// 3. if store doesn't support bulk get, fallback to call get() method one by one
	limiter := common.NewLimiter(int(in.Parallelism))
	for i := 0; i < len(reqs); i++ {
		fn := func(param interface{}) {
			req := param.(*state.GetRequest)
			r, err := store.Get(req)
			item := &runtimev1pb.BulkStateItem{
				Key: runtime_state.GetOriginalStateKey(req.Key),
			}
			if err != nil {
				item.Error = err.Error()
			} else if r != nil {
				item.Data = r.Data
				item.Etag = stringValueOrEmpty(r.ETag)
				item.Metadata = r.Metadata
			}
			bulkResp.Items = append(bulkResp.Items, item)
		}
		limiter.Execute(fn, &reqs[i])
	}
	limiter.Wait()

	return bulkResp, nil
}

func (a *api) SaveState(ctx context.Context, in *runtimev1pb.SaveStateRequest) (*emptypb.Empty, error) {
	// 1. get store
	store, err := a.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}

	// 2. convert requests
	reqs := []state.SetRequest{}
	for _, s := range in.States {
		key, err1 := runtime_state.GetModifiedStateKey(s.Key, in.StoreName, a.appId)
		if err1 != nil {
			return &emptypb.Empty{}, err1
		}
		req := state.SetRequest{
			Key:      key,
			Metadata: s.Metadata,
			Value:    s.Value,
		}
		if s.Etag != nil {
			req.ETag = &s.Etag.Value
		}
		if s.Options != nil {
			req.Options = state.SetStateOption{
				Consistency: stateConsistencyToString(s.Options.Consistency),
				Concurrency: stateConcurrencyToString(s.Options.Concurrency),
			}
		}
		reqs = append(reqs, req)
	}

	// 3. query
	err = store.BulkSet(reqs)
	// 4. check result
	if err != nil {
		err = a.stateErrorResponse(err, messages.ErrStateSave, in.StoreName, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// stateErrorResponse takes a state store error, format and args and returns a status code encoded gRPC error
func (a *api) stateErrorResponse(err error, format string, args ...interface{}) error {
	e, ok := err.(*state.ETagError)
	if !ok {
		return status.Errorf(codes.Internal, format, args...)
	}
	switch e.Kind() {
	case state.ETagMismatch:
		return status.Errorf(codes.Aborted, format, args...)
	case state.ETagInvalid:
		return status.Errorf(codes.InvalidArgument, format, args...)
	}

	return status.Errorf(codes.Internal, format, args...)
}

func (a *api) DeleteState(ctx context.Context, in *runtimev1pb.DeleteStateRequest) (*emptypb.Empty, error) {
	// 1. get store
	store, err := a.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 2. generate the actual key
	key, err := runtime_state.GetModifiedStateKey(in.Key, in.StoreName, a.appId)
	if err != nil {
		return &empty.Empty{}, err
	}
	// 3. convert request
	req := state.DeleteRequest{
		Key:      key,
		Metadata: in.Metadata,
	}
	if in.Etag != nil {
		req.ETag = &in.Etag.Value
	}
	if in.Options != nil {
		req.Options = state.DeleteStateOption{
			Concurrency: stateConcurrencyToString(in.Options.Concurrency),
			Consistency: stateConsistencyToString(in.Options.Consistency),
		}
	}
	// 4. send request
	err = store.Delete(&req)
	// 5. check result
	if err != nil {
		err = a.stateErrorResponse(err, messages.ErrStateDelete, in.Key, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteState] error: %v", err)
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (a *api) DeleteBulkState(ctx context.Context, in *runtimev1pb.DeleteBulkStateRequest) (*empty.Empty, error) {
	// 1. get store
	store, err := a.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteBulkState] error: %v", err)
		return &empty.Empty{}, err
	}
	// 2. convert request
	reqs := make([]state.DeleteRequest, 0, len(in.States))
	for _, item := range in.States {
		key, err1 := runtime_state.GetModifiedStateKey(item.Key, in.StoreName, a.appId)
		if err1 != nil {
			return &empty.Empty{}, err1
		}
		req := state.DeleteRequest{
			Key:      key,
			Metadata: item.Metadata,
		}
		if item.Etag != nil {
			req.ETag = &item.Etag.Value
		}
		if item.Options != nil {
			req.Options = state.DeleteStateOption{
				Concurrency: stateConcurrencyToString(item.Options.Concurrency),
				Consistency: stateConsistencyToString(item.Options.Consistency),
			}
		}
		reqs = append(reqs, req)
	}
	// 3. send request
	err = store.BulkDelete(reqs)
	// 4. check result
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteBulkState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (a *api) ExecuteStateTransaction(ctx context.Context, in *runtimev1pb.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	// 1. check params
	if a.stateStores == nil || len(a.stateStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	storeName := in.StoreName
	if a.stateStores[storeName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, storeName)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 2. find store
	transactionalStore, ok := a.transactionalStateStores[storeName]
	if !ok {
		err := status.Errorf(codes.Unimplemented, messages.ErrStateStoreNotSupported, storeName)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 3. convert request
	operations := []state.TransactionalStateOperation{}
	for _, inputReq := range in.Operations {
		// 3.1. extract and validate fields
		var operation state.TransactionalStateOperation
		var req = inputReq.Request
		// tolerant npe
		if req == nil {
			log.DefaultLogger.Warnf("[runtime] [grpc.ExecuteStateTransaction] one of TransactionalStateOperation.Request is nil")
			continue
		}
		key, err := runtime_state.GetModifiedStateKey(req.Key, in.StoreName, a.appId)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		// 3.2. prepare TransactionalStateOperation struct according to the operation type
		switch state.OperationType(inputReq.OperationType) {
		case state.Upsert:
			operation = state.TransactionalStateOperation{
				Operation: state.Upsert,
				Request:   prepareSetRequest(req, key),
			}
		case state.Delete:
			operation = state.TransactionalStateOperation{
				Operation: state.Delete,
				Request:   prepareDelRequest(req, key),
			}
		default:
			err := status.Errorf(codes.Unimplemented, messages.ErrNotSupportedStateOperation, inputReq.OperationType)
			log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
			return &emptypb.Empty{}, err
		}
		operations = append(operations, operation)
	}
	// 4. submit transactional request
	err := transactionalStore.Multi(&state.TransactionalStateRequest{
		Operations: operations,
		Metadata:   in.Metadata,
	})
	// 5. check result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateTransaction, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func prepareDelRequest(req *runtimev1pb.StateItem, key string) state.DeleteRequest {
	hasEtag, etag := extractEtag(req)
	delReq := state.DeleteRequest{
		Key:      key,
		Metadata: req.Metadata,
	}

	if hasEtag {
		delReq.ETag = &etag
	}
	if req.Options != nil {
		delReq.Options = state.DeleteStateOption{
			Concurrency: stateConcurrencyToString(req.Options.Concurrency),
			Consistency: stateConsistencyToString(req.Options.Consistency),
		}
	}
	return delReq
}

func prepareSetRequest(req *runtimev1pb.StateItem, key string) state.SetRequest {
	hasEtag, etag := extractEtag(req)
	setReq := state.SetRequest{
		Key: key,
		// Limitation:
		// components that cannot handle byte array need to deserialize/serialize in
		// component specific way in components-contrib repo.
		Value:    req.Value,
		Metadata: req.Metadata,
	}

	if hasEtag {
		setReq.ETag = &etag
	}
	if req.Options != nil {
		setReq.Options = state.SetStateOption{
			Concurrency: stateConcurrencyToString(req.Options.Concurrency),
			Consistency: stateConsistencyToString(req.Options.Consistency),
		}
	}
	return setReq
}

func extractEtag(req *runtimev1pb.StateItem) (bool, string) {
	if req.Etag != nil {
		return true, req.Etag.Value
	}
	return false, ""
}
