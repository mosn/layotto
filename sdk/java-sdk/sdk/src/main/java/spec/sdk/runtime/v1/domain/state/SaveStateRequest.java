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
 *
 */
package spec.sdk.runtime.v1.domain.state;

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
