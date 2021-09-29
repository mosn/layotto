package diagnostics

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
	ltrace "mosn.io/layotto/components/trace"
	"mosn.io/mosn/pkg/trace"
)

// UnaryInterceptorFilter is an implementation of grpc.UnaryServerInterceptor
func UnaryInterceptorFilter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if !trace.IsEnabled() {
		resp, err = handler(ctx, req)
		return resp, err
	}
	tracer := trace.Tracer("layotto")
	span := tracer.Start(ctx, req, time.Now())
	defer span.FinishSpan()
	span.SetTag(ltrace.LAYOTTO_METHOD_NAME, info.FullMethod)
	span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "0")
	ctx = GetNewContext(ctx, span)
	resp, err = handler(ctx, req)
	if err != nil {
		span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "1")
	}
	return
}

func StreamInterceptorFilter(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if !trace.IsEnabled() {
		err := handler(srv, ss)
		return err
	}
	tracer := trace.Tracer("layotto")
	ctx := ss.Context()
	span := tracer.Start(ctx, nil, time.Now())
	defer span.FinishSpan()
	span.SetTag(ltrace.LAYOTTO_METHOD_NAME, info.FullMethod)
	span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "0")
	wrapped := grpc_middleware.WrapServerStream(ss)
	ctx = GetNewContext(ctx, span)
	wrapped.WrappedContext = ctx
	err := handler(srv, wrapped)
	if err != nil {
		span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "1")
	}
	return err
}
