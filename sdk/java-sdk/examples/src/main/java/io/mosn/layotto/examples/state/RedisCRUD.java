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
import spec.sdk.runtime.v1.domain.state.GetBulkStateRequest;
import spec.sdk.runtime.v1.domain.state.State;

import java.util.ArrayList;
import java.util.List;

public class RedisCRUD {

    public static void main(String[] args) {
        // build RuntimeClient
        RuntimeClient client = new RuntimeClientBuilder()
                .withPort(RuntimeProperties.DEFAULT_PORT)
                .build();
        // saveState
        client.saveState("redis", "key1", "v11");
        // getState
        State<String> state = client.getState("redis", "key1", String.class);
        assertEquals(state.getKey(), "key1");
        assertEquals(state.getValue(), "v11");
        System.out.println("get state key:" + state.getKey() + "  value:" + state.getValue());

        // deleteState
        client.deleteState("redis", "key1");

        // getState
        state = client.getState("redis", "key1", String.class);
        assertEquals(state.getKey(), "key1");
        assertEquals(state.getValue(), null);
        System.out.println("get state after delete. key:" + state.getKey() + "  value:" + state.getValue());

        // saveBulkState
        List<State<?>> list = new ArrayList<>();
        State<?> state1 = new State<>("key11", "v1", null, null);
        State<?> state2 = new State<>("key22", "v2", null, null);
        list.add(state2);
        list.add(state1);
        client.saveBulkState("redis", list);

        //    getBulkState
        List<String> keys = new ArrayList<>();
        keys.add("key1");
        keys.add("key2");
        GetBulkStateRequest req = new GetBulkStateRequest("redis", keys);
        List<State<byte[]>> bulkState = client.getBulkState(req);
        assertTrue(bulkState.size() == 2);
        for (State<byte[]> st : bulkState) {
            String key = st.getKey();
            if (key.equals("key1")) {
                assertEquals(new String(st.getValue()), "v1");
            } else if (key.equals("key2")) {
                assertEquals(new String(st.getValue()), "v2");
            } else {
                throw new RuntimeException("Unexpected key:" + key);
            }
        }
    }

    private static void assertTrue(boolean b) {
        if (!b) {
            throw new RuntimeException("Assertion fail");
        }
    }

    private static void assertEquals(String actualResult, String expected) {
        if (actualResult == expected || actualResult.equals(expected)) {
            return;
        }
        throw new RuntimeException("Unexpected result:" + actualResult);
    }
}
