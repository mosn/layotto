package io.mosn.layotto.v1.domain;

public class SayHelloRequest {

    public String serviceName;
    public String name;

    public String getServiceName() {
        return serviceName;
    }

    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
