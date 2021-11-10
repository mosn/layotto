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
 *
 */
// CODE ATTRIBUTION: https://github.com/dapr/java-sdk
// Modified some test cases to test layotto's code
package io.mosn.layotto.v1;

import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.junit.Assert;
import org.junit.Test;
import spec.sdk.runtime.v1.client.RuntimeClient;

import static org.mockito.Mockito.mock;

public class RuntimeClientBuilderTest {

    @Test
    public void build() {
        ObjectSerializer stateSerializer = mock(ObjectSerializer.class);
        RuntimeClientBuilder builder = new RuntimeClientBuilder();
        builder.withStateSerializer(stateSerializer);
        RuntimeClient client = builder.build();
        Assert.assertNotNull(client);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noLogger() {
        new RuntimeClientBuilder().withLogger(null);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noTimeout() {
        new RuntimeClientBuilder().withTimeout(0);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noPort() {
        new RuntimeClientBuilder().withPort(0);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noStateSerializer() {
        new RuntimeClientBuilder().withStateSerializer(null);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noApiProtocal() {
        new RuntimeClientBuilder().withApiProtocol(null);
    }

}
