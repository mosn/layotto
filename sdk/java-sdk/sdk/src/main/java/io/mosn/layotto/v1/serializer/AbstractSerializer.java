package io.mosn.layotto.v1.serializer;

import com.alibaba.fastjson.JSONObject;

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