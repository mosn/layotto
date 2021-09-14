package trace

import (
	"context"
	"sync"
)

var (
	generators sync.Map
)

type Generator interface {
	GetTraceId(ctx context.Context) string
	GetSpanId(ctx context.Context) string
	GenerateNewContext(ctx context.Context, span *Span) context.Context
	GetParentSpanId(ctx context.Context) string
}

func RegisterGenerator(name string, ge Generator) {
	generators.Store(name, ge)
}

func GetGenerator(name string) Generator {
	g, ok := generators.Load(name)
	if ok {
		return g.(Generator)
	}
	return nil
}
