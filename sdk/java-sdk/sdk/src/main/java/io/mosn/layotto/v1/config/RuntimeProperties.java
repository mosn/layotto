package io.mosn.layotto.v1.config;

import io.mosn.layotto.v1.domain.ApiProtocol;

import java.util.function.Supplier;

public class RuntimeProperties {

    /**
     * Layotto Runtimes default use of GRPC.
     */
    public static final ApiProtocol DEFAULT_API_PROTOCOL = ApiProtocol.GRPC;

    /**
     * Layotto Runtimes default start on "127.0.0.1".
     */
    public static final String DEFAULT_IP = "127.0.0.1";

    /**
     * Layotto Runtimes default start on 34904 port.
     */
    public static final Integer DEFAULT_PORT = 34904;

    /**
     * Layotto Runtimes default timeout in 1000ms for GRPC client reads.
     */
    public static final Integer DEFAULT_TIMEOUT_MS = 1000;

    /**
     * Layotto Runtimes default pubsub content type.
     */
    public static final String DEFAULT_PUBSUB_CONTENT_TYPE = "";

}
