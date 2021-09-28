package io.mosn.layotto.v1.serializer;

import com.alibaba.fastjson.JSONObject;

import java.io.IOException;

public class JSONSerializer implements ObjectSerializer {

    /**
     * {@inheritDoc}
     */
    @Override
    public byte[] serialize(Object o) throws IOException {
        return JSONObject.toJSONBytes(o);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public <T> T deserialize(byte[] data, Class<T> clazz) throws IOException {
        return JSONObject.parseObject(data, clazz);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getContentType() {
        return "application/json";
    }
}
