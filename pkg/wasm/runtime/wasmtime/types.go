package wasmtime

import (
	"reflect"

	wasmtimego "github.com/bytecodealliance/wasmtime-go"
	"mosn.io/mosn/pkg/log"
)

func convertFromGoType(t reflect.Type) *wasmtimego.ValType {
	switch t.Kind() {
	case reflect.Int32:
		return wasmtimego.NewValType(wasmtimego.KindI32)
	case reflect.Int64:
		return wasmtimego.NewValType(wasmtimego.KindI64)
	case reflect.Float32:
		return wasmtimego.NewValType(wasmtimego.KindF32)
	case reflect.Float64:
		return wasmtimego.NewValType(wasmtimego.KindF64)
	default:
		log.DefaultLogger.Errorf("[wasmtimego][type] convertFromGoType unsupported type: %v", t.Kind().String())
	}

	return nil
}

func convertToGoTypes(in wasmtimego.Val) reflect.Value {
	switch in.Kind() {
	case wasmtimego.KindI32:
		return reflect.ValueOf(in.I32())
	case wasmtimego.KindI64:
		return reflect.ValueOf(in.I64())
	case wasmtimego.KindF32:
		return reflect.ValueOf(in.F32())
	case wasmtimego.KindF64:
		return reflect.ValueOf(in.F64())
	}

	return reflect.Value{}
}

func convertFromGoValue(val reflect.Value) wasmtimego.Val {
	switch val.Kind() {
	case reflect.Int32:
		return wasmtimego.ValI32(int32(val.Int()))
	case reflect.Int64:
		return wasmtimego.ValI64(int64(val.Int()))
	case reflect.Float32:
		return wasmtimego.ValF32(float32(val.Float()))
	case reflect.Float64:
		return wasmtimego.ValF64(float64(val.Float()))
	default:
		log.DefaultLogger.Errorf("[wasmtimego][type] convertFromGoValue unsupported val type: %v", val.Kind().String())
	}

	return wasmtimego.Val{}
}
