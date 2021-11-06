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
package io.mosn.layotto.v1;

import com.google.protobuf.Empty;
import io.grpc.ManagedChannel;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.stub.StreamObserver;
import io.grpc.testing.GrpcCleanupRule;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.state.GetBulkStateRequest;
import spec.sdk.runtime.v1.domain.state.GetStateRequest;
import spec.sdk.runtime.v1.domain.state.State;
import spec.sdk.runtime.v1.domain.state.StateOptions;
import spec.sdk.runtime.v1.domain.state.TransactionalStateOperation;

import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Deque;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.mockito.AdditionalAnswers.delegatesTo;
import static org.mockito.Mockito.mock;

@RunWith(JUnit4.class)
public class StateTest {
    @Rule
    public final GrpcCleanupRule grpcCleanup = new GrpcCleanupRule();

    private final RuntimeGrpc.RuntimeImplBase serviceImpl =
            mock(RuntimeGrpc.RuntimeImplBase.class, delegatesTo(
                    new RuntimeGrpc.RuntimeImplBase() {
                        private Map<String, RuntimeProto.StateItem> store = new HashMap<>();

                        @Override
                        public void getState(RuntimeProto.GetStateRequest request,
                                             StreamObserver<RuntimeProto.GetStateResponse> responseObserver) {
                            RuntimeProto.StateItem item = store.get(request.getKey());
                            if (item == null) {
                                responseObserver.onNext(null);
                                responseObserver.onCompleted();
                                return;
                            }
                            RuntimeProto.GetStateResponse resp = RuntimeProto.GetStateResponse.newBuilder()
                                    .setData(item.getValue())
                                    .setEtag(item.getEtag().getValue())
                                    .putAllMetadata(item.getMetadataMap())
                                    .build();
                            responseObserver.onNext(resp);
                            responseObserver.onCompleted();
                        }

                        @Override
                        public void getBulkState(RuntimeProto.GetBulkStateRequest request,
                                                 StreamObserver<RuntimeProto.GetBulkStateResponse> responseObserver) {
                            RuntimeProto.GetBulkStateResponse.Builder builder = RuntimeProto.GetBulkStateResponse.newBuilder();
                            for (int i = 0; i < request.getKeysCount(); i++) {
                                RuntimeProto.BulkStateItem.Builder itemBuilder = RuntimeProto.BulkStateItem.newBuilder().setKey(
                                        request.getKeys(i));
                                RuntimeProto.StateItem item = store.get(request.getKeys(i));
                                if (item != null) {
                                    itemBuilder = itemBuilder.setData(item.getValue()).setEtag(item.getEtag().getValue()).putAllMetadata(
                                            item.getMetadataMap());
                                }

                                builder = builder.addItems(itemBuilder.build());
                            }
                            responseObserver.onNext(builder.build());
                            responseObserver.onCompleted();
                        }

                        @Override
                        public void saveState(RuntimeProto.SaveStateRequest request, StreamObserver<Empty> responseObserver) {
                            for (int i = 0; i < request.getStatesCount(); i++) {
                                store.put(request.getStates(i).getKey(), request.getStates(i));
                            }
                            responseObserver.onNext(null);
                            responseObserver.onCompleted();
                        }

                        @Override
                        public void deleteState(RuntimeProto.DeleteStateRequest request, StreamObserver<Empty> responseObserver) {
                            store.remove(request.getKey());
                            responseObserver.onNext(null);
                            responseObserver.onCompleted();
                        }

                        @Override
                        public void deleteBulkState(RuntimeProto.DeleteBulkStateRequest request, StreamObserver<Empty> responseObserver) {
                            for (int i = 0; i < request.getStatesCount(); i++) {
                                store.remove(request.getStates(i).getKey());
                            }
                            responseObserver.onNext(null);
                            responseObserver.onCompleted();
                        }

                        @Override
                        public void executeStateTransaction(RuntimeProto.ExecuteStateTransactionRequest request,
                                                            StreamObserver<Empty> responseObserver) {
                            List<RuntimeProto.TransactionalStateOperation> list = request.getOperationsList();
                            Map<String, RuntimeProto.StateItem> newStore = new HashMap<>(store);
                            try {
                                for (RuntimeProto.TransactionalStateOperation tso : list) {
                                    String type = tso.getOperationType();
                                    RuntimeProto.StateItem req = tso.getRequest();
                                    if ("upsert".equals(type)) {
                                        newStore.put(req.getKey(), req);
                                    } else if ("delete".equals(type)) {
                                        newStore.remove(req.getKey());
                                    } else {
                                        throw new RuntimeException("illegal type" + type);
                                    }
                                }
                                store = newStore;
                                responseObserver.onNext(null);
                                responseObserver.onCompleted();
                            } catch (Exception e) {
                                responseObserver.onError(e);
                                responseObserver.onCompleted();
                            }
                        }
                    }));

    private RuntimeClient client;

    @Before
    public void setUp() throws Exception {
        String serverName = InProcessServerBuilder.generateName();
        grpcCleanup.register(InProcessServerBuilder
                .forName(serverName).directExecutor()
                .addService(serviceImpl)
                .build().start());
        ManagedChannel channel = grpcCleanup.register(
                InProcessChannelBuilder.forName(serverName).directExecutor().build());
        client = new RuntimeClientBuilder()
                .buildGrpcWithExistingChannel(channel);
    }

    @Test
    public void testStateCrud() {
        String storeName = "redis";

        // saveState
        client.saveState(storeName, "foo", "bar".getBytes());
        GetStateRequest req = new GetStateRequest(storeName, "foo");

        // getState
        State<String> resp = client.getState(req, String.class);
        Assert.assertEquals(resp.getValue(), "bar");

        // delete
        client.deleteState(storeName, "foo");

        // getState
        req = new GetStateRequest(storeName, "foo");
        resp = client.getState(req, String.class);
        Assert.assertEquals(resp.getValue().length(), 0);

        // saveState
        client.saveState(storeName, "key1", "bar1".getBytes());
        client.saveState(storeName, "key2", "bar2".getBytes());

        GetBulkStateRequest br = new GetBulkStateRequest(storeName, Arrays.asList("key1", "key2"));
        List<State<String>> bulkResp = client.getBulkState(br, String.class);
        Assert.assertEquals(bulkResp.get(0).getValue(), "bar1");
        Assert.assertEquals(bulkResp.get(1).getValue(), "bar2");

        // deleteState
        client.deleteState(storeName, "key1", null
                , new StateOptions(StateOptions.Consistency.STRONG, StateOptions.Concurrency.FIRST_WRITE));
        client.deleteState(storeName, "key2", null
                , new StateOptions(null, null));

        br = new GetBulkStateRequest(storeName, Arrays.asList("key1", "key2"));
        bulkResp = client.getBulkState(br, String.class);
        Assert.assertEquals(bulkResp.get(0).getValue().length(), 0);
        Assert.assertEquals(bulkResp.get(1).getValue().length(), 0);
    }

    @Test
    public void testTransaction() {
        String storeName = "redis";
        String key1 = "key11";
        String key2 = "key22";
        String key3 = "key33";

        // execute transaction
        List<TransactionalStateOperation<?>> operationList = new ArrayList<>();
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key1, new TestClass(key1), "")));
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key2, new TestClass(key2), "")));
        client.executeStateTransaction(storeName, operationList);

        //    getBulkState
        List<String> keys = new ArrayList<>();
        keys.add(key2);
        keys.add(key1);
        GetBulkStateRequest bulkReq = new GetBulkStateRequest(storeName, keys);
        List<State<TestClass>> bulkState = client.getBulkState(bulkReq, TestClass.class);
        Assert.assertTrue(bulkState.size() == 2);
        for (State<TestClass> st : bulkState) {
            String key = st.getKey();
            if (key.equals(key1)) {
                Assert.assertEquals(st.getValue().getName(), key1);
            } else if (key.equals(key2)) {
                Assert.assertEquals(st.getValue().getName(), key2);
            } else {
                throw new RuntimeException("Unexpected key:" + key);
            }
        }

        // execute transaction,update key3 and delete key1
        operationList = new ArrayList<>();
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key3, new TestClass(key3), "", null
                        , new StateOptions(StateOptions.Consistency.STRONG, StateOptions.Concurrency.LAST_WRITE))));
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.DELETE,
                new State<>(key2, new TestClass(key2), "")));
        client.executeStateTransaction(storeName, operationList);

        //    getBulkState
        keys = new ArrayList<>();
        keys.add(key2);
        keys.add(key1);
        keys.add(key3);
        bulkReq = new GetBulkStateRequest(storeName, keys);
        bulkState = client.getBulkState(bulkReq, TestClass.class);
        Assert.assertTrue(bulkState.size() == 3);
        for (State<TestClass> st : bulkState) {
            String key = st.getKey();
            if (key.equals(key1)) {
                Assert.assertEquals(st.getValue().getName(), key1);
            } else if (key.equals(key2)) {
                Assert.assertEquals(st.getValue(), null);
            } else if (key.equals(key3)) {
                Assert.assertEquals(st.getValue().getName(), key3);
            } else {
                throw new RuntimeException("Unexpected key:" + key);
            }
        }
    }

    @Test(expected = IllegalArgumentException.class)
    public void testTransactionNullOperation_thenIllegal() {
        String storeName = "redis";
        String key1 = "key11";
        String key2 = "key22";

        // execute transaction
        List<TransactionalStateOperation<?>> operationList = new ArrayList<>();
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key1, new TestClass(key1), "")));
        // null
        operationList.add(new TransactionalStateOperation<>(null,
                new State<>(key2, new TestClass(key2), "")));
        client.executeStateTransaction(storeName, operationList);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testTransactionEmptyKey_thenIllegal() {
        String storeName = "redis";
        String key1 = "";
        String key2 = "key22";

        // execute transaction
        List<TransactionalStateOperation<?>> operationList = new ArrayList<>();
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key1, new TestClass(key1), "")));
        operationList.add(new TransactionalStateOperation<>(TransactionalStateOperation.OperationType.UPSERT,
                new State<>(key2, new TestClass(key2), "")));
        client.executeStateTransaction(storeName, operationList);
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
