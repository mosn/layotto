/*
 * Copyright (c) Microsoft Corporation and Dapr Contributors.
 * Licensed under the MIT License.
 */

package io.mosn.layotto.v1.value;

/**
 * Common headers for GRPC and HTTP communication.
 */
public class Headers {

    /**
     * OpenCensus's metadata for GRPC.
     */
    public static final String GRPC_TRACE_BIN = "grpc-trace-bin";

    /**
     * Token for authentication from Application to Layotto runtime.
     */
    public static final String DAPR_API_TOKEN = "layotto-api-token";
}
