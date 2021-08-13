package proxywasm

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

//export proxy_id
func proxyID(contextID uint32) types.Action {
	//ctx, ok := currentState.httpStreams[contextID]
	//if !ok {
	//	panic("invalid context on proxy_id")
	//}
	//currentState.setActiveContextID(contextID)
	//return ctx.OnHttpRequestBody(bodySize, endOfStream)
}
