package io.mosn.layotto.v1.serializer;

import java.io.IOException;

/**
 * Serializes and deserializes application's objects.
 */
public interface ObjectSerializer {

    /**
     * Serializes the given object as byte[].
     *
     * @param o Object to be serialized.
     * @return Serialized object.
     * @throws IOException If cannot serialize.
     */
    byte[] serialize(Object o) throws IOException;

    /**
     * Deserializes the given byte[] into a object.
     *
     * @param data Data to be deserialized.
     * @param clazz Class of object to be deserialized.
     * @param <T>  Type of object to be deserialized.
     * @return Deserialized object.
     * @throws IOException If cannot deserialize object.
     */
    public <T> T deserialize(byte[] data, Class<T> clazz) throws IOException;

    /**
     * Returns the content type of the request.
     *
     * @return content type of the request
     */
    String getContentType();
}
