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
// CODE ATTRIBUTION: https://github.com/dapr/java-sdk
// Modified some test cases to test layotto's code
package io.mosn.layotto.v1.serializer;

import org.junit.Assert;
import org.junit.Test;
import spec.proto.runtime.v1.RuntimeProto;

import java.io.IOException;
import java.io.Serializable;
import java.util.Base64;
import java.util.List;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

public class JSONSerializerTest {

    private static final ObjectSerializer SERIALIZER = new JSONSerializer();

    @Test
    public void getContentType() {
        assertEquals(new JSONSerializer().getContentType(), "application/json");
    }

    public static class MyObjectTestToSerialize implements Serializable {
        private String  stringValue;
        private int     intValue;
        private boolean boolValue;
        private char    charValue;
        private byte    byteValue;
        private short   shortValue;
        private long    longValue;
        private float   floatValue;
        private double  doubleValue;

        public String getStringValue() {
            return stringValue;
        }

        public void setStringValue(String stringValue) {
            this.stringValue = stringValue;
        }

        public int getIntValue() {
            return intValue;
        }

        public void setIntValue(int intValue) {
            this.intValue = intValue;
        }

        public boolean isBoolValue() {
            return boolValue;
        }

        public void setBoolValue(boolean boolValue) {
            this.boolValue = boolValue;
        }

        public char getCharValue() {
            return charValue;
        }

        public void setCharValue(char charValue) {
            this.charValue = charValue;
        }

        public byte getByteValue() {
            return byteValue;
        }

        public void setByteValue(byte byteValue) {
            this.byteValue = byteValue;
        }

        public short getShortValue() {
            return shortValue;
        }

        public void setShortValue(short shortValue) {
            this.shortValue = shortValue;
        }

        public long getLongValue() {
            return longValue;
        }

        public void setLongValue(long longValue) {
            this.longValue = longValue;
        }

        public float getFloatValue() {
            return floatValue;
        }

        public void setFloatValue(float floatValue) {
            this.floatValue = floatValue;
        }

        public double getDoubleValue() {
            return doubleValue;
        }

        public void setDoubleValue(double doubleValue) {
            this.doubleValue = doubleValue;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) {
                return true;
            }
            if (!(o instanceof MyObjectTestToSerialize)) {
                return false;
            }

            MyObjectTestToSerialize that = (MyObjectTestToSerialize) o;

            if (getIntValue() != that.getIntValue()) {
                return false;
            }
            if (isBoolValue() != that.isBoolValue()) {
                return false;
            }
            if (getCharValue() != that.getCharValue()) {
                return false;
            }
            if (getByteValue() != that.getByteValue()) {
                return false;
            }
            if (getShortValue() != that.getShortValue()) {
                return false;
            }
            if (getLongValue() != that.getLongValue()) {
                return false;
            }
            if (Float.compare(that.getFloatValue(), getFloatValue()) != 0) {
                return false;
            }
            if (Double.compare(that.getDoubleValue(), getDoubleValue()) != 0) {
                return false;
            }
            if (getStringValue() != null ? !getStringValue().equals(that.getStringValue())
                : that.getStringValue() != null) {
                return false;
            }

            return true;
        }

        @Override
        public int hashCode() {
            int result;
            long temp;
            result = getStringValue() != null ? getStringValue().hashCode() : 0;
            result = 31 * result + getIntValue();
            result = 31 * result + (isBoolValue() ? 1 : 0);
            result = 31 * result + (int) getCharValue();
            result = 31 * result + (int) getByteValue();
            result = 31 * result + (int) getShortValue();
            result = 31 * result + (int) (getLongValue() ^ (getLongValue() >>> 32));
            result = 31 * result + (getFloatValue() != +0.0f ? Float.floatToIntBits(getFloatValue()) : 0);
            temp = Double.doubleToLongBits(getDoubleValue());
            result = 31 * result + (int) (temp ^ (temp >>> 32));
            return result;
        }

        @Override
        public String toString() {
            return "MyObjectTestToSerialize{" +
                "stringValue='" + stringValue + '\'' +
                ", intValue=" + intValue +
                ", boolValue=" + boolValue +
                ", charValue=" + charValue +
                ", byteValue=" + byteValue +
                ", shortValue=" + shortValue +
                ", longValue=" + longValue +
                ", floatValue=" + floatValue +
                ", doubleValue=" + doubleValue +
                '}';
        }
    }

    @Test
    public void serializeStringObjectTest() {
        MyObjectTestToSerialize obj = new MyObjectTestToSerialize();
        obj.setStringValue("A String");
        obj.setIntValue(2147483647);
        obj.setBoolValue(true);
        obj.setCharValue('a');
        obj.setByteValue((byte) 65);
        obj.setShortValue((short) 32767);
        obj.setLongValue(9223372036854775807L);
        obj.setFloatValue(1.0f);
        obj.setDoubleValue(1000.0);
        String expectedResult = "{\"boolValue\":true,\"byteValue\":65,\"charValue\":\"a\",\"doubleValue\":1000.0,\"floatValue\":1.0,"
            + "\"intValue\":2147483647,\"longValue\":9223372036854775807,\"shortValue\":32767,\"stringValue\":\"A String\"}";

        String serializedValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(obj));
            assertEquals("FOUND:[[" + serializedValue + "]] \n but was EXPECTING: [[" + expectedResult + "]]",
                expectedResult,
                serializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeObjectTest() {
        MyObjectTestToSerialize obj = new MyObjectTestToSerialize();
        obj.setStringValue("A String");
        obj.setIntValue(2147483647);
        obj.setBoolValue(true);
        obj.setCharValue('a');
        obj.setByteValue((byte) 65);
        obj.setShortValue((short) 32767);
        obj.setLongValue(9223372036854775807L);
        obj.setFloatValue(1.0f);
        obj.setDoubleValue(1000.0);
        //String expectedResult = "{\"stringValue\":\"A String\",\"intValue\":2147483647,\"boolValue\":true,\"charValue\":\"a\",
        // \"byteValue\":65,\"shortValue\":32767,\"longValue\":9223372036854775807,\"floatValue\":1.0,\"doubleValue\":1000.0}";

        byte[] serializedValue;
        try {
            serializedValue = SERIALIZER.serialize(obj);
            Assert.assertNotNull(serializedValue);
            MyObjectTestToSerialize deserializedValue = SERIALIZER.deserialize(serializedValue,
                MyObjectTestToSerialize.class);
            assertEquals(obj, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }

        try {
            serializedValue = SERIALIZER.serialize(obj);
            Assert.assertNotNull(serializedValue);
            MyObjectTestToSerialize deserializedValue = SERIALIZER.deserialize(serializedValue,
                MyObjectTestToSerialize.class);
            assertEquals(obj, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeNullTest() {

        byte[] byteSerializedValue;
        try {
            byteSerializedValue = SERIALIZER.serialize(null);
            Assert.assertNull(byteSerializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeEmptyByteArrayTest() {

        byte[] byteSerializedValue;
        try {
            byteSerializedValue = SERIALIZER.serialize(new byte[] {});
            Assert.assertTrue(byteSerializedValue != null && byteSerializedValue.length == 0);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
        try {
            byteSerializedValue = SERIALIZER.deserialize(null, byte[].class);
            Assert.assertNull(byteSerializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
        try {
            byteSerializedValue = SERIALIZER.deserialize(new byte[] {}, byte[].class);
            Assert.assertTrue(byteSerializedValue != null && byteSerializedValue.length == 0);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
        try {
            MyObjectTestToSerialize de = SERIALIZER.deserialize(new byte[] {}, MyObjectTestToSerialize.class);
            Assert.assertNull(de);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
        try {
            MyObjectTestToSerialize de = SERIALIZER.deserialize(null, MyObjectTestToSerialize.class);
            Assert.assertNull(de);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }

    }

    @Test
    public void serializeStringTest() {
        String valueToSerialize = "A String";

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize));
            assertEquals(valueToSerialize, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            String deserializedValue = SERIALIZER.deserialize(byteValue, String.class);
            assertEquals(valueToSerialize, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeIntTest() {
        Integer valueToSerialize = 1;
        String expectedResult = valueToSerialize.toString();

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize.intValue()));
            assertEquals(expectedResult, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            Integer deserializedValue = SERIALIZER.deserialize(byteValue, Integer.class);
            assertEquals(valueToSerialize, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeShortTest() {
        Short valueToSerialize = 1;
        String expectedResult = valueToSerialize.toString();

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize.shortValue()));
            assertEquals(expectedResult, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            Short deserializedValue = SERIALIZER.deserialize(byteValue, Short.class);
            assertEquals(valueToSerialize, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeLongTest() {
        Long valueToSerialize = Long.MAX_VALUE;
        String expectedResult = valueToSerialize.toString();

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize.longValue()));
            assertEquals(expectedResult, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            Long deserializedValue = SERIALIZER.deserialize(byteValue, Long.class);
            assertEquals(valueToSerialize, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeFloatTest() {
        Float valueToSerialize = -1.23456f;
        String expectedResult = valueToSerialize.toString();

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize.floatValue()));
            assertEquals(expectedResult, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            Float deserializedValue = SERIALIZER.deserialize(byteValue, Float.class);
            assertEquals(valueToSerialize, deserializedValue, 0.00000000001);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeDoubleTest() {
        Double valueToSerialize = 1.0;
        String expectedResult = valueToSerialize.toString();

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize.doubleValue()));
            assertEquals(expectedResult, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            Double deserializedValue = SERIALIZER.deserialize(byteValue, Double.class);
            assertEquals(valueToSerialize, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void serializeBooleanTest() {
        Boolean valueToSerialize = true;
        String expectedResult = valueToSerialize.toString();

        String serializedValue;
        byte[] byteValue;
        try {
            serializedValue = new String(SERIALIZER.serialize(valueToSerialize.booleanValue()));
            assertEquals(expectedResult, serializedValue);
            byteValue = SERIALIZER.serialize(valueToSerialize);
            Assert.assertNotNull(byteValue);
            Boolean deserializedValue = SERIALIZER.deserialize(byteValue, Boolean.class);
            assertEquals(valueToSerialize, deserializedValue);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

    @Test
    public void deserializeObjectTest() {
        String jsonToDeserialize = "{\"stringValue\":\"A String\",\"intValue\":2147483647,\"boolValue\":true,\"charValue\":\"a\",\"byteValue\":65,"
            + "\"shortValue\":32767,\"longValue\":9223372036854775807,\"floatValue\":1.0,\"doubleValue\":1000.0}";
        MyObjectTestToSerialize expectedResult = new MyObjectTestToSerialize();
        expectedResult.setStringValue("A String");
        expectedResult.setIntValue(2147483647);
        expectedResult.setBoolValue(true);
        expectedResult.setCharValue('a');
        expectedResult.setByteValue((byte) 65);
        expectedResult.setShortValue((short) 32767);
        expectedResult.setLongValue(9223372036854775807L);
        expectedResult.setFloatValue(1.0f);
        expectedResult.setDoubleValue(1000.0);
        MyObjectTestToSerialize result;

        try {
            result = SERIALIZER.deserialize(jsonToDeserialize.getBytes(), MyObjectTestToSerialize.class);
            assertEquals("The expected value is different than the actual result", expectedResult, result);
        } catch (IOException exception) {
            fail(exception.getMessage());
        }
    }

}
