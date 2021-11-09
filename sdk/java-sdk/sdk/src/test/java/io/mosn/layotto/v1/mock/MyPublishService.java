package io.mosn.layotto.v1.mock;

import com.google.protobuf.Empty;
import io.grpc.stub.StreamObserver;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;

public class MyPublishService extends RuntimeGrpc.RuntimeImplBase {
    RuntimeProto.PublishEventRequest lastReq = null;

    @Override
    public void publishEvent(RuntimeProto.PublishEventRequest request, StreamObserver<Empty> responseObserver) {
        lastReq = request;
        responseObserver.onNext(null);
        responseObserver.onCompleted();
    }

    /**
     * Getter method for property <tt>lastReq</tt>.
     *
     * @return property value of lastReq
     */
    public RuntimeProto.PublishEventRequest getLastReq() {
        return lastReq;
    }
}