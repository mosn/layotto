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

package skywalking

import (
	"context"
	"time"

	"mosn.io/layotto/diagnostics/grpc"

	"github.com/SkyAPM/go2sky"
	language_agent "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/trace/skywalking"
	"mosn.io/mosn/pkg/types"

	ltrace "mosn.io/layotto/components/trace"
)

func NewGrpcSkyTracer(_ map[string]interface{}) (api.Tracer, error) {
	return &grpcSkyTracer{}, nil
}

type grpcSkyTracer struct {
	*go2sky.Tracer
}

func (tracer *grpcSkyTracer) SetGO2SkyTracer(t *go2sky.Tracer) {
	tracer.Tracer = t
}

func (tracer *grpcSkyTracer) Start(ctx context.Context, request interface{}, _ time.Time) api.Span {
	info, ok := request.(*grpc.RequestInfo)
	if !ok {
		log.DefaultLogger.Debugf("[SkyWalking] [tracer] [layotto] unable to get request header, downstream trace ignored")
		return skywalking.NoopSpan
	}

	// create entry span (downstream)
	entry, nCtx, err := tracer.CreateEntrySpan(ctx, info.FullMethod, func() (sw8 string, err error) {
		return
	})

	if err != nil {
		log.DefaultLogger.Errorf("[SkyWalking] [tracer] [http1] create entry span error, err: %v", err)
		return skywalking.NoopSpan
	}
	entry.Tag(go2sky.TagHTTPMethod, "POST")
	entry.Tag(go2sky.TagURL, info.FullMethod)
	entry.SetComponent(skywalking.MOSNComponentID)
	entry.SetSpanLayer(language_agent.SpanLayer_Http)

	return &grpcSkySpan{
		tracer: tracer,
		ctx:    nCtx,
		carrier: &skywalking.SpanCarrier{
			EntrySpan: entry,
		},
		Span: &ltrace.Span{},
	}
}

type grpcSkySpan struct {
	*ltrace.Span
	tracer  *grpcSkyTracer
	ctx     context.Context
	carrier *skywalking.SpanCarrier
}

func (h *grpcSkySpan) TraceId() string {
	return go2sky.TraceID(h.ctx)
}

func (h *grpcSkySpan) InjectContext(requestHeaders types.HeaderMap, requestInfo api.RequestInfo) {
}

func (h *grpcSkySpan) SetRequestInfo(requestInfo api.RequestInfo) {
}

func (h *grpcSkySpan) FinishSpan() {
	entry := h.carrier.EntrySpan
	if h.Tag(ltrace.LAYOTTO_REQUEST_RESULT) == "1" {
		entry.Error(time.Now(), skywalking.ErrorLog)
		entry.Tag(go2sky.TagStatusCode, "500")
	} else {
		entry.Tag(go2sky.TagStatusCode, "200")
	}

	entry.End()
}
