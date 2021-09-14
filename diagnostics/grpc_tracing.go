package diagnostics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/mosn/pkg/trace"
)

// UnaryInterceptorFilter is an implementation of grpc.UnaryServerInterceptor
func UnaryInterceptorFilter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tracer := trace.Tracer("layotto")
	span := tracer.Start(ctx, req, time.Now())
	defer span.FinishSpan()
	span.SetTag(ltrace.LAYOTTO_METHOD_NAME, info.FullMethod)
	span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "0")
	resp, err = handler(ctx, req)
	if err != nil {
		span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "1")
	}
	return
}
