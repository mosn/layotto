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
package transport_protocol

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/header"

	"mosn.io/layotto/components/rpc"
)

func Test_boltCommon_FromFrame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := &bolt.Response{
			Content: buffer.NewIoBuffer(100),
			ResponseHeader: bolt.ResponseHeader{
				BytesHeader: header.BytesHeader{
					Kvs: []header.BytesKV{
						{
							Key:   []byte("key1"),
							Value: []byte("value1"),
						},
					},
				},
			},
		}
		resp.Content.Write([]byte("hello"))
		resp.ResponseStatus = bolt.ResponseStatusSuccess
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": "bolt",
		}
		err := b.Init(conf)
		assert.Nil(t, err)

		rsp, err := b.FromFrame(resp)
		assert.Nil(t, err)
		assert.Equal(t, "hello", string(rsp.Data))
		assert.Equal(t, "value1", rsp.Header.Get("key1"))
	})

	t.Run("fail", func(t *testing.T) {
		resp := &bolt.Response{}
		resp.ResponseStatus = bolt.ResponseStatusError
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": "bolt",
		}
		err := b.Init(conf)
		assert.Nil(t, err)

		f, err := b.FromFrame(resp)
		assert.Nil(t, err)
		assert.Equal(t, false, f.Success)
		assert.True(t, strings.Contains(f.Error.Error(), "bolt error code 1"))
	})

	t.Run("fail", func(t *testing.T) {
		resp := &bolt.Response{}
		resp.ResponseStatus = bolt.ResponseStatusServerDeserialException
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": "bolt",
		}
		err := b.Init(conf)
		assert.Nil(t, err)

		f, err := b.FromFrame(resp)
		assert.Nil(t, err)
		assert.Equal(t, false, f.Success)
		assert.True(t, strings.Contains(f.Error.Error(), "bolt error code 18"))
	})

	t.Run("fail", func(t *testing.T) {
		resp := &bolt.Response{}
		resp.ResponseStatus = bolt.ResponseStatusServerSerialException
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": "bolt",
		}
		err := b.Init(conf)
		assert.Nil(t, err)

		f, err := b.FromFrame(resp)
		assert.Nil(t, err)
		assert.Equal(t, false, f.Success)
		assert.True(t, strings.Contains(f.Error.Error(), "bolt error code 17"))
	})

	t.Run("fail", func(t *testing.T) {
		resp := &bolt.Response{}
		resp.ResponseStatus = bolt.ResponseStatusCodecException
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": "bolt",
		}
		err := b.Init(conf)
		assert.Nil(t, err)

		f, err := b.FromFrame(resp)
		assert.Nil(t, err)
		assert.Equal(t, false, f.Success)
		assert.True(t, strings.Contains(f.Error.Error(), "bolt error code 9"))
	})
}

func Test_boltCommon_Init(t *testing.T) {
	t.Run("empty conf", func(t *testing.T) {
		b := &boltCommon{}
		err := b.Init(nil)
		assert.NotNil(t, err)
		assert.Equal(t, "missing bolt classname", err.Error())
	})

	t.Run("class not exist", func(t *testing.T) {
		b := &boltCommon{}
		conf := map[string]interface{}{
			"key": "value",
		}
		err := b.Init(conf)
		assert.NotNil(t, err)
		assert.Equal(t, "bolt need class", err.Error())
	})

	t.Run("bolt class not string", func(t *testing.T) {
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": 1,
		}
		err := b.Init(conf)
		assert.NotNil(t, err)
		assert.Equal(t, "bolt class not string", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		b := &boltCommon{}
		conf := map[string]interface{}{
			"class": "bolt",
		}
		err := b.Init(conf)
		assert.Nil(t, err)
		assert.Equal(t, "bolt", b.className)
	})

}

func Test_boltProtocol_ToFrame(t *testing.T) {
	b := newBoltProtocol()
	req := &rpc.RPCRequest{
		Ctx:         context.Background(),
		Id:          "1",
		Timeout:     100,
		Method:      "Hello",
		ContentType: "application/json",
		Header: rpc.RPCHeader{
			"env":  []string{"test"},
			"name": []string{"bolt"},
		},
		Data: []byte("hello world"),
	}
	frame := b.ToFrame(req)
	assert.NotNil(t, frame)
	assert.Equal(t, uint64(0), frame.GetRequestId())
	assert.Equal(t, int32(100), frame.GetTimeout())
	assert.Equal(t, "hello world", frame.GetData().String())
	var headers []string
	frame.GetHeader().Range(func(key, value string) bool {
		headers = append(headers, key+":"+value)
		return true
	})
	sort.Slice(headers, func(i, j int) bool {
		return headers[i] < headers[j]
	})
	assert.Equal(t, "env:test,name:bolt", strings.Join(headers, ","))

}

func Test_boltv2Protocol_ToFrame(t *testing.T) {
	b := newBoltV2Protocol()
	req := &rpc.RPCRequest{
		Ctx:         context.Background(),
		Id:          "1",
		Timeout:     100,
		Method:      "Hello",
		ContentType: "application/json",
		Header: rpc.RPCHeader{
			"env":  []string{"test"},
			"name": []string{"bolt"},
		},
		Data: []byte("hello world"),
	}
	frame := b.ToFrame(req)
	assert.NotNil(t, frame)
	assert.Equal(t, uint64(0), frame.GetRequestId())
	assert.Equal(t, int32(100), frame.GetTimeout())
	assert.Equal(t, "hello world", frame.GetData().String())
	var headers []string
	frame.GetHeader().Range(func(key, value string) bool {
		headers = append(headers, key+":"+value)
		return true
	})
	sort.Slice(headers, func(i, j int) bool {
		return headers[i] < headers[j]
	})
	assert.Equal(t, "env:test,name:bolt", strings.Join(headers, ","))
}
