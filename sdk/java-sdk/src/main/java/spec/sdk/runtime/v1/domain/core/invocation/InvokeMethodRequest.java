/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.invocation;

import java.util.Map;

/**
 * A request to invoke a service.
 */
public class InvokeMethodRequest {

    private final String appId;

    private final String method;

    private Object body;

    private spec.sdk.runtime.v1.domain.core.invocation.HttpExtension httpExtension;

    private String contentType;

    private Map<String, String> metadata;

    public InvokeMethodRequest(String appId, String method) {
        this.appId = appId;
        this.method = method;
    }

    public String getAppId() {
        return appId;
    }

    public String getMethod() {
        return method;
    }

    public Object getBody() {
        return body;
    }

    public InvokeMethodRequest setBody(Object body) {
        this.body = body;
        return this;
    }

    public spec.sdk.runtime.v1.domain.core.invocation.HttpExtension getHttpExtension() {
        return httpExtension;
    }

    public InvokeMethodRequest setHttpExtension(HttpExtension httpExtension) {
        this.httpExtension = httpExtension;
        return this;
    }

    public String getContentType() {
        return contentType;
    }

    public InvokeMethodRequest setContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }
}
