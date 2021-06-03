package callback

import (
	"testing"

	"github.com/layotto/layotto/components/rpc"
	"github.com/stretchr/testify/assert"
)

func TestCallback(t *testing.T) {
	cb := NewCallback()

	phase := []string{}

	cb.AddBeforeInvoke(func(request *rpc.RPCRequest) (*rpc.RPCRequest, error) {
		phase = append(phase, "before")
		return request, nil
	})

	cb.AddAfterInvoke(func(response *rpc.RPCResponse) (*rpc.RPCResponse, error) {
		phase = append(phase, "after")
		return response, nil
	})

	cb.BeforeInvoke(nil)
	cb.AfterInvoke(nil)

	assert.Equal(t, "before", phase[0])
	assert.Equal(t, "after", phase[1])
}
