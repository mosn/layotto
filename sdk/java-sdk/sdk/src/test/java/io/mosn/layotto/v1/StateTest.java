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

import java.util.Arrays;
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
                            responseObserver.onError(new Exception("not support"));
                            responseObserver.onCompleted();
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

        client.saveState(storeName, "foo", "bar".getBytes());
        GetStateRequest req = new GetStateRequest(storeName, "foo");
        State<byte[]> resp = client.getState(req);
        Assert.assertEquals(new String(resp.getValue()), "bar");

        client.deleteState(storeName, "foo");

        req = new GetStateRequest(storeName, "foo");
        resp = client.getState(req);
        Assert.assertEquals(resp.getValue().length, 0);

        client.saveState(storeName, "key1", "bar1".getBytes());
        client.saveState(storeName, "key2", "bar2".getBytes());

        GetBulkStateRequest br = new GetBulkStateRequest(storeName, Arrays.asList("key1", "key2"));
        List<State<byte[]>> bulkResp = client.getBulkState(br);
        Assert.assertEquals(new String(bulkResp.get(0).getValue()), "bar1");
        Assert.assertEquals(new String(bulkResp.get(1).getValue()), "bar2");

        client.deleteState(storeName, "key1");
        client.deleteState(storeName, "key2");

        br = new GetBulkStateRequest(storeName, Arrays.asList("key1", "key2"));
        bulkResp = client.getBulkState(br);
        Assert.assertEquals(bulkResp.get(0).getValue().length, 0);
        Assert.assertEquals(bulkResp.get(1).getValue().length, 0);
    }

}
