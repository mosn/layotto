//go:build wasmtime
// +build wasmtime

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

package wasmtime

import (
	wasmtimego "github.com/bytecodealliance/wasmtime-go"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
)

func init() {
	wasm.RegisterWasmEngine("wasmtime", NewwasmtimegoVM())
}

type VM struct {
	engine *wasmtimego.Engine
}

func NewwasmtimegoVM() types.WasmVM {
	vm := &VM{}
	vm.Init()

	return vm
}

func (w *VM) Name() string {
	return "wasmtime"
}

func (w *VM) Init() {
	w.engine = wasmtimego.NewEngine()
}

func (w *VM) NewModule(wasmBytes []byte) types.WasmModule {
	if len(wasmBytes) == 0 {
		return nil
	}
	m, err := wasmtimego.NewModule(w.engine, wasmBytes)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmtimego][vm] fail to new module, err: %v", err)
		return nil
	}

	return NewwasmtimegoModule(w, m, wasmBytes)
}
