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
	"time"

	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/rpc/callback"
	"mosn.io/layotto/components/rpc/invoker/mosn/channel"
	_ "mosn.io/mosn/pkg/filter/network/proxy"
	"mosn.io/pkg/log"
)

const (
	Name = "mosn"
)

type mosnInvoker struct {
	channel rpc.Channel
	cb      rpc.Callback
}

type mosnConfig struct {
	Before  []rpc.CallbackFunc      `json:"before_invoke"`
	After   []rpc.CallbackFunc      `json:"after_invoke"`
	Channel []channel.ChannelConfig `json:"channel"`
}

func NewMosnInvoker() rpc.Invoker {
	invoker := &mosnInvoker{cb: callback.NewCallback()}
	return invoker
}

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

func (m *mosnInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (resp *rpc.RPCResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[runtime][rpc]mosn invoker panic: %v", r)
			log.DefaultLogger.Errorf("%v", err)
		}
	}()

	if req.Timeout == 0 {
		req.Timeout = 3000
	}
	req.Ctx = ctx
	startTime := time.Now()
	log.DefaultLogger.Debugf("[runtime][rpc]request %+v", req)
	req, err = m.cb.BeforeInvoke(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]before filter error %s", err.Error())
		return nil, err
	}
	beforeInvokeTime := time.Now()
	resp, err = m.channel.Do(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]error %s", err.Error())
		return nil, err
	}

	resp.Ctx = req.Ctx
	requestTime := time.Now()
	resp, err = m.cb.AfterInvoke(resp)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]after filter error %s", err.Error())
	}
	afterInvokeTime := time.Now()
	if afterInvokeTime.Sub(startTime) > 100 {
		log.DefaultLogger.Infof("[Layotto] rpc request spend time is: beforeInvokeTime: %+v, requestTime: %+v,afterInvoke: %+v, tootle:%+v",
			strconv.FormatInt(beforeInvokeTime.Sub(startTime).Nanoseconds()/1000000, 10),
			strconv.FormatInt(requestTime.Sub(beforeInvokeTime).Nanoseconds()/1000000, 10),
			strconv.FormatInt(afterInvokeTime.Sub(requestTime).Nanoseconds()/1000000, 10),
			strconv.FormatInt(afterInvokeTime.Sub(startTime).Nanoseconds()/1000000, 10),
		)
	}
	return resp, err
}
