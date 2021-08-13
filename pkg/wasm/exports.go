package wasm

import "mosn.io/proxy-wasm-go-host/proxywasm"

type Exports interface {
	proxywasm.Exports

	// ProxyGetID return the id
	ProxyGetID() (string, error)
}
