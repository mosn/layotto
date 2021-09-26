package io.mosn.layotto.v1.client.reactor;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import spec.sdk.reactor.v1.client.CloudRuntimesClient;
import spec.sdk.reactor.v1.domain.core.configuration.ConfigurationItem;
import spec.sdk.reactor.v1.domain.core.configuration.ConfigurationRequestItem;
import spec.sdk.reactor.v1.domain.core.configuration.SaveConfigurationRequest;
import spec.sdk.reactor.v1.domain.core.configuration.SubConfigurationResp;
import spec.sdk.reactor.v1.domain.core.invocation.HttpExtension;
import spec.sdk.reactor.v1.domain.core.invocation.InvokeMethodRequest;
import spec.sdk.reactor.v1.domain.core.pubsub.PublishEventRequest;
import spec.sdk.reactor.v1.domain.core.state.*;
import spec.sdk.reactor.v1.utils.TypeRef;

import java.util.List;
import java.util.Map;

public interface LayottoReactorClient extends CloudRuntimesClient {

    @Override
    Mono<Void> waitForSidecar(int timeoutInMilliseconds);

    @Override
    Mono<Void> shutdown();

    @Override
    void close() throws Exception;

    @Override
    <T> Mono<List<ConfigurationItem<T>>> getConfiguration(ConfigurationRequestItem configurationRequestItem, TypeRef<T> type);

    @Override
    Mono<Void> saveConfiguration(SaveConfigurationRequest saveConfigurationRequest);

    @Override
    Mono<Void> deleteConfiguration(ConfigurationRequestItem configurationRequestItem);

    @Override
    <T> Flux<SubConfigurationResp<T>> subscribeConfiguration(ConfigurationRequestItem configurationRequestItem, TypeRef<T> type);

    @Override
    <T> Mono<T> invokeMethod(String appId, String methodName, Object data, HttpExtension httpExtension, Map<String, String> metadata, TypeRef<T> type);

    @Override
    <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, Map<String, String> metadata, Class<T> clazz);

    @Override
    <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, TypeRef<T> type);

    @Override
    <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, Class<T> clazz);

    @Override
    <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata, TypeRef<T> type);

    @Override
    <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata, Class<T> clazz);

    @Override
    Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension, Map<String, String> metadata);

    @Override
    Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension);

    @Override
    Mono<Void> invokeMethod(String appId, String methodName, HttpExtension httpExtension, Map<String, String> metadata);

    @Override
    Mono<byte[]> invokeMethod(String appId, String methodName, byte[] request, HttpExtension httpExtension, Map<String, String> metadata);

    @Override
    <T> Mono<T> invokeMethod(InvokeMethodRequest invokeMethodRequest, TypeRef<T> type);

    @Override
    Mono<Void> publishEvent(String pubsubName, String topicName, Object data);

    @Override
    Mono<Void> publishEvent(String pubsubName, String topicName, Object data, Map<String, String> metadata);

    @Override
    Mono<Void> publishEvent(PublishEventRequest request);

    @Override
    <T> Mono<State<T>> getState(String storeName, State<T> state, TypeRef<T> type);

    @Override
    <T> Mono<State<T>> getState(String storeName, State<T> state, Class<T> clazz);

    @Override
    <T> Mono<State<T>> getState(String storeName, String key, TypeRef<T> type);

    @Override
    <T> Mono<State<T>> getState(String storeName, String key, Class<T> clazz);

    @Override
    <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, TypeRef<T> type);

    @Override
    <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, Class<T> clazz);

    @Override
    <T> Mono<State<T>> getState(GetStateRequest request, TypeRef<T> type);

    @Override
    <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, TypeRef<T> type);

    @Override
    <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, Class<T> clazz);

    @Override
    <T> Mono<List<State<T>>> getBulkState(GetBulkStateRequest request, TypeRef<T> type);

    @Override
    Mono<Void> executeStateTransaction(String storeName, List<TransactionalStateOperation<?>> operations);

    @Override
    Mono<Void> executeStateTransaction(ExecuteStateTransactionRequest request);

    @Override
    Mono<Void> saveBulkState(String storeName, List<State<?>> states);

    @Override
    Mono<Void> saveBulkState(SaveStateRequest request);

    @Override
    Mono<Void> saveState(String storeName, String key, Object value);

    @Override
    Mono<Void> saveState(String storeName, String key, String etag, Object value, StateOptions options);

    @Override
    Mono<Void> deleteState(String storeName, String key);

    @Override
    Mono<Void> deleteState(String storeName, String key, String etag, StateOptions options);

    @Override
    Mono<Void> deleteState(DeleteStateRequest request);
}
