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

package wasm

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	mockwasm "mosn.io/layotto/pkg/mock/wasm"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/mock"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/header"
	"mosn.io/proxy-wasm-go-host/proxywasm/common"
	v1 "mosn.io/proxy-wasm-go-host/proxywasm/v1"
)

func TestMapEncodeAndDecode(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	b := common.EncodeMap(m)
	n := common.DecodeMap(b)
	assert.Equal(t, m, n)
}

func TestFilter_Append(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		senderFilterHandler api.StreamSenderFilterHandler
		responseBuffer      api.IoBuffer
	}
	tests := []struct {
		name         string
		fields       fields
		mockAndCheck func(t *testing.T, f *Filter, ctrl *gomock.Controller)
	}{
		{
			name: "normal",
			fields: fields{
				senderFilterHandler: mock.NewMockStreamSenderFilterHandler(ctrl),
				responseBuffer:      buffer.NewIoBufferString("test"),
			},
			mockAndCheck: func(t *testing.T, f *Filter, ctrl *gomock.Controller) {
				f.senderFilterHandler.(*mock.MockStreamSenderFilterHandler).EXPECT().SetResponseData(f.responseBuffer).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				senderFilterHandler: tt.fields.senderFilterHandler,
				responseBuffer:      tt.fields.responseBuffer,
			}

			tt.mockAndCheck(t, f, ctrl)
			assert.Equal(t, api.StreamFilterContinue, f.Append(context.Background(), nil, nil, nil))
		})
	}
}

func TestFilterGetMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userData := map[string]string{
		"plugin_config": "plugin_config",
	}

	type fields struct {
		factory               *FilterConfigFactory
		contextID             int32
		pluginUsed            *WasmPlugin
		receiverFilterHandler api.StreamReceiverFilterHandler
		senderFilterHandler   api.StreamSenderFilterHandler
		requestBuffer         api.IoBuffer
		responseBuffer        api.IoBuffer
	}
	type wants struct {
		rootContextID       int32
		httpRequestHeader   common.HeaderMap
		httpRequestBody     common.IoBuffer
		httpRequestTrailer  common.HeaderMap
		httpResponseHeader  common.HeaderMap
		httpResponseBody    common.IoBuffer
		httpResponseTrailer common.HeaderMap
	}
	tests := []struct {
		name         string
		fields       fields
		wants        wants
		mockAndCheck func(ctrl *gomock.Controller, f *Filter)
	}{
		{
			name: "normal",
			fields: fields{
				factory: &FilterConfigFactory{
					RootContextID: 1,
				},
				pluginUsed: &WasmPlugin{
					pluginName: "plugin_1",
					plugin:     mock.NewMockWasmPlugin(ctrl),
					config: &filterConfigItem{
						UserData: userData,
					},
				},
				receiverFilterHandler: mock.NewMockStreamReceiverFilterHandler(ctrl),
				senderFilterHandler:   mock.NewMockStreamSenderFilterHandler(ctrl),
				requestBuffer:         buffer.NewIoBufferString("request body"),
				responseBuffer:        buffer.NewIoBufferString("response body"),
			},
			wants: wants{
				rootContextID: 1,
				httpRequestHeader: &proxywasm010.HeaderMapWrapper{
					HeaderMap: &header.CommonHeader{"request": "header"},
				},
				httpRequestBody: buffer.NewIoBufferString("request body"),
				httpRequestTrailer: &proxywasm010.HeaderMapWrapper{
					HeaderMap: &header.CommonHeader{"request": "trailer"},
				},
				httpResponseHeader: &proxywasm010.HeaderMapWrapper{
					HeaderMap: &header.CommonHeader{"response": "header"},
				},
				httpResponseBody: buffer.NewIoBufferString("response body"),
				httpResponseTrailer: &proxywasm010.HeaderMapWrapper{
					HeaderMap: &header.CommonHeader{"response": "trailer"},
				},
			},
			mockAndCheck: func(ctrl *gomock.Controller, f *Filter) {
				f.pluginUsed.plugin.(*mock.MockWasmPlugin).EXPECT().
					GetVmConfig().Return(v2.WasmVmConfig{Engine: "test"}).Times(1)
				f.receiverFilterHandler.(*mock.MockStreamReceiverFilterHandler).EXPECT().
					GetRequestHeaders().Return(&header.CommonHeader{"request": "header"}).Times(1)
				f.receiverFilterHandler.(*mock.MockStreamReceiverFilterHandler).EXPECT().
					GetRequestTrailers().Return(&header.CommonHeader{"request": "trailer"}).Times(1)
				f.senderFilterHandler.(*mock.MockStreamSenderFilterHandler).EXPECT().
					GetResponseHeaders().Return(&header.CommonHeader{"response": "header"}).Times(1)
				f.senderFilterHandler.(*mock.MockStreamSenderFilterHandler).EXPECT().
					GetResponseTrailers().Return(&header.CommonHeader{"response": "trailer"}).Times(1)
			},
		},
		{
			name: "get nil",
			fields: fields{
				factory: &FilterConfigFactory{
					RootContextID: 1,
				},
				pluginUsed: &WasmPlugin{
					pluginName: "plugin_1",
					plugin:     mock.NewMockWasmPlugin(ctrl),
					config: &filterConfigItem{
						UserData: userData,
					},
				},
			},
			wants: wants{
				rootContextID:       1,
				httpRequestHeader:   nil,
				httpRequestBody:     nil,
				httpRequestTrailer:  nil,
				httpResponseHeader:  nil,
				httpResponseBody:    nil,
				httpResponseTrailer: nil,
			},
			mockAndCheck: func(ctrl *gomock.Controller, f *Filter) {
				f.pluginUsed.plugin.(*mock.MockWasmPlugin).EXPECT().
					GetVmConfig().Return(v2.WasmVmConfig{Engine: "test"}).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				factory:               tt.fields.factory,
				contextID:             tt.fields.contextID,
				pluginUsed:            tt.fields.pluginUsed,
				receiverFilterHandler: tt.fields.receiverFilterHandler,
				senderFilterHandler:   tt.fields.senderFilterHandler,
				requestBuffer:         tt.fields.requestBuffer,
				responseBuffer:        tt.fields.responseBuffer,
			}

			tt.mockAndCheck(ctrl, f)
			assert.Equalf(t, tt.wants.rootContextID, f.GetRootContextID(), "GetRootContextID()")
			assert.NotNilf(t, f.GetVmConfig(), "GetVmConfig()")
			assert.NotNilf(t, f.GetPluginConfig(), "GetPluginConfig()")

			assert.Equalf(t, tt.wants.httpRequestHeader, f.GetHttpRequestHeader(), "GetHttpRequestHeader()")
			assert.Equalf(t, tt.wants.httpRequestTrailer, f.GetHttpRequestTrailer(), "GetHttpRequestTrailer()")
			if f.GetHttpRequestBody() != nil {
				assert.Equalf(t, tt.wants.httpRequestBody.Len(), f.GetHttpRequestBody().Len(), "GetHttpRequestBody()")
			}

			assert.Equalf(t, tt.wants.httpResponseHeader, f.GetHttpResponseHeader(), "GetHttpResponseHeader()")
			assert.Equalf(t, tt.wants.httpResponseTrailer, f.GetHttpResponseTrailer(), "GetHttpResponseTrailer()")
			if f.GetHttpResponseBody() != nil {
				assert.Equalf(t, tt.wants.httpResponseBody.Len(), f.GetHttpResponseBody().Len(), "GetHttpRequestBody()")
			}
		})
	}
}

func TestFilter_OnReceive(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	mockAbiFunc := func(ctrl *gomock.Controller, instance *mock.MockWasmInstance) *mock.MockABI {
		a := mock.NewMockABI(ctrl)
		abiFactory := func(instance types.WasmInstance) types.ABI {
			return a
		}
		abi.RegisterABI(AbiV2, abiFactory)
		return a
	}

	type fields struct {
		router *Router
	}
	type args struct {
		ctx      context.Context
		headers  api.HeaderMap
		buf      buffer.IoBuffer
		trailers api.HeaderMap
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		mockAndCheck func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap)
		want         api.StreamFilterStatus
	}{
		{
			name: "get id fail",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				headers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("", false).Times(1)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "f.router.GetRandomPluginByID error",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{},
				},
			},
			args: args{
				ctx:     context.Background(),
				headers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "pluginABI is nil",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				headers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				gomock.InOrder(
					headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().GetInstance().Return(nil).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().ReleaseInstance(gomock.Any()).Times(1),
				)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "exports.ProxyOnContextCreate error",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				headers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				instance := mock.NewMockWasmInstance(ctrl)
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)

				gomock.InOrder(
					headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().GetInstance().Return(instance).Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports),
					instance.EXPECT().Lock(a).Times(1),
					exports.EXPECT().ProxyOnContextCreate(gomock.Any(), gomock.Any()).
						Return(errors.New("exports.ProxyOnContextCreate error")).Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "exports.ProxyOnRequestHeaders result is not ActionContinue",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				headers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				instance := mock.NewMockWasmInstance(ctrl)
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)

				gomock.InOrder(
					headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().GetInstance().Return(instance).Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports),
					instance.EXPECT().Lock(a).Times(1),
					exports.EXPECT().ProxyOnContextCreate(gomock.Any(), gomock.Any()).Return(nil).Times(1),
					headers.(*mock.MockHeaderMap).EXPECT().Range(gomock.Any()).Times(1),
					exports.EXPECT().ProxyOnRequestHeaders(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(v1.ActionPause, nil).Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "exports.ProxyOnRequestBody result is not ActionContinue",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				headers: mock.NewMockHeaderMap(ctrl),
				buf:     buffer.NewIoBufferString("test"),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				instance := mock.NewMockWasmInstance(ctrl)
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)

				gomock.InOrder(
					headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().GetInstance().Return(instance).Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports),
					instance.EXPECT().Lock(a).Times(1),
					exports.EXPECT().ProxyOnContextCreate(gomock.Any(), gomock.Any()).Return(nil).Times(1),
					headers.(*mock.MockHeaderMap).EXPECT().Range(gomock.Any()).Times(1),
					exports.EXPECT().ProxyOnRequestHeaders(gomock.Any(), gomock.Any(), gomock.Any()).Return(v1.ActionContinue, nil).Times(1),
					exports.EXPECT().ProxyOnRequestBody(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(v1.ActionPause, nil).Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "exports.ProxyOnRequestTrailers result is not ActionContinue",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				headers:  mock.NewMockHeaderMap(ctrl),
				buf:      buffer.NewIoBufferString("test"),
				trailers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				instance := mock.NewMockWasmInstance(ctrl)
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)

				gomock.InOrder(
					headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().GetInstance().Return(instance).Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports),
					instance.EXPECT().Lock(a).Times(1),
					exports.EXPECT().ProxyOnContextCreate(gomock.Any(), gomock.Any()).Return(nil).Times(1),
					headers.(*mock.MockHeaderMap).EXPECT().Range(gomock.Any()).Times(1),
					exports.EXPECT().ProxyOnRequestHeaders(gomock.Any(), gomock.Any(), gomock.Any()).Return(v1.ActionContinue, nil).Times(1),
					exports.EXPECT().ProxyOnRequestBody(gomock.Any(), gomock.Any(), gomock.Any()).Return(v1.ActionContinue, nil).Times(1),
					trailers.(*mock.MockHeaderMap).EXPECT().Range(gomock.Any()).Times(1),
					exports.EXPECT().ProxyOnRequestTrailers(gomock.Any(), gomock.Any()).
						Return(v1.ActionPause, nil).Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
			want: api.StreamFilterStop,
		},
		{
			name: "normal",
			fields: fields{
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				headers:  mock.NewMockHeaderMap(ctrl),
				buf:      buffer.NewIoBufferString("test"),
				trailers: mock.NewMockHeaderMap(ctrl),
			},
			mockAndCheck: func(headers api.HeaderMap, plugin types.WasmPlugin, f *Filter, trailers api.HeaderMap) {
				instance := mock.NewMockWasmInstance(ctrl)
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)

				gomock.InOrder(
					headers.(*mock.MockHeaderMap).EXPECT().Get("id").Return("id_1", true).Times(1),
					plugin.(*mock.MockWasmPlugin).EXPECT().GetInstance().Return(instance).Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports),
					instance.EXPECT().Lock(a).Times(1),
					exports.EXPECT().ProxyOnContextCreate(gomock.Any(), gomock.Any()).Return(nil).Times(1),
					headers.(*mock.MockHeaderMap).EXPECT().Range(gomock.Any()).Times(1),
					exports.EXPECT().ProxyOnRequestHeaders(gomock.Any(), gomock.Any(), gomock.Any()).Return(v1.ActionContinue, nil).Times(1),
					exports.EXPECT().ProxyOnRequestBody(gomock.Any(), gomock.Any(), gomock.Any()).Return(v1.ActionContinue, nil).Times(1),
					trailers.(*mock.MockHeaderMap).EXPECT().Range(gomock.Any()).Times(1),
					exports.EXPECT().ProxyOnRequestTrailers(gomock.Any(), gomock.Any()).Return(v1.ActionContinue, nil).Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
			want: api.StreamFilterContinue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				router: tt.fields.router,
			}

			layottoPlugin, _ := f.router.GetRandomPluginByID("id_1")
			var plugin types.WasmPlugin
			if layottoPlugin != nil {
				plugin = layottoPlugin.plugin
			}
			tt.mockAndCheck(tt.args.headers, plugin, f, tt.args.trailers)

			assert.Equalf(t, tt.want, f.OnReceive(tt.args.ctx, tt.args.headers, tt.args.buf, tt.args.trailers), "OnReceive(%v, %v, %v, %v)", tt.args.ctx, tt.args.headers, tt.args.buf, tt.args.trailers)
		})
	}
}

func TestFilter_SetReceiveFilterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		handler api.StreamReceiverFilterHandler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handler is nil",
			args: args{
				handler: nil,
			},
		},
		{
			name: "handler is not nil",
			args: args{
				handler: mock.NewMockStreamReceiverFilterHandler(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{}
			f.SetReceiveFilterHandler(tt.args.handler)
		})
	}
}

func TestFilter_SetSenderFilterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		handler api.StreamSenderFilterHandler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handler is nil",
			args: args{
				handler: nil,
			},
		},
		{
			name: "handler is not nil",
			args: args{
				handler: mock.NewMockStreamSenderFilterHandler(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{}
			f.SetSenderFilterHandler(tt.args.handler)
		})
	}
}

func TestFilter_OnDestroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := &Filter{}
	f.OnDestroy()
}

func TestFilter_releaseUsedInstance(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	type fields struct {
		contextID  int32
		pluginUsed *WasmPlugin
		instance   types.WasmInstance
		abi        types.ABI
		exports    Exports
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      assert.ErrorAssertionFunc
		mockAndCheck func(ctrl *gomock.Controller, f *Filter)
	}{
		{
			name: "p.instance is nil",
			fields: fields{
				pluginUsed: mockLayottoWasmPlugin(
					"plugin_1", 2, mock.NewMockWasmPlugin(ctrl),
				),
				instance: nil,
			},
			wantErr:      assert.NoError,
			mockAndCheck: func(ctrl *gomock.Controller, f *Filter) {},
		},
		{
			name: "p.exports.ProxyOnDone error",
			fields: fields{
				pluginUsed: mockLayottoWasmPlugin(
					"plugin_1", 2, mock.NewMockWasmPlugin(ctrl),
				),
				instance:  mock.NewMockWasmInstance(ctrl),
				abi:       mock.NewMockABI(ctrl),
				exports:   mockwasm.NewMockExports(ctrl),
				contextID: 1,
			},
			wantErr: assert.Error,
			mockAndCheck: func(ctrl *gomock.Controller, f *Filter) {
				gomock.InOrder(
					f.instance.(*mock.MockWasmInstance).EXPECT().Lock(f.abi).Times(1),
					f.exports.(*mockwasm.MockExports).EXPECT().ProxyOnDone(f.contextID).Return(int32(0), errors.New("error")).Times(1),
				)
			},
		},
		{
			name: "p.exports.ProxyOnDelete error",
			fields: fields{
				pluginUsed: mockLayottoWasmPlugin(
					"plugin_1", 2, mock.NewMockWasmPlugin(ctrl),
				),
				instance:  mock.NewMockWasmInstance(ctrl),
				abi:       mock.NewMockABI(ctrl),
				exports:   mockwasm.NewMockExports(ctrl),
				contextID: 1,
			},
			wantErr: assert.Error,
			mockAndCheck: func(ctrl *gomock.Controller, f *Filter) {
				gomock.InOrder(
					f.instance.(*mock.MockWasmInstance).EXPECT().Lock(f.abi).Times(1),
					f.exports.(*mockwasm.MockExports).EXPECT().ProxyOnDone(f.contextID).Return(int32(0), nil).Times(1),
					f.exports.(*mockwasm.MockExports).EXPECT().ProxyOnDelete(f.contextID).Return(errors.New("error")).Times(1),
				)
			},
		},
		{
			name: "releaseUsedInstance success",
			fields: fields{
				pluginUsed: mockLayottoWasmPlugin(
					"plugin_1", 2, mock.NewMockWasmPlugin(ctrl),
				),
				instance:  mock.NewMockWasmInstance(ctrl),
				abi:       mock.NewMockABI(ctrl),
				exports:   mockwasm.NewMockExports(ctrl),
				contextID: 1,
			},
			wantErr: assert.NoError,
			mockAndCheck: func(ctrl *gomock.Controller, f *Filter) {
				gomock.InOrder(
					f.instance.(*mock.MockWasmInstance).EXPECT().Lock(f.abi).Times(1),
					f.exports.(*mockwasm.MockExports).EXPECT().ProxyOnDone(f.contextID).Return(int32(0), nil).Times(1),
					f.exports.(*mockwasm.MockExports).EXPECT().ProxyOnDelete(f.contextID).Return(nil).Times(1),
					f.instance.(*mock.MockWasmInstance).EXPECT().Unlock().Times(1),
					f.pluginUsed.plugin.(*mock.MockWasmPlugin).EXPECT().ReleaseInstance(f.instance).Times(1),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				contextID:  tt.fields.contextID,
				pluginUsed: tt.fields.pluginUsed,
				instance:   tt.fields.instance,
				abi:        tt.fields.abi,
				exports:    tt.fields.exports,
			}

			tt.mockAndCheck(ctrl, f)
			tt.wantErr(t, f.releaseUsedInstance(), "releaseUsedInstance()")
		})
	}
}

func TestNewFilter(t *testing.T) {
	type args struct {
		ctx     context.Context
		factory *FilterConfigFactory
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(res *Filter)
	}{
		{
			name: "normal",
			args: args{
				ctx:     context.Background(),
				factory: &FilterConfigFactory{},
			},
			checkFunc: func(res *Filter) {
				assert.NotNil(t, res)
				assert.Equal(t, int32(1), res.contextID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.checkFunc(NewFilter(tt.args.ctx, tt.args.factory))
		})
	}
}

func TestWasmPlugin_GetPluginConfig(t *testing.T) {
	type fields struct {
		pluginName        string
		plugin            types.WasmPlugin
		rootContextID     int32
		config            *filterConfigItem
		vmConfigBytes     buffer.IoBuffer
		pluginConfigBytes buffer.IoBuffer
	}
	tests := []struct {
		name   string
		fields fields
		want   common.IoBuffer
	}{
		{
			name: "p.pluginConfigBytes is not nil",
			fields: fields{
				pluginConfigBytes: buffer.NewIoBuffer(1),
			},
			want: buffer.NewIoBuffer(1),
		},
		{
			name: "p.config.UserData is empty",
			fields: fields{
				config: &filterConfigItem{
					UserData: map[string]string{},
				},
			},
			want: nil,
		},
		{
			name: "p.config.UserData is not empty",
			fields: fields{
				config: &filterConfigItem{
					UserData: map[string]string{
						"plugin_config": "plugin_config",
					},
				},
			},
			want: buffer.NewIoBufferBytes(common.EncodeMap(map[string]string{
				"plugin_config": "plugin_config",
			})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &WasmPlugin{
				pluginName:        tt.fields.pluginName,
				plugin:            tt.fields.plugin,
				rootContextID:     tt.fields.rootContextID,
				config:            tt.fields.config,
				vmConfigBytes:     tt.fields.vmConfigBytes,
				pluginConfigBytes: tt.fields.pluginConfigBytes,
			}
			assert.Equalf(t, tt.want, p.GetPluginConfig(), "GetPluginConfig()")
		})
	}
}

func TestWasmPlugin_GetVmConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		pluginName        string
		plugin            types.WasmPlugin
		rootContextID     int32
		config            *filterConfigItem
		vmConfigBytes     buffer.IoBuffer
		pluginConfigBytes buffer.IoBuffer
	}
	tests := []struct {
		name      string
		fields    fields
		mockFunc  func(plugin types.WasmPlugin)
		checkFunc func(res common.IoBuffer)
	}{
		{
			name: "p.vmConfigBytes is not nil",
			fields: fields{
				vmConfigBytes: buffer.NewIoBuffer(1),
			},
			mockFunc: func(plugin types.WasmPlugin) {},
			checkFunc: func(res common.IoBuffer) {
				assert.Equal(t, buffer.NewIoBuffer(1), res)
			},
		},
		{
			name: "normal",
			fields: fields{
				plugin: mock.NewMockWasmPlugin(ctrl),
			},
			mockFunc: func(plugin types.WasmPlugin) {
				vmConfig := v2.WasmVmConfig{
					Engine: "wasmtime",
					Path:   "no_file",
				}
				plugin.(*mock.MockWasmPlugin).EXPECT().GetVmConfig().Return(vmConfig).Times(1)
			},
			checkFunc: func(res common.IoBuffer) {
				assert.NotNil(t, res)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &WasmPlugin{
				pluginName:        tt.fields.pluginName,
				plugin:            tt.fields.plugin,
				rootContextID:     tt.fields.rootContextID,
				config:            tt.fields.config,
				vmConfigBytes:     tt.fields.vmConfigBytes,
				pluginConfigBytes: tt.fields.pluginConfigBytes,
			}

			tt.mockFunc(tt.fields.plugin)
			res := p.GetVmConfig()
			tt.checkFunc(res)
		})
	}
}
