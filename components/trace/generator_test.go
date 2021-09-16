package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/api"
)

type MockGenerator struct {
}

func (m *MockGenerator) GetTraceId(ctx context.Context) string {
	return "mock"
}
func (m *MockGenerator) GetSpanId(ctx context.Context) string {
	return "mock"
}

func (m *MockGenerator) GenerateNewContext(ctx context.Context, span api.Span) context.Context {
	return ctx
}

func (m *MockGenerator) GetParentSpanId(ctx context.Context) string {
	return "mock"
}

func TestGenerator(t *testing.T) {
	RegisterGenerator("mock", &MockGenerator{})
	ge := GetGenerator("mock")
	assert.Equal(t, ge.GetSpanId(context.TODO()), "mock")
	assert.Equal(t, ge.GetTraceId(context.TODO()), "mock")
	assert.Equal(t, ge.GetParentSpanId(context.TODO()), "mock")
}
