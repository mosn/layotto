package spec.sdk.runtime.v1;

public interface HelloClient {

    String sayHello(String name);

    String sayHello(String name, int timeout);
}
