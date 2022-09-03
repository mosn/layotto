package wasmtime

import 	wasmtimego "github.com/bytecodealliance/wasmtime-go"

type Function struct {
	vm           *VM
	function  *wasmtimego.Func
}

func (self *Function) Call(parameters ...interface{}) (interface{}, error) {
	return self.function.Call(self.vm.store, parameters...)
}

