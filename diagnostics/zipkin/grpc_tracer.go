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

package zipkin

import (
	"context"
	"time"

	"github.com/openzipkin/zipkin-go"
	"mosn.io/api"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/layotto/diagnostics/protocol"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/types"
)

const (
	PORT = "9005"

	SERVICE_NAME              = "layotto"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

type grpcZipTracer struct {
	*zipkin.Tracer
}

func init() {
	trace.RegisterTracerBuilder("ZipKin", protocol.Layotto, NewGrpcZipTracer)
}

func NewGrpcZipTracer(_ map[string]interface{}) (api.Tracer, error) {
	return nil, nil
}

func (tracer *grpcZipTracer) Start(ctx context.Context, request interface{}, _ time.Time) api.Span {
	//info, ok := request.(*grpc.RequestInfo)

	return nil
}

type grpcZipSpan struct {
	*ltrace.Span
	tracer *grpcZipTracer
	ctx    context.Context
	span   zipkin.Span
}

func (h *grpcZipSpan) TraceId() string {
	return ""
}

func (h *grpcZipSpan) InjectContext(requestHeaders types.HeaderMap, requestInfo api.RequestInfo) {
}

func (h *grpcZipSpan) SetRequestInfo(requestInfo api.RequestInfo) {
}

func (h *grpcZipSpan) FinishSpan() {
	h.span.Finish()
}
