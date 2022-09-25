// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
