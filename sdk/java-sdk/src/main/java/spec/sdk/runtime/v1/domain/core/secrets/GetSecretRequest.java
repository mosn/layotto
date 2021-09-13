/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.secrets;

import java.util.Collections;
import java.util.Map;

/**
 * A request to get a secret by key.
 */
public class GetSecretRequest {

    private final String storeName;

    private final String key;

    public GetSecretRequest(String storeName, String key) {
        this.storeName = storeName;
        this.key = key;
    }

    private Map<String, String> metadata;

    public String getStoreName() {
        return storeName;
    }

    public String getKey() {
        return key;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public GetSecretRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
