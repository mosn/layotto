package diagnostics

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	"mosn.io/api"
	mtrace "mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/trace/sofa"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/variable"

	"mosn.io/layotto/components/trace"
)

func init() {
	trace.RegisterGenerator("mosntracing", &OpenGenerator{})
}

// OpenGenerator is the default implementation of Generator
type OpenGenerator struct {
	md metadata.MD
}

func (o *OpenGenerator) Init(ctx context.Context) {
	o.md = map[string][]string{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		o.md = md
	}
}
func (o *OpenGenerator) GetTraceId(ctx context.Context) string {
	var traceId string
	if v, ok := o.md[strings.ToLower(sofa.TRACER_ID_KEY)]; ok {
		traceId = v[0]
	} else {
		traceId = mtrace.IdGen().GenerateTraceId()
	}
	return traceId
}

func (o *OpenGenerator) GetSpanId(ctx context.Context) string {
	var spanId string
	if v, ok := o.md[strings.ToLower(sofa.RPC_ID_KEY)]; ok {
		spanId = v[0]
	} else {
		spanId = "0"
	}
	return spanId
}

// GetParentSpanId returns the same id as GetSpanId.
// It's because currently Layotto don't know the parent id.
func (o *OpenGenerator) GetParentSpanId(ctx context.Context) string {
	// TODO: need some design to get the parent id
	var spanId string
	if v, ok := o.md[strings.ToLower(sofa.RPC_ID_KEY)]; ok {
		spanId = v[0]
	} else {
		spanId = "0"
	}
	return spanId
}

func (o *OpenGenerator) GenerateNewContext(ctx context.Context, span api.Span) context.Context {
	newMd := o.md.Copy()
	newMd[strings.ToLower(sofa.TRACER_ID_KEY)] = []string{span.TraceId()}
	newMd[strings.ToLower(sofa.RPC_ID_KEY)] = []string{span.SpanId()}
	if v, ok := o.md[strings.ToLower(sofa.APP_NAME_KEY)]; ok && len(v) > 0 {
		span.SetTag(trace.LAYOTTO_APP_NAME, v[0])
	}
	if v, ok := o.md[strings.ToLower(sofa.SOFA_TRACE_BAGGAGE_DATA)]; ok && len(v) > 0 {
		span.SetTag(trace.LAYOTTO_ATTRS_CONTENT, v[0])
	}
	ctx = metadata.NewIncomingContext(ctx, newMd)
	_ = variable.Set(ctx, types.VariableTraceSpan, span)
	return ctx
}
