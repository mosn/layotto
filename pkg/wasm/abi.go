package wasm

import (
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
	v1 "mosn.io/mosn/pkg/wasm/abi/proxywasm010"
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
