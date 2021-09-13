package spec.sdk.runtime.v1.domain.core.configuration;

import java.util.List;
import java.util.Map;

/**
 * ConfigurationRequestItem used for GET,DEL,SUB request
 */
public class ConfigurationRequestItem {

    /**
     * The name of configuration store.
     */
    private String storeName;
    /**
     * The application id which
     * Only used for admin, Ignored and reset for normal client
     */
    private String appId;
    /**
     * The group of keys.
     */
    private String group;
    /**
     * The label for keys.
     */
    private String label;
    /**
     * The keys to get.
     */
    private List<String> keys;
    /**
     * The metadata which will be sent to configuration store components.
     */
    private Map<String, String> metadata;

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

    public List<String> getKeys() {
        return keys;
    }

    public void setKeys(List<String> keys) {
        this.keys = keys;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }
}
