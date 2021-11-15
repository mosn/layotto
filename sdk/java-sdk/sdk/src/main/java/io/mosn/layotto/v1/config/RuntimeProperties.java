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
    public static final Integer DEFAULT_TIMEOUT_MS = 3000;

    /**
     * Layotto Runtimes default pubsub content type.
     */
    public static final String DEFAULT_PUBSUB_CONTENT_TYPE = "";

}
