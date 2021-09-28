package spec.sdk.runtime.v1.domain.state;

import java.util.Map;

/**
 * A request to get a state by key.
 */
public class GetStateRequest {

  private final String storeName;

  private final String key;

  private Map<String, String> metadata;

  private StateOptions stateOptions;

  public GetStateRequest(String storeName, String key) {
    this.storeName = storeName;
    this.key = key;
  }

  public String getStoreName() {
    return storeName;
  }

  public String getKey() {
    return key;
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

  /**
   * Getter method for property <tt>stateOptions</tt>.
   *
   * @return property value of stateOptions
   */
  public StateOptions getStateOptions() {
    return stateOptions;
  }

  /**
   * Setter method for property <tt>stateOptions</tt>.
   *
   * @param stateOptions value to be assigned to property stateOptions
   */
  public void setStateOptions(StateOptions stateOptions) {
    this.stateOptions = stateOptions;
  }
}
