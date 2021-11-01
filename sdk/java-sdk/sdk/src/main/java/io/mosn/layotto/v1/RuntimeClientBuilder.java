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
 *
 */
package io.mosn.layotto.v1;

import com.google.errorprone.annotations.DoNotCall;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.domain.ApiProtocol;
import io.mosn.layotto.v1.grpc.stub.PooledStubFactory;
import io.mosn.layotto.v1.grpc.stub.SingleStubFactory;
import io.mosn.layotto.v1.grpc.stub.StubCreator;
import io.mosn.layotto.v1.grpc.stub.StubFactory;
import io.mosn.layotto.v1.serializer.JSONSerializer;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.sdk.runtime.v1.client.RuntimeClient;

/**
 * A builder for the RuntimeClient,
 */
public class RuntimeClientBuilder {

    private static final Logger DEFAULT_LOGGER = LoggerFactory.getLogger(RuntimeClient.class.getName());

    private int timeoutMs = RuntimeProperties.DEFAULT_TIMEOUT_MS;

    private String ip = RuntimeProperties.DEFAULT_IP;

    private int port = RuntimeProperties.DEFAULT_PORT;

    private ApiProtocol protocol = RuntimeProperties.DEFAULT_API_PROTOCOL;

    private Logger logger = DEFAULT_LOGGER;

    private ObjectSerializer stateSerializer = new JSONSerializer();

    private boolean useConnectionPool = false;
    private int     poolSize;

    // TODO add rpc serializer

    /**
     * Creates a constructor for RuntimeClient.
     */
    public RuntimeClientBuilder() {
    }

    /**
     * Setter method for property <tt>useConnectionPool</tt> and <tt>poolSize</tt>.
     *
     * @param useConnectionPool value to be assigned to property useConnectionPool
     */
    public void setUseConnectionPool(boolean useConnectionPool, int poolSize) {
        this.useConnectionPool = useConnectionPool;
        if (useConnectionPool && poolSize <= 0) {
            throw new IllegalArgumentException("Invalid poolSize.");
        }
        this.poolSize = poolSize;
    }

    public RuntimeClientBuilder withIp(String ip) {
        if (ip == null || ip.isEmpty()) {
            throw new IllegalArgumentException("Invalid ip.");
        }
        this.ip = ip;
        return this;
    }

    public RuntimeClientBuilder withPort(int port) {
        if (port <= 0) {
            throw new IllegalArgumentException("Invalid port.");
        }
        this.port = port;
        return this;
    }

    @DoNotCall
    public RuntimeClientBuilder withApiProtocol(ApiProtocol protocol) {
        if (protocol == null) {
            throw new IllegalArgumentException("Invalid protocol.");
        }
        this.protocol = protocol;
        return this;
    }

    public RuntimeClientBuilder withTimeout(int timeoutMillisecond) {
        if (timeoutMillisecond <= 0) {
            throw new IllegalArgumentException("Invalid timeout.");
        }
        this.timeoutMs = timeoutMillisecond;
        return this;
    }

    public RuntimeClientBuilder withLogger(Logger logger) {
        if (logger == null) {
            throw new IllegalArgumentException("Invalid logger.");
        }
        this.logger = logger;
        return this;
    }

    /**
     * Sets the serializer for objects to be persisted.
     *
     * @param stateSerializer Serializer for objects to be persisted.
     * @return This instance.
     */
    public RuntimeClientBuilder withStateSerializer(ObjectSerializer stateSerializer) {
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
     * @throws IllegalStateException if any required field is missing
     */
    public RuntimeClient build() {
        if (protocol == null) {
            throw new IllegalStateException("Protocol is required.");
        }

        switch (protocol) {
            case GRPC:
                return buildGrpc();
            default:
                throw new IllegalStateException("Unsupported protocol: " + protocol.name());
        }
    }

    private RuntimeClient buildGrpc() {
        if (port <= 0) {
            throw new IllegalArgumentException("Invalid port.");
        }
        // construct stubFactory
        StubFactory<RuntimeGrpc.RuntimeStub, RuntimeGrpc.RuntimeBlockingStub> stubFactory;
        if (useConnectionPool) {
            stubFactory = new PooledStubFactory<>(ip, port, poolSize, new StubCreatorImpl());
        } else {
            ManagedChannel channel = ManagedChannelBuilder.forAddress(ip, port)
                    .usePlaintext()
                    .build();
            stubFactory = new SingleStubFactory(channel, new StubCreatorImpl());
        }
        return new RuntimeClientGrpc(
                logger,
                timeoutMs,
                stateSerializer,
                stubFactory);
    }

    public static class StubCreatorImpl implements StubCreator<RuntimeGrpc.RuntimeStub, RuntimeGrpc.RuntimeBlockingStub> {

        @Override
        public RuntimeGrpc.RuntimeStub createAsyncStub(ManagedChannel channel) {
            return RuntimeGrpc.newStub(channel);
        }

        @Override
        public RuntimeGrpc.RuntimeBlockingStub createBlockingStub(ManagedChannel channel) {
            return RuntimeGrpc.newBlockingStub(channel);
        }
    }
}
