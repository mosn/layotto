package trace

import (
	"time"

	"mosn.io/api"

	"mosn.io/mosn/pkg/trace/sofa/xprotocol"
	"mosn.io/mosn/pkg/types"
)

type Span struct {
	StartTime     time.Time
	EndTime       time.Time
	traceId       string
	spanId        string
	parentSpanId  string
	tags          [xprotocol.TRACE_END]string
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
