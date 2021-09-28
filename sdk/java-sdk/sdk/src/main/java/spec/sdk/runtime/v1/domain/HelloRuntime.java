package spec.sdk.runtime.v1.domain;

public interface HelloRuntime {

    String sayHello(String name);

    String sayHello(String name, int timeoutMillisecond);
}
