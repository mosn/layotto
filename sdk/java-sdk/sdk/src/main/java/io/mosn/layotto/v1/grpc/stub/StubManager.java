package io.mosn.layotto.v1.grpc.stub;

import io.grpc.stub.AbstractAsyncStub;
import io.grpc.stub.AbstractBlockingStub;

/**
 * An abstraction to manage grpc stub.
 * It can be used to implement connection pool pattern
 */
public interface StubManager<A extends AbstractAsyncStub, B extends AbstractBlockingStub> {

    A getAsyncStub();

    B getBlockingStub();

    void destroy();
}
