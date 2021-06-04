package wasm

import (
	"github.com/stretchr/testify/assert"
	"mosn.io/proxy-wasm-go-host/proxywasm"
	"testing"
)

func TestImportsHandler(t *testing.T) {
	d := &LayottoHandler{}
	assert.Equal(t, d.Log(proxywasm.LogLevelCritical, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelError, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelWarn, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelInfo, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelDebug, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelTrace, "msg"), proxywasm.WasmResultOk)
}
