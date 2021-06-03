package rpc

import (
	"context"
	"encoding/json"
	"strings"
)

type RPCHeader map[string][]string

func (r RPCHeader) Range(f func(key string, value string) bool) {
	if len(r) == 0 {
		return
	}

	for k, values := range r {
		if ok := f(k, strings.Join(values, ",")); !ok {
			break
		}
	}
}

func (r RPCHeader) Get(key string) string {
	if r == nil {
		return ""
	}
	values, ok := r[key]
	if !ok {
		return ""
	}
	return strings.Join(values, ",")
}

type RPCRequest struct {
	Ctx         context.Context
	Id          string
	Timeout     int32
	Method      string
	ContentType string
	Header      RPCHeader
	Data        []byte
}

type RPCResponse struct {
	Header      RPCHeader
	ContentType string
	Data        []byte
}

type RpcConfig struct {
	Config json.RawMessage
}

type Invoker interface {
	Init(config RpcConfig) error
	Invoke(ctx context.Context, req *RPCRequest) (*RPCResponse, error)
}

type Callback interface {
	AddBeforeInvoke(CallbackFunc)
	AddAfterInvoke(CallbackFunc)

	BeforeInvoke(*RPCRequest) (*RPCRequest, error)
	AfterInvoke(*RPCResponse) (*RPCResponse, error)
}

type CallbackFunc struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

type Channel interface {
	Do(*RPCRequest) (*RPCResponse, error)
}
