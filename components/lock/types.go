//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package lock

type Feature string

type Config struct {
	Metadata map[string]string `json:"metadata"`
}

type Metadata struct {
	Properties map[string]string `json:"properties"`
}

type TryLockRequest struct {
	ResourceId string
	LockOwner  string
	Expire     int32
}

type TryLockResponse struct {
	Success bool
}

type UnlockRequest struct {
	ResourceId string
	LockOwner  string
}

type UnlockResponse struct {
	Status LockStatus
}

type LockStatus int32

const (
	SUCCESS               LockStatus = 0
	LOCK_UNEXIST          LockStatus = 1
	LOCK_BELONG_TO_OTHERS LockStatus = 2
	INTERNAL_ERROR        LockStatus = 3
)
