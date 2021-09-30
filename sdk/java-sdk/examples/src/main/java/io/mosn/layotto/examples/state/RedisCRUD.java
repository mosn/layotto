package io.mosn.layotto.examples.state;

import io.mosn.layotto.v1.RuntimeClientBuilder;
import io.mosn.layotto.v1.config.RuntimeProperties;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.state.State;

public class RedisCRUD {

    public static void main(String[] args) {
        RuntimeClient client = new RuntimeClientBuilder()
                .withPort(RuntimeProperties.PORT.get())
                .build();

        client.saveState("redis", "key1", "v11");
        State<String> state = client.getState("redis", "key1", String.class);
        System.out.println("get state key:" + state.getKey() + "  value:" + state.getValue());
        client.deleteState("redis", "key1");
        state = client.getState("redis", "key1", String.class);
        System.out.println("get state after delete. key:" + state.getKey() + "  value:" + state.getValue());
    }
}
