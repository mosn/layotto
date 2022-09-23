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
package trace

import (
	"time"

	"mosn.io/api"

	"mosn.io/mosn/pkg/trace/sofa"
	"mosn.io/mosn/pkg/types"
)

type Span struct {
	StartTime     time.Time
	EndTime       time.Time
	traceId       string
	spanId        string
	parentSpanId  string
	tags          [sofa.TRACE_END]string
	operationName string
}

func (span *Span) SetTraceId(id string) {
	span.traceId = id
}

func (span *Span) TraceId() string {
	return span.traceId
}

func (span *Span) SetSpanId(id string) {
	span.spanId = id
}

func (span *Span) SpanId() string {
	return span.spanId
}

func (span *Span) SetParentSpanId(id string) {
	span.parentSpanId = id
}

func (span *Span) ParentSpanId() string {
	return span.parentSpanId
}

func (span *Span) SetOperation(operation string) {
	span.operationName = operation
}

func (span *Span) SetTag(key uint64, value string) {
	span.tags[key] = value
}

func (span *Span) SetRequestInfo(reqInfo types.RequestInfo) {

}

func (span *Span) Tag(key uint64) string {
	return span.tags[key]
}

func (span *Span) FinishSpan() {
	span.EndTime = time.Now()
	for _, name := range activeExporters {
		exporter := GetExporter(name)
		if exporter == nil {
			return
		}
		exporter.ExportSpan(span)
	}
}
func (span *Span) InjectContext(requestHeaders types.HeaderMap, requestInfo types.RequestInfo) {

}

func (span *Span) SpawnChild(operationName string, startTime time.Time) api.Span {
	return nil
}
