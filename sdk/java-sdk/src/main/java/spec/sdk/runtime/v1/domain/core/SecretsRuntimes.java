package spec.sdk.runtime.v1.domain.core;

import spec.sdk.runtime.v1.domain.core.secrets.GetBulkSecretRequest;
import spec.sdk.runtime.v1.domain.core.secrets.GetSecretRequest;
import reactor.core.publisher.Mono;

import java.util.Map;

/**
 * Secrets Runtimes standard API defined.
 */
public interface SecretsRuntimes {

    /**
     * Fetches a secret from the configured vault.
     *
     * @param storeName  Name of vault component in CloudRuntimes.
     * @param secretName Secret to be fetched.
     * @param metadata   Optional metadata.
     * @return Key-value pairs for the secret.
     */
    Mono<Map<String, String>> getSecret(String storeName, String secretName, Map<String, String> metadata);

    /**
     * Fetches a secret from the configured vault.
     *
     * @param storeName  Name of vault component in CloudRuntimes.
     * @param secretName Secret to be fetched.
     * @return Key-value pairs for the secret.
     */
    Mono<Map<String, String>> getSecret(String storeName, String secretName);

    /**
     * Fetches a secret from the configured vault.
     *
     * @param request Request to fetch secret.
     * @return Key-value pairs for the secret.
     */
    Mono<Map<String, String>> getSecret(GetSecretRequest request);

    /**
     * Fetches all secrets from the configured vault.
     *
     * @param storeName Name of vault component in CloudRuntimes.
     * @return Key-value pairs for all the secrets in the state store.
     */
    Mono<Map<String, Map<String, String>>> getBulkSecret(String storeName);

    /**
     * Fetches all secrets from the configured vault.
     *
     * @param storeName Name of vault component in CloudRuntimes.
     * @param metadata  Optional metadata.
     * @return Key-value pairs for all the secrets in the state store.
     */
    Mono<Map<String, Map<String, String>>> getBulkSecret(String storeName, Map<String, String> metadata);

    /**
     * Fetches all secrets from the configured vault.
     *
     * @param request Request to fetch secret.
     * @return Key-value pairs for the secret.
     */
    Mono<Map<String, Map<String, String>>> getBulkSecret(GetBulkSecretRequest request);
}
