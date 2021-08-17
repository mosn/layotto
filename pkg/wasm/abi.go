/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wasm

import (
	sdkTypes "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
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
	ih := &ImportsHandler{
		DefaultImportsHandler: v1.DefaultImportsHandler{},
		ioBuffer:              &common.CommonBuffer{},
	}
	abi.SetImports(ih)
	abi.importsHandler = ih
	return abi
}

// easy for extension
type AbiV2Impl struct {
	v1.ABIContext

	importsHandler *ImportsHandler
}

var (
	_ types.ABIHandler = &AbiV2Impl{}
	_ Exports          = &AbiV2Impl{}
)

func (a *AbiV2Impl) Name() string {
	return AbiV2
}

func (a *AbiV2Impl) GetABIExports() interface{} {
	return a
}

func (a *AbiV2Impl) ProxyGetID() (string, error) {
	ff, err := a.Instance.GetExportsFunc("proxy_get_id")
	if err != nil {
		return "", err
	}
	res, err := ff.Call()
	if err != nil {
		a.Instance.HandleError(err)
		return "", err
	}
	a.Imports.Wait()

	status := sdkTypes.Status(res.(int32))
	if err := sdkTypes.StatusToError(status); err != nil {
		a.Instance.HandleError(err)
		return "", err
	}

	return string(a.importsHandler.ioBuffer.Bytes()), nil
}

type ImportsHandler struct {
	v1.DefaultImportsHandler

	ioBuffer common.IoBuffer
}

var _ proxywasm.ImportsHandler = &ImportsHandler{}

func (h *ImportsHandler) GetFuncCallData() common.IoBuffer {
	h.ioBuffer = &common.CommonBuffer{}
	return h.ioBuffer
}
