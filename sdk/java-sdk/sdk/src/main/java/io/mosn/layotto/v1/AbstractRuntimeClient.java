package io.mosn.layotto.v1;

import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.invocation.InvokeResponse;
import spec.sdk.runtime.v1.domain.state.DeleteStateRequest;
import spec.sdk.runtime.v1.domain.state.ExecuteStateTransactionRequest;
import spec.sdk.runtime.v1.domain.state.GetBulkStateRequest;
import spec.sdk.runtime.v1.domain.state.GetStateRequest;
import spec.sdk.runtime.v1.domain.state.SaveStateRequest;
import spec.sdk.runtime.v1.domain.state.State;
import spec.sdk.runtime.v1.domain.state.StateOptions;
import spec.sdk.runtime.v1.domain.state.TransactionalStateOperation;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public abstract class AbstractRuntimeClient implements RuntimeClient {

    protected final String DEFAULT_PUBSUB_CONTENT_TYPE = "";

    protected Logger logger;

    private int timeoutMs;

    protected ObjectSerializer stateSerializer;

    AbstractRuntimeClient(Logger logger, int timeoutMs, ObjectSerializer stateSerializer) {
        this.logger = logger;
        this.timeoutMs = timeoutMs;
        this.stateSerializer = stateSerializer;
    }

    public String sayHello(String name) {
        return sayHello(name, getTimeoutMs());
    }

    public InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header) {
        return invokeMethod(appId, methodName, data, header, getTimeoutMs());
    }
    // TODO add some methods that serialize data before invoking method

    public void publishEvent(String pubsubName, String topicName, byte[] data) {
        publishEvent(pubsubName, topicName, data, DEFAULT_PUBSUB_CONTENT_TYPE, new HashMap<>());
    }

    @Override
    public void publishEvent(String pubsubName, String topicName, byte[] data, Map<String, String> metadata) {
        publishEvent(pubsubName, topicName, data, DEFAULT_PUBSUB_CONTENT_TYPE, metadata);
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
        return getState(storeName, key, null, clazz);
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
        deleteState(request);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveState(String storeName, String key, Object value) {
        this.saveState(storeName, key, null, value, null, new HashMap<>());
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void saveState(String storeName, String key, String etag, Object value, StateOptions options, Map<String, String> metadata) {
        State<?> state = new State<>(key, value, etag, metadata, options);
        this.saveBulkState(storeName, Collections.singletonList(state));
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
    public void executeStateTransaction(String storeName,
                                        List<TransactionalStateOperation<?>> operations) {
        ExecuteStateTransactionRequest request = new ExecuteStateTransactionRequest(storeName);
        request.setOperations(operations);
        executeStateTransaction(request);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> List<State<T>> getBulkState(String storeName, List<String> keys, Class<T> clazz) {
        return this.getBulkState(new GetBulkStateRequest(storeName, keys), clazz);
    }

    /**
     * Getter method for property <tt>timeoutMs</tt>.
     *
     * @return property value of timeoutMs
     */
    public int getTimeoutMs() {
        return timeoutMs;
    }
}
