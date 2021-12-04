package spec.sdk.runtime.v1.domain;

public interface HelloRuntimes {

    String sayHello(String name);

    String sayHello(String name, int timeout);
}
