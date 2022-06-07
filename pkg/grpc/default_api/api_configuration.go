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

package default_api

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/pkg/log"
	"mosn.io/pkg/utils"

	"mosn.io/layotto/components/configstores"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

// GetConfiguration gets configuration from configuration store.
func (a *api) GetConfiguration(ctx context.Context, req *runtimev1pb.GetConfigurationRequest) (*runtimev1pb.GetConfigurationResponse, error) {
	resp := &runtimev1pb.GetConfigurationResponse{}
	// check store type supported or not
	store, ok := a.configStores[req.StoreName]
	if !ok {
		return nil, fmt.Errorf("configure store [%+v] don't support now", req.StoreName)
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
		return nil, fmt.Errorf("get configuration failed with error: %+v", err)
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
		return nil, fmt.Errorf("configure store [%+v] don't support now", req.StoreName)
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
		return nil, fmt.Errorf("configure store [%+v] don't support now", req.StoreName)
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
	// TODO currently this goroutine model is error-prone,and it should be refactored after new version of configuration API being accepted
	// 1. start a reader goroutine
	utils.GoWithRecover(func() {
		defer wg.Done()
		for {
			// 1.1. read stream
			req, err := sub.Recv()
			// 1.2. if an error happens,stop all the subscribers
			if err != nil {
				log.DefaultLogger.Errorf("occur error in subscribe, err: %+v", err)
				// stop all the subscribers
				for _, store := range subscribedStore {
					// TODO this method will stop subscribers created by other connections.Should be refactored
					store.StopSubscribe()
				}
				subErr = err
				// stop writer goroutine
				close(recvExitCh)
				return
			}
			// 1.3. else find the component and delegate to it
			store, ok := a.configStores[req.StoreName]
			// 1.3.1. stop if StoreName is not supported
			if !ok {
				log.DefaultLogger.Errorf("configure store [%+v] don't support now", req.StoreName)
				// stop all the subscribers
				for _, store := range subscribedStore {
					store.StopSubscribe()
				}
				subErr = fmt.Errorf("configure store [%+v] don't support now", req.StoreName)
				// stop writer goroutine
				close(recvExitCh)
				return
			}
			// 1.3.2. use default settings if blank
			if strings.ReplaceAll(req.Group, " ", "") == "" {
				req.Group = store.GetDefaultGroup()
			}
			if strings.ReplaceAll(req.Label, " ", "") == "" {
				req.Label = store.GetDefaultLabel()
			}
			// 1.3.3. delegate to the component
			store.Subscribe(&configstores.SubscribeReq{AppId: req.AppId, Group: req.Group, Label: req.Label, Keys: req.Keys, Metadata: req.Metadata}, respCh)
			subscribedStore = append(subscribedStore, store)
		}
	}, nil)
	// 2. start a writer goroutine
	utils.GoWithRecover(func() {
		defer wg.Done()
		for {
			select {
			// read response from components
			case resp, ok := <-respCh:
				if !ok {
					return
				}
				items := make([]*runtimev1pb.ConfigurationItem, 0, 10)
				for _, item := range resp.Items {
					items = append(items, &runtimev1pb.ConfigurationItem{Group: item.Group, Label: item.Label, Key: item.Key, Content: item.Content, Tags: item.Tags, Metadata: item.Metadata})
				}
				// write to response stream
				sub.Send(&runtimev1pb.SubscribeConfigurationResponse{StoreName: resp.StoreName, AppId: resp.AppId, Items: items})
			//	read exit signal
			case <-recvExitCh:
				return
			}
		}
	}, nil)
	wg.Wait()
	log.DefaultLogger.Warnf("subscribe gorountine exit")
	return subErr
}
