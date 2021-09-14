package spec.sdk.runtime.v1.domain;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * A request to save states to state store.
 */
public class SaveStateRequest {

  private final String storeName;

  private List<State<?>> states;

  public SaveStateRequest(String storeName) {
    this.storeName = storeName;
  }

  public String getStoreName() {
    return storeName;
  }

  /**
   * Getter method for property <tt>states</tt>.
   *
   * @return property value of states
   */
  public List<State<?>> getStates() {
    return states;
  }

  /**
   * Setter method for property <tt>states</tt>.
   *
   * @param states value to be assigned to property states
   */
  public void setStates(List<State<?>> states) {
    this.states = states;
  }
}
