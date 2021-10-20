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
package io.mosn.layotto.examples.state;

import io.mosn.layotto.v1.RuntimeClientBuilder;
import io.mosn.layotto.v1.config.RuntimeProperties;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.state.State;

public class RedisCRUD {

    public static void main(String[] args) {
        RuntimeClient client = new RuntimeClientBuilder()
                .withPort(RuntimeProperties.DEFAULT_PORT)
                .build();

        client.saveState("redis", "key1", "v11");
        State<String> state = client.getState("redis", "key1", String.class);
        System.out.println("get state key:" + state.getKey() + "  value:" + state.getValue());
        client.deleteState("redis", "key1");
        state = client.getState("redis", "key1", String.class);
        System.out.println("get state after delete. key:" + state.getKey() + "  value:" + state.getValue());
    }
}
