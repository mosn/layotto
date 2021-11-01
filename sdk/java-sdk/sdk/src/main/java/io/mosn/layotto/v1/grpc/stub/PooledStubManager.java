package io.mosn.layotto.v1.grpc.stub;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.AbstractAsyncStub;
import io.grpc.stub.AbstractBlockingStub;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;

public class PooledStubManager<A extends AbstractAsyncStub, B extends AbstractBlockingStub>
        implements StubManager<A, B> {

    private volatile ManagedChannel[] channels;
    private final    RRPool<A>        asyncRuntimePool;
    private final    RRPool<B>        runtimePool;

    /**
     * Construct a new PooledStubManager with same channels in `other` PooledStubManager
     *
     * @param other
     * @param sc
     */
    public PooledStubManager(PooledStubManager other, StubCreator<A, B> sc) {
        // 1. validate
        if (other == null || other.channels == null) {
            throw new IllegalArgumentException("Invalid other");
        }
        if (sc == null) {
            throw new IllegalArgumentException("Invalid StubCreator");
        }
        int size = other.channels.length;
        this.channels = new ManagedChannel[size];
        List<A> asyncStubs = new ArrayList<>();
        List<B> blockingStubs = new ArrayList<>();
        // 2. construct channels and stubs
        for (int i = 0; i < size; i++) {
            // change the order of channels to avoid unbalanced load
            this.channels[i] = other.channels[size - 1 - i];
            asyncStubs.add(sc.createAsyncStub(channels[i]));
            blockingStubs.add(sc.createBlockingStub(channels[i]));
        }
        // 3. construct pools
        asyncRuntimePool = new RRPool<>(new CopyOnWriteArrayList<>(asyncStubs));
        runtimePool = new RRPool<>(new CopyOnWriteArrayList<>(blockingStubs));
        // 4. init connections
        init();
    }

    public PooledStubManager(String host, int port, int size,
                             StubCreator<A, B> sc) {
        channels = new ManagedChannel[size];
        List<A> asyncStubs = new ArrayList<>();
        List<B> blockingStubs = new ArrayList<>();
        // construct channels and stubs
        for (int i = 0; i < size; i++) {
            channels[i] = ManagedChannelBuilder.forAddress(host, port).usePlaintext().build();
            asyncStubs.add(sc.createAsyncStub(channels[i]));
            blockingStubs.add(sc.createBlockingStub(channels[i]));
        }
        // construct pools
        asyncRuntimePool = new RRPool<>(new CopyOnWriteArrayList<>(asyncStubs));
        runtimePool = new RRPool<>(new CopyOnWriteArrayList<>(blockingStubs));

        // init connections
        init();
    }

    protected void init() {
        // TODO establish connection
    }

    @Override
    public void destroy() {
        // get reference
        ManagedChannel[] chs = channels;
        // shutdown
        if (chs != null) {
            for (ManagedChannel c : chs) {
                c.shutdown();
            }
            channels = null;
        }
    }

    @Override
    public A getAsyncStub() {
        return asyncRuntimePool.next();
    }

    @Override
    public B getBlockingStub() {
        return runtimePool.next();
    }

    @Override
    public ManagedChannel[] getChannels() {
        return Arrays.copyOf(channels, channels.length);
    }
}
