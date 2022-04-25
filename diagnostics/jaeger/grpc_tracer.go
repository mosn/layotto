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

package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaegerc "github.com/uber/jaeger-client-go"
	"mosn.io/api"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/layotto/diagnostics/protocol"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/trace/jaeger"
	"mosn.io/mosn/pkg/types"
	"strings"
	"time"
)

type httpJaegerTracer struct {
	*jaegerc.Tracer
}

type httpJaegerSpan struct {
	*ltrace.Span
	ctx     context.Context
	trace   *httpJaegerTracer
	spanCtx jaegerc.SpanContext
}

func init() {
	trace.RegisterTracerBuilder(jaeger.DriverName, protocol.Layotto, NewHttpJaegerTracer)
}

func NewHttpJaegerTracer(_ map[string]interface{}) (api.Tracer, error) {
	return &httpJaegerTracer{}, nil
}

func (t *httpJaegerTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	header, ok := request.(http.RequestHeader)
	if !ok {
		log.DefaultLogger.Debugf("[jaeger] [tracer] [layotto] unable to get request header, downstream trace ignored")
		return &jaeger.Span{}
	}

	sp, spanCtx := t.getSpan(ctx, header, startTime)

	ext.HTTPMethod.Set(sp, string(header.Method()))
	ext.HTTPUrl.Set(sp, string(header.RequestURI()))

	return &httpJaegerSpan{
		trace:   t,
		ctx:     ctx,
		Span:    &ltrace.Span{},
		spanCtx: spanCtx,
	}
}

func (t *httpJaegerTracer) getSpan(ctx context.Context, header http.RequestHeader, startTime time.Time) (opentracing.Span, jaegerc.SpanContext) {
	httpHeaderPropagator := jaegerc.NewHTTPHeaderPropagator(newDefaultHeadersConfig(), *jaegerc.NewNullMetrics())
	// extract metadata
	spanCtx, _ := httpHeaderPropagator.Extract(jaeger.HTTPHeadersCarrier(header))
	sp, _ := opentracing.StartSpanFromContextWithTracer(ctx, t.Tracer, getOperationName(header.RequestURI()), opentracing.ChildOf(spanCtx), opentracing.StartTime(startTime))

	//renew span context
	newSpanCtx, ok := sp.Context().(jaegerc.SpanContext)
	if !ok {
		return sp, spanCtx
	}

	return sp, newSpanCtx
}

func getOperationName(uri []byte) string {
	arr := strings.Split(string(uri), "?")
	return arr[0]
}

func newDefaultHeadersConfig() *jaegerc.HeadersConfig {
	return &jaegerc.HeadersConfig{
		JaegerDebugHeader:        jaegerc.JaegerDebugHeader,
		JaegerBaggageHeader:      jaegerc.JaegerBaggageHeader,
		TraceContextHeaderName:   jaegerc.TraceContextHeaderName,
		TraceBaggageHeaderPrefix: jaegerc.TraceBaggageHeaderPrefix,
	}
}

func (s *httpJaegerSpan) TraceId() string {
	return string(s.spanCtx.SpanID())
}

func (s *httpJaegerSpan) InjectContext(requestHeaders types.HeaderMap, requestInfo api.RequestInfo) {
}

func (s *httpJaegerSpan) SetRequestInfo(requestInfo api.RequestInfo) {
}

func (s *httpJaegerSpan) FinishSpan() {
	s.Span.FinishSpan()
}
