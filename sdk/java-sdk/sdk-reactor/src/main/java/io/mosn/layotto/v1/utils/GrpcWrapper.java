/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package io.mosn.layotto.v1.utils;

import io.grpc.CallOptions;
import io.grpc.Channel;
import io.grpc.ClientCall;
import io.grpc.ClientInterceptor;
import io.grpc.ForwardingClientCall;
import io.grpc.Metadata;
import io.grpc.MethodDescriptor;
import io.mosn.layotto.v1.config.Property;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import reactor.util.context.Context;
import spec.proto.runtime.v1.RuntimeGrpc;

/**
 * Wraps a Layotto gRPC stub with telemetry interceptor.
 */
public final class GrpcWrapper {

    private static final Logger               LOGGER             = LoggerFactory.getLogger(Property.class.getName());

    private static final Metadata.Key<byte[]> GRPC_TRACE_BIN_KEY =
                                                                         Metadata.Key.of("grpc-trace-bin",
                                                                             Metadata.BINARY_BYTE_MARSHALLER);

    private static final Metadata.Key<String> TRACEPARENT_KEY    =
                                                                         Metadata.Key.of("traceparent",
                                                                             Metadata.ASCII_STRING_MARSHALLER);

    private static final Metadata.Key<String> TRACESTATE_KEY     =
                                                                         Metadata.Key.of("tracestate",
                                                                             Metadata.ASCII_STRING_MARSHALLER);

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