package io.mosn.layotto.v1.grpc.stub;

import io.grpc.ManagedChannel;
import io.grpc.stub.AbstractAsyncStub;
import io.grpc.stub.AbstractBlockingStub;

public class SingleStubManager<A extends AbstractAsyncStub, B extends AbstractBlockingStub>
        implements StubManager<A, B> {

    private       ManagedChannel channel;
    private final A              asyncStub;
    private final B              blockingStub;

    public SingleStubManager(ManagedChannel channel,
                             StubCreator<A, B> sc) {
        this.channel = channel;
        asyncStub = sc.createAsyncStub(channel);
        blockingStub = sc.createBlockingStub(channel);
    }

    @Override
    public B getBlockingStub() {
        return blockingStub;
    }

    @Override
    public A getAsyncStub() {
        return asyncStub;
    }

    @Override
    public void destroy() {
        if (channel != null) {
            channel.shutdown();
            channel = null;
        }
    }
}
