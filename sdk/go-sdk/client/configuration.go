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

package client

import (
	"context"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

type WatchChan <-chan WatchResponse

// ConfigurationRequestItem used for GET,DEL,SUB request
type ConfigurationRequestItem struct {
	// The name of configuration store.
	StoreName string
	// The application id which
	// Only used for admin, Ignored and reset for normal client
	AppId string
	// The group of keys.
	Group string
	// The label for keys.
	Label string
	// The keys to get.
	Keys []string
	// The metadata which will be sent to configuration store components.
	Metadata map[string]string
}

type ConfigurationItem struct {
	// Required. The key of configuration item
	Key string
	// The content of configuration item
	// Empty if the configuration is not set, including the case that the configuration is changed from value-set to value-not-set.
	Content string
	// The group of configuration item.
	Group string
	// The label of configuration item.
	Label string
	// The tag list of configuration item.
	Tags map[string]string
	// The metadata which will be passed to configuration store component.
	Metadata map[string]string
}

type SaveConfigurationRequest struct {
	// The name of configuration store.
	StoreName string
	// The application id which
	// Only used for admin, ignored and reset for normal client
	AppId string
	// The list of configuration items to save.
	// To delete a exist item, set the key (also label) and let content to be empty
	Items []*ConfigurationItem
	// The metadata which will be sent to configuration store components.
	Metadata map[string]string
}

type SubConfigurationResp struct {
	// The name of configuration store.
	StoreName string
	// The application id which
	// Only used for admin, ignored and reset for normal client
	AppId string
	// The list of configuration items to save.
	// To delete a exist item, set the key (also label) and let content to be empty
	Items []*ConfigurationItem
}

type WatchResponse struct {
	Item *SubConfigurationResp
	Err  error
}

func (c *GRPCClient) GetConfiguration(ctx context.Context, in *ConfigurationRequestItem) ([]*ConfigurationItem, error) {
	req := &runtimev1pb.GetConfigurationRequest{StoreName: in.StoreName, AppId: in.AppId, Group: in.Group, Label: in.Label, Keys: in.Keys, Metadata: in.Metadata}
	resp, err := c.protoClient.GetConfiguration(ctx, req)
	if err != nil {
		return nil, err
	}
	items := make([]*ConfigurationItem, 0, len(resp.Items))
	for _, v := range resp.Items {
		c := &ConfigurationItem{Group: v.Group, Label: v.Label, Key: v.Key, Content: v.Content, Tags: v.Tags, Metadata: v.Metadata}
		items = append(items, c)
	}
	return items, nil
}

// SaveConfiguration saves configuration into configuration store.
func (c *GRPCClient) SaveConfiguration(ctx context.Context, in *SaveConfigurationRequest) error {
	req := &runtimev1pb.SaveConfigurationRequest{StoreName: in.StoreName, AppId: in.AppId, Metadata: in.Metadata}
	for _, v := range in.Items {
		c := &runtimev1pb.ConfigurationItem{Group: v.Group, Label: v.Label, Key: v.Key, Content: v.Content, Tags: v.Tags, Metadata: v.Metadata}
		req.Items = append(req.Items, c)
	}
	_, err := c.protoClient.SaveConfiguration(ctx, req)
	return err
}

// DeleteConfiguration deletes configuration from configuration store.
func (c *GRPCClient) DeleteConfiguration(ctx context.Context, in *ConfigurationRequestItem) error {
	req := &runtimev1pb.DeleteConfigurationRequest{StoreName: in.StoreName, AppId: in.AppId, Group: in.Group, Label: in.Label, Keys: in.Keys, Metadata: in.Metadata}
	_, err := c.protoClient.DeleteConfiguration(ctx, req)
	return err
}

// SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
func (c *GRPCClient) SubscribeConfiguration(ctx context.Context, in *ConfigurationRequestItem) WatchChan {
	cli, err := c.protoClient.SubscribeConfiguration(ctx)
	res := WatchResponse{}
	resCh := make(chan WatchResponse, 1)
	if err != nil {
		res.Err = err
		resCh <- res
		close(resCh)
		return resCh
	}
	request := &runtimev1pb.SubscribeConfigurationRequest{StoreName: in.StoreName, AppId: in.AppId, Group: in.Group, Label: in.Label, Keys: in.Keys, Metadata: in.Metadata}
	err = cli.Send(request)
	if err != nil {
		res.Err = err
		resCh <- res
		close(resCh)
		return resCh
	}
	GoWithRecover(func() {
		for {
			resp, err := cli.Recv()
			if err != nil {
				res.Err = err
				resCh <- res
				close(resCh)
				return
			}
			item := &SubConfigurationResp{}
			item.StoreName = resp.StoreName
			item.AppId = resp.AppId
			for _, v := range resp.Items {
				c := &ConfigurationItem{}
				c.Metadata = v.Metadata
				c.Label = v.Label
				c.Group = v.Group
				c.Key = v.Key
				c.Tags = v.Tags
				c.Content = v.Content
				item.Items = append(item.Items, c)
			}
			res.Item = item
			res.Err = nil
			resCh <- res
		}
	}, nil)
	return resCh
}
