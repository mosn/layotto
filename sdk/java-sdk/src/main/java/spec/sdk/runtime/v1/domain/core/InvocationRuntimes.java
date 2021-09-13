package spec.sdk.runtime.v1.domain.core;

import spec.sdk.runtime.v1.domain.core.invocation.HttpExtension;
import spec.sdk.runtime.v1.domain.core.invocation.InvokeMethodRequest;
import spec.sdk.runtime.v1.utils.TypeRef;
import reactor.core.publisher.Mono;

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
    <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata,
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
    <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata,
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
