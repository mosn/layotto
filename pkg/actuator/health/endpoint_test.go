/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package health

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockIndicator struct {
}

func (m mockIndicator) Report() (Status, map[string]interface{}) {
	h := NewHealth(DOWN)
	h.SetDetail("reason", "mock")
	return h.Status, h.Details
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
	health, ok := handle["components"].(map[string]Health)["test"]
	assert.True(t, ok)
	assert.True(t, health.Status == DOWN)
	assert.True(t, health.GetDetail("reason") == "mock")
}
