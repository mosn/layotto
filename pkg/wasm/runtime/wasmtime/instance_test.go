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
	"reflect"
	"testing"
	"time"

	"fmt"

	wasmtimego "github.com/bytecodealliance/wasmtime-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mosn.io/mosn/pkg/mock"
	"mosn.io/mosn/pkg/types"
)

func TestRegisterFunc(t *testing.T) {
	vm := NewwasmtimegoVM()
	assert.Equal(t, vm.Name(), "wasmtime")

	wasm, _ := wasmtimego.Wat2Wasm(`(module (func (export "_start")))`)
	module := vm.NewModule(wasm)
	ins := module.NewInstance()

	// invalid namespace
	assert.Equal(t, ins.RegisterFunc("", "funcName", nil), ErrInvalidParam)

	// nil f
	assert.Equal(t, ins.RegisterFunc("TestRegisterFuncNamespace", "funcName", nil), ErrInvalidParam)

	var testStruct struct{}

	// f is not func
	assert.Equal(t, ins.RegisterFunc("TestRegisterFuncNamespace", "funcName", &testStruct), ErrRegisterNotFunc)

	// f is func with 0 args
	assert.Equal(t, ins.RegisterFunc("TestRegisterFuncNamespace", "funcName", func() {}), ErrRegisterArgNum)

	// wrong number of imports
	assert.Nil(t, ins.RegisterFunc("TestRegisterFuncNamespace", "funcName", func(f types.WasmInstance) {}))

	assert.Nil(t, ins.Start())

	assert.Equal(t, ins.RegisterFunc("TestRegisterFuncNamespace", "funcName", func(f types.WasmInstance) {}), ErrInstanceAlreadyStart)
}

func TestRegisterFuncRecoverPanic(t *testing.T) {
	vm := NewwasmtimegoVM()
	wasm, _ := wasmtimego.Wat2Wasm(`
			(module
				(import "TestRegisterFuncRecover" "somePanic" (func $somePanic (result i32)))
				(func (export "_start"))
				(func (export "panicTrigger") (param) (result i32)
					call $somePanic))
	`)
	module := vm.NewModule(wasm)
	ins := module.NewInstance()

	assert.Nil(t, ins.RegisterFunc("TestRegisterFuncRecover", "somePanic", func(instance types.WasmInstance) int32 {
		panic("some panic")
	}))

	assert.Nil(t, ins.Start())

	f, err := ins.GetExportsFunc("panicTrigger")
	assert.Nil(t, err)

	_, err = f.Call()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "panic [some panic] when calling func [somePanic]\nwasm backtrace:\n    0:   0x61 - <unknown>!<wasm function 2>\n")
}

func TestInstanceMalloc(t *testing.T) {
	vm := NewwasmtimegoVM()
	wasm, err := wasmtimego.Wat2Wasm(`
			(module
				(func (export "_start"))
				(func (export "malloc") (param i32) (result i32) i32.const 10))
	`)
	if err != nil {
		panic("Wat2Wasm failed")
	}
	module := vm.NewModule(wasm)
	ins := module.NewInstance()

	assert.Nil(t, ins.RegisterFunc("TestRegisterFuncRecover", "somePanic", func(instance types.WasmInstance) int32 {
		panic("some panic")
	}))

	assert.Nil(t, ins.Start())

	addr, err := ins.Malloc(100)
	assert.Nil(t, err)
	assert.Equal(t, addr, uint64(10))
}

func TestInstanceMem(t *testing.T) {
	vm := NewwasmtimegoVM()
	wasm, err := wasmtimego.Wat2Wasm(`(module (memory (export "memory") 1) (func (export "_start")))`)
	assert.Nil(t, err)
	module := vm.NewModule(wasm)
	ins := module.NewInstance()
	assert.Nil(t, ins.Start())

	m, err := ins.GetExportsMem("memory")
	assert.Nil(t, err)
	// A WebAssembly page has a constant size of 65,536 bytes, i.e., 64KiB
	assert.Equal(t, len(m), 1<<16)

	assert.Nil(t, ins.PutByte(uint64(100), 'a'))
	b, err := ins.GetByte(uint64(100))
	assert.Nil(t, err)
	assert.Equal(t, b, byte('a'))

	assert.Nil(t, ins.PutUint32(uint64(200), 99))
	u, err := ins.GetUint32(uint64(200))
	assert.Nil(t, err)
	assert.Equal(t, u, uint32(99))

	assert.Nil(t, ins.PutMemory(uint64(300), 10, []byte("1111111111")))
	bs, err := ins.GetMemory(uint64(300), 10)
	assert.Nil(t, err)
	assert.Equal(t, string(bs), "1111111111")
}

func TestInstanceData(t *testing.T) {
	vm := NewwasmtimegoVM()
	wasm, _ := wasmtimego.Wat2Wasm(`
			(module
				(func (export "_start")))
	`)
	module := vm.NewModule(wasm)
	ins := module.NewInstance()
	assert.Nil(t, ins.Start())

	var data int = 1
	ins.SetData(data)
	assert.Equal(t, ins.GetData().(int), 1)

	for i := 0; i < 10; i++ {
		ins.Lock(i)
		assert.Equal(t, ins.GetData().(int), i)
		ins.Unlock()
	}
}

func TestWasmtimegoTypes(t *testing.T) {
	testDatas := []struct {
		refType     reflect.Type
		refValue    reflect.Value
		refValKind  reflect.Kind
		wasmValKind wasmtimego.ValKind
	}{
		{reflect.TypeOf(int32(0)), reflect.ValueOf(int32(0)), reflect.Int32, wasmtimego.KindI32},
		{reflect.TypeOf(int64(0)), reflect.ValueOf(int64(0)), reflect.Int64, wasmtimego.KindI64},
		{reflect.TypeOf(float32(0)), reflect.ValueOf(float32(0)), reflect.Float32, wasmtimego.KindF32},
		{reflect.TypeOf(float64(0)), reflect.ValueOf(float64(0)), reflect.Float64, wasmtimego.KindF64},
	}

	for _, tc := range testDatas {
		assert.Equal(t, convertFromGoType(tc.refType).Kind(), tc.wasmValKind)
		assert.Equal(t, convertToGoTypes(convertFromGoValue(tc.refValue)).Kind(), tc.refValKind)
	}
}

func TestRefCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	destroyCount := 0
	abi := mock.NewMockABI(ctrl)
	abi.EXPECT().OnInstanceDestroy(gomock.Any()).DoAndReturn(func(types.WasmInstance) {
		destroyCount++
	})

	vm := NewwasmtimegoVM()
	wasm, _ := wasmtimego.Wat2Wasm(`
			(module
				(func (export "_start")))
	`)
	module := vm.NewModule(wasm)
	ins := NewwasmtimegoInstance(vm.(*VM), module.(*Module))

	ins.abiList = []types.ABI{abi}

	assert.False(t, ins.Acquire())

	ins.started = 1
	for i := 0; i < 100; i++ {
		assert.True(t, ins.Acquire())
	}
	assert.Equal(t, ins.refCount, 100)

	ins.Stop()
	ins.Stop() // double stop
	time.Sleep(time.Second)
	assert.Equal(t, ins.started, uint32(1))

	for i := 0; i < 100; i++ {
		ins.Release()
	}

	time.Sleep(time.Second)
	assert.False(t, ins.Acquire())
	assert.Equal(t, ins.started, uint32(0))
	assert.Equal(t, ins.refCount, 0)
	assert.Equal(t, destroyCount, 1)
}

const TextWat = `
(module
    ;; Import the required fd_write WASI function which will write the given io vectors to stdout
    ;; The function signature for fd_write is:
    ;; (File Descriptor, *iovs, iovs_len, nwritten) -> Returns number of bytes written
    (import "wasi_snapshot_preview1" "fd_write" (func $fd_write (param i32 i32 i32 i32) (result i32)))

    (memory 1)
    (export "memory" (memory 0))

    ;; Write 'hello world\n' to memory at an offset of 8 bytes
    ;; Note the trailing newline which is required for the text to appear
    (data (i32.const 8) "hello world\n")

    (func $main (export "_start")
        ;; Creating a new io vector within linear memory
        (i32.store (i32.const 0) (i32.const 8))  ;; iov.iov_base - This is a pointer to the start of the 'hello world\n' string
        (i32.store (i32.const 4) (i32.const 12))  ;; iov.iov_len - The length of the 'hello world\n' string

        (call $fd_write
            (i32.const 1) ;; file_descriptor - 1 for stdout
            (i32.const 0) ;; *iovs - The pointer to the iov array, which is stored at memory location 0
            (i32.const 1) ;; iovs_len - We're printing 1 string stored in an iov - so one.
            (i32.const 20) ;; nwritten - A place in memory to store the number of bytes written
        )
        drop ;; Discard the number of bytes written from the top of the stack
    )
)
`

func TestStartWithWASI(t *testing.T) {
	wasm, err := wasmtimego.Wat2Wasm(TextWat)
	if err != nil {
		fmt.Println(err)
	}
	vm := NewwasmtimegoVM()
	module := vm.NewModule(wasm)
	ins := module.NewInstance()
	assert.Nil(t, ins.Start())
}
