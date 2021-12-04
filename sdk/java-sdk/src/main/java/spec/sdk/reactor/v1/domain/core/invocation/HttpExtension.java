/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.reactor.v1.domain.core.invocation;

import okhttp3.HttpUrl;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Optional;

/**
 * HTTP Extension class.
 * This class is only needed if the app you are calling is listening on HTTP.
 * It contains properties that represent data that may be populated for an HTTP receiver.
 */
public final class HttpExtension {

    /**
     * Convenience HttpExtension object for {@link HttpMethods#NONE} with empty queryString.
     */
    public static final HttpExtension NONE = new HttpExtension(HttpMethods.NONE);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#GET} Verb with empty queryString.
     */
    public static final HttpExtension GET = new HttpExtension(HttpMethods.GET);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#PUT} Verb with empty queryString.
     */
    public static final HttpExtension PUT = new HttpExtension(HttpMethods.PUT);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#POST} Verb with empty queryString.
     */
    public static final HttpExtension POST = new HttpExtension(HttpMethods.POST);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#DELETE} Verb with empty queryString.
     */
    public static final HttpExtension DELETE = new HttpExtension(HttpMethods.DELETE);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#HEAD} Verb with empty queryString.
     */
    public static final HttpExtension HEAD = new HttpExtension(HttpMethods.HEAD);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#CONNECT} Verb with empty queryString.
     */
    public static final HttpExtension CONNECT = new HttpExtension(HttpMethods.CONNECT);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#OPTIONS} Verb with empty queryString.
     */
    public static final HttpExtension OPTIONS = new HttpExtension(HttpMethods.OPTIONS);
    /**
     * Convenience HttpExtension object for the {@link HttpMethods#TRACE} Verb with empty queryString.
     */
    public static final HttpExtension TRACE = new HttpExtension(HttpMethods.TRACE);

    /**
     * HTTP verb.
     */
    private HttpMethods method;

    /**
     * HTTP query params.
     */
    private Map<String, List<String>> queryParams;

    /**
     * HTTP headers.
     */
    private Map<String, String> headers;

    /**
     * Construct a HttpExtension object.
     *
     * @param method      Required value denoting the HttpMethod.
     * @param queryParams map for the query parameters the HTTP call.
     * @param headers     map to set HTTP headers.
     * @throws IllegalArgumentException on null method or queryString.
     * @see HttpMethods for supported methods.
     */
    public HttpExtension(HttpMethods method,
                         Map<String, List<String>> queryParams,
                         Map<String, String> headers) {
        if (method == null) {
            throw new IllegalArgumentException("HttpExtension method cannot be null");
        }

        this.method = method;
        this.queryParams = Collections.unmodifiableMap(queryParams == null ? Collections.emptyMap() : queryParams);
        this.headers = Collections.unmodifiableMap(headers == null ? Collections.emptyMap() : headers);
    }

    /**
     * Construct a HttpExtension object.
     *
     * @param method Required value denoting the HttpMethod.
     * @throws IllegalArgumentException on null method or queryString.
     * @see HttpMethods for supported methods.
     */
    public HttpExtension(HttpMethods method) {
        this(method, null, null);
    }

    public HttpMethods getMethod() {
        return method;
    }

    public Map<String, List<String>> getQueryParams() {
        return queryParams;
    }

    public Map<String, String> getHeaders() {
        return headers;
    }

    /**
     * Encodes the query string for the HTTP request.
     *
     * @return Encoded HTTP query string.
     */
    public String encodeQueryString() {
        if ((this.queryParams == null) || (this.queryParams.isEmpty())) {
            return "";
        }

        HttpUrl.Builder urlBuilder = new HttpUrl.Builder();
        // Setting required values but we only need query params in the end.
        urlBuilder.scheme("http").host("localhost");
        Optional.ofNullable(this.queryParams).orElse(Collections.emptyMap()).entrySet().stream()
                .forEach(urlParameter ->
                        Optional.ofNullable(urlParameter.getValue()).orElse(Collections.emptyList()).stream()
                                .forEach(urlParameterValue ->
                                        urlBuilder.addQueryParameter(urlParameter.getKey(), urlParameterValue)));
        return urlBuilder.build().encodedQuery();
    }

    /**
     * HTTP Methods supported.
     */
    public enum HttpMethods {
        NONE,
        GET,
        PUT,
        POST,
        DELETE,
        HEAD,
        CONNECT,
        OPTIONS,
        TRACE
    }
}
