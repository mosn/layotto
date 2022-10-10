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
package actuators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeAllIndicators(t *testing.T) {
	readinessIndicator := NewHealthIndicator()
	livenessIndicator := NewHealthIndicator()
	indicators := &ComponentsIndicator{ReadinessIndicator: readinessIndicator, LivenessIndicator: livenessIndicator}
	SetComponentsIndicator("test", indicators)
	idc := GetIndicatorWithName("test")
	assert.Equal(t, indicators, idc)
	SetComponentsIndicator("test2", indicators)
	cnt := 0
	RangeAllIndicators(func(key string, value *ComponentsIndicator) bool {
		cnt++
		assert.Equal(t, indicators, value)
		return true
	})
	assert.Equal(t, cnt, 2)
}
