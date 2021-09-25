package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	mosnctx "mosn.io/mosn/pkg/context"
	"mosn.io/mosn/pkg/types"
)

func TestSetExtraComponentInfo(t *testing.T) {
	var span Span
	ctx := mosnctx.WithValue(context.TODO(), types.ContextKeyActiveSpan, &span)
	SetExtraComponentInfo(ctx, "hello")
	v := span.Tag(LAYOTTO_COMPONENT_DETAIL)
	assert.Equal(t, v, "hello")
}
