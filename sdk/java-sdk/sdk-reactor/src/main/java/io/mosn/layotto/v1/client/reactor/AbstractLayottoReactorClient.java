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
package io.mosn.layotto.v1.client.reactor;

import io.mosn.layotto.v1.serializer.LayottoObjectSerializer;
import reactor.core.publisher.Mono;
import spec.sdk.reactor.v1.domain.core.invocation.HttpExtension;
import spec.sdk.reactor.v1.domain.core.invocation.InvokeMethodRequest;
import spec.sdk.reactor.v1.domain.core.pubsub.PublishEventRequest;
import spec.sdk.reactor.v1.domain.core.state.*;
import spec.sdk.reactor.v1.utils.TypeRef;

import java.util.Collections;
import java.util.List;
import java.util.Map;

abstract class AbstractLayottoReactorClient implements LayottoReactorClient {

    /**
     * A utility class for serialize and deserialize the transient objects.
     */
    protected LayottoObjectSerializer objectSerializer;

    /**
     * A utility class for serialize and deserialize state objects.
     */
    protected LayottoObjectSerializer stateSerializer;

    /**
     * Common constructor for implementations of this class.
     *
     * @param objectSerializer Serializer for transient request/response objects.
     * @param stateSerializer  Serializer for state objects.
     */
    AbstractLayottoReactorClient(LayottoObjectSerializer objectSerializer,
                                 LayottoObjectSerializer stateSerializer) {
        this.objectSerializer = objectSerializer;
        this.stateSerializer = stateSerializer;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> publishEvent(String pubsubName, String topicName, Object data) {
        return this.publishEvent(pubsubName, topicName, data, null);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> publishEvent(String pubsubName, String topicName, Object data, Map<String, String> metadata) {
        PublishEventRequest req = new PublishEventRequest(pubsubName, topicName, data)
            .setMetadata(metadata);
        return this.publishEvent(req).then();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<T> invokeMethod(String appId,
                                    String methodName,
                                    Object data,
                                    HttpExtension httpExtension,
                                    Map<String, String> metadata,
                                    TypeRef<T> type) {
        InvokeMethodRequest req = new InvokeMethodRequest(appId, methodName)
            .setBody(data)
            .setHttpExtension(httpExtension)
            .setContentType(objectSerializer.getContentType());
        return this.invokeMethod(req, type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<T> invokeMethod(String appId,
                                    String methodName,
                                    Object request,
                                    HttpExtension httpExtension,
                                    Map<String, String> metadata,
                                    Class<T> clazz) {
        return this.invokeMethod(appId, methodName, request, httpExtension, metadata, TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension,
                                    Map<String, String> metadata, TypeRef<T> type) {
        return this.invokeMethod(appId, methodName, null, httpExtension, metadata, type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<T> invokeMethod(String appId, String methodName, HttpExtension httpExtension,
                                    Map<String, String> metadata, Class<T> clazz) {
        return this.invokeMethod(appId, methodName, null, httpExtension, metadata, TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                                    TypeRef<T> type) {
        return this.invokeMethod(appId, methodName, request, httpExtension, null, type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<T> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                                    Class<T> clazz) {
        return this.invokeMethod(appId, methodName, request, httpExtension, null, TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension) {
        return this.invokeMethod(appId, methodName, request, httpExtension, null, TypeRef.BYTE_ARRAY).then();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> invokeMethod(String appId, String methodName, Object request, HttpExtension httpExtension,
                                   Map<String, String> metadata) {
        return this.invokeMethod(appId, methodName, request, httpExtension, metadata, TypeRef.BYTE_ARRAY).then();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> invokeMethod(String appId, String methodName, HttpExtension httpExtension,
                                   Map<String, String> metadata) {
        return this.invokeMethod(appId, methodName, null, httpExtension, metadata, TypeRef.BYTE_ARRAY).then();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<byte[]> invokeMethod(String appId, String methodName, byte[] request, HttpExtension httpExtension,
                                     Map<String, String> metadata) {
        return this.invokeMethod(appId, methodName, request, httpExtension, metadata, TypeRef.BYTE_ARRAY);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<State<T>> getState(String storeName, State<T> state, TypeRef<T> type) {
        return this.getState(storeName, state.getKey(), state.getOptions(), type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<State<T>> getState(String storeName, State<T> state, Class<T> clazz) {
        return this.getState(storeName, state.getKey(), state.getOptions(), TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<State<T>> getState(String storeName, String key, TypeRef<T> type) {
        return this.getState(storeName, key, null, type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<State<T>> getState(String storeName, String key, Class<T> clazz) {
        return this.getState(storeName, key, null, TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, TypeRef<T> type) {
        GetStateRequest request = new GetStateRequest(storeName, key)
            .setStateOptions(options);
        return this.getState(request, type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, Class<T> clazz) {
        return this.getState(storeName, key, options, TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, TypeRef<T> type) {
        return this.getBulkState(new GetBulkStateRequest(storeName, keys), type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, Class<T> clazz) {
        return this.getBulkState(storeName, keys, TypeRef.get(clazz));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> executeStateTransaction(String storeName, List<TransactionalStateOperation<?>> operations) {
        ExecuteStateTransactionRequest request = new ExecuteStateTransactionRequest(storeName)
            .setOperations(operations);
        return executeStateTransaction(request).then();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> saveBulkState(String storeName, List<State<?>> states) {
        SaveStateRequest request = new SaveStateRequest(storeName)
            .setStates(states);
        return this.saveBulkState(request).then();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> saveState(String storeName, String key, Object value) {
        return this.saveState(storeName, key, null, value, null);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> saveState(String storeName, String key, String etag, Object value, StateOptions options) {
        State<?> state = new State<>(key, value, etag, options);
        return this.saveBulkState(storeName, Collections.singletonList(state));
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> deleteState(String storeName, String key) {
        return this.deleteState(storeName, key, null, null);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Mono<Void> deleteState(String storeName, String key, String etag, StateOptions options) {
        DeleteStateRequest request = new DeleteStateRequest(storeName, key)
            .setEtag(etag)
            .setStateOptions(options);
        return deleteState(request).then();
    }
}
