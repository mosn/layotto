package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/pubsub"
	jsoniter "github.com/json-iterator/go"
	"github.com/layotto/layotto/pkg/messages"
	runtime_pubsub "github.com/layotto/layotto/pkg/runtime/pubsub"
	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/pkg/services/hello"
	"github.com/layotto/layotto/spec/proto/runtime/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/pkg/log"
	"strings"
	"sync"
)

var (
	ErrNoInstance = errors.New("no instance found")
)

type API interface {
	SayHello(ctx context.Context, in *runtime.SayHelloRequest) (*runtime.SayHelloResponse, error)
	// GetConfiguration gets configuration from configuration store.
	GetConfiguration(context.Context, *runtime.GetConfigurationRequest) (*runtime.GetConfigurationResponse, error)
	// SaveConfiguration saves configuration into configuration store.
	SaveConfiguration(context.Context, *runtime.SaveConfigurationRequest) (*empty.Empty, error)
	// DeleteConfiguration deletes configuration from configuration store.
	DeleteConfiguration(context.Context, *runtime.DeleteConfigurationRequest) (*empty.Empty, error)
	// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
	SubscribeConfiguration(runtime.Runtime_SubscribeConfigurationServer) error
	// Publishes events to the specific topic.
	PublishEvent(context.Context, *runtime.PublishEventRequest) (*empty.Empty, error)
}

// api is a default implementation for MosnRuntimeServer.
type api struct {
	hellos       map[string]hello.HelloService
	configStores map[string]configstores.Store
	pubSubs      map[string]pubsub.PubSub
}

func NewAPI(
	hellos map[string]hello.HelloService,
	configStores map[string]configstores.Store,
	pubSubs map[string]pubsub.PubSub,
) API {
	return &api{
		hellos:       hellos,
		configStores: configStores,
		pubSubs:      pubSubs,
	}
}

func (a *api) SayHello(ctx context.Context, in *runtime.SayHelloRequest) (*runtime.SayHelloResponse, error) {
	h, err := a.getHello(in.ServiceName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.say_hello] get hello error: %v", err)
		return nil, err
	}
	// create hello request based on pb.go struct
	req := &hello.HelloRequest{}
	resp, err := h.Hello(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.say_hello] request hello error: %v", err)
		return nil, err
	}
	// create response base on hello.Response
	return &runtime.SayHelloResponse{
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

// GetConfiguration gets configuration from configuration store.
func (a *api) GetConfiguration(ctx context.Context, req *runtime.GetConfigurationRequest) (*runtime.GetConfigurationResponse, error) {
	resp := &runtime.GetConfigurationResponse{}
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
		resp.Items = append(resp.Items, &runtime.ConfigurationItem{Group: item.Group, Label: item.Label, Key: item.Key, Content: item.Content, Tags: item.Tags, Metadata: item.Metadata})
	}
	return resp, err
}

// SaveConfiguration saves configuration into configuration store.
func (a *api) SaveConfiguration(ctx context.Context, req *runtime.SaveConfigurationRequest) (*empty.Empty, error) {
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
	return &empty.Empty{}, err
}

// DeleteConfiguration deletes configuration from configuration store.
func (a *api) DeleteConfiguration(ctx context.Context, req *runtime.DeleteConfigurationRequest) (*empty.Empty, error) {
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
	return &empty.Empty{}, err
}

// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
func (a *api) SubscribeConfiguration(sub runtime.Runtime_SubscribeConfigurationServer) error {
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
				items := make([]*runtime.ConfigurationItem, 0, 10)
				for _, item := range resp.Items {
					items = append(items, &runtime.ConfigurationItem{Group: item.Group, Label: item.Label, Key: item.Key, Content: item.Content, Tags: item.Tags, Metadata: item.Metadata})
				}
				sub.Send(&runtime.SubscribeConfigurationResponse{StoreName: resp.StoreName, AppId: resp.StoreName, Items: items})
			case <-recvExitCh:
				return
			}
		}
	}()
	wg.Wait()
	log.DefaultLogger.Warnf("subscribe gorountine exit")
	return subErr
}

func (a *api) PublishEvent(ctx context.Context, in *runtime.PublishEventRequest) (*emptypb.Empty, error) {
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
	envelope, err := runtime_pubsub.NewCloudEvent(&runtime_pubsub.CloudEvent{
		Topic:           topic,
		DataContentType: contentType,
		Data:            data,
		Pubsub:          pubsubName,
	})
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, messages.ErrPubsubCloudEventCreation, err.Error())
		return &emptypb.Empty{}, err
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
		if errors.As(err, &runtime_pubsub.NotAllowedError{}) {
			nerr = status.Errorf(codes.PermissionDenied, err.Error())
		}

		if errors.As(err, &runtime_pubsub.NotFoundError{}) {
			nerr = status.Errorf(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, nerr
	}
	return &emptypb.Empty{}, nil
}
