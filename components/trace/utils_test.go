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

	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/variable"
)

func TestSetExtraComponentInfo(t *testing.T) {
	var span Span
	ctx := variable.NewVariableContext(context.TODO())
	_ = variable.Set(ctx, types.VariableTraceSpan, &span)
	SetExtraComponentInfo(ctx, "hello")
	v := span.Tag(LAYOTTO_COMPONENT_DETAIL)
	assert.Equal(t, v, "hello")
}

func TestSetterAndGetter(t *testing.T) {
	var span Span
	// ParentSpanId
	span.SetParentSpanId("par")
	v := span.ParentSpanId()
	assert.Equal(t, v, "par")
	// traceId
	span.SetTraceId("traceId")
	v = span.TraceId()
	assert.Equal(t, v, "traceId")
	// span id
	span.SetSpanId("spanId")
	v = span.SpanId()
	assert.Equal(t, v, "spanId")
}
