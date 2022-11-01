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
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"

	wasmtimego "github.com/bytecodealliance/wasmtime-go"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/pkg/utils"
)

var (
	ErrAddrOverflow         = errors.New("addr overflow")
	ErrInstanceNotStart     = errors.New("instance has not started")
	ErrInstanceAlreadyStart = errors.New("instance has already started")
	ErrInvalidParam         = errors.New("invalid param")
	ErrRegisterNotFunc      = errors.New("register a non-func object")
	ErrRegisterArgNum       = errors.New("register func with invalid arg num")
	ErrRegisterArgType      = errors.New("register func with invalid arg type")
)

type Instance struct {
	vm           *VM
	store        *wasmtimego.Store
	module       *Module
	importObject []wasmtimego.AsExtern
	instance     *wasmtimego.Instance
	linker       *wasmtimego.Linker
	debug        *dwarfInfo
	abiList      []types.ABI

	lock     sync.Mutex
	started  uint32
	refCount int
	stopCond *sync.Cond

	// for cache
	memory    *wasmtimego.Memory
	funcCache sync.Map // string -> *wasmtimego.Function

	// user-defined data
	data interface{}
}

type InstanceOptions func(instance *Instance)

func InstanceWithDebug(debug *dwarfInfo) InstanceOptions {
	return func(instance *Instance) {
		if debug != nil {
			instance.debug = debug
		}
	}
}

func NewwasmtimegoInstance(vm *VM, module *Module, options ...InstanceOptions) *Instance {
	ins := &Instance{
		vm:     vm,
		module: module,
		lock:   sync.Mutex{},
	}
	ins.stopCond = sync.NewCond(&ins.lock)

	for _, option := range options {
		option(ins)
	}

	ins.importObject = make([]wasmtimego.AsExtern, 0)
	ins.linker = wasmtimego.NewLinker(vm.engine)
	err := ins.linker.DefineWasi()
	if err != nil {
		log.DefaultLogger.Errorf("[wasmtime][instance] DefineWasi failed, err: %v", err)
		return nil
	}

	wasiConfig := wasmtimego.NewWasiConfig()
	ins.store = wasmtimego.NewStore(vm.engine)
	ins.store.SetWasi(wasiConfig)
	return ins
}

func (w *Instance) GetData() interface{} {
	return w.data
}

func (w *Instance) SetData(data interface{}) {
	w.data = data
}

func (w *Instance) Acquire() bool {
	w.lock.Lock()
	defer w.lock.Unlock()

	if !w.checkStart() {
		return false
	}

	w.refCount++

	return true
}

func (w *Instance) Release() {
	w.lock.Lock()
	w.refCount--

	if w.refCount <= 0 {
		w.stopCond.Broadcast()
	}
	w.lock.Unlock()
}

func (w *Instance) Lock(data interface{}) {
	w.lock.Lock()
	w.data = data
}

func (w *Instance) Unlock() {
	w.data = nil
	w.lock.Unlock()
}

func (w *Instance) GetModule() types.WasmModule {
	return w.module
}

func (w *Instance) Start() error {
	w.abiList = abi.GetABIList(w)

	for _, abi := range w.abiList {
		abi.OnInstanceCreate(w)
	}

	ins, err := w.linker.Instantiate(w.store, w.module.module)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmtime][instance] Start fail to new wasmtimego instance, err: %v", err)
		return err
	}

	w.instance = ins

	f := w.instance.GetFunc(w.store, "_start")
	if f == nil {
		log.DefaultLogger.Errorf("[wasmtime][instance] Start fail to get export func: _start, err: %v", err)
		return err
	}

	_, err = f.Call(w.store)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmtime][instance] Start fail to call _start func, err: %v", err)
		w.HandleError(err)
		return err
	}

	for _, abi := range w.abiList {
		abi.OnInstanceStart(w)
	}

	atomic.StoreUint32(&w.started, 1)

	return nil
}

func (w *Instance) Stop() {
	utils.GoWithRecover(func() {
		w.lock.Lock()
		for w.refCount > 0 {
			w.stopCond.Wait()
		}
		swapped := atomic.CompareAndSwapUint32(&w.started, 1, 0)
		w.lock.Unlock()

		if swapped {
			for _, abi := range w.abiList {
				abi.OnInstanceDestroy(w)
			}
		}
	}, nil)
}

// return true is Instance is started, false if not started.
func (w *Instance) checkStart() bool {
	return atomic.LoadUint32(&w.started) == 1
}

type Val struct {
	ValType *wasmtimego.ValType
}

func (w *Instance) RegisterFunc(namespace string, funcName string, f interface{}) error {
	if w.checkStart() {
		log.DefaultLogger.Errorf("[wasmtime][instance] RegisterFunc not allow to register func after instance started, namespace: %v, funcName: %v",
			namespace, funcName)
		return ErrInstanceAlreadyStart
	}

	if namespace == "" || funcName == "" {
		log.DefaultLogger.Errorf("[wasmtime][instance] RegisterFunc invalid param, namespace: %v, funcName: %v", namespace, funcName)
		return ErrInvalidParam
	}

	if f == nil || reflect.ValueOf(f).IsNil() {
		log.DefaultLogger.Errorf("[wasmtime][instance] RegisterFunc f is nil")
		return ErrInvalidParam
	}

	if reflect.TypeOf(f).Kind() != reflect.Func {
		log.DefaultLogger.Errorf("[wasmtime][instance] RegisterFunc f is not func, actual type: %v", reflect.TypeOf(f))
		return ErrRegisterNotFunc
	}

	funcType := reflect.TypeOf(f)

	argsNum := funcType.NumIn()
	if argsNum < 1 {
		log.DefaultLogger.Errorf("[wasmtime][instance] RegisterFunc invalid args num: %v, must >= 1", argsNum)
		return ErrRegisterArgNum
	}

	argsKind := make([]*wasmtimego.ValType, argsNum-1)
	for i := 1; i < argsNum; i++ {
		argsKind[i-1] = convertFromGoType(funcType.In(i))
	}

	retsNum := funcType.NumOut()
	retsKind := make([]*wasmtimego.ValType, retsNum)
	for i := 0; i < retsNum; i++ {
		retsKind[i] = convertFromGoType(funcType.Out(i))
	}

	fwasmtimego := wasmtimego.NewFunc(
		w.store,
		wasmtimego.NewFuncType(argsKind, retsKind),
		func(caller *wasmtimego.Caller, args []wasmtimego.Val) (callRes []wasmtimego.Val, trap *wasmtimego.Trap) {
			defer func() {
				if r := recover(); r != nil {
					log.DefaultLogger.Errorf("[wasmtime][instance] RegisterFunc recover func call: %v, r: %v, stack: %v",
						funcName, r, string(debug.Stack()))
					callRes = nil
					err := fmt.Sprintf("panic [%v] when calling func [%v]", r, funcName)
					trap = wasmtimego.NewTrap(err)
				}
			}()

			callArgs := make([]reflect.Value, 1+len(args))
			callArgs[0] = reflect.ValueOf(w)

			for i, arg := range args {
				callArgs[i+1] = convertToGoTypes(arg)
			}

			callResult := reflect.ValueOf(f).Call(callArgs)

			ret := convertFromGoValue(callResult[0])

			return []wasmtimego.Val{ret}, nil
		},
	)

	w.linker.Define(namespace, funcName, fwasmtimego)

	//w.importObject = append(w.importObject, fwasmtimego)

	return nil
}

func (w *Instance) Malloc(size int32) (uint64, error) {
	if !w.checkStart() {
		log.DefaultLogger.Errorf("[wasmtime][instance] call malloc before starting instance")
		return 0, ErrInstanceNotStart
	}

	malloc, err := w.GetExportsFunc("malloc")
	if err != nil {
		return 0, err
	}

	addr, err := malloc.Call(size)
	if err != nil {
		w.HandleError(err)
		return 0, err
	}

	return uint64(addr.(int32)), nil
}

func (w *Instance) GetExportsFunc(funcName string) (types.WasmFunction, error) {
	if !w.checkStart() {
		log.DefaultLogger.Errorf("[wasmtime][instance] call GetExportsFunc before starting instance")
		return nil, ErrInstanceNotStart
	}

	if v, ok := w.funcCache.Load(funcName); ok {
		return v.(*Function), nil
	}

	f := w.instance.GetFunc(w.store, funcName)
	if f == nil {
		return nil, errors.New("func" + funcName + " is not exist")
	}

	ff := &Function{
		ins:      w,
		function: f,
	}

	w.funcCache.Store(funcName, ff)
	return ff, nil
}

func (w *Instance) GetExportsMem(memName string) ([]byte, error) {
	if !w.checkStart() {
		log.DefaultLogger.Errorf("[wasmtime][instance] call GetExportsMem before starting instance")
		return nil, ErrInstanceNotStart
	}

	if w.memory == nil {
		m := w.instance.GetExport(w.store, memName).Memory()
		if m == nil {
			return nil, errors.New("mem " + memName + " is not exist")
		}

		w.memory = m
	}

	return w.memory.UnsafeData(w.store), nil
}

func (w *Instance) GetMemory(addr uint64, size uint64) ([]byte, error) {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return nil, err
	}

	if int(addr) > len(mem) || int(addr+size) > len(mem) {
		return nil, ErrAddrOverflow
	}

	return mem[addr : addr+size], nil
}

func (w *Instance) PutMemory(addr uint64, size uint64, content []byte) error {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if int(addr) > len(mem) || int(addr+size) > len(mem) {
		return ErrAddrOverflow
	}

	copySize := uint64(len(content))
	if size < copySize {
		copySize = size
	}

	copy(mem[addr:], content[:copySize])

	return nil
}

func (w *Instance) GetByte(addr uint64) (byte, error) {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if int(addr) > len(mem) {
		return 0, ErrAddrOverflow
	}

	return mem[addr], nil
}

func (w *Instance) PutByte(addr uint64, b byte) error {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if int(addr) > len(mem) {
		return ErrAddrOverflow
	}

	mem[addr] = b

	return nil
}

func (w *Instance) GetUint32(addr uint64) (uint32, error) {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if int(addr) > len(mem) || int(addr+4) > len(mem) {
		return 0, ErrAddrOverflow
	}

	return binary.LittleEndian.Uint32(mem[addr:]), nil
}

func (w *Instance) PutUint32(addr uint64, value uint32) error {
	mem, err := w.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if int(addr) > len(mem) || int(addr+4) > len(mem) {
		return ErrAddrOverflow
	}

	binary.LittleEndian.PutUint32(mem[addr:], value)

	return nil
}

func (w *Instance) HandleError(err error) {
	var trapError *wasmtimego.Trap
	if !errors.As(err, &trapError) {
		return
	}

	trace := trapError.Frames()
	if trace == nil {
		return
	}

	log.DefaultLogger.Errorf("[wasmtime][instance] HandleError err: %v, trace:", err)

	if w.debug == nil {
		// do not have dwarf debug info
		for _, t := range trace {
			log.DefaultLogger.Errorf("[wasmtime][instance]\t funcIndex: %v, funcOffset: 0x%08x, moduleOffset: 0x%08x",
				t.FuncIndex(), t.FuncOffset(), t.ModuleOffset())
		}
	} else {
		for _, t := range trace {
			pc := uint64(t.ModuleOffset())
			line := w.debug.SeekPC(pc)
			if line != nil {
				log.DefaultLogger.Errorf("[wasmtime][instance]\t funcIndex: %v, funcOffset: 0x%08x, pc: 0x%08x %v:%v",
					t.FuncIndex(), t.FuncOffset(), pc, line.File.Name, line.Line)
			} else {
				log.DefaultLogger.Errorf("[wasmtime][instance]\t funcIndex: %v, funcOffset: 0x%08x, pc: 0x%08x fail to seek pc",
					t.FuncIndex(), t.FuncOffset(), t.ModuleOffset())
			}
		}
	}
}
