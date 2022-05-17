/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ref

//RefConfig is ref json config
type RefConfig struct {
	ComponentRef     []*ComponentRef
	SecretRef        []*RefItem
	ConfigurationRef []*RefItem
}

type ComponentRef struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type RefItem struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
