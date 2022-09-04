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
	store  *wasmtimego.Store
}

func NewwasmtimegoVM() types.WasmVM {
	vm := &VM{}
	vm.Init()

	return vm
}

func (w *VM) Name() string {
	return "wasmtimego"
}

func (w *VM) Init() {
	w.engine = wasmtimego.NewEngine()
	w.store = wasmtimego.NewStore(w.engine)
}

func (w *VM) NewModule(wasmBytes []byte) types.WasmModule {
	if len(wasmBytes) == 0 {
		return nil
	}
	//wasm, err := wasmtimego.wat2wasm(string(wasmBytes))
	m, err := wasmtimego.NewModule(w.engine, wasmBytes)
	if err != nil {
		log.DefaultLogger.Errorf("[wasmtimego][vm] fail to new module, err: %v", err)
		return nil
	}

	return NewwasmtimegoModule(w, m, wasmBytes)
}
/*

config := wasmtimego.NewConfig()
	config.SetConsumeFuel(true)
	engine := wasmtimego.NewEngineWithConfig(config)
	store := wasmtimego.NewStore(engine)
	err := store.AddFuel(10)
	if err != nil {
		panic(err)
	}

	wasm, err := wasmtimego.Wat2Wasm(`
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
	module, err := wasmtimego.NewModule(store.Engine, wasm)
	if err != nil {
		panic(err)
	}
	instance, err := wasmtimego.NewInstance(store, module, []wasmtimego.AsExtern{})
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
 */