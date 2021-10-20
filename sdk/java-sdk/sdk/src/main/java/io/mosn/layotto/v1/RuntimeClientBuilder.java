package io.mosn.layotto.v1;

import com.google.errorprone.annotations.DoNotCall;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.domain.ApiProtocol;
import io.mosn.layotto.v1.serializer.JSONSerializer;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.sdk.runtime.v1.client.RuntimeClient;

import java.io.Closeable;
import java.util.function.Supplier;

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

    // TODO add rpc serializer

    /**
     * Creates a constructor for RuntimeClient.
     */
    public RuntimeClientBuilder() {
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
        ManagedChannel channel = ManagedChannelBuilder.forAddress(ip, port)
                .usePlaintext()
                .build();
        Closeable closeable = () -> {
            if (channel != null && !channel.isShutdown()) {
                channel.shutdown();
            }
        };
        RuntimeGrpc.RuntimeBlockingStub blockingStub = RuntimeGrpc.newBlockingStub(channel);
        return new RuntimeClientGrpc(
                logger,
                timeoutMs,
                stateSerializer,
                closeable,
                blockingStub);
    }
}
