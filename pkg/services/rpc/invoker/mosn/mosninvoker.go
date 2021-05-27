package mosn

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"github.com/layotto/layotto/pkg/services/rpc"
	"github.com/layotto/layotto/pkg/services/rpc/callback"
	"github.com/layotto/layotto/pkg/services/rpc/invoker/mosn/channel"
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
		if f := callback.GetBefore(before); f != nil {
			m.cb.AddBeforeInvoke(f)
		}
	}

	for _, after := range config.After {
		if f := callback.GetAfter(after); f != nil {
			m.cb.AddAfterInvoke(f)
		}
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

func (i *mosnInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	if req.Timeout == 0 {
		req.Timeout = 3000
	}
	log.DefaultLogger.Debugf("[runtime][rpc]request %+v", req)
	req, err := i.cb.BeforeInvoke(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]before filter error %s", err.Error())
		return nil, err
	}

	resp, err := i.getChannel().Do(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]error %s", err.Error())
		return nil, err
	}

	resp, err = i.cb.AfterInvoke(resp)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]after filter error %s", err.Error())
	}
	return resp, err
}

func (i *mosnInvoker) getChannel() rpc.Channel {
	idx := atomic.AddUint32(&i.rrIdx, 1) % uint32(len(i.channels))
	return i.channels[idx]
}
