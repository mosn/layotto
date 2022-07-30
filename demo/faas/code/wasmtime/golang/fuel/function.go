package main

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

func main() {
	config := wasmtime.NewConfig()
	config.SetConsumeFuel(true)
	engine := wasmtime.NewEngineWithConfig(config)
	store := wasmtime.NewStore(engine)
	err := store.AddFuel(10)
	if err != nil {
		panic(err)
	}

	wasm, err := wasmtime.Wat2Wasm(`
	(module
         (func $addition (param i32 i32) (result i32)
             get_local 0
             get_local 1
             i32.add
         )
         (func (export "arithmetic") (param $n i32) ( param $m i32) (result i32)
             get_local $n
             get_local $m
             call $addition
             i32.const 3
             i32.mul
         )
    )`)

	if err != nil {
		panic(err)
	}
	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		panic(err)
	}
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{})
	if err != nil {
		panic(err)
	}

	arithmetic := instance.GetFunc(store, "arithmetic")
	if arithmetic == nil {
		panic("failed to find function `arithmetic`")
	}
	for n := 0; n < 10; n++ {
		fuelBefore, _ := store.FuelConsumed()
		output, err := arithmetic.Call(store, n, n+1)
		if err != nil {
			fmt.Println(err)
			break
		}
		fuelAfter, _ := store.FuelConsumed()
		fmt.Println(fmt.Sprintf("arithmetic(%d, %d) = %d [consumed %d fuel]\n", n, n+1, output, fuelAfter-fuelBefore))
		err = store.AddFuel(fuelAfter - fuelBefore)
		if err != nil {
			panic(err)
		}
	}
}
