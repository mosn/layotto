/*
 * Copyright (c) Microsoft Corporation and Dapr Contributors.
 * Licensed under the MIT License.
 */

package io.mosn.layotto.v1.config;

import java.util.function.Function;

/**
 * Configuration property for any type.
 */
public class GenericProperty<T> extends Property<T> {

    private final Function<String, T> parser;

    /**
     * {@inheritDoc}
     */
    GenericProperty(String name, String envName, T defaultValue, Function<String, T> parser) {
        super(name, envName, defaultValue);
        this.parser = parser;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    protected T parse(String value) {
        return parser.apply(value);
    }
}
