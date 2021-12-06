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
package spec.sdk.reactor.v1.domain.core;

import reactor.core.publisher.Mono;
import spec.sdk.reactor.v1.domain.core.invocation.HttpExtension;
import spec.sdk.reactor.v1.domain.core.invocation.InvokeMethodRequest;
import spec.sdk.reactor.v1.utils.TypeRef;

import java.util.Map;

/**
 * Service-to-Service Invocation Runtimes standard API defined.
 */
public interface InvocationRuntimes {

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param data          The data to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in data.
     * @param type          The Type needed as return for the call.
     * @param <T>           The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(String appId, String methodName, Object data, HttpExtension httpExtension,
                             Map<String, String> metadata, TypeRef<T> type);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param request       The request to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in request.
     * @param clazz         The type needed as return for the call.
     * @param <T>           The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                             Map<String, String> metadata, Class<T> clazz);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param request       The request to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param type          The Type needed as return for the call.
     * @param <T>           The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                             TypeRef<T> type);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param request       The request to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param clazz         The type needed as return for the call.
     * @param <T>           The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                             Class<T> clazz);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in request.
     * @param type          The Type needed as return for the call.
     * @param <T>           The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension,
                             Map<String, String> metadata,
                             TypeRef<T> type);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in request.
     * @param clazz         The type needed as return for the call.
     * @param <T>           The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension,
                             Map<String, String> metadata,
                             Class<T> clazz);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param request       The request to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in request.
     * @return A Mono Plan of type Void.
     */
    Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                            Map<String, String> metadata);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param request       The request to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @return A Mono Plan of type Void.
     */
    Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension);

    /**
     * Invoke a service method, using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in request.
     * @return A Mono Plan of type Void.
     */
    Mono<Void> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata);

    /**
     * Invoke a service method, without using serialization.
     *
     * @param appId         The Application ID where the service is.
     * @param methodName    The actual Method to be call in the application.
     * @param request       The request to be sent to invoke the service, use byte[] to skip serialization.
     * @param httpExtension Additional fields that are needed if the receiving app is listening on
     *                      HTTP, {@link HttpExtension#NONE} otherwise.
     * @param metadata      Metadata (in GRPC) or headers (in HTTP) to be sent in request.
     * @return A Mono Plan of type byte[].
     */
    Mono<byte[]> invokeMethod(String appId, String methodName, byte[] request, HttpExtension httpExtension,
                              Map<String, String> metadata);

    /**
     * Invoke a service method.
     *
     * @param invokeMethodRequest Request object.
     * @param type                The Type needed as return for the call.
     * @param <T>                 The Type of the return, use byte[] to skip serialization.
     * @return A Mono Plan of type T.
     */
    <T> Mono<T> invokeMethod(InvokeMethodRequest invokeMethodRequest, TypeRef<T> type);
}
