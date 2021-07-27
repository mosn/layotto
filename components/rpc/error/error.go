package error

import (
	"encoding/json"
	"fmt"
)

const (
	TimeoutCode int = iota
	UnavailebleCode
	InternalCode
)

type RpcError struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	marshaled string
}

func Error(code int, msg string) *RpcError {
	return &RpcError{Code: code, Msg: msg}
}

func Errorf(code int, format string, a ...interface{}) *RpcError {
	return &RpcError{Code: code, Msg: fmt.Sprintf(format, a...)}
}

func (e *RpcError) Error() string {
	if len(e.marshaled) > 0 {
		return e.marshaled
	}
	bytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	e.marshaled = string(bytes)
	return e.marshaled
}
