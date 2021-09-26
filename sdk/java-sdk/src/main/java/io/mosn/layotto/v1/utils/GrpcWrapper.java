/*
 * Copyright (c) Microsoft Corporation and Layotto Contributors.
 * Licensed under the MIT License.
 */

package io.mosn.layotto.v1.utils;

import io.grpc.*;
import io.mosn.layotto.v1.config.Property;
import reactor.util.context.Context;
import spec.proto.runtime.v1.RuntimeGrpc;

import java.util.logging.Logger;

/**
 * Wraps a Layotto gRPC stub with telemetry interceptor.
 */
public final class GrpcWrapper {

    private static final Logger LOGGER = Logger.getLogger(Property.class.getName());

    private static final Metadata.Key<byte[]> GRPC_TRACE_BIN_KEY =
            Metadata.Key.of("grpc-trace-bin", Metadata.BINARY_BYTE_MARSHALLER);

    private static final Metadata.Key<String> TRACEPARENT_KEY =
            Metadata.Key.of("traceparent", Metadata.ASCII_STRING_MARSHALLER);

    private static final Metadata.Key<String> TRACESTATE_KEY =
            Metadata.Key.of("tracestate", Metadata.ASCII_STRING_MARSHALLER);

    private GrpcWrapper() {
    }

    /**
     * Populates GRPC client with interceptors.
     *
     * @param context Reactor's context.
     * @param client  GRPC client for Layotto.
     * @return Client after adding interceptors.
     */
    public static RuntimeGrpc.RuntimeStub intercept(final Context context, RuntimeGrpc.RuntimeStub client) {
        ClientInterceptor interceptor = new ClientInterceptor() {
            @Override
            public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(MethodDescriptor<ReqT, RespT> methodDescriptor,
                                                                       CallOptions callOptions,
                                                                       Channel channel) {
                ClientCall<ReqT, RespT> clientCall = channel.newCall(methodDescriptor, callOptions);
                return new ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(clientCall) {
                    @Override
                    public void start(final Listener<RespT> responseListener, final Metadata metadata) {
                        // FIXME: 2021/9/26 Refer to Dapr
                        super.start(responseListener, metadata);
                    }
                };
            }
        };
        return client.withInterceptors(interceptor);
    }
}