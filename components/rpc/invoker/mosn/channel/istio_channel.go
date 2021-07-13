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

package channel

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/golang/protobuf/ptypes"
	"github.com/valyala/fasthttp"
	"mosn.io/layotto/components/rpc"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/router"
	clusterAdapter "mosn.io/mosn/pkg/upstream/cluster"
	"mosn.io/mosn/pkg/xds/conv"
	v2 "mosn.io/mosn/pkg/xds/v2"
	"net"
	"net/http"
	"time"
)

func init() {
	RegistChannel("istio", newIstioChannel)
	v2.RegisterTypeURLHandleFunc(v2.EnvoyRouteConfiguration, HandleEnvoyRouteConfiguration)
}

func HandleEnvoyRouteConfiguration(client *v2.ADSClient, resp *envoy_api_v2.DiscoveryResponse) {
	log.DefaultLogger.Tracef("get rds resp,handle it")
	routes := handleRoutesResp(resp)
	log.DefaultLogger.Infof("get %d routes from RDS", len(routes))
	ConvertAddOrUpdateRouters(routes)
	v2.AckResponse(client.StreamClient, resp)
}

func handleRoutesResp(resp *envoy_api_v2.DiscoveryResponse) []*envoy_api_v2.RouteConfiguration {
	routes := make([]*envoy_api_v2.RouteConfiguration, 0, len(resp.Resources))
	for _, res := range resp.Resources {
		route := &envoy_api_v2.RouteConfiguration{}
		if err := ptypes.UnmarshalAny(res, route); err != nil {
			log.DefaultLogger.Errorf("ADSClient unmarshal route fail: %v", err)
		}
		routes = append(routes, route)
	}
	return routes
}

var domainToCluster = make(map[string]string)

func ConvertAddOrUpdateRouters(routers []*envoy_api_v2.RouteConfiguration) {
	if routersMngIns := router.GetRoutersMangerInstance(); routersMngIns == nil {
		log.DefaultLogger.Errorf("xds OnAddOrUpdateRouters error: router manager in nil")
	} else {
		for _, router := range routers {
			log.DefaultLogger.Debugf("xds convert router config: %+v", router)

			mosnRouter, _ := conv.ConvertRouterConf("", router)

			for _, h := range mosnRouter.VirtualHosts {
				for _, d := range h.Domains {
					if len(h.Routers) == 0 {
						break
					}
					domainToCluster[d] = h.Routers[0].Route.ClusterName
				}
			}

			if err := routersMngIns.AddOrUpdateRouters(mosnRouter); err != nil {
				log.DefaultLogger.Errorf("xds client  routersMngIns.AddOrUpdateRouters error: %v", err)
			}
		}
	}
}

type istioChannel struct {
	httpChannel
}

func newIstioChannel(config ChannelConfig) (rpc.Channel, error) {
	return &istioChannel{
		httpChannel{},
	}, nil
}

func (i *istioChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	timeout := time.Duration(req.Timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(req.Ctx, timeout)
	defer cancel()

	clusterName, ok := domainToCluster[req.Id]
	if !ok {
		return nil, errors.New("no available service found")
	}
	snapshot := clusterAdapter.GetClusterMngAdapterInstance().GetClusterSnapshot(ctx, clusterName)
	if snapshot == nil || len(snapshot.HostSet().Hosts()) == 0 {
		return nil, errors.New("no available server address")
	}

	conn, err := net.Dial("tcp", snapshot.HostSet().Hosts()[0].AddressString())
	if err != nil {
		return nil, err
	}

	deadline, _ := ctx.Deadline()
	conn.SetDeadline(deadline)

	httpReq := i.httpChannel.constructReq(req)
	defer fasthttp.ReleaseRequest(httpReq)

	if _, err = httpReq.WriteTo(conn); err != nil {
		return nil, err
	}

	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpResp)
	if err = httpResp.Read(bufio.NewReader(conn)); err != nil {
		return nil, err
	}
	body := httpResp.Body()
	if httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http response code %d, body: %s", httpResp.StatusCode(), string(body))
	}

	data := make([]byte, len(body))
	copy(data, body)
	rpcResp := &rpc.RPCResponse{
		ContentType: string(httpResp.Header.ContentType()),
		Data:        data,
		Header:      map[string][]string{},
	}
	httpResp.Header.VisitAll(func(key, value []byte) {
		rpcResp.Header[string(key)] = []string{string(value)}
	})
	return rpcResp, nil

}


