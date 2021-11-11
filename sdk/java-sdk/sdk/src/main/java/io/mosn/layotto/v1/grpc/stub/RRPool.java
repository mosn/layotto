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

import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * A utility to generate the next index in a pool
 *
 * @param <T>
 */
public class RRPool<T> {
    private final List<T> stubs;
    private final RRIndex idx;

    /**
     * Construct a RRPool
     *
     * @param stubs should be a concurrent safe list
     */
    public RRPool(List<T> stubs) {
        this.stubs = stubs;
        this.idx = new RRIndex(stubs.size());
    }

    public T next() {
        return stubs.get(idx.next());
    }

    static class RRIndex {
        private final static int mask = 0x7FFFFFFF;
        AtomicInteger idx;
        private final int size;

        RRIndex(int size) {
            this.idx = new AtomicInteger(-1);
            this.size = size;
        }

        RRIndex(int idx, int size) {
            this.idx = new AtomicInteger(idx);
            this.size = size;
        }

        int next() {
            if (size == 1) {
                return 0;
            }
            int n = idx.incrementAndGet();
            n = n & mask;
            n = n % size;
            return n;
        }
    }
}
