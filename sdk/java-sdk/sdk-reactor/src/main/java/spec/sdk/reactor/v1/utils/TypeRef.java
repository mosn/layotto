/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.reactor.v1.utils;

import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;

/**
 * Used to reference a type.
 *
 * <p>Usage: new TypeRef&lt;MyClass&gt;(){}</p>
 *
 * @param <T> Type to be deserialized.
 */
public abstract class TypeRef<T> {

    public static final TypeRef<String> STRING = new TypeRef<String>() {
    };

    public static final TypeRef<Boolean> BOOLEAN = new TypeRef<Boolean>(boolean.class) {
    };

    public static final TypeRef<Integer> INT = new TypeRef<Integer>(int.class) {
    };

    public static final TypeRef<Long> LONG = new TypeRef<Long>(long.class) {
    };

    public static final TypeRef<Character> CHAR = new TypeRef<Character>(char.class) {
    };

    public static final TypeRef<Byte> BYTE = new TypeRef<Byte>(byte.class) {
    };

    public static final TypeRef<Void> VOID = new TypeRef<Void>(void.class) {
    };

    public static final TypeRef<Float> FLOAT = new TypeRef<Float>(float.class) {
    };

    public static final TypeRef<Double> DOUBLE = new TypeRef<Double>(double.class) {
    };

    public static final TypeRef<byte[]> BYTE_ARRAY = new TypeRef<byte[]>() {
    };

    public static final TypeRef<int[]> INT_ARRAY = new TypeRef<int[]>() {
    };

    public static final TypeRef<String[]> STRING_ARRAY = new TypeRef<String[]>() {
    };

    private final Type type;

    /**
     * Constructor.
     */
    public TypeRef() {
        Type superClass = this.getClass().getGenericSuperclass();
        if (superClass instanceof Class) {
            throw new IllegalArgumentException("TypeReference requires type.");
        }

        this.type = ((ParameterizedType) superClass).getActualTypeArguments()[0];
    }

    /**
     * Constructor for reflection.
     *
     * @param type Type to be referenced.
     */
    private TypeRef(Type type) {
        this.type = type;
    }

    /**
     * Gets the type referenced.
     *
     * @return type referenced.
     */
    public Type getType() {
        return this.type;
    }

    /**
     * Creates a reference to a given class type.
     *
     * @param clazz Class type to be referenced.
     * @param <T>   Type to be referenced.
     * @return Class type reference.
     */
    public static <T> TypeRef<T> get(Class<T> clazz) {
        if (clazz == String.class) {
            return (TypeRef<T>) STRING;
        }
        if (clazz == boolean.class) {
            return (TypeRef<T>) BOOLEAN;
        }
        if (clazz == int.class) {
            return (TypeRef<T>) INT;
        }
        if (clazz == long.class) {
            return (TypeRef<T>) LONG;
        }
        if (clazz == char.class) {
            return (TypeRef<T>) CHAR;
        }
        if (clazz == byte.class) {
            return (TypeRef<T>) BYTE;
        }
        if (clazz == void.class) {
            return (TypeRef<T>) VOID;
        }
        if (clazz == float.class) {
            return (TypeRef<T>) FLOAT;
        }
        if (clazz == double.class) {
            return (TypeRef<T>) DOUBLE;
        }
        if (clazz == byte[].class) {
            return (TypeRef<T>) BYTE_ARRAY;
        }
        if (clazz == int[].class) {
            return (TypeRef<T>) INT_ARRAY;
        }
        if (clazz == String[].class) {
            return (TypeRef<T>) STRING_ARRAY;
        }

        return new TypeRef<T>(clazz) {
        };
    }

    /**
     * Creates a reference to a given class type.
     *
     * @param type Type to be referenced.
     * @param <T>  Type to be referenced.
     * @return Class type reference.
     */
    public static <T> TypeRef<T> get(Type type) {
        if (type instanceof Class) {
            Class clazz = (Class) type;
            return get(clazz);
        }

        return new TypeRef<T>(type) {
        };
    }
}
