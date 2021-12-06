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

import java.io.IOException;

public abstract class AbstractSerializer implements ObjectSerializer {
    /**
     * {@inheritDoc}
     */
    @Override
    public byte[] serialize(Object o) throws IOException {
        if (o == null) {
            return null;
        }
        if (o instanceof byte[]) {
            return (byte[]) o;
        }
        if (o instanceof String) {
            return ((String) o).getBytes();
        }

        return doSerialize(o);
    }

    protected abstract byte[] doSerialize(Object o) throws IOException;

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> T deserialize(byte[] data, Class<T> clazz) throws IOException {
        if (data == null) {
            return null;
        }
        if (clazz == byte[].class) {
            return (T) data;
        }
        if (clazz == String.class) {
            return (T) new String(data);
        }
        return doDeserialize(data, clazz);
    }

    protected abstract <T> T doDeserialize(byte[] data, Class<T> clazz) throws IOException;

}