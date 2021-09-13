package spec.sdk.runtime.v1.client;

import spec.sdk.runtime.v1.domain.core.binding.InvokeBindingRequest;
import spec.sdk.runtime.v1.domain.core.configuration.ConfigurationItem;
import spec.sdk.runtime.v1.domain.core.configuration.ConfigurationRequestItem;
import spec.sdk.runtime.v1.domain.core.configuration.SaveConfigurationRequest;
import spec.sdk.runtime.v1.domain.core.configuration.SubConfigurationResp;
import spec.sdk.runtime.v1.domain.core.invocation.HttpExtension;
import spec.sdk.runtime.v1.domain.core.invocation.InvokeMethodRequest;
import spec.sdk.runtime.v1.domain.core.pubsub.PublishEventRequest;
import spec.sdk.runtime.v1.domain.core.secrets.GetBulkSecretRequest;
import spec.sdk.runtime.v1.domain.core.secrets.GetSecretRequest;
import spec.sdk.runtime.v1.domain.core.state.*;
import spec.sdk.runtime.v1.utils.TypeRef;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.Map;

/**
 * An exception is thrown by default, and some functions can be implemented on demand.
 */
public interface DefaultCloudRuntimesClient extends CloudRuntimesClient {

    @Override
    default Mono<Void> waitForSidecar(int timeoutInMilliseconds) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> invokeBinding(String bindingName, String operation, Object data) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<byte[]> invokeBinding(String bindingName, String operation, byte[] data, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, Map<String, String> metadata, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeBinding(String bindingName, String operation, Object data, Map<String, String> metadata, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeBinding(InvokeBindingRequest request, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(String appId, String methodName, Object data, HttpExtension httpExtension, Map<String, String> metadata, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, Map<String, String> metadata, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<byte[]> invokeMethod(String appId, String methodName, byte[] request, HttpExtension httpExtension, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<T> invokeMethod(InvokeMethodRequest invokeMethodRequest, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> publishEvent(String pubsubName, String topicName, Object data) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> publishEvent(String pubsubName, String topicName, Object data, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> publishEvent(PublishEventRequest request) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Map<String, String>> getSecret(String storeName, String secretName, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Map<String, String>> getSecret(String storeName, String secretName) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Map<String, String>> getSecret(GetSecretRequest request) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Map<String, Map<String, String>>> getBulkSecret(String storeName) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Map<String, Map<String, String>>> getBulkSecret(String storeName, Map<String, String> metadata) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Map<String, Map<String, String>>> getBulkSecret(GetBulkSecretRequest request) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(String storeName, State<T> state, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(String storeName, State<T> state, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(String storeName, String key, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(String storeName, String key, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<State<T>> getState(GetStateRequest request, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, Class<T> clazz) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<List<State<T>>> getBulkState(GetBulkStateRequest request, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> executeStateTransaction(String storeName, List<TransactionalStateOperation<?>> operations) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> executeStateTransaction(ExecuteStateTransactionRequest request) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> saveBulkState(String storeName, List<State<?>> states) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> saveBulkState(SaveStateRequest request) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> saveState(String storeName, String key, Object value) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> saveState(String storeName, String key, String etag, Object value, StateOptions options) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> deleteState(String storeName, String key) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> deleteState(String storeName, String key, String etag, StateOptions options) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> deleteState(DeleteStateRequest request) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Mono<List<ConfigurationItem<T>>> getConfiguration(ConfigurationRequestItem configurationRequestItem, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> saveConfiguration(SaveConfigurationRequest saveConfigurationRequest) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> deleteConfiguration(ConfigurationRequestItem configurationRequestItem) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default <T> Flux<SubConfigurationResp<T>> subscribeConfiguration(ConfigurationRequestItem configurationRequestItem, TypeRef<T> type) {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    default Mono<Void> shutdown() {
        throw new UnsupportedOperationException("CloudRuntimes Operate Unsupported.");
    }

    @Override
    void close();
}
