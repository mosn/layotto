package mosn

import (
	"context"
	"encoding/json"
	"sync/atomic"

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
	channels []rpc.Channel
	rrIdx    uint32
	cb       rpc.Callback
}

type mosnConfig struct {
	Before        []rpc.CallbackFunc    `json:"before_invoke"`
	After         []rpc.CallbackFunc    `json:"after_invoke"`
	TotalChannels int                   `json:"total_channels"`
	Channel       channel.ChannelConfig `json:"channel"`
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

	if config.TotalChannels <= 0 {
		config.TotalChannels = 1
	}
	for i := 0; i < config.TotalChannels; i++ {
		channel, err := channel.GetChannel(config.Channel)
		if err != nil {
			return err
		}
		m.channels = append(m.channels, channel)
	}
	return nil
}

func (m *mosnInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.DefaultLogger.Errorf("[runtime][rpc]mosn invoker panic: %v", r)
		}
	}()

	if req.Timeout == 0 {
		req.Timeout = 3000
	}
	log.DefaultLogger.Debugf("[runtime][rpc]request %+v", req)
	req, err := m.cb.BeforeInvoke(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]before filter error %s", err.Error())
		return nil, err
	}

	resp, err := m.getChannel().Do(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]error %s", err.Error())
		return nil, err
	}

	resp.Ctx = req.Ctx
	resp, err = m.cb.AfterInvoke(resp)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]after filter error %s", err.Error())
	}
	return resp, err
}

func (m *mosnInvoker) getChannel() rpc.Channel {
	idx := atomic.AddUint32(&m.rrIdx, 1) % uint32(len(m.channels))
	return m.channels[idx]
}
