package io.mosn.layotto.v1.config;

import io.mosn.layotto.v1.domain.ApiProtocol;

import java.util.function.Supplier;

public class RuntimeProperties {

    /**
     * Layotto Runtimes default use of GRPC.
     */
    private static final ApiProtocol DEFAULT_API_PROTOCOL = ApiProtocol.GRPC;

    public static final Supplier<ApiProtocol> API_PROTOCOL = () -> DEFAULT_API_PROTOCOL;

    /**
     * Layotto Runtimes default start on "127.0.0.1".
     */
    private static final String DEFAULT_IP = "127.0.0.1";

    public static final Supplier<String> IP = () -> DEFAULT_IP;

    /**
     * Layotto Runtimes default start on 34904 port.
     */
    private static final Integer DEFAULT_PORT = 34904;

    public static final Supplier<Integer> PORT = () -> DEFAULT_PORT;

    /**
     * Layotto Runtimes default timeout in 1000ms for GRPC client reads.
     */
    private static final Integer DEFAULT_TIMEOUT_MS = 1000;

    public static final Supplier<Integer> TIMEOUT_MS = () -> DEFAULT_TIMEOUT_MS;

    /**
     * Layotto Runtimes default pubsub content type.
     */
    private static final String DEFAULT_PUBSUB_CONTENT_TYPE = "";

    public static final Supplier<String> PUBSUB_CONTENT_TYPE = () -> DEFAULT_PUBSUB_CONTENT_TYPE;
}
