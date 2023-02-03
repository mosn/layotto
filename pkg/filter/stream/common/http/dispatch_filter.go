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

package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
	"mosn.io/api"
	mosnhttp "mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/variable"
)

type DispatchFilter struct {
	filterType     string
	requestHandler RequestHandler
	handler        api.StreamReceiverFilterHandler
}

func (dis *DispatchFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	dis.handler = handler
}

func (dis *DispatchFilter) OnDestroy() {}

func (dis *DispatchFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	// 1. log
	log.DefaultLogger.Debugf("[%v] receive %v pkt", dis.filterType, dis.filterType)
	path, err := variable.GetString(ctx, types.VarHttpRequestPath)
	if err != nil {
		dis.write404()
		return api.StreamFilterStop
	}
	log.DefaultLogger.Debugf("[%v] path: %v", dis.filterType, path)
	// 2. validate path
	resolver := NewPathResolver(path)
	// http path must be /{dis.filterType}/{endpoint_name}/{params}
	// So we can return 404 directly if it does not start with {dis.filterType}
	if resolver.Next() != dis.filterType {
		// illegal
		dis.write404()
		return api.StreamFilterStop
	}
	// 3. process request
	requestData := dis.handler.GetRequestData()
	if requestData != nil {
		ctx = context.WithValue(ctx, ContextKeyRequestData{}, requestData.Bytes())
	}
	epName := resolver.Next()
	endpoint, ok := dis.requestHandler.GetEndpoint(epName)
	if !ok {
		// illegal
		dis.write404()
		return api.StreamFilterStop
	}
	json, err := endpoint.Handle(ctx, resolver)
	var code int
	if err != nil {
		code = http.StatusInternalServerError
	} else {
		code = http.StatusOK
	}
	dis.writeJsonResult(json, code)
	return api.StreamFilterStop
}

func (dis *DispatchFilter) write404() {
	dis.writeJsonResult(nil, http.StatusNotFound)
}

func (dis *DispatchFilter) writeJsonResult(jsonObject map[string]interface{}, code int) {
	if code == 0 {
		code = http.StatusOK
	}
	// 0. marshal
	var byteSlice []byte
	if jsonObject != nil {
		var err error
		byteSlice, err = json.Marshal(jsonObject)
		if err != nil {
			log.DefaultLogger.Errorf("[%v][dispatch_filter]error when marshal result:%v", dis.filterType, err)
			code = http.StatusInternalServerError
		}
	}
	// 1. header
	fastHttpHeader := &fasthttp.ResponseHeader{}
	rspHeader := mosnhttp.ResponseHeader{
		ResponseHeader: fastHttpHeader,
	}
	rspHeader.Set("Content-Type", "application/json")
	rspHeader.SetStatusCode(code)
	// 2. body
	data := buffer.NewIoBufferBytes(byteSlice)
	// 3. write response
	dis.handler.SendDirectResponse(rspHeader, data, nil)
}
