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

import io.mosn.layotto.v1.grpc.stub.RRPool;
import org.junit.Assert;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

@RunWith(JUnit4.class)
public class RRPoolTest {

    @Test
    public void testRR() {
        RRPool.RRIndex idx = new RRPool.RRIndex(1);
        Assert.assertEquals(0, idx.next());
        Assert.assertEquals(0, idx.next());
        Assert.assertEquals(0, idx.next());

        int count = 32;
        idx = new RRPool.RRIndex(count);
        for (int i = 0; i < count; i++) {
            Assert.assertEquals(i, idx.next());
        }
        for (int i = 0; i < count; i++) {
            Assert.assertEquals(i, idx.next());
        }
        for (int i = 0; i < count; i++) {
            Assert.assertEquals(i, idx.next());
        }
    }

    @Test
    public void testRROverflow() {
        RRPool.RRIndex idx = new RRPool.RRIndex(Integer.MAX_VALUE, 10);
        Assert.assertEquals(0, idx.next());
        Assert.assertEquals(1, idx.next());
        Assert.assertEquals(2, idx.next());
    }
}
