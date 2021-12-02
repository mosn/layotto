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
package io.mosn.layotto.v1.config;

/**
 * String configuration property.
 */
public class StringProperty extends Property<String> {

    /**
     * {@inheritDoc}
     */
    StringProperty(String name, String envName, String defaultValue) {
        super(name, envName, defaultValue);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    protected String parse(String value) {
        return value;
    }
}
