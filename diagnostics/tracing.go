package diagnostics

import (
	"context"
	"encoding/json"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/api"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/mosn/pkg/trace"
)

const (
	Generator        = "generator"
	Exporter         = "exporter"
	DefaultGenerator = "mosntracing"
)

func init() {
	//register with mosn
	trace.RegisterTracerBuilder("SOFATracer", "layotto", NewTracer)
}

//grpcTracer  is used to start a new Span
type grpcTracer struct {
	config map[string]interface{}
}

func NewTracer(config map[string]interface{}) (api.Tracer, error) {
	v := getActiveExportersFromConfig(config)
	ltrace.SetActiveExporters(v)
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

func NewSpan(ctx context.Context, startTime time.Time, config map[string]interface{}) api.Span {
	span := &ltrace.Span{StartTime: startTime}
	generator := DefaultGenerator
	if v, ok := config[Generator]; ok {
		generator = v.(string)
	}
	ge := ltrace.GetGenerator(generator)
	if ge == nil {
		log.DefaultLogger.Errorf("not support trace type: %+v", generator)
		return nil
	}
	spanId := ge.GetSpanId(ctx)
	traceId := ge.GetTraceId(ctx)
	parentSpanId := ge.GetParentSpanId(ctx)
	span.SetSpanId(spanId)
	span.SetTraceId(traceId)
	span.SetParentSpanId(parentSpanId)
	span.SetTag(ltrace.LAYOTTO_GENERATOR_TYPE, generator)
	return span
}

func GetNewContext(ctx context.Context, span api.Span) context.Context {
	genType := span.Tag(ltrace.LAYOTTO_GENERATOR_TYPE)
	ge := ltrace.GetGenerator(genType)
	//if no implement generator, return old ctx
	if ge == nil {
		return ctx
	}
	newCtx := ge.GenerateNewContext(ctx, span)
	return newCtx
}
