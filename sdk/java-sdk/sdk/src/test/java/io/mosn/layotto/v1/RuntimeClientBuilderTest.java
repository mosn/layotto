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

}
