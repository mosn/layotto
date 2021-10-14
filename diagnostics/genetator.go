package diagnostics

import (
	"context"
	"strings"

	"mosn.io/api"

	"mosn.io/mosn/pkg/types"

	"google.golang.org/grpc/metadata"
	"mosn.io/layotto/components/trace"
	mosnctx "mosn.io/mosn/pkg/context"
	mtrace "mosn.io/mosn/pkg/trace"
	"mosn.io/mosn/pkg/trace/sofa"
)

func init() {
	trace.RegisterGenerator("mosntracing", &OpenGenerator{})
}

type OpenGenerator struct {
}

func (o *OpenGenerator) GetTraceId(ctx context.Context) string {
	var traceId string
	md, _ := metadata.FromIncomingContext(ctx)
	if v, ok := md[strings.ToLower(sofa.TRACER_ID_KEY)]; ok {
		traceId = v[0]
	} else {
		traceId = mtrace.IdGen().GenerateTraceId()
	}
	return traceId
}

func (o *OpenGenerator) GetSpanId(ctx context.Context) string {
	var spanId string
	md, _ := metadata.FromIncomingContext(ctx)
	if v, ok := md[strings.ToLower(sofa.RPC_ID_KEY)]; ok {
		spanId = v[0]
	} else {
		spanId = "0"
	}
	return spanId
}

func (o *OpenGenerator) GetParentSpanId(ctx context.Context) string {
	var spanId string
	md, _ := metadata.FromIncomingContext(ctx)
	if v, ok := md[strings.ToLower(sofa.RPC_ID_KEY)]; ok {
		spanId = v[0]
	} else {
		spanId = "0"
	}
	return spanId
}

func (o *OpenGenerator) GenerateNewContext(ctx context.Context, span api.Span) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	newMd := md.Copy()
	newMd[strings.ToLower(sofa.TRACER_ID_KEY)] = []string{span.TraceId()}
	newMd[strings.ToLower(sofa.RPC_ID_KEY)] = []string{span.SpanId()}
	if v, ok := md[strings.ToLower(sofa.APP_NAME)]; ok && len(v) > 0 {
		span.SetTag(trace.LAYOTTO_APP_NAME, v[0])
	}
	if v, ok := md[strings.ToLower(sofa.SOFA_TRACE_BAGGAGE_DATA)]; ok && len(v) > 0 {
		span.SetTag(trace.LAYOTTO_ATTRS_CONTENT, v[0])
	}
	ctx = metadata.NewIncomingContext(ctx, newMd)
	ctx = mosnctx.WithValue(ctx, types.ContextKeyActiveSpan, span)
	return ctx
}
