package callback

import "github.com/layotto/layotto/components/rpc"

var (
	beforeInvokeRegistry = map[string]func(*rpc.RPCRequest) (*rpc.RPCRequest, error){}
	afterInvokeRegistry  = map[string]func(*rpc.RPCResponse) (*rpc.RPCResponse, error){}
)

func RegistBefore(callbackFunc string, f func(*rpc.RPCRequest) (*rpc.RPCRequest, error)) {
	beforeInvokeRegistry[callbackFunc] = f
}

func RegistAfter(callbackFunc string, f func(response *rpc.RPCResponse) (*rpc.RPCResponse, error)) {
	afterInvokeRegistry[callbackFunc] = f
}

func GetBefore(callbackFunc rpc.CallbackFunc) func(*rpc.RPCRequest) (*rpc.RPCRequest, error) {
	return beforeInvokeRegistry[callbackFunc.Name]
}

func GetAfter(callbackFunc rpc.CallbackFunc) func(*rpc.RPCResponse) (*rpc.RPCResponse, error) {
	return afterInvokeRegistry[callbackFunc.Name]
}

func NewCallback() rpc.Callback {
	return &callback{}
}

type callback struct {
	beforeInvoke []func(*rpc.RPCRequest) (*rpc.RPCRequest, error)
	afterInvoke  []func(*rpc.RPCResponse) (*rpc.RPCResponse, error)
}

func (c *callback) AddBeforeInvoke(f func(request *rpc.RPCRequest) (req *rpc.RPCRequest, err error)) {
	c.beforeInvoke = append(c.beforeInvoke, f)
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

func (c *callback) AddAfterInvoke(f func(response *rpc.RPCResponse) (resp *rpc.RPCResponse, reterr error)) {
	c.afterInvoke = append(c.afterInvoke, f)
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
