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

package state

const (
	// StateOptions_CONSISTENCY
	// Unspecified
	StateOptions_CONSISTENCY_UNSPECIFIED int32 = 0
	//  The API server assumes data stores are eventually consistent by default.A state store should:
	// - For read requests, the state store can return data from any of the replicas
	// - For write request, the state store should asynchronously replicate updates to configured quorum after acknowledging the update request.
	StateOptions_CONSISTENCY_EVENTUAL int32 = 1
	// When a strong consistency hint is attached, a state store should:
	// - For read requests, the state store should return the most up-to-date data consistently across replicas.
	// - For write/delete requests, the state store should synchronisely replicate updated data to configured quorum before completing the write request.
	StateOptions_CONSISTENCY_STRONG int32 = 2

	// StateOptions_CONCURRENCY
	// Unspecified
	StateOptions_CONCURRENCY_UNSPECIFIED int32 = 0
	// First write wins
	StateOptions_CONCURRENCY_FIRST_WRITE int32 = 1
	// Last write wins
	StateOptions_CONCURRENCY_LAST_WRITE int32 = 2
)

func StateConsistencyToString(c int32) string {
	switch c {
	case StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}
	return ""
}

func StateConcurrencyToString(c int32) string {
	switch c {
	case StateOptions_CONCURRENCY_FIRST_WRITE:
		return "first-write"
	case StateOptions_CONCURRENCY_LAST_WRITE:
		return "last-write"
	}

	return ""
}
