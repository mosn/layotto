package io.mosn.layotto.v1;

import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.invocation.InvokeResponse;
import spec.sdk.runtime.v1.domain.state.*;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public abstract class AbstractRuntimeClient implements RuntimeClient {

    /**
     * Runtime client logger.
     */
    protected final Logger logger;
    /**
     * Serializer used for state objects.
     */
    protected final ObjectSerializer stateSerializer;
    /**
     * Runtime invocation timeout ms.
     */
    private final int timeoutMs;

    AbstractRuntimeClient(Logger logger, int timeoutMs, ObjectSerializer stateSerializer) {
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
        this.publishEvent(pubsubName, topicName, data, RuntimeProperties.PUBSUB_CONTENT_TYPE.get(), metadata);
    }

    @Override
    public void publishEvent(String pubsubName, String topicName, byte[] data, Map<String, String> metadata) {
        this.publishEvent(pubsubName, topicName, data, RuntimeProperties.PUBSUB_CONTENT_TYPE.get(), metadata);
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

    /**
     * Getter method for property <tt>timeoutMs</tt>.
     *
     * @return property value of timeoutMs
     */
    public int getTimeoutMs() {
        return this.timeoutMs;
    }
}
