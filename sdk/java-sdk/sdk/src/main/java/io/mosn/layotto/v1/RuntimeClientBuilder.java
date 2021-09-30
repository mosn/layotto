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

    private Supplier<String> ipSupplier = RuntimeProperties.IP;
    private Supplier<Integer> portSupplier = RuntimeProperties.PORT;
    private Supplier<Integer> timeoutMsSupplier = RuntimeProperties.TIMEOUT_MS;

    private Supplier<ApiProtocol> protocolSupplier = RuntimeProperties.API_PROTOCOL;

    private Supplier<Logger> loggerSupplier = () -> LoggerFactory.getLogger(RuntimeClient.class.getName());
    private Supplier<ObjectSerializer> stateSerializerSupplier = JSONSerializer::new;

    /**
     * Creates a constructor for RuntimeClient.
     */
    public RuntimeClientBuilder() {
    }

    public RuntimeClientBuilder withIp(String ip) {
        if (ip == null || ip.isEmpty()) {
            throw new IllegalArgumentException("Invalid ip.");
        }
        this.ipSupplier = () -> ip;
        return this;
    }

    public RuntimeClientBuilder withPort(int port) {
        if (port <= 0) {
            throw new IllegalArgumentException("Invalid port.");
        }
        this.portSupplier = () -> port;
        return this;
    }

    @DoNotCall
    public RuntimeClientBuilder withApiProtocol(ApiProtocol protocol) {
        if (protocol == null) {
            throw new IllegalArgumentException("Invalid protocol.");
        }
        this.protocolSupplier = () -> protocol;
        return this;
    }

    public RuntimeClientBuilder withTimeout(int timeoutMillisecond) {
        if (timeoutMillisecond <= 0) {
            throw new IllegalArgumentException("Invalid timeout.");
        }
        this.timeoutMsSupplier = () -> timeoutMillisecond;
        return this;
    }

    public RuntimeClientBuilder withLogger(Logger logger) {
        if (logger == null) {
            throw new IllegalArgumentException("Invalid logger.");
        }
        this.loggerSupplier = () -> logger;
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

        this.stateSerializerSupplier = () -> stateSerializer;
        return this;
    }

    /**
     * Build an instance of the Client based on the provided setup.
     *
     * @return an instance of the setup Client
     * @throws IllegalStateException if any required field is missing
     */
    public RuntimeClient build() {
        final ApiProtocol apiProtocol = protocolSupplier.get();
        if (apiProtocol == null) {
            throw new IllegalStateException("Protocol is required.");
        }
        switch (apiProtocol) {
            case GRPC:
                return buildGrpc();
            default:
                throw new IllegalStateException("Unsupported protocol: " + apiProtocol.name());
        }
    }

    private RuntimeClient buildGrpc() {
        ManagedChannel channel = ManagedChannelBuilder.forAddress(ipSupplier.get(), portSupplier.get())
                .usePlaintext()
                .build();
        Closeable closeable = () -> {
            if (channel != null && !channel.isShutdown()) {
                channel.shutdown();
            }
        };
        RuntimeGrpc.RuntimeBlockingStub blockingStub = RuntimeGrpc.newBlockingStub(channel);
        return new RuntimeClientGrpc(
                loggerSupplier.get(),
                timeoutMsSupplier.get(),
                stateSerializerSupplier.get(),
                closeable,
                blockingStub);
    }
}
