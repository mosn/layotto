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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/api"
)

type MockGenerator struct {
}

func (m *MockGenerator) Init(ctx context.Context) {}
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
