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
package spec.sdk.runtime.v1.domain;

import spec.sdk.runtime.v1.domain.state.*;

import java.util.List;
import java.util.Map;

public interface StateRuntime {

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param clazz     The type of State needed as return.
     * @param <T>       The type of the return.
     */
    <T> State<T> getState(String storeName, String key, Class<T> clazz);

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param options   Optional settings for retrieve operation.
     * @param clazz     The Type of State needed as return.
     * @param <T>       The Type of the return.
     */
    <T> State<T> getState(String storeName, String key, StateOptions options, Class<T> clazz);

    /**
     * Retrieve a State based on their key.
     *
     * @param request The request to get state.
     * @param clazz   The Class of State needed as return.
     * @param <T>     The Type of the return.
     * @return The requested State.
     */
    <T> State<T> getState(GetStateRequest request, Class<T> clazz);

    /**
     * Retrieve a State based on their key.
     *
     * @param request   The request to get state.
     * @param clazz     The Class of State needed as return.
     * @param timeoutMs
     * @param <T>       The Type of the return.
     * @return
     */
    <T> State<T> getState(GetStateRequest request, Class<T> clazz, int timeoutMs);

    /**
     * Retrieve bulk States based on their keys.
     *
     * @param storeName The name of the state store.
     * @param keys      The keys of the State to be retrieved.
     * @param clazz     The type of State needed as return.
     * @param <T>       The type of the return.
     */
    <T> List<State<T>> getBulkState(String storeName, List<String> keys, Class<T> clazz);

    /**
     * Retrieve bulk States based on their keys.
     *
     * @param request The request to get state.
     * @param clazz   The Class of State needed as return.
     * @param <T>     The Type of the return.
     * @return The requested State.
     */
    <T> List<State<T>> getBulkState(GetBulkStateRequest request, Class<T> clazz);

    /**
     * Execute a transaction.
     *
     * @param storeName  The name of the state store.
     * @param operations The operations to be performed.
     */
    void executeStateTransaction(String storeName,
                                 List<TransactionalStateOperation<?>> operations);

    /**
     * Execute a transaction.
     *
     * @param request Request to execute transaction.
     */
    void executeStateTransaction(ExecuteStateTransactionRequest request);

    /**
     * Save/Update a list of states.
     *
     * @param storeName The name of the state store.
     * @param states    The States to be saved.
     */
    void saveBulkState(String storeName, List<State<?>> states);

    /**
     * Save/Update a list of states.
     *
     * @param request Request to save states.
     */
    void saveBulkState(SaveStateRequest request);

    /**
     * Save/Update a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the state.
     * @param value     The value of the state.
     */
    void saveState(String storeName, String key, Object value);

    /**
     * Save/Update a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the state.
     * @param etag      The etag to be used.
     * @param value     The value of the state.
     * @param options   The Options to use for each state.
     */
    void saveState(String storeName, String key, String etag, Object value, StateOptions options, Map<String, String> metadata);

    /**
     * Delete a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be removed.
     */
    void deleteState(String storeName, String key);

    /**
     * Delete a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be removed.
     * @param etag      Optional etag for conditional delete.
     * @param options   Optional settings for state operation.
     */
    void deleteState(String storeName, String key, String etag, StateOptions options);

    /**
     * Delete a state.
     *
     * @param request Request to delete a state.
     */
    void deleteState(DeleteStateRequest request);
}
