package wasmtime

import 	wasmtimego "github.com/bytecodealliance/wasmtime-go"

type Function struct {
	ins           *Instance
	function  *wasmtimego.Func
}

func (self *Function) Call(parameters ...interface{}) (interface{}, error) {
	return self.function.Call(self.ins.store, parameters...)
}

