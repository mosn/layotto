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
package spec.sdk.runtime.v1.domain.lock;

enum UnlockResponseStatus {
    SUCCESS(0),
    LOCK_UNEXIST(1),
    LOCK_BELONG_TO_OTHERS(2),
    INTERNAL_ERROR(3);
    private final int value;

    UnlockResponseStatus(int value) {
        this.value = value;
    }
}

public class UnlockResponse {
    public UnlockResponseStatus status;
}
