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
package spec.sdk.runtime.v1.domain.file;

import java.util.HashMap;
import java.util.Map;

public class GetMeteResponse {
    private long                  size;
    private String                lastModified;
    private Map<String, String[]> meta;

    public long getSize() {
        return size;
    }

    public void setSize(long size) {
        this.size = size;
    }

    public String getLastModified() {
        return lastModified;
    }

    public void setLastModified(String lastModified) {
        this.lastModified = lastModified;
    }

    public Map<String, String[]> getMeta() {

        if (meta == null) {
            meta = new HashMap<>();
        }

        return meta;
    }

    public void setMeta(Map<String, String[]> meta) {
        this.meta = meta;
    }
}
