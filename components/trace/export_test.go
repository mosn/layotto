package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockExporter struct {
	runtime int
}

func (m *MockExporter) ExportSpan(s *Span) {
	m.runtime++
}

func TestExport(t *testing.T) {
	active := []string{"mock"}
	m := &MockExporter{}
	span := &Span{}
	SetActiveExporters(active)
	RegisterExporter("mock", m)
	span.FinishSpan()
	assert.Equal(t, m.runtime, 1)

	UnregisterExporter("mock")
	span.FinishSpan()
	assert.Equal(t, m.runtime, 1)
}
