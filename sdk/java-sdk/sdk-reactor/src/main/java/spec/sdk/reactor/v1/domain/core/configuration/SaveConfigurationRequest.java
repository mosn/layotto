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

public class SaveConfigurationRequest {

    /**
     * The name of configuration store.
     */
    private String                          storeName;
    /**
     * The application id which
     * Only used for admin, Ignored and reset for normal client
     */
    private String                          appId;
    /**
     * The list of configuration items to save.
     * To delete a exist item, set the key (also label) and let content to be empty
     */
    private List<ConfigurationItem<Object>> items;
    /**
     * The metadata which will be sent to configuration store components.
     */
    private Map<String, String>             metadata;

    public String getStoreName() {
        return storeName;
    }

    public void setStoreName(String storeName) {
        this.storeName = storeName;
    }

    public String getAppId() {
        return appId;
    }

    public void setAppId(String appId) {
        this.appId = appId;
    }

    public List<ConfigurationItem<Object>> getItems() {
        return items;
    }

    public void setItems(List<ConfigurationItem<Object>> items) {
        this.items = items;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }
}
