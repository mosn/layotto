package spec.sdk.runtime.v1.domain.core;

import spec.sdk.runtime.v1.domain.core.binding.InvokeBindingRequest;
import spec.sdk.runtime.v1.utils.TypeRef;
import reactor.core.publisher.Mono;

import java.util.Map;

/**
 * Resource Bindings Runtimes standard API defined.
 */
public interface BindingRuntimes {

    /**
     * Invokes a Binding operation.
     *
     * @param bindingName The bindingName of the biding to call.
     * @param operation   The operation to be performed by the binding request processor.
     * @param data        The data to be processed, use byte[] to skip serialization.
     * @return an empty Mono.
     */
    Mono<Void> invokeBinding(String bindingName, String operation, Object data);

    /**
     * Invokes a Binding operation, skipping serialization.
     *
     * @param bindingName The name of the biding to call.
     * @param operation   The operation to be performed by the binding request processor.
     * @param data        The data to be processed, skipping serialization.
     * @param metadata    The metadata map.
     * @return a Mono plan of type byte[].
     */
    Mono<byte[]> invokeBinding(String bindingName, String operation, byte[] data, Map<String, String> metadata);

    /**
     * Invokes a Binding operation.
     *
     * @param bindingName The name of the biding to call.
     * @param operation   The operation to be performed by the binding request processor.
     * @param data        The data to be processed, use byte[] to skip serialization.
     * @param type        The type being returned.
     * @param <T>         The type of the return
     * @return a Mono plan of type T.
     */
    <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, TypeRef<T> type);

    /**
     * Invokes a Binding operation.
     *
     * @param bindingName The name of the biding to call.
     * @param operation   The operation to be performed by the binding request processor.
     * @param data        The data to be processed, use byte[] to skip serialization.
     * @param clazz       The type being returned.
     * @param <T>         The type of the return
     * @return a Mono plan of type T.
     */
    <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, Class<T> clazz);

    /**
     * Invokes a Binding operation.
     *
     * @param bindingName The name of the biding to call.
     * @param operation   The operation to be performed by the binding request processor.
     * @param data        The data to be processed, use byte[] to skip serialization.
     * @param metadata    The metadata map.
     * @param type        The type being returned.
     * @param <T>         The type of the return
     * @return a Mono plan of type T.
     */
    <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, Map<String, String> metadata,
                              TypeRef<T> type);

    /**
     * Invokes a Binding operation.
     *
     * @param bindingName The name of the biding to call.
     * @param operation   The operation to be performed by the binding request processor.
     * @param data        The data to be processed, use byte[] to skip serialization.
     * @param metadata    The metadata map.
     * @param clazz       The type being returned.
     * @param <T>         The type of the return
     * @return a Mono plan of type T.
     */
    <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, Map<String, String> metadata,
                              Class<T> clazz);

    /**
     * Invokes a Binding operation.
     *
     * @param request The binding invocation request.
     * @param type    The type being returned.
     * @param <T>     The type of the return
     * @return a Mono plan of type T.
     */
    <T> Mono<T> invokeBinding(InvokeBindingRequest request, TypeRef<T> type);
}
