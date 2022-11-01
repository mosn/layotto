// +build wasmer

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
	"strings"

	wasmtimego "github.com/bytecodealliance/wasmtime-go"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

type Module struct {
	vm          *VM
	module      *wasmtimego.Module
	abiNameList []string
	debug       *dwarfInfo
	rawBytes    []byte
}

func NewwasmtimegoModule(vm *VM, module *wasmtimego.Module, wasmBytes []byte) *Module {
	m := &Module{
		vm:       vm,
		module:   module,
		rawBytes: wasmBytes,
	}

	m.Init()

	return m
}

func (w *Module) Init() {
	log.DefaultLogger.Infof("[wasmtime][module] Init module")

	w.abiNameList = w.GetABINameList()

	// parse dwarf info from wasm data bytes
	if debug := parseDwarf(w.rawBytes); debug != nil {
		w.debug = debug
	}

	// release raw bytes, the parsing of dwarf info is the only place that uses module raw bytes
	w.rawBytes = nil
}

func (w *Module) NewInstance() types.WasmInstance {
	if w.debug != nil {
		return NewwasmtimegoInstance(w.vm, w, InstanceWithDebug(w.debug))
	}

	return NewwasmtimegoInstance(w.vm, w)
}

func (w *Module) GetABINameList() []string {
	abiNameList := make([]string, 0)

	exportList := w.module.Exports()

	for _, export := range exportList {
		//if export.Type() == wasmtimego.FuncType{
		if strings.HasPrefix(export.Name(), "proxy_abi") {
			abiNameList = append(abiNameList, export.Name())
		}
		//}
	}

	return abiNameList
}
