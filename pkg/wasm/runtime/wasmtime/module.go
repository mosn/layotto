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
