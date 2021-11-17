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
package spec.sdk.runtime.v1.domain.state;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * A request to get bulk state by keys.
 */
public class GetBulkStateRequest {

    private final String storeName;

    private final List<String> keys;

    private Map<String, String> metadata = new HashMap<>();

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

    /**
     * Getter method for property <tt>metadata</tt>.
     *
     * @return property value of metadata
     */
    public Map<String, String> getMetadata() {
        return metadata;
    }

    /**
     * Setter method for property <tt>metadata</tt>.
     *
     * @param metadata value to be assigned to property metadata
     */
    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }

    public void putMetadata(String key, String value) {
        if (key == null) { throw new java.lang.NullPointerException(); }
        if (value == null) { throw new java.lang.NullPointerException(); }
        if (metadata == null) {
            metadata = new HashMap<>();
        }
        metadata.put(key, value);
    }

    /**
     * Getter method for property <tt>parallelism</tt>.
     *
     * @return property value of parallelism
     */
    public int getParallelism() {
        return parallelism;
    }

    /**
     * Setter method for property <tt>parallelism</tt>.
     *
     * @param parallelism value to be assigned to property parallelism
     */
    public void setParallelism(int parallelism) {
        this.parallelism = parallelism;
    }
}
