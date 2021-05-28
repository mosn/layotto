package health

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockIndicator struct {
}

func (m mockIndicator) Report() Health {
	h := NewHealth(DOWN)
	h.SetDetail("reason", "mock")
	return h
}

type mockScanner struct {
	cnt int
}

func newMockScanner() *mockScanner {
	return &mockScanner{}
}
func (m *mockScanner) Next() string {
	if m.cnt < 1 {
		m.cnt++
		return "readiness"
	}
	return ""
}

func (m *mockScanner) HasNext() bool {
	return m.cnt < 1
}

func TestEndpoint_WhenNoIndicator(t *testing.T) {
	ep := NewEndpoint()
	handle, err := ep.Handle(context.Background(), nil)
	assert.True(t, err != nil)
	assert.True(t, len(handle) == 0)

	AddReadinessIndicator("test", nil)
	AddReadinessIndicatorFunc("test", nil)
	handle, err = ep.Handle(context.Background(), nil)
	assert.True(t, err != nil)
	assert.True(t, len(handle) == 0)

	AddReadinessIndicator("test", mockIndicator{})
	handle, err = ep.Handle(context.Background(), newMockScanner())
	assert.True(t, err != nil)
	assert.True(t, len(handle) == 2)
	assert.True(t, handle["status"] == DOWN)
	health := handle["components"].(map[string]Health)["test"]
	assert.True(t, health.Status == DOWN)
	assert.True(t, health.GetDetail("reason") == "mock")
}
