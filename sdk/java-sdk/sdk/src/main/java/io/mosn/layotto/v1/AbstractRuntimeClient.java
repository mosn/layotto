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
package io.mosn.layotto.v1;

import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.exceptions.RuntimeClientException;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.invocation.InvokeResponse;
import spec.sdk.runtime.v1.domain.state.*;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public abstract class AbstractRuntimeClient implements RuntimeClient {

    /**
     * Runtime client logger.
     */
    protected final Logger           logger;
    /**
     * Serializer used for state objects.
     */
    protected final ObjectSerializer stateSerializer;
    /**
     * Runtime invocation timeout ms.
     */
    private final   int              timeoutMs;

    AbstractRuntimeClient(Logger logger, int timeoutMs, ObjectSerializer stateSerializer) {
        if (logger == null) {
            throw new IllegalArgumentException("logger shouldn't be null");
        }
        if (stateSerializer == null) {
            throw new IllegalArgumentException("stateSerializer shouldn't be null");
        }
        this.logger = logger;
        this.timeoutMs = timeoutMs;
        this.stateSerializer = stateSerializer;
    }

    @Override
    public String sayHello(String name) {
        return this.sayHello(name, this.getTimeoutMs());
    }

    @Override
    public InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header) {
        return this.invokeMethod(appId, methodName, data, header, this.getTimeoutMs());
    }
    // TODO add some methods that serialize data before invoking method

    @Override
    public void publishEvent(String pubsubName, String topicName, byte[] data) {
        Map<String, String> metadata = new HashMap<>(2, 1);
        this.publishEvent(pubsubName, topicName, data, RuntimeProperties.DEFAULT_PUBSUB_CONTENT_TYPE, metadata);
    }

    @Override
    public void publishEvent(String pubsubName, String topicName, byte[] data, Map<String, String> metadata) {
        this.publishEvent(pubsubName, topicName, data, RuntimeProperties.DEFAULT_PUBSUB_CONTENT_TYPE, metadata);
    }

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param clazz     The type of State needed as return.
     */
    @Override
    public <T> State<T> getState(String storeName, String key, Class<T> clazz) {
        return this.getState(storeName, key, null, clazz);
    }

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param options   Optional settings for retrieve operation.
     * @param clazz     The Type of State needed as return.
     */
    @Override
    public <T> State<T> getState(String storeName, String key, StateOptions options, Class<T> clazz) {
        GetStateRequest request = new GetStateRequest(storeName, key);
        request.setStateOptions(options);
        return this.getState(request, clazz);
    }

    /**
     * Retrieve a State based on their key.
     *
     * @param request The request to get state.
     * @param clazz   The Class of State needed as return.
     * @return The requested State.
     */
    @Override
    public <T> State<T> getState(GetStateRequest request, Class<T> clazz) {
        return getState(request, clazz, getTimeoutMs());
    }

    @Override
    public <T> State<T> getState(GetStateRequest request, Class<T> clazz, int timeoutMs) {
        // 1. validate
        if (clazz == null) {
            throw new IllegalArgumentException("clazz cannot be null.");
        }
        final String stateStoreName = request.getStoreName();
        if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
            throw new IllegalArgumentException("State store name cannot be null or empty.");
        }
        final String key = request.getKey();
        if ((key == null) || (key.trim().isEmpty())) {
            throw new IllegalArgumentException("Key cannot be null or empty.");
        }
        // 2. invoke
        State<byte[]> state = doGetState(request, timeoutMs);
        try {
            // 3. deserialize
            T value = null;
            byte[] data = state.getValue();
            if (data != null) {
                value = stateSerializer.deserialize(data, clazz);
            }
            return new State<>(state.getKey(), value, state.getEtag(), state.getMetadata(), state.getOptions());
        } catch (Exception e) {
            logger.error("getState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    protected abstract State<byte[]> doGetState(GetStateRequest request, int timeoutMs);

    /**
     * {@inheritDoc}
     */
    @Override
    public void deleteState(String storeName, String key) {
        this.deleteState(storeName, key, null, null);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void deleteState(String storeName, String key, String etag, StateOptions options) {
        DeleteStateRequest request = new DeleteStateRequest(storeName, key);
        request.setEtag(etag);
        request.setStateOptions(options);
        this.deleteState(request);
    }

    /**
     * Delete a state.
     *
     * @param request Request to delete a state.
     */
    @Override
    public void deleteState(DeleteStateRequest request) {
        deleteState(request, getTimeoutMs());
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveState(String storeName, String key, Object value) {
        Map<String, String> metadata = new HashMap<>(2, 1);
        this.saveState(storeName, key, null, value, null, metadata);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveState(String storeName, String key, String etag, Object value, StateOptions options, Map<String, String> metadata) {
        State<?> state = new State<>(key, value, etag, metadata, options);
        List<State<?>> states = Collections.singletonList(state);
        this.saveBulkState(storeName, states);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveBulkState(String storeName, List<State<?>> states) {
        SaveStateRequest request = new SaveStateRequest(storeName);
        request.setStates(states);
        this.saveBulkState(request);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveBulkState(SaveStateRequest request) {
        saveBulkState(request, getTimeoutMs());
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void executeStateTransaction(String storeName, List<TransactionalStateOperation<?>> operations) {
        ExecuteStateTransactionRequest request = new ExecuteStateTransactionRequest(storeName);
        request.setOperations(operations);
        this.executeStateTransaction(request);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> List<State<T>> getBulkState(String storeName, List<String> keys, Class<T> clazz) {
        GetBulkStateRequest request = new GetBulkStateRequest(storeName, keys);
        return this.getBulkState(request, clazz);
    }

    @Override
    public <T> List<State<T>> getBulkState(GetBulkStateRequest request, Class<T> clazz) {
        // 1. validate
        if (clazz == null) {
            throw new IllegalArgumentException("clazz cannot be null.");
        }
        try {
            // 2. invoke
            List<State<byte[]>> bulkState = getBulkState(request, getTimeoutMs());
            // 3. deserialize
            List<State<T>> result = new ArrayList<>(bulkState.size());
            for (State<byte[]> state : bulkState) {
                byte[] value = state.getValue();
                T deValue = null;
                if (value != null) {
                    deValue = stateSerializer.deserialize(value, clazz);
                }
                State<T> tState = new State<>(state.getKey(), deValue, state.getEtag(), state.getMetadata(), state.getOptions());
                result.add(tState);
            }
            return result;
        } catch (Exception e) {
            logger.error("getBulkState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public List<State<byte[]>> getBulkState(GetBulkStateRequest request) {
        return getBulkState(request, getTimeoutMs());
    }

    /**
     * Getter method for property <tt>timeoutMs</tt>.
     *
     * @return property value of timeoutMs
     */
    public int getTimeoutMs() {
        return this.timeoutMs;
    }
}
