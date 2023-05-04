package diagnostics

import (
	"context"
	"encoding/json"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/api"

	ltrace "mosn.io/layotto/components/trace"
)

const (
	generatorConfigKey = "generator"
	exporterConfigKey  = "exporter"
	defaultGenerator   = "mosntracing"
)

// grpcTracer  is used to start a new Span
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
	if v, ok := config[exporterConfigKey]; ok {
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

// NewSpan constructs a span and tag it with span/trace/parentSpan IDs.
// These IDs are generated using the Generator
func NewSpan(ctx context.Context, startTime time.Time, config map[string]interface{}) api.Span {
	// construct span
	span := &ltrace.Span{StartTime: startTime}
	// get generator according to configuration
	generatorName := defaultGenerator
	if v, ok := config[generatorConfigKey]; ok {
		generatorName = v.(string)
	}
	ge := ltrace.GetGenerator(generatorName)
	if ge == nil {
		log.DefaultLogger.Errorf("not support trace type: %+v", generatorName)
		return nil
	}
	ge.Init(ctx)
	// use generator to extract the span/trace/parentSpan IDs
	spanId := ge.GetSpanId(ctx)
	traceId := ge.GetTraceId(ctx)
	parentSpanId := ge.GetParentSpanId(ctx)
	span.SetSpanId(spanId)
	span.SetTraceId(traceId)
	span.SetParentSpanId(parentSpanId)
	// tagging generator type
	span.SetTag(ltrace.LAYOTTO_GENERATOR_TYPE, generatorName)
	return span
}

func GetNewContext(ctx context.Context, span api.Span) context.Context {
	genType := span.Tag(ltrace.LAYOTTO_GENERATOR_TYPE)
	ge := ltrace.GetGenerator(genType)
	//if no implement generator, return old ctx
	if ge == nil {
		return ctx
	}
	ge.Init(ctx)
	newCtx := ge.GenerateNewContext(ctx, span)
	return newCtx
}
