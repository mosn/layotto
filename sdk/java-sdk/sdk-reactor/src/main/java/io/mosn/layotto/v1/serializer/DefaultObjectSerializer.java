/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
