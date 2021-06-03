package callback

import (
	"encoding/json"

	"github.com/layotto/layotto/components/rpc"
	"mosn.io/pkg/log"
)

func RegisterBeforeInvoke(f BeforeFactory) {
	beforeInvokeRegistry[f.Name()] = f
}

func RegisterAfterInvoke(f AfterFactory) {
	afterInvokeRegistry[f.Name()] = f
}

type BeforeFactory interface {
	Name() string
	Init(json.RawMessage) error
	Create() func(*rpc.RPCRequest) (*rpc.RPCRequest, error)
}

type AfterFactory interface {
	Name() string
	Init(json.RawMessage) error
	Create() func(*rpc.RPCResponse) (*rpc.RPCResponse, error)
}

var (
	beforeInvokeRegistry = map[string]BeforeFactory{}
	afterInvokeRegistry  = map[string]AfterFactory{}
)

func NewCallback() rpc.Callback {
	return &callback{}
}

type callback struct {
	beforeInvoke []func(*rpc.RPCRequest) (*rpc.RPCRequest, error)
	afterInvoke  []func(*rpc.RPCResponse) (*rpc.RPCResponse, error)
}

func (c *callback) AddBeforeInvoke(conf rpc.CallbackFunc) {
	f, ok := beforeInvokeRegistry[conf.Name]
	if !ok {
		log.DefaultLogger.Errorf("[runtime][rpc]can't find before filter %s", conf.Name)
		return
	}
	if err := f.Init(conf.Config); err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]init before filter err %s", err.Error())
		return
	}
	c.beforeInvoke = append(c.beforeInvoke, f.Create())
}

func (c *callback) AddAfterInvoke(conf rpc.CallbackFunc) {
	f, ok := afterInvokeRegistry[conf.Name]
	if !ok {
		log.DefaultLogger.Errorf("[runtime][rpc]can't find after filter %s", conf.Name)
		return
	}
	if err := f.Init(conf.Config); err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]init after filter err %s", err.Error())
		return
	}
	c.afterInvoke = append(c.afterInvoke, f.Create())
}

func (c *callback) BeforeInvoke(request *rpc.RPCRequest) (*rpc.RPCRequest, error) {
	var err error
	for _, cb := range c.beforeInvoke {
		if request, err = cb(request); err != nil {
			return nil, err
		}
	}
	return request, err
}

func (c *callback) AfterInvoke(response *rpc.RPCResponse) (*rpc.RPCResponse, error) {
	var err error
	for _, cb := range c.afterInvoke {
		if response, err = cb(response); err != nil {
			return nil, err
		}
	}
	return response, err
}
