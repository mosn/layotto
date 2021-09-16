package trace

import (
	"context"

	"mosn.io/mosn/pkg/types"

	mosnctx "mosn.io/mosn/pkg/context"
)

func SetExtraComponentInfo(ctx context.Context, info string) {
	span := mosnctx.Get(ctx, types.ContextKeyActiveSpan).(*Span)
	if span == nil {
		return
	}
	span.SetTag(LAYOTTO_COMPONENT_DETAIL, info)
}
