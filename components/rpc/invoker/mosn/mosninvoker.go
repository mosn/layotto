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

package mosn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	// bridge to mosn
	_ "mosn.io/mosn/pkg/filter/network/proxy"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/rpc/callback"
	"mosn.io/layotto/components/rpc/invoker/mosn/channel"
)

const (
	Name = "mosn"
)

// mosnInvoker is Invoker implement
type mosnInvoker struct {
	channel rpc.Channel
	cb      rpc.Callback
}

// mosnConfig is mosn config
type mosnConfig struct {
	Before  []rpc.CallbackFunc      `json:"before_invoke"`
	After   []rpc.CallbackFunc      `json:"after_invoke"`
	Channel []channel.ChannelConfig `json:"channel"`
}

// NewMosnInvoker is init mosnInvoker
func NewMosnInvoker() rpc.Invoker {
	invoker := &mosnInvoker{cb: callback.NewCallback()}
	return invoker
}

// Init is init mosn RpcConfig
func (m *mosnInvoker) Init(conf rpc.RpcConfig) error {
	var config mosnConfig
	if err := json.Unmarshal(conf.Config, &config); err != nil {
		return err
	}

	for _, before := range config.Before {
		m.cb.AddBeforeInvoke(before)
	}

	for _, after := range config.After {
		m.cb.AddAfterInvoke(after)
	}

	if len(config.Channel) == 0 {
		return errors.New("missing channel config")
	}

	// todo support multiple channel
	channel, err := channel.GetChannel(config.Channel[0])
	if err != nil {
		return err
	}
	m.channel = channel
	return nil
}

// Invoke is invoke mosn RPCRequest and Context to RPCResponse
func (m *mosnInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (resp *rpc.RPCResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[runtime][rpc]mosn invoker panic: %v", r)
			log.DefaultLogger.Errorf("%v", err)
		}
	}()

	// 1. validate request
	if req.Timeout == 0 {
		req.Timeout = rpc.DefaultRequestTimeoutMs
		if ts, ok := req.Header[rpc.RequestTimeoutMs]; ok && len(ts) > 0 {
			t, err := strconv.ParseInt(ts[0], 10, 32)
			if err == nil && t != 0 {
				req.Timeout = int32(t)
			}
		}
	}
	req.Ctx = ctx
	log.DefaultLogger.Debugf("[runtime][rpc]request %+v", req)
	// 2. beforeInvoke callback
	req, err = m.cb.BeforeInvoke(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]before filter error %s", err.Error())
		return nil, err
	}
	// 3. do invocation
	resp, err = m.channel.Do(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]error %s", err.Error())
		return nil, err
	}
	resp.Ctx = req.Ctx
	// 4. afterInvoke callback
	resp, err = m.cb.AfterInvoke(resp)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]after filter error %s", err.Error())
		return nil, err
	}
	return resp, nil
}
