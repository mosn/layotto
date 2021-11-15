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
import spec.sdk.runtime.v1.domain.state.TransactionalStateOperation;

import java.util.ArrayList;
import java.util.List;

public class RedisCRUD {
    static String storeName = "redis";
    static String key1      = "key1";
    static String key2      = "key2";
    static String key3      = "key3";

    public static void main(String[] args) {
        // build RuntimeClient
        RuntimeClient client = new RuntimeClientBuilder()
                .withPort(RuntimeProperties.DEFAULT_PORT)
                .build();
        // saveState
        client.saveState(storeName, key1, "v11");
        // getState
        State<String> state = client.getState(storeName, key1, String.class);
        assertEquals(state.getKey(), key1);
        assertEquals(state.getValue(), "v11");
        System.out.println("get state key:" + state.getKey() + "  value:" + state.getValue());

        // deleteState
        client.deleteState(storeName, key1);

        // getState
        state = client.getState(storeName, key1, String.class);
        assertEquals(state.getKey(), key1);
        // TODO: currently Redis component can't tell the difference between null and 'non exist'
        //assertEquals(state.getValue(), null);
        assertEquals(state.getValue(), "");
        System.out.println("get state after delete. key:" + state.getKey() + "  value:" + state.getValue());

        // saveBulkState
        List<State<?>> list = new ArrayList<>();
        State<?> state1 = new State<>(key1, "v1", null, null);
        State<?> state2 = new State<>(key2, "v2", null, null);
        list.add(state2);
        list.add(state1);
        client.saveBulkState(storeName, list);

        // execute transaction
        List<TransactionalStateOperation<?>> operationList = new ArrayList<>();
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key2, new TestClass(key2), "")));

        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key3, "v3", "")));

        client.executeStateTransaction(storeName, operationList);

        //    getBulkState
        List<String> keys = new ArrayList<>();
        keys.add(key3);
        keys.add(key1);
        GetBulkStateRequest req = new GetBulkStateRequest(storeName, keys);
        List<State<String>> bulkState = client.getBulkState(req, String.class);
        assertTrue(bulkState.size() == 2);
        for (State<String> st : bulkState) {
            String key = st.getKey();
            if (key.equals(key1)) {
                assertEquals(st.getValue(), "v1");
            } else if (key.equals(key3)) {
                assertEquals(st.getValue(), "v3");
            } else {
                throw new RuntimeException("Unexpected key:" + key);
            }
        }

        keys = new ArrayList<>();
        keys.add(key2);
        req = new GetBulkStateRequest(storeName, keys);
        List<State<TestClass>> resp = client.getBulkState(req, TestClass.class);
        assertTrue(resp.size() == 1);
        assertEquals(resp.get(0).getValue().name, key2);

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

    public static class TestClass {
        String name;

        public TestClass(String name) {
            this.name = name;
        }

        /**
         * Getter method for property <tt>name</tt>.
         *
         * @return property value of name
         */
        public String getName() {
            return name;
        }
    }
}
