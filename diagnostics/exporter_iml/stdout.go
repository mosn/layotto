package exporter_iml

import (
	"strconv"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/trace"
)

func init() {
	trace.RegisterExporter("stdout", &StdoutExporter{})
}

// StdoutExporter is the implementation of Exporter, export span information to log
type StdoutExporter struct{}

var _ trace.Exporter = &StdoutExporter{}

const msg = "%s, AppName: %+s, Method: %s, TraceId: %s, SpanId: %s, ParentSpanId:%s, Time: [%s ->  %s], processTime: %+v, result: %+v, extraInfo:[ %+v ]"

// ExportSpan implements the open census exporter interface.
func (e *StdoutExporter) ExportSpan(sd *trace.Span) {
	processingTime := strconv.FormatInt(sd.EndTime.Sub(sd.StartTime).Nanoseconds()/1000000, 10)
	log.DefaultLogger.Infof(msg, time.Now().Format("2006-01-02 15:04:05.999"), sd.Tag(trace.LAYOTTO_APP_NAME), sd.Tag(trace.LAYOTTO_METHOD_NAME),
		sd.TraceId(), sd.SpanId(), sd.ParentSpanId(), sd.StartTime, sd.EndTime, processingTime, sd.Tag(trace.LAYOTTO_REQUEST_RESULT), sd.Tag(trace.LAYOTTO_COMPONENT_DETAIL))
}
