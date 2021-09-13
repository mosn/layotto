/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.state;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * A request to get bulk state by keys.
 */
public class GetBulkStateRequest {

    private final String storeName;

    private final List<String> keys;

    private Map<String, String> metadata;

    private int parallelism = 1;

    public GetBulkStateRequest(String storeName, List<String> keys) {
        this.storeName = storeName;
        this.keys = keys == null ? null : Collections.unmodifiableList(keys);
    }

    public GetBulkStateRequest(String storeName, String... keys) {
        this.storeName = storeName;
        this.keys = keys == null ? null : Collections.unmodifiableList(Arrays.asList(keys));
    }

    public String getStoreName() {
        return storeName;
    }

    public List<String> getKeys() {
        return keys;
    }

    public int getParallelism() {
        return parallelism;
    }

    public GetBulkStateRequest setParallelism(int parallelism) {
        this.parallelism = parallelism;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public GetBulkStateRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
