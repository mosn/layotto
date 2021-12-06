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
package spec.sdk.reactor.v1.domain.core.configuration;

import java.util.List;
import java.util.Map;

/**
 * Builds a request to invoke configuration.
 */
public class ConfigurationRequestItemBuilder {

    private final String        storeName;

    private final String        appId;

    private String              group;

    private String              label;

    private List<String>        keys;

    private Map<String, String> metadata;

    public ConfigurationRequestItemBuilder(String storeName, String appId) {
        this.storeName = storeName;
        this.appId = appId;
    }

    public ConfigurationRequestItemBuilder withGroup(String group) {
        this.group = group;
        return this;
    }

    public ConfigurationRequestItemBuilder withLabel(String label) {
        this.label = label;
        return this;
    }

    public ConfigurationRequestItemBuilder withKeys(List<String> keys) {
        this.keys = keys;
        return this;
    }

    public ConfigurationRequestItemBuilder withMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
        return this;
    }

    /**
     * Builds a request object.
     *
     * @return Request object.
     */
    public ConfigurationRequestItem build() {
        ConfigurationRequestItem request = new ConfigurationRequestItem();
        request.setStoreName(this.appId);
        request.setAppId(this.appId);
        request.setGroup(this.group);
        request.setLabel(this.label);
        request.setKeys(this.keys);
        request.setMetadata(this.metadata);
        return request;
    }
}
