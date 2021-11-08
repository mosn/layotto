package io.mosn.layotto.v1;

import com.google.protobuf.Empty;
import io.grpc.stub.StreamObserver;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class MyStateService extends RuntimeGrpc.RuntimeImplBase {
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
}