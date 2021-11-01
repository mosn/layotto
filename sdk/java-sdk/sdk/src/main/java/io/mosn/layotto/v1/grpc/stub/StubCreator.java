package io.mosn.layotto.v1.grpc.stub;

import io.grpc.ManagedChannel;
import io.grpc.stub.AbstractAsyncStub;
import io.grpc.stub.AbstractBlockingStub;

public interface StubCreator<A extends AbstractAsyncStub, B extends AbstractBlockingStub> {

    A createAsyncStub(ManagedChannel channel);

    B createBlockingStub(ManagedChannel channel);

}