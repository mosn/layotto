package wasm

type Exports interface {
	ProxyGetID() (string, error)
}
