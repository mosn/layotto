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
	"fmt"
	jaegerc "github.com/uber/jaeger-client-go"
	"mosn.io/api"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/layotto/diagnostics/grpc"
	"mosn.io/layotto/diagnostics/protocol"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/trace/jaeger"
	"mosn.io/mosn/pkg/types"
	"os"
	"time"
)

const (
	serviceName            = "service_name"
	agentHost              = "agent_host"
	defaultServiceName     = "layotto"
	defaultJaegerAgentHost = "0.0.0.0:6831"
	jaegerAgentHostKey     = "TRACE"
	appIDKey               = "APP_ID"
)

type grpcJaegerTracer struct {
	*jaegerc.Tracer
}

type grpcJaegerSpan struct {
	*ltrace.Span
	ctx        context.Context
	trace      *grpcJaegerTracer
	jaegerSpan opentracing.Span
	spanCtx    jaegerc.SpanContext
}

func init() {
	trace.RegisterTracerBuilder(jaeger.DriverName, protocol.Layotto, NewGrpcJaegerTracer)
}

func NewGrpcJaegerTracer(traceCfg map[string]interface{}) (api.Tracer, error) {
	cfg := config.Configuration{
		Disabled: false,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  getAgentHost(traceCfg),
		},
	}

	cfg.ServiceName = getServiceName(traceCfg)
	tracer, _, err := cfg.NewTracer()

	log.DefaultLogger.Infof("[jaeger] [tracer] jaeger agent host:%s, report service name:%s",
		getAgentHost(traceCfg), getServiceName(traceCfg))

	if err != nil {
		log.DefaultLogger.Errorf("[jaeger] [tracer] [http1] cannot initialize Jaeger Tracer")
		return nil, err
	}

	return &grpcJaegerTracer{}, nil
}

func getServiceName(traceCfg map[string]interface{}) string {
	if service, ok := traceCfg[serviceName]; ok {
		return service.(string)
	}

	//if service_name is not set, get it from the env variable
	if appID := os.Getenv(appIDKey); appID != "" {
		return fmt.Sprintf("%s_sidecar", appID)
	}

	return defaultServiceName
}

func (t *grpcJaegerTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	header, ok := request.(*grpc.RequestInfo)
	if !ok {
		log.DefaultLogger.Debugf("[jaeger] [tracer] [layotto] unable to get request header, downstream trace ignored")
		return &jaeger.Span{}
	}

	fmt.Println(t)
	fmt.Println(header)

	//create entry span (downstream)
	//sp, _ := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, header.FullMethod)

	//renew span context
	//newSpanCtx, ok := sp.Context().(jaegerc.SpanContext)

	//sp, spanCtx := t.getSpan(ctx, header, startTime)
	//
	//ext.HTTPMethod.Set(sp, string(header.Method()))
	//ext.HTTPUrl.Set(sp, string(header.RequestURI()))

	return &grpcJaegerSpan{
		trace: t,
		ctx:   ctx,
		Span:  &ltrace.Span{},
		//spanCtx:    newSpanCtx,
		//jaegerSpan: sp,
	}
}

//func (t *grpcJaegerTracer) getSpan(ctx context.Context, header http.RequestHeader, startTime time.Time) (opentracing.Span, jaegerc.SpanContext) {
//	httpHeaderPropagator := jaegerc.NewHTTPHeaderPropagator(newDefaultHeadersConfig(), *jaegerc.NewNullMetrics())
//	// extract metadata
//	spanCtx, _ := httpHeaderPropagator.Extract(jaeger.HTTPHeadersCarrier(header))
//	sp, _ := opentracing.StartSpanFromContextWithTracer(ctx, t.trace, getOperationName(header.RequestURI()), opentracing.ChildOf(spanCtx), opentracing.StartTime(startTime))
//
//	//renew span context
//	newSpanCtx, ok := sp.Context().(jaegerc.SpanContext)
//	if !ok {
//		return sp, spanCtx
//	}
//
//	return sp, newSpanCtx
//}
//
//func getOperationName(uri []byte) string {
//	arr := strings.Split(string(uri), "?")
//	return arr[0]
//}
//
//func newDefaultHeadersConfig() *jaegerc.HeadersConfig {
//	return &jaegerc.HeadersConfig{
//		JaegerDebugHeader:        jaegerc.JaegerDebugHeader,
//		JaegerBaggageHeader:      jaegerc.JaegerBaggageHeader,
//		TraceContextHeaderName:   jaegerc.TraceContextHeaderName,
//		TraceBaggageHeaderPrefix: jaegerc.TraceBaggageHeaderPrefix,
//	}
//}

func (s *grpcJaegerSpan) TraceId() string {
	return s.spanCtx.TraceID().String()
}

func (s *grpcJaegerSpan) InjectContext(requestHeaders types.HeaderMap, requestInfo api.RequestInfo) {
}

func (s *grpcJaegerSpan) SetRequestInfo(requestInfo api.RequestInfo) {
}

func (s *grpcJaegerSpan) FinishSpan() {
	s.jaegerSpan.Finish()
}
