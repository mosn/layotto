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
package spec.sdk.runtime.v1.domain.pubsub;

// TopicEventResponseStatus allows apps to have finer control over handling of the message.
public enum TopicEventResponseStatus {
    // SUCCESS is the default behavior: message is acknowledged and not retried or logged.
    SUCCESS(0),
    // RETRY status signals runtime to retry the message as part of an expected scenario (no warning is logged).
    RETRY(1),
    // DROP status signals runtime to drop the message as part of an unexpected scenario (warning is logged).
    DROP(2);

    int idx;

    TopicEventResponseStatus(int idx) {
        this.idx = idx;
    }

    /**
     * Getter method for property <tt>idx</tt>.
     *
     * @return property value of idx
     */
    public int getIdx() {
        return idx;
    }
}
