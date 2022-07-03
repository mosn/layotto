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
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/opentracing/opentracing-go"
	jaegerc "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/trace/jaeger"
	"mosn.io/mosn/pkg/types"

	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/layotto/diagnostics/grpc"
)

const (
	serviceName              = "service_name"
	strategy                 = "strategy"
	agentHost                = "agent_host"
	collectorEndpoint        = "collector_endpoint"
	defaultServiceName       = "layotto"
	defaultJaegerAgentHost   = "127.0.0.1:6831"
	jaegerAgentHostKey       = "TRACE"
	appIDKey                 = "APP_ID"
	defaultCollectorEndpoint = "http://127.0.0.1:14268/api/traces"
	defaultStrategy          = "collector"
)

type grpcJaegerTracer struct {
	tracer opentracing.Tracer
}

type grpcJaegerSpan struct {
	*ltrace.Span
	ctx        context.Context
	trace      *grpcJaegerTracer
	jaegerSpan opentracing.Span
	spanCtx    jaegerc.SpanContext
}

func NewGrpcJaegerTracer(traceCfg map[string]interface{}) (api.Tracer, error) {
	// 1. construct the ReporterConfig, which is used to communicate with jaeger
	var reporter *config.ReporterConfig

	// Determining whether to start the agent
	strategy, err := getStrategy(traceCfg)

	if err != nil {
		return nil, err
	}

	if strategy == defaultStrategy {
		reporter = &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   getCollectorEndpoint(traceCfg),
		}
	} else {
		reporter = &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  getAgentHost(traceCfg),
		}
	}
	// 2. construct the Configuration
	cfg := config.Configuration{
		Disabled: false,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: reporter,
	}

	cfg.ServiceName = getServiceName(traceCfg)

	// 3. use the Configuration to construct a new tracer
	tracer, _, err := cfg.NewTracer()

	log.DefaultLogger.Infof("[layotto] [jaeger] [tracer] report service name:%s", getServiceName(traceCfg))

	if err != nil {
		log.DefaultLogger.Errorf("[layotto] [jaeger] [tracer] cannot initialize Jaeger Tracer")
		return nil, err
	}

	// 4. adapt to the `api.Tracer`
	return &grpcJaegerTracer{
		tracer: tracer,
	}, nil
}

func getAgentHost(traceCfg map[string]interface{}) string {
	if agentHost, ok := traceCfg[agentHost]; ok {
		return agentHost.(string)
	}

	//if TRACE is not set, get it from the env variable
	if host := os.Getenv(jaegerAgentHostKey); host != "" {
		return host
	}

	return defaultJaegerAgentHost
}

func getStrategy(traceCfg map[string]interface{}) (string, error) {
	if k, ok := traceCfg[strategy]; ok {
		if ok && (k.(string) == defaultStrategy || k.(string) == "agent") {
			return k.(string), nil
		} else if ok {
			return "", errors.New("Unknown Strategy")
		}
	}

	return defaultStrategy, nil
}

func getCollectorEndpoint(traceCfg map[string]interface{}) string {
	if collectorEndpoint, ok := traceCfg[collectorEndpoint]; ok {
		return collectorEndpoint.(string)
	}

	return defaultCollectorEndpoint
}

func getServiceName(traceCfg map[string]interface{}) string {
	if service, ok := traceCfg[serviceName]; ok {
		return service.(string)
	}

	//if service_name is not set, get it from the env variable
	if appID := os.Getenv(appIDKey); appID != "" {
		return appID + "_sidecar"
	}

	return defaultServiceName
}

func (t *grpcJaegerTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	header, ok := request.(*grpc.RequestInfo)
	if !ok {
		log.DefaultLogger.Debugf("[layotto] [jaeger] [tracer] unable to get request header, downstream trace ignored")
		return &jaeger.Span{}
	}

	//create entry span (downstream)
	sp, _ := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, header.FullMethod)

	//renew span context
	newSpanCtx, _ := sp.Context().(jaegerc.SpanContext)

	return &grpcJaegerSpan{
		trace:      t,
		ctx:        ctx,
		Span:       &ltrace.Span{},
		spanCtx:    newSpanCtx,
		jaegerSpan: sp,
	}
}

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
