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

package info

// Contributor holds info and err for Endpoint.
type Contributor interface {
	GetInfo() (info interface{}, err error)
}

type ContributorAdapter func() (interface{}, error)

func (ca ContributorAdapter) GetInfo() (interface{}, error) {
	return ca()
}
