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
package io.mosn.layotto.v1.config;


import io.mosn.layotto.v1.value.LayottoApiProtocol;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

/**
 * Global properties for Layotto's SDK, using Supplier so they are dynamically resolved.
 */
public class Properties {

    /**
     * Layotto's default IP for gRPC communication.
     */
    private static final String DEFAULT_SIDECAR_IP = "127.0.0.1";

    /**
     * Layotto's default gRPC port.
     */
    private static final Integer DEFAULT_GRPC_PORT = 50001;

    /**
     * Layotto's default use of gRPC.
     */
    private static final LayottoApiProtocol DEFAULT_API_PROTOCOL = LayottoApiProtocol.GRPC;

    /**
     * Layotto's default String encoding: UTF-8.
     */
    private static final Charset DEFAULT_STRING_CHARSET = StandardCharsets.UTF_8;

    /**
     * IP for Layotto's sidecar.
     */
    public static final Property<String> SIDECAR_IP = new StringProperty(
            "layotto.sidecar.ip",
            "LAYOTTO_SIDECAR_IP",
            DEFAULT_SIDECAR_IP);

    /**
     * GRPC port for Layotto after checking system property and environment variable.
     */
    public static final Property<Integer> GRPC_PORT = new IntegerProperty(
            "layotto.grpc.port",
            "LAYOTTO_GRPC_PORT",
            DEFAULT_GRPC_PORT);

    /**
     * Determines if Layotto client will use gRPC to talk to Layotto's sidecar.
     */
    public static final Property<LayottoApiProtocol> API_PROTOCOL = new GenericProperty<>(
            "layotto.api.protocol",
            "LAYOTTO_API_PROTOCOL",
            DEFAULT_API_PROTOCOL,
            (s) -> LayottoApiProtocol.valueOf(s.toUpperCase()));

    /**
     * API token for authentication between App and Layotto's sidecar.
     */
    public static final Property<String> API_TOKEN = new StringProperty(
            "layotto.api.token",
            "LAYOTTO_API_TOKEN",
            null);

    /**
     * Determines which string encoding is used in Layotto's Java SDK.
     */
    public static final Property<Charset> STRING_CHARSET = new GenericProperty<>(
            "layotto.string.charset",
            "LAYOTTO_STRING_CHARSET",
            DEFAULT_STRING_CHARSET,
            (s) -> Charset.forName(s));
}
