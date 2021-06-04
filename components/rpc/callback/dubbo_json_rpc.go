package callback

import (
	"encoding/json"

	"github.com/layotto/layotto/components/rpc"
)

func init() {
	RegisterBeforeInvoke(&beforeFactory{})
}

type beforeFactory struct {
}

func (b *beforeFactory) Name() string {
	return "dubbo_json_rpc"
}

func (b *beforeFactory) Init(json.RawMessage) error {
	return nil
}

func (b *beforeFactory) Create() func(*rpc.RPCRequest) (*rpc.RPCRequest, error) {
	return func(request *rpc.RPCRequest) (*rpc.RPCRequest, error) {
		request.Header["x-services"] = []string{request.Id}
		request.Header["x-method"] = []string{request.Method}
		request.Header["content-type"] = []string{"application/json"}
		request.Header["accept"] = []string{"application/json"}

		request.Method = request.Id
		return request, nil
	}
}
