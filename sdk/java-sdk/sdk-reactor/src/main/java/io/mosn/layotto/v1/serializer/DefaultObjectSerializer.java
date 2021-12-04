/*
 * Copyright (c) Microsoft Corporation and Dapr Contributors.
 * Licensed under the MIT License.
 */

package io.mosn.layotto.v1.serializer;


import spec.sdk.reactor.v1.utils.TypeRef;

import java.io.IOException;

/**
 * Default serializer/deserializer for request/response objects and for state objects too.
 */
public class DefaultObjectSerializer extends ObjectSerializer implements LayottoObjectSerializer {

    /**
     * {@inheritDoc}
     */
    @Override
    public byte[] serialize(Object o) throws IOException {
        return super.serialize(o);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> T deserialize(byte[] data, TypeRef<T> type) throws IOException {
        return super.deserialize(data, type);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getContentType() {
        return "application/json";
    }
}
