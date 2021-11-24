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
package spec.sdk.runtime.v1.domain.state;

public class StateOptions {
    private final Consistency consistency;
    private final Concurrency concurrency;

    /**
     * Represents options for a state API call.
     * @param consistency The consistency mode.
     * @param concurrency The concurrency mode.
     */
    public StateOptions(Consistency consistency, Concurrency concurrency) {
        this.consistency = consistency;
        this.concurrency = concurrency;
    }

    public Concurrency getConcurrency() {
        return concurrency;
    }

    public Consistency getConsistency() {
        return consistency;
    }

    public enum Consistency {
        EVENTUAL("eventual"),
        STRONG("strong");

        private final String value;

        Consistency(String value) {
            this.value = value;
        }

        public String getValue() {
            return this.value;
        }

        public static Consistency fromValue(String value) {
            return Consistency.valueOf(value);
        }
    }

    public enum Concurrency {
        FIRST_WRITE("first-write"),
        LAST_WRITE("last-write");

        private final String value;

        Concurrency(String value) {
            this.value = value;
        }

        public String getValue() {
            return this.value;
        }

        public static Concurrency fromValue(String value) {
            return Concurrency.valueOf(value);
        }
    }
}
