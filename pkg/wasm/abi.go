package wasm

import (
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
	v1 "mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/proxy-wasm-go-host/common"
	"mosn.io/proxy-wasm-go-host/proxywasm"
)

const AbiV2 = "proxy_abi_version_0_2_0"

func init() {
	abi.RegisterABI(AbiV2, abiImplFactory)
}

func abiImplFactory(instance types.WasmInstance) types.ABI {
	abi := &AbiV2Impl{}
	abi.SetInstance(instance)
	return abi
}

// easy for extension
type AbiV2Impl struct {
	v1.ABIContext
}

func (a *AbiV2Impl) Name() string {
	return AbiV2
}

func (a *AbiV2Impl) OnInstanceCreate(instance types.WasmInstance) {
	proxywasm.RegisterImports(instance)
	_ = instance.RegisterFunc("env", "proxy_get_configuration", ProxyGetConfiguration)
	_ = instance.RegisterFunc("env", "proxy_continue_request", ProxyContinueRequest)
	_ = instance.RegisterFunc("env", "proxy_continue_response", ProxyContinueResponse)
	_ = instance.RegisterFunc("env", "proxy_send_local_response", ProxyContinueResponse)
	_ = instance.RegisterFunc("env", "proxy_clear_route_cache", ProxyContinueResponse)
	_ = instance.RegisterFunc("env", "proxy_grpc_stream", ProxyContinueResponse)
	_ = instance.RegisterFunc("env", "proxy_grpc_send", ProxyContinueResponse)
	_ = instance.RegisterFunc("env", "proxy_grpc_cancel", ProxyContinueResponse)
	_ = instance.RegisterFunc("env", "proxy_grpc_close", ProxyContinueResponse)
}

func ProxyGetConfiguration(instance common.WasmInstance, returnValueData int32, returnValueSize int32) int32 {
	return proxywasm.WasmResultOk.Int32()
}

func ProxyContinueRequest(instance common.WasmInstance) int32 {
	return proxywasm.WasmResultOk.Int32()
}

func ProxyContinueResponse(instance common.WasmInstance) int32 {
	return proxywasm.WasmResultOk.Int32()
}
