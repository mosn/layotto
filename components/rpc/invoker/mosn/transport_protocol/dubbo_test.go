/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package transport_protocol

import (
	"context"
	"sort"
	"strings"
	"testing"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/stretchr/testify/assert"
	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"
	"mosn.io/pkg/buffer"

	"mosn.io/layotto/components/rpc"
)

func Test_dubboProtocol_FromFrame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := buffer.NewIoBuffer(1000)
		data.Write(buildDubboRequestData(1))
		resp := dubbo.NewRpcResponse(nil, data)
		resp.Header.Set("key1", "value1")
		resp.Header.Status = dubbo.RespStatusOK
		d := newDubboProtocol()

		rsp, err := d.FromFrame(resp)
		assert.Nil(t, err)
		assert.Equal(t, "value1", rsp.Header.Get("key1"))
	})

	t.Run("fail", func(t *testing.T) {
		data := buffer.NewIoBuffer(1000)
		data.Write(buildDubboRequestData(1))
		resp := dubbo.NewRpcResponse(nil, data)
		resp.Header.Status = dubbo.RespStatusBadRequest
		d := newDubboProtocol()

		_, err := d.FromFrame(resp)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "dubbo error code 40"))
	})
}

func Test_dubboProtocol_Init(t *testing.T) {
	d := newDubboProtocol()
	err := d.Init(nil)
	assert.Nil(t, err)
}

func buildDubboRequestData(requestId uint64) []byte {
	service := hessian.Service{
		Path:      "io.mosn.layotto",
		Interface: "test",
		Group:     "test",
		Version:   "v1",
		Method:    "Call",
	}
	codec := hessian.NewHessianCodec(nil)
	header := hessian.DubboHeader{
		SerialID: 2,
		Type:     hessian.PackageRequest,
		ID:       int64(requestId),
	}
	body := hessian.NewRequest([]interface{}{}, nil)
	reqData, err := codec.Write(service, header, body)
	if err != nil {
		return nil
	}
	return reqData
}

func Test_dubboProtocol_ToFrame(t *testing.T) {
	d := newDubboProtocol()
	req := &rpc.RPCRequest{
		Ctx:         context.Background(),
		Id:          "1",
		Timeout:     100,
		Method:      "Hello",
		ContentType: "",
		Header: rpc.RPCHeader{
			"env":  []string{"test"},
			"name": []string{"bolt"},
		},
		Data: buildDubboRequestData(1),
	}
	frame := d.ToFrame(req)
	assert.NotNil(t, frame)
	assert.Equal(t, uint64(1), frame.GetRequestId())
	var headers []string
	frame.GetHeader().Range(func(key, value string) bool {
		// key dubbo's value is relative to version, ignore this key
		if key != "dubbo" {
			headers = append(headers, key+":"+value)
		}
		return true
	})
	sort.Slice(headers, func(i, j int) bool {
		return headers[i] < headers[j]
	})
	assert.Equal(t, "env:test,method:Call,name:bolt,service:io.mosn.layotto,version:v1", strings.Join(headers, ","))
}
