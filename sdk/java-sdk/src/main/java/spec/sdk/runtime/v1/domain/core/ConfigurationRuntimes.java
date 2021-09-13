package spec.sdk.runtime.v1.domain.core;

import spec.sdk.runtime.v1.domain.core.configuration.ConfigurationItem;
import spec.sdk.runtime.v1.domain.core.configuration.ConfigurationRequestItem;
import spec.sdk.runtime.v1.domain.core.configuration.SaveConfigurationRequest;
import spec.sdk.runtime.v1.domain.core.configuration.SubConfigurationResp;
import spec.sdk.runtime.v1.utils.TypeRef;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.List;

/**
 * Configuration Runtimes standard API defined.
 */
public interface ConfigurationRuntimes {

    /**
     * Gets configuration from configuration store
     *
     * @param <T>                      The Type of the return.
     * @param configurationRequestItem Request object.
     * @param type                     The Type needed as return for the call.
     * @return A Mono Plan of response with type T.
     */
    <T> Mono<List<ConfigurationItem<T>>> getConfiguration(ConfigurationRequestItem configurationRequestItem, TypeRef<T> type);

    /**
     * Saves configuration into configuration store.
     *
     * @param saveConfigurationRequest Request object.
     * @return A Mono Plan of invocation.
     */
    Mono<Void> saveConfiguration(SaveConfigurationRequest saveConfigurationRequest);

    /**
     * Deletes configuration from configuration store.
     *
     * @param configurationRequestItem Request object.
     * @return A Mono Plan of invocation.
     */
    Mono<Void> deleteConfiguration(ConfigurationRequestItem configurationRequestItem);

    /**
     * Gets configuration from configuration store and subscribe the updates.
     *
     * @param <T>                      The Type of the return.
     * @param configurationRequestItem Request object.
     * @param type                     The Type needed as return for the call.
     * @return A Flux Plan of response with type T. Subscribe update listener.
     */
    <T> Flux<SubConfigurationResp<T>> subscribeConfiguration(ConfigurationRequestItem configurationRequestItem, TypeRef<T> type);
}
