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
package spec.sdk.reactor.v1;

import spec.sdk.reactor.v1.domain.core.ConfigurationRuntimes;
import spec.sdk.reactor.v1.domain.core.InvocationRuntimes;
import spec.sdk.reactor.v1.domain.core.PubSubRuntimes;
import spec.sdk.reactor.v1.domain.core.StateRuntimes;

/**
 * Core Cloud Runtimes standard API defined.
 */
public interface CoreCloudRuntimes extends
                                  InvocationRuntimes,
                                  PubSubRuntimes,
                                  // BindingRuntimes,
                                  StateRuntimes,
                                  // SecretsRuntimes,
                                  ConfigurationRuntimes {
}
