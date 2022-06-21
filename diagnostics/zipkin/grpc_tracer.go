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
	"fmt"
	"time"

	"mosn.io/layotto/diagnostics/grpc"

	"github.com/openzipkin/zipkin-go"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"mosn.io/api"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/log"

	ltrace "mosn.io/layotto/components/trace"
)

const (
	service_name       = "service_name"
	reporter_endpoint  = "reporter_endpoint"
	recorder_host_post = "recorder_host_post"
)

type grpcZipTracer struct {
	*zipkin.Tracer
}

type grpcZipSpan struct {
	*ltrace.Span
	tracer *grpcZipTracer
	ctx    context.Context
	span   zipkin.Span
}

func NewGrpcZipTracer(traceCfg map[string]interface{}) (api.Tracer, error) {
	point, err := getReporterEndpoint(traceCfg)
	if err != nil {
		return nil, err
	}

	reporter := reporterhttp.NewReporter(point)

	name, err := getServerName(traceCfg)
	if err != nil {
		return nil, err
	}

	host_post, err := getRecorderHostPort(traceCfg)
	if err != nil {
		return nil, err
	}

	endpoint, err := zipkin.NewEndpoint(name, host_post)
	if err != nil {
		log.DefaultLogger.Errorf("[layotto] [zipkin] [tracer] unable to create zipkin reporter endpoint")
		return nil, err
	}

	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithTraceID128Bit(true))
	if err != nil {
		log.DefaultLogger.Errorf("[layotto] [zipkin] [tracer] cannot initialize zipkin Tracer")
		return nil, err
	}

	log.DefaultLogger.Infof("[layotto] [zipkin] [tracer] create success")

	return &grpcZipTracer{
		tracer,
	}, nil
}

func getRecorderHostPort(traceCfg map[string]interface{}) (string, error) {
	if recorder, ok := traceCfg[recorder_host_post]; ok {
		return recorder.(string), nil
	}

	return "", fmt.Errorf("[layotto] [zipkin] [tracer] no config zipkin server host and port")
}

func getReporterEndpoint(traceCfg map[string]interface{}) (string, error) {
	if point, ok := traceCfg[reporter_endpoint]; ok {
		return point.(string), nil
	}

	return "", fmt.Errorf("[layotto] [zipkin] [tracer] no config zipkin reporter endpoint")
}

func getServerName(traceCfg map[string]interface{}) (string, error) {
	if name, ok := traceCfg[service_name]; ok {
		return name.(string), nil
	}

	return "", fmt.Errorf("[layotto] [zipkin] [tracer] no config zipkin server name")
}

func (t *grpcZipTracer) Start(ctx context.Context, request interface{}, _ time.Time) api.Span {
	info, ok := request.(*grpc.RequestInfo)
	if !ok {
		log.DefaultLogger.Debugf("[layotto] [zipkin] [tracer] unable to get request header, downstream trace ignored")
		return nil
	}

	// start span
	span := t.StartSpan(info.FullMethod)

	return &grpcZipSpan{
		tracer: t,
		ctx:    ctx,
		Span:   &ltrace.Span{},
		span:   span,
	}
}

func (s *grpcZipSpan) TraceId() string {
	return s.span.Context().TraceID.String()
}

func (s *grpcZipSpan) InjectContext(requestHeaders types.HeaderMap, requestInfo api.RequestInfo) {
}

func (s *grpcZipSpan) SetRequestInfo(requestInfo api.RequestInfo) {
}

func (s *grpcZipSpan) FinishSpan() {
	s.span.Finish()
}
