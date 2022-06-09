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
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"mosn.io/api"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/layotto/diagnostics/protocol"
	"mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/types"
)

const (
	PORT = "9005"

	service_name              = "layotto"
	defaultReporterEndpoint   = "http://127.0.0.1:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
	configs                   = "config"
	reporter_endpoint         = "reporter_endpoint"
)

type grpcZipTracer struct {
	*zipkin.Tracer
}

func init() {
	trace.RegisterTracerBuilder("ZipKin", protocol.Layotto, NewGrpcZipTracer)
}

func NewGrpcZipTracer(traceCfg map[string]interface{}) (api.Tracer, error) {
	reporter := reporterhttp.NewReporter(getReporterEndpoint(traceCfg))
	tracer, err := zipkin.NewTracer(reporter)
	if err != nil {

	}

	return &grpcZipTracer{
		tracer,
	}, nil
}

func getReporterEndpoint(traceCfg map[string]interface{}) string {
	if cfg, ok := traceCfg[configs]; ok {
		endpoint := cfg.(map[string]interface{})
		if point, ok := endpoint[reporter_endpoint]; ok {
			return point.(string)
		}
	}

	return defaultReporterEndpoint
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
