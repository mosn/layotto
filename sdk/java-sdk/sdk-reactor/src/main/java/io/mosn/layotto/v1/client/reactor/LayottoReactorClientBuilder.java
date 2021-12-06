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
package io.mosn.layotto.v1.client.reactor;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.mosn.layotto.v1.config.Properties;
import io.mosn.layotto.v1.serializer.DefaultObjectSerializer;
import io.mosn.layotto.v1.serializer.LayottoObjectSerializer;
import io.mosn.layotto.v1.value.LayottoApiProtocol;
import spec.proto.runtime.v1.RuntimeGrpc;

import java.io.Closeable;

/**
 * A builder for the LayottoClient, Currently only gRPC Client will be supported.
 */
public class LayottoReactorClientBuilder {

    /**
     * Determine if this builder will create GRPC clients instead of HTTP clients.
     */
    private final LayottoApiProtocol apiProtocol;

    /**
     * Serializer used for request and response objects in LayottoClient.
     */
    private LayottoObjectSerializer  objectSerializer;

    /**
     * Serializer used for state objects in LayottoClient.
     */
    private LayottoObjectSerializer  stateSerializer;

    /**
     * Creates a constructor for LayottoClient.
     * <p>
     * {@link DefaultObjectSerializer} is used for object and state serializers by default but is not recommended
     * for production scenarios.
     */
    public LayottoReactorClientBuilder() {
        this.objectSerializer = new DefaultObjectSerializer();
        this.stateSerializer = new DefaultObjectSerializer();
        this.apiProtocol = Properties.API_PROTOCOL.get();
    }

    /**
     * Sets the serializer for objects to be sent and received from Layotto.
     * See {@link DefaultObjectSerializer} as possible serializer for non-production scenarios.
     *
     * @param objectSerializer Serializer for objects to be sent and received from Layotto.
     * @return This instance.
     */
    public LayottoReactorClientBuilder withObjectSerializer(LayottoObjectSerializer objectSerializer) {
        if (objectSerializer == null) {
            throw new IllegalArgumentException("Object serializer is required");
        }

        if (objectSerializer.getContentType() == null || objectSerializer.getContentType().isEmpty()) {
            throw new IllegalArgumentException("Content Type should not be null or empty");
        }

        this.objectSerializer = objectSerializer;
        return this;
    }

    /**
     * Sets the serializer for objects to be persisted.
     * See {@link DefaultObjectSerializer} as possible serializer for non-production scenarios.
     *
     * @param stateSerializer Serializer for objects to be persisted.
     * @return This instance.
     */
    public LayottoReactorClientBuilder withStateSerializer(LayottoObjectSerializer stateSerializer) {
        if (stateSerializer == null) {
            throw new IllegalArgumentException("State serializer is required");
        }

        this.stateSerializer = stateSerializer;
        return this;
    }

    /**
     * Build an instance of the Client based on the provided setup.
     *
     * @return an instance of the setup Client
     * @throws java.lang.IllegalStateException if any required field is missing
     */
    public LayottoReactorClient build() {
        return buildLayottoClient(this.apiProtocol);
    }

    /**
     * Creates an instance of a Layotto Client based on the chosen protocol.
     *
     * @param protocol Layotto API's protocol.
     * @return the GRPC Client.
     * @throws java.lang.IllegalStateException if either host is missing or if port is missing or a negative number.
     */
    private LayottoReactorClient buildLayottoClient(LayottoApiProtocol protocol) {
        if (protocol == null) {
            throw new IllegalStateException("Protocol is required.");
        }

        switch (protocol) {
            case GRPC:
                return buildLayottoClientGrpc();
            default:
                throw new IllegalStateException("Unsupported protocol: " + protocol.name());
        }
    }

    /**
     * Creates an instance of the GPRC Client.
     *
     * @return the GRPC Client.
     * @throws java.lang.IllegalStateException if either host is missing or if port is missing or a negative number.
     */
    private LayottoReactorClient buildLayottoClientGrpc() {
        int port = Properties.GRPC_PORT.get();
        if (port <= 0) {
            throw new IllegalArgumentException("Invalid port.");
        }
        ManagedChannel channel = ManagedChannelBuilder
                .forAddress(Properties.SIDECAR_IP.get(), port)
                .usePlaintext()
                .build();
        Closeable closeableChannel = () -> {
            if (channel != null && !channel.isShutdown()) {
                channel.shutdown();
            }
        };
        RuntimeGrpc.RuntimeStub asyncStub = RuntimeGrpc.newStub(channel);
        return new LayottoReactorClientGrpc(this.objectSerializer, this.stateSerializer, closeableChannel, asyncStub);
    }
}
