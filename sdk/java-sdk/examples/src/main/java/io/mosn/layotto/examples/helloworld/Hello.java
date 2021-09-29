package io.mosn.layotto.examples.helloworld;

import io.mosn.layotto.v1.RuntimeClientBuilder;
import spec.sdk.runtime.v1.client.RuntimeClient;

public class Hello {
    public static void main(String[] args) {
        RuntimeClient client = new RuntimeClientBuilder().withPort(34904).build();
        String resp = client.sayHello("helloworld");
        System.out.println(resp);
    }
}
