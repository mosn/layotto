package diagnostics

import (
	"context"
	"encoding/json"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/api"
	trace2 "mosn.io/layotto/components/trace"
	"mosn.io/mosn/pkg/trace"
)

const (
	Generator = "generator"
	Exporter  = "exporter"
)

func init() {
	trace.RegisterTracerBuilder("SOFATracer", "layotto", NewTracer)
}

type grpcTracer struct {
	config map[string]interface{}
}

func NewTracer(config map[string]interface{}) (api.Tracer, error) {
	v := getActiveExportersFromConfig(config)
	trace2.SetActiveExporters(v)
	return &grpcTracer{config: config}, nil
}

func getActiveExportersFromConfig(config map[string]interface{}) []string {
	var exporters []string
	if v, ok := config[Exporter]; ok {
		data, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(data, &exporters)
		if err != nil {
			return nil
		}
	}
	return exporters
}

func (tracer *grpcTracer) Start(ctx context.Context, request interface{}, startTime time.Time) api.Span {
	span := NewSpan(ctx, startTime, tracer.config)
	return span
}

func NewSpan(ctx context.Context, startTime time.Time, config map[string]interface{}) *trace2.Span {
	span := &trace2.Span{StartTime: startTime}
	generator := "mosntracing"
	if v, ok := config[Generator]; ok {
		generator = v.(string)
	}
	g := trace2.GetGenerator(generator)
	if g == nil {
		log.DefaultLogger.Errorf("not support trace type: %+v", generator)
		return nil
	}
	spanId := g.GetSpanId(ctx)
	traceId := g.GetTraceId(ctx)
	parentSpanId := g.GetParentSpanId(ctx)
	span.SetSpanId(spanId)
	span.SetTraceId(traceId)
	span.SetParentSpanId(parentSpanId)
	span.Ctx = g.GenerateNewContext(ctx, span)
	return span
}
