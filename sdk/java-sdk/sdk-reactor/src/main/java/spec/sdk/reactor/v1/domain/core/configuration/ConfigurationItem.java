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

import java.util.Map;

public class ConfigurationItem<T> {

    /**
     * Required. The key of configuration item
     */
    private String              key;
    /**
     * The content of configuration item
     * Empty if the configuration is not set, including the case that the configuration is changed from value-set to value-not-set.
     */
    private T                   content;
    /**
     * The group of configuration item.
     */
    private String              group;
    /**
     * The label of configuration item.
     */
    private String              label;
    /**
     * The tag list of configuration item.
     */
    private Map<String, String> tags;
    /**
     * The metadata which will be passed to configuration store component.
     */
    private Map<String, String> metadata;

    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public T getContent() {
        return content;
    }

    public void setContent(T content) {
        this.content = content;
    }

    public String getGroup() {
        return group;
    }

    public void setGroup(String group) {
        this.group = group;
    }

    public String getLabel() {
        return label;
    }

    public void setLabel(String label) {
        this.label = label;
    }

    public Map<String, String> getTags() {
        return tags;
    }

    public void setTags(Map<String, String> tags) {
        this.tags = tags;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }
}
