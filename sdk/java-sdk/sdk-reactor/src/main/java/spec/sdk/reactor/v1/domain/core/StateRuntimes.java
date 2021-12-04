package spec.sdk.reactor.v1.domain.core;

import spec.sdk.reactor.v1.domain.core.state.DeleteStateRequest;
import spec.sdk.reactor.v1.domain.core.state.ExecuteStateTransactionRequest;
import spec.sdk.reactor.v1.domain.core.state.GetBulkStateRequest;
import spec.sdk.reactor.v1.domain.core.state.GetStateRequest;
import spec.sdk.reactor.v1.domain.core.state.SaveStateRequest;
import spec.sdk.reactor.v1.domain.core.state.State;
import spec.sdk.reactor.v1.domain.core.state.StateOptions;
import spec.sdk.reactor.v1.domain.core.state.TransactionalStateOperation;
import spec.sdk.reactor.v1.domain.core.state.*;
import spec.sdk.reactor.v1.utils.TypeRef;
import reactor.core.publisher.Mono;

import java.util.List;

/**
 * State Management Runtimes standard API defined.
 */
public interface StateRuntimes {

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param state     State to be re-retrieved.
     * @param type      The type of State needed as return.
     * @param <T>       The type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(String storeName, State<T> state, TypeRef<T> type);

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param state     State to be re-retrieved.
     * @param clazz     The type of State needed as return.
     * @param <T>       The type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(String storeName, State<T> state, Class<T> clazz);

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param type      The type of State needed as return.
     * @param <T>       The type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(String storeName, String key, TypeRef<T> type);

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param clazz     The type of State needed as return.
     * @param <T>       The type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(String storeName, String key, Class<T> clazz);

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param options   Optional settings for retrieve operation.
     * @param type      The Type of State needed as return.
     * @param <T>       The Type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, TypeRef<T> type);

    /**
     * Retrieve a State based on their key.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be retrieved.
     * @param options   Optional settings for retrieve operation.
     * @param clazz     The Type of State needed as return.
     * @param <T>       The Type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(String storeName, String key, StateOptions options, Class<T> clazz);

    /**
     * Retrieve a State based on their key.
     *
     * @param request The request to get state.
     * @param type    The Type of State needed as return.
     * @param <T>     The Type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<State<T>> getState(GetStateRequest request, TypeRef<T> type);

    /**
     * Retrieve bulk States based on their keys.
     *
     * @param storeName The name of the state store.
     * @param keys      The keys of the State to be retrieved.
     * @param type      The type of State needed as return.
     * @param <T>       The type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, TypeRef<T> type);

    /**
     * Retrieve bulk States based on their keys.
     *
     * @param storeName The name of the state store.
     * @param keys      The keys of the State to be retrieved.
     * @param clazz     The type of State needed as return.
     * @param <T>       The type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<List<State<T>>> getBulkState(String storeName, List<String> keys, Class<T> clazz);

    /**
     * Retrieve bulk States based on their keys.
     *
     * @param request The request to get state.
     * @param type    The Type of State needed as return.
     * @param <T>     The Type of the return.
     * @return A Mono Plan for the requested State.
     */
    <T> Mono<List<State<T>>> getBulkState(GetBulkStateRequest request, TypeRef<T> type);

    /**
     * Execute a transaction.
     *
     * @param storeName  The name of the state store.
     * @param operations The operations to be performed.
     * @return a Mono plan of type Void
     */
    Mono<Void> executeStateTransaction(String storeName,
                                       List<TransactionalStateOperation<?>> operations);

    /**
     * Execute a transaction.
     *
     * @param request Request to execute transaction.
     * @return a Mono plan of type Response Void
     */
    Mono<Void> executeStateTransaction(ExecuteStateTransactionRequest request);

    /**
     * Save/Update a list of states.
     *
     * @param storeName The name of the state store.
     * @param states    The States to be saved.
     * @return a Mono plan of type Void.
     */
    Mono<Void> saveBulkState(String storeName, List<State<?>> states);

    /**
     * Save/Update a list of states.
     *
     * @param request Request to save states.
     * @return a Mono plan of type Void.
     */
    Mono<Void> saveBulkState(SaveStateRequest request);

    /**
     * Save/Update a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the state.
     * @param value     The value of the state.
     * @return a Mono plan of type Void.
     */
    Mono<Void> saveState(String storeName, String key, Object value);

    /**
     * Save/Update a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the state.
     * @param etag      The etag to be used.
     * @param value     The value of the state.
     * @param options   The Options to use for each state.
     * @return a Mono plan of type Void.
     */
    Mono<Void> saveState(String storeName, String key, String etag, Object value, StateOptions options);

    /**
     * Delete a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be removed.
     * @return a Mono plan of type Void.
     */
    Mono<Void> deleteState(String storeName, String key);

    /**
     * Delete a state.
     *
     * @param storeName The name of the state store.
     * @param key       The key of the State to be removed.
     * @param etag      Optional etag for conditional delete.
     * @param options   Optional settings for state operation.
     * @return a Mono plan of type Void.
     */
    Mono<Void> deleteState(String storeName, String key, String etag, StateOptions options);

    /**
     * Delete a state.
     *
     * @param request Request to delete a state.
     * @return a Mono plan of type Void.
     */
    Mono<Void> deleteState(DeleteStateRequest request);
}
