package diagnostics

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
	"mosn.io/mosn/pkg/trace"

	ltrace "mosn.io/layotto/components/trace"
	lgrpc "mosn.io/layotto/diagnostics/grpc"
	"mosn.io/layotto/diagnostics/protocol"
)

// UnaryInterceptorFilter is an implementation of grpc.UnaryServerInterceptor
func UnaryInterceptorFilter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if !trace.IsEnabled() {
		resp, err = handler(ctx, req)
		return resp, err
	}
	// get tracer
	tracer := trace.Tracer(protocol.Layotto)
	// start a span
	span := tracer.Start(ctx, &lgrpc.RequestInfo{
		FullMethod: info.FullMethod,
	}, time.Now())
	defer span.FinishSpan()

	span.SetTag(ltrace.LAYOTTO_METHOD_NAME, info.FullMethod)
	span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "0")

	// construct a new context which contains the span
	ctx = GetNewContext(ctx, span)

	// handle request
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
	// get tracer
	tracer := trace.Tracer(protocol.Layotto)
	ctx := ss.Context()
	// start a span
	span := tracer.Start(ctx, &lgrpc.RequestInfo{
		FullMethod: info.FullMethod,
	}, time.Now())
	defer span.FinishSpan()

	span.SetTag(ltrace.LAYOTTO_METHOD_NAME, info.FullMethod)
	span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "0")

	// construct a new context which contains the span
	wrapped := grpc_middleware.WrapServerStream(ss)
	ctx = GetNewContext(ctx, span)
	wrapped.WrappedContext = ctx
	// handle request
	err := handler(srv, wrapped)
	if err != nil {
		span.SetTag(ltrace.LAYOTTO_REQUEST_RESULT, "1")
	}

	return err
}
