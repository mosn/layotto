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
package io.mosn.layotto.v1.grpc.stub;

import io.grpc.ManagedChannel;
import io.grpc.stub.AbstractAsyncStub;
import io.grpc.stub.AbstractBlockingStub;

public class SingleStubManager<A extends AbstractAsyncStub, B extends AbstractBlockingStub>
        implements StubManager<A, B> {

    private volatile ManagedChannel channel;
    private final    A              asyncStub;
    private final    B              blockingStub;

    /**
     * Construct a new SingleStubManager with the existing channel in `other` SingleStubManager
     *
     * @param other
     * @param sc
     */
    public SingleStubManager(SingleStubManager other,
                             StubCreator<A, B> sc) {
        // 1. validate
        if (other == null || other.channel == null) {
            throw new IllegalArgumentException("Invalid other");
        }
        if (sc == null) {
            throw new IllegalArgumentException("Invalid StubCreator");
        }
        // 2. set fields
        this.channel = other.channel;
        asyncStub = sc.createAsyncStub(channel);
        blockingStub = sc.createBlockingStub(channel);
    }

    public SingleStubManager(ManagedChannel channel,
                             StubCreator<A, B> sc) {
        // 1. validate
        if (channel == null) {
            throw new IllegalArgumentException("Invalid channel");
        }
        if (sc == null) {
            throw new IllegalArgumentException("Invalid StubCreator");
        }
        // 2. set fields
        this.channel = channel;
        asyncStub = sc.createAsyncStub(channel);
        blockingStub = sc.createBlockingStub(channel);
    }

    @Override
    public B getBlockingStub() {
        return blockingStub;
    }

    @Override
    public ManagedChannel[] getChannels() {
        ManagedChannel[] chs = new ManagedChannel[1];
        chs[0] = channel;
        return chs;
    }

    @Override
    public A getAsyncStub() {
        return asyncStub;
    }

    @Override
    public void destroy() {
        // 1. get reference
        ManagedChannel ch = this.channel;
        // 2. shutdown
        if (ch != null) {
            ch.shutdown();
            this.channel = null;
        }
    }
}
