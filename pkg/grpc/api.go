package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/layotto/L8-components/configstores"
	"github.com/layotto/L8-components/hello"
	"strings"
	"sync"

	empty "google.golang.org/protobuf/types/known/emptypb"

	runtimev1pb "github.com/layotto/layotto/proto/runtime/v1"
	"mosn.io/pkg/log"
)

var (
	ErrNoInstance = errors.New("no instance found")
)

type API interface {
	SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error)
	// GetConfiguration gets configuration from configuration store.
	GetConfiguration(context.Context, *runtimev1pb.GetConfigurationRequest) (*runtimev1pb.GetConfigurationResponse, error)
	// SaveConfiguration saves configuration into configuration store.
	SaveConfiguration(context.Context, *runtimev1pb.SaveConfigurationRequest) (*empty.Empty, error)
	// DeleteConfiguration deletes configuration from configuration store.
	DeleteConfiguration(context.Context, *runtimev1pb.DeleteConfigurationRequest) (*empty.Empty, error)
	// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
	SubscribeConfiguration(runtimev1pb.MosnRuntime_SubscribeConfigurationServer) error
}

// api is a default implementation for MosnRuntimeServer.
type api struct {
	hellos       map[string]hello.HelloService
	configStores map[string]configstores.Store
}

func NewAPI(
	hellos map[string]hello.HelloService,
	configStores map[string]configstores.Store,
) API {
	return &api{
		hellos:       hellos,
		configStores: configStores,
	}
}

func (a *api) SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error) {
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
func (a *api) SaveConfiguration(ctx context.Context, req *runtimev1pb.SaveConfigurationRequest) (*empty.Empty, error) {
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
func (a *api) DeleteConfiguration(ctx context.Context, req *runtimev1pb.DeleteConfigurationRequest) (*empty.Empty, error) {
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
func (a *api) SubscribeConfiguration(sub runtimev1pb.MosnRuntime_SubscribeConfigurationServer) error {
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
