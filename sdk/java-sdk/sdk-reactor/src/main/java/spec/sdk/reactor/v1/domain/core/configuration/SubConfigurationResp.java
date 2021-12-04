package spec.sdk.reactor.v1.domain.core.configuration;


import java.util.List;

public class SubConfigurationResp<T> {

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
     * The list of configuration items to save.
     * To delete a exist item, set the key (also label) and let content to be empty
     */
    private List<ConfigurationItem<T>> items;

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

    public List<ConfigurationItem<T>> getItems() {
        return items;
    }

    public void setItems(List<ConfigurationItem<T>> items) {
        this.items = items;
    }
}
