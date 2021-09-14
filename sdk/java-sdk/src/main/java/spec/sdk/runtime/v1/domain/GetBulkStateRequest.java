package spec.sdk.runtime.v1.domain;

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
