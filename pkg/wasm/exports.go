package wasm

import (
	proxywasm "mosn.io/proxy-wasm-go-host/proxywasm/v1"
)

type Exports interface {
	proxywasm.Exports

	// ProxyGetID return the id
	ProxyGetID() (string, error)
}
