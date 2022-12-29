// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wasm

import (
	"errors"

	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
	v1 "mosn.io/mosn/pkg/wasm/abi/proxywasm010"
)

const AbiV2 = "proxy_abi_version_0_1_0"

// registers an abi
func init() {
	abi.RegisterABI(AbiV2, abiImplFactory)
}

// initialization an abi
func abiImplFactory(instance types.WasmInstance) types.ABI {
	abi := &AbiV2Impl{}
	abi.SetInstance(instance)
	return abi
}

// easy for extension
type AbiV2Impl struct {
	v1.ABIContext
}

var (
	_ types.ABIHandler = &AbiV2Impl{}
	_ Exports          = &AbiV2Impl{}
)

// Get abi name
func (a *AbiV2Impl) Name() string {
	return AbiV2
}

// Get abi
func (a *AbiV2Impl) GetABIExports() interface{} {
	return a
}

// Get id
func (a *AbiV2Impl) ProxyGetID() (string, error) {
	// store the funcName and common.WasmFunction, then return the common.WasmFunction
	ff, err := a.Instance.GetExportsFunc("proxy_get_id")
	if err != nil {
		return "", err
	}
	_, err = ff.Call()
	if err != nil {
		a.Instance.HandleError(err)
		return "", err
	}
	a.Imports.Wait()
	// generate functionId
	functionId := string(a.Imports.GetFuncCallData().Bytes())
	if functionId == "" {
		return "", errors.New("")
	}

	return functionId, nil
}
