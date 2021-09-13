/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.state;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

public class StateOptions {

    private final Consistency consistency;
    private final Concurrency concurrency;

    /**
     * Represents options for a CloudRuntimes state API call.
     *
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

    /**
     * Returns state options as a Map of option name to value.
     *
     * @return A map of state options.
     */
    public Map<String, String> getStateOptionsAsMap() {
        Map<String, String> mapOptions = new HashMap<>();
        if (this.getConsistency() != null) {
            mapOptions.put("consistency", this.getConsistency().getValue());
        }
        if (this.getConcurrency() != null) {
            mapOptions.put("concurrency", this.getConcurrency().getValue());
        }
        return Collections.unmodifiableMap(Optional.ofNullable(mapOptions).orElse(Collections.EMPTY_MAP));
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
